package database

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const batchSize = 500

// SeedWilayah membaca file CSV wilayah Indonesia dan memasukkan data ke database secara streaming dan transaksional.
func SeedWilayah(ctx context.Context, db *sql.DB, csvBaseDir string) error {
	slog.Info("Memulai seeding data Wilayah Indonesia...")

	steps := []struct {
		name      string
		tableName string
		columns   string
		csvFile   string
		mapper    func([]string) ([]any, error)
	}{
		{
			name:      "Provinsi",
			tableName: "provinces",
			columns:   "(id, name)",
			csvFile:   "province.csv",
			mapper: func(record []string) ([]any, error) {
				if len(record) < 2 {
					return nil, fmt.Errorf("kolom kurang dari 2")
				}
				id, err := strconv.Atoi(strings.TrimSpace(record[0]))
				if err != nil {
					return nil, err
				}
				name := strings.TrimSpace(record[1])
				return []any{id, name}, nil
			},
		},
		{
			name:      "Kota/Kabupaten",
			tableName: "cities",
			columns:   "(id, name, province_id)",
			csvFile:   "city.csv",
			mapper: func(record []string) ([]any, error) {
				if len(record) < 3 {
					return nil, fmt.Errorf("kolom kurang dari 3")
				}
				id, err := strconv.Atoi(strings.TrimSpace(record[0]))
				if err != nil {
					return nil, err
				}
				name := strings.TrimSpace(record[1])
				provID, err := strconv.Atoi(strings.TrimSpace(record[2]))
				if err != nil {
					return nil, err
				}
				return []any{id, name, provID}, nil
			},
		},
		{
			name:      "Kecamatan",
			tableName: "districts",
			columns:   "(id, name, city_id)",
			csvFile:   "district.csv",
			mapper: func(record []string) ([]any, error) {
				if len(record) < 3 {
					return nil, fmt.Errorf("kolom kurang dari 3")
				}
				id, err := strconv.Atoi(strings.TrimSpace(record[0]))
				if err != nil {
					return nil, err
				}
				name := strings.TrimSpace(record[1])
				cityID, err := strconv.Atoi(strings.TrimSpace(record[2]))
				if err != nil {
					return nil, err
				}
				return []any{id, name, cityID}, nil
			},
		},
		{
			name:      "Kelurahan/Desa",
			tableName: "subdistricts",
			columns:   "(id, name, district_id)",
			csvFile:   "subdistrict.csv",
			mapper: func(record []string) ([]any, error) {
				if len(record) < 3 {
					return nil, fmt.Errorf("kolom kurang dari 3")
				}
				id, err := strconv.Atoi(strings.TrimSpace(record[0]))
				if err != nil {
					return nil, err
				}
				name := strings.TrimSpace(record[1])
				distID, err := strconv.Atoi(strings.TrimSpace(record[2]))
				if err != nil {
					return nil, err
				}
				return []any{id, name, distID}, nil
			},
		},
		{
			name:      "Kode Pos",
			tableName: "postal_codes",
			columns:   "(id, subdistrict_id, postal_code)",
			csvFile:   "postal_code.csv",
			mapper: func(record []string) ([]any, error) {
				if len(record) < 3 {
					return nil, fmt.Errorf("kolom kurang dari 3")
				}
				id, err := strconv.Atoi(strings.TrimSpace(record[0]))
				if err != nil {
					return nil, err
				}
				subdistID, err := strconv.Atoi(strings.TrimSpace(record[1]))
				if err != nil {
					return nil, err
				}
				postalCode, err := strconv.Atoi(strings.TrimSpace(record[2]))
				if err != nil {
					return nil, err
				}
				return []any{id, subdistID, postalCode}, nil
			},
		},
	}

	for _, step := range steps {
		csvPath := filepath.Join(csvBaseDir, step.csvFile)

		// Cek apakah sudah ada data di tabel
		var exists bool
		query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s LIMIT 1)", step.tableName)
		err := db.QueryRowContext(ctx, query).Scan(&exists)
		if err != nil {
			return fmt.Errorf("gagal cek tabel %s: %w", step.tableName, err)
		}

		if exists {
			slog.Info(fmt.Sprintf("✓ %s: sudah ada data, dilewati.", step.name))
			continue
		}

		inserted, err := genericSeed(ctx, db, step.tableName, step.columns, csvPath, step.mapper)
		if err != nil {
			return fmt.Errorf("gagal seed %s: %w", step.name, err)
		}
		slog.Info(fmt.Sprintf("✓ %s: %d data berhasil dimasukkan.", step.name, inserted))
	}

	slog.Info("Seeding data Wilayah Indonesia selesai!")
	return nil
}

func genericSeed(
	ctx context.Context,
	db *sql.DB,
	tableName string,
	columns string,
	csvPath string,
	mapper func([]string) ([]any, error),
) (int, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return 0, fmt.Errorf("gagal membuka file %s: %w", csvPath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Mulai Transaksi
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("gagal memulai transaksi untuk %s: %w", tableName, err)
	}
	defer tx.Rollback()

	total := 0
	var batchArgs []any
	var placeholders []string

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			slog.Warn("Gagal membaca baris CSV, baris dilewati", "file", filepath.Base(csvPath), "error", err)
			continue
		}

		args, err := mapper(record)
		if err != nil {
			// Skip quietly if it's the header row. But for visibility, we log warning.
			slog.Warn("Gagal parsing baris CSV, baris dilewati", "file", filepath.Base(csvPath), "record", record, "error", err)
			continue
		}

		qs := make([]string, len(args))
		for j := range qs {
			qs[j] = "?"
		}
		placeholders = append(placeholders, "("+strings.Join(qs, ", ")+")")
		batchArgs = append(batchArgs, args...)

		if len(placeholders) >= batchSize {
			if err := insertBatch(ctx, tx, tableName, columns, placeholders, batchArgs); err != nil {
				return total, err
			}
			total += len(placeholders)
			placeholders = nil
			batchArgs = nil
		}
	}

	// Insert sisa batch
	if len(placeholders) > 0 {
		if err := insertBatch(ctx, tx, tableName, columns, placeholders, batchArgs); err != nil {
			return total, err
		}
		total += len(placeholders)
	}

	if err := tx.Commit(); err != nil {
		return total, fmt.Errorf("gagal commit transaksi %s: %w", tableName, err)
	}

	return total, nil
}

func insertBatch(ctx context.Context, tx *sql.Tx, tableName string, columns string, placeholders []string, args []any) error {
	query := fmt.Sprintf("INSERT IGNORE INTO %s %s VALUES %s", tableName, columns, strings.Join(placeholders, ", "))
	
	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("batch insert ke %s gagal: %w", tableName, err)
	}
	return nil
}
