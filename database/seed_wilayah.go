package database

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const batchSize = 500

// SeedWilayah membaca file CSV wilayah Indonesia dan memasukkan data ke database.
// csvBaseDir adalah path ke folder CSV (contoh: "../Wilayah-Indonesia-Beserta-Kode-Pos/CSV")
func SeedWilayah(db *sql.DB, csvBaseDir string) error {
	log.Println("Memulai seeding data Wilayah Indonesia...")

	// Urutan seeding penting karena ada foreign key constraint
	steps := []struct {
		name    string
		csvFile string
		seedFn  func(*sql.DB, string) (int, error)
	}{
		{"Provinsi", "province.csv", seedProvinces},
		{"Kota/Kabupaten", "city.csv", seedCities},
		{"Kecamatan", "district.csv", seedDistricts},
		{"Kelurahan/Desa", "subdistrict.csv", seedSubdistricts},
		{"Kode Pos", "postal_code.csv", seedPostalCodes},
	}

	for _, step := range steps {
		csvPath := filepath.Join(csvBaseDir, step.csvFile)

		// Cek apakah sudah ada data di tabel
		var count int
		tableName := tableNameFromCSV(step.csvFile)
		err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)).Scan(&count)
		if err != nil {
			return fmt.Errorf("gagal cek tabel %s: %w", tableName, err)
		}

		if count > 0 {
			log.Printf("  ✓ %s: sudah ada %d data, dilewati.", step.name, count)
			continue
		}

		inserted, err := step.seedFn(db, csvPath)
		if err != nil {
			return fmt.Errorf("gagal seed %s: %w", step.name, err)
		}
		log.Printf("  ✓ %s: %d data berhasil dimasukkan.", step.name, inserted)
	}

	log.Println("Seeding data Wilayah Indonesia selesai!")
	return nil
}

// tableNameFromCSV mengembalikan nama tabel berdasarkan nama file CSV
func tableNameFromCSV(csvFile string) string {
	switch csvFile {
	case "province.csv":
		return "provinces"
	case "city.csv":
		return "cities"
	case "district.csv":
		return "districts"
	case "subdistrict.csv":
		return "subdistricts"
	case "postal_code.csv":
		return "postal_codes"
	default:
		return ""
	}
}

// seedProvinces memasukkan data provinsi dari CSV
func seedProvinces(db *sql.DB, csvPath string) (int, error) {
	records, err := readCSV(csvPath)
	if err != nil {
		return 0, err
	}

	total := 0
	for i := 0; i < len(records); i += batchSize {
		end := i + batchSize
		if end > len(records) {
			end = len(records)
		}
		batch := records[i:end]

		valueStrings := make([]string, 0, len(batch))
		valueArgs := make([]interface{}, 0, len(batch)*2)

		for _, record := range batch {
			if len(record) < 2 {
				continue
			}
			id, err := strconv.Atoi(strings.TrimSpace(record[0]))
			if err != nil {
				continue
			}
			name := strings.TrimSpace(record[1])

			valueStrings = append(valueStrings, "(?, ?)")
			valueArgs = append(valueArgs, id, name)
		}

		if len(valueStrings) == 0 {
			continue
		}

		query := fmt.Sprintf("INSERT IGNORE INTO provinces (id, name) VALUES %s",
			strings.Join(valueStrings, ", "))

		_, err := db.Exec(query, valueArgs...)
		if err != nil {
			return total, fmt.Errorf("batch insert provinces gagal: %w", err)
		}
		total += len(valueStrings)
	}

	return total, nil
}

// seedCities memasukkan data kota/kabupaten dari CSV
func seedCities(db *sql.DB, csvPath string) (int, error) {
	records, err := readCSV(csvPath)
	if err != nil {
		return 0, err
	}

	total := 0
	for i := 0; i < len(records); i += batchSize {
		end := i + batchSize
		if end > len(records) {
			end = len(records)
		}
		batch := records[i:end]

		valueStrings := make([]string, 0, len(batch))
		valueArgs := make([]interface{}, 0, len(batch)*3)

		for _, record := range batch {
			if len(record) < 3 {
				continue
			}
			id, err := strconv.Atoi(strings.TrimSpace(record[0]))
			if err != nil {
				continue
			}
			name := strings.TrimSpace(record[1])
			provID, err := strconv.Atoi(strings.TrimSpace(record[2]))
			if err != nil {
				continue
			}

			valueStrings = append(valueStrings, "(?, ?, ?)")
			valueArgs = append(valueArgs, id, name, provID)
		}

		if len(valueStrings) == 0 {
			continue
		}

		query := fmt.Sprintf("INSERT IGNORE INTO cities (id, name, province_id) VALUES %s",
			strings.Join(valueStrings, ", "))

		_, err := db.Exec(query, valueArgs...)
		if err != nil {
			return total, fmt.Errorf("batch insert cities gagal: %w", err)
		}
		total += len(valueStrings)
	}

	return total, nil
}

// seedDistricts memasukkan data kecamatan dari CSV
func seedDistricts(db *sql.DB, csvPath string) (int, error) {
	records, err := readCSV(csvPath)
	if err != nil {
		return 0, err
	}

	total := 0
	for i := 0; i < len(records); i += batchSize {
		end := i + batchSize
		if end > len(records) {
			end = len(records)
		}
		batch := records[i:end]

		valueStrings := make([]string, 0, len(batch))
		valueArgs := make([]interface{}, 0, len(batch)*3)

		for _, record := range batch {
			if len(record) < 3 {
				continue
			}
			id, err := strconv.Atoi(strings.TrimSpace(record[0]))
			if err != nil {
				continue
			}
			name := strings.TrimSpace(record[1])
			cityID, err := strconv.Atoi(strings.TrimSpace(record[2]))
			if err != nil {
				continue
			}

			valueStrings = append(valueStrings, "(?, ?, ?)")
			valueArgs = append(valueArgs, id, name, cityID)
		}

		if len(valueStrings) == 0 {
			continue
		}

		query := fmt.Sprintf("INSERT IGNORE INTO districts (id, name, city_id) VALUES %s",
			strings.Join(valueStrings, ", "))

		_, err := db.Exec(query, valueArgs...)
		if err != nil {
			return total, fmt.Errorf("batch insert districts gagal: %w", err)
		}
		total += len(valueStrings)
	}

	return total, nil
}

// seedSubdistricts memasukkan data kelurahan/desa dari CSV
func seedSubdistricts(db *sql.DB, csvPath string) (int, error) {
	records, err := readCSV(csvPath)
	if err != nil {
		return 0, err
	}

	total := 0
	for i := 0; i < len(records); i += batchSize {
		end := i + batchSize
		if end > len(records) {
			end = len(records)
		}
		batch := records[i:end]

		valueStrings := make([]string, 0, len(batch))
		valueArgs := make([]interface{}, 0, len(batch)*3)

		for _, record := range batch {
			if len(record) < 3 {
				continue
			}
			id, err := strconv.Atoi(strings.TrimSpace(record[0]))
			if err != nil {
				continue
			}
			name := strings.TrimSpace(record[1])
			districtID, err := strconv.Atoi(strings.TrimSpace(record[2]))
			if err != nil {
				continue
			}

			valueStrings = append(valueStrings, "(?, ?, ?)")
			valueArgs = append(valueArgs, id, name, districtID)
		}

		if len(valueStrings) == 0 {
			continue
		}

		query := fmt.Sprintf("INSERT IGNORE INTO subdistricts (id, name, district_id) VALUES %s",
			strings.Join(valueStrings, ", "))

		_, err := db.Exec(query, valueArgs...)
		if err != nil {
			return total, fmt.Errorf("batch insert subdistricts gagal: %w", err)
		}
		total += len(valueStrings)
	}

	return total, nil
}

// seedPostalCodes memasukkan data kode pos dari CSV
func seedPostalCodes(db *sql.DB, csvPath string) (int, error) {
	records, err := readCSV(csvPath)
	if err != nil {
		return 0, err
	}

	total := 0
	for i := 0; i < len(records); i += batchSize {
		end := i + batchSize
		if end > len(records) {
			end = len(records)
		}
		batch := records[i:end]

		valueStrings := make([]string, 0, len(batch))
		valueArgs := make([]interface{}, 0, len(batch)*6)

		for _, record := range batch {
			if len(record) < 6 {
				continue
			}
			id, err := strconv.Atoi(strings.TrimSpace(record[0]))
			if err != nil {
				continue
			}
			subdisID, err := strconv.Atoi(strings.TrimSpace(record[1]))
			if err != nil {
				continue
			}
			disID, err := strconv.Atoi(strings.TrimSpace(record[2]))
			if err != nil {
				continue
			}
			cityID, err := strconv.Atoi(strings.TrimSpace(record[3]))
			if err != nil {
				continue
			}
			provID, err := strconv.Atoi(strings.TrimSpace(record[4]))
			if err != nil {
				continue
			}
			postalCode := strings.TrimSpace(record[5])

			valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?)")
			valueArgs = append(valueArgs, id, subdisID, disID, cityID, provID, postalCode)
		}

		if len(valueStrings) == 0 {
			continue
		}

		query := fmt.Sprintf(
			"INSERT IGNORE INTO postal_codes (id, subdistrict_id, district_id, city_id, province_id, postal_code) VALUES %s",
			strings.Join(valueStrings, ", "))

		_, err := db.Exec(query, valueArgs...)
		if err != nil {
			return total, fmt.Errorf("batch insert postal_codes gagal: %w", err)
		}
		total += len(valueStrings)
	}

	return total, nil
}

// readCSV membaca file CSV dan mengembalikan semua record (tanpa header)
func readCSV(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	// Skip header
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("gagal membaca header %s: %w", filePath, err)
	}

	var records [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			// Skip baris yang bermasalah, lanjutkan
			continue
		}
		records = append(records, record)
	}

	return records, nil
}
