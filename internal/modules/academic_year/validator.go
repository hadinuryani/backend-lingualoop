package academic_year

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// parsedDates menyimpan hasil parsing dan validasi seluruh tanggal dari request.
// Struct ini adalah output dari validateAndParseRequest dan digunakan oleh buildEntity.
type parsedDates struct {
	StartDate time.Time
	EndDate   time.Time
	Ganjil    parsedSemester
	Genap     parsedSemester
}

type parsedSemester struct {
	Start      sql.NullTime
	EndKBM     sql.NullTime
	Assessment sql.NullTime
}

var yearFormatRegex = regexp.MustCompile(`^\d{4}/\d{4}$`)

// validateAndParseRequest adalah satu pintu masuk validasi bisnis untuk AcademicYearRequest.
// Mengembalikan parsedDates yang siap dipakai untuk membangun entity, atau error jika ada
// pelanggaran aturan bisnis.
func validateAndParseRequest(req AcademicYearRequest) (*parsedDates, error) {
	// 1. Validasi format Year (harus YYYY/YYYY)
	if !yearFormatRegex.MatchString(req.Year) {
		return nil, ErrInvalidYearFormat
	}

	// 2. Validasi tahun kedua = tahun pertama + 1
	parts := strings.SplitN(req.Year, "/", 2)
	startYear, _ := strconv.Atoi(parts[0])
	endYear, _ := strconv.Atoi(parts[1])

	if endYear != startYear+1 {
		return nil, ErrInvalidYearSequence
	}

	// 3. Parse StartDate dan EndDate
	startDate, err := parseDate(req.StartDate)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}
	endDate, err := parseDate(req.EndDate)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}

	// 4. EndDate harus setelah StartDate
	if endDate.Before(startDate) {
		return nil, ErrInvalidDate
	}

	// 5. StartDate.Year() harus cocok dengan tahun pertama
	if startDate.Year() != startYear {
		return nil, ErrStartDateYearMismatch
	}

	// 6. EndDate.Year() harus cocok dengan tahun kedua
	if endDate.Year() != endYear {
		return nil, ErrEndDateYearMismatch
	}

	// 7. Rentang akademik maks 13 bulan
	maxEnd := startDate.AddDate(1, 1, 0) // + 13 bulan
	if endDate.After(maxEnd) {
		return nil, ErrAcademicRangeTooLong
	}

	// 8. Parse semester dates
	ganjil, err := parseSemesterDates(req.SemesterGanjilStart, req.SemesterGanjilKbm, req.SemesterGanjilAssessment)
	if err != nil {
		return nil, err
	}
	genap, err := parseSemesterDates(req.SemesterGenapStart, req.SemesterGenapKbm, req.SemesterGenapAssessment)
	if err != nil {
		return nil, err
	}

	// 9. Validasi tanggal semester berada dalam rentang akademik [startDate, endDate]
	if err := validateSemesterInRange(ganjil, startDate, endDate); err != nil {
		return nil, err
	}
	if err := validateSemesterInRange(genap, startDate, endDate); err != nil {
		return nil, err
	}

	// 10. Validasi urutan internal semester: Start <= EndKBM <= Assessment
	if err := validateSemesterOrder(ganjil); err != nil {
		return nil, err
	}
	if err := validateSemesterOrder(genap); err != nil {
		return nil, err
	}

	// 11. Ganjil harus lebih dulu dari Genap (jika keduanya punya StartDate)
	if ganjil.Start.Valid && genap.Start.Valid {
		if !ganjil.Start.Time.Before(genap.Start.Time) {
			return nil, ErrOddBeforeEven
		}
	}

	return &parsedDates{
		StartDate: startDate,
		EndDate:   endDate,
		Ganjil:    ganjil,
		Genap:     genap,
	}, nil
}

// parseSemesterDates mem-parsing 3 string tanggal semester menjadi parsedSemester.
// Menghilangkan duplikasi 6x pemanggilan parseNullDate di Create dan Update.
func parseSemesterDates(startStr, endKBMStr, assessmentStr string) (parsedSemester, error) {
	start, err := parseNullDate(startStr)
	if err != nil {
		return parsedSemester{}, err
	}
	endKBM, err := parseNullDate(endKBMStr)
	if err != nil {
		return parsedSemester{}, err
	}
	assessment, err := parseNullDate(assessmentStr)
	if err != nil {
		return parsedSemester{}, err
	}

	return parsedSemester{
		Start:      start,
		EndKBM:     endKBM,
		Assessment: assessment,
	}, nil
}

// validateSemesterInRange memastikan semua tanggal semester berada di dalam rentang akademik.
func validateSemesterInRange(sem parsedSemester, rangeStart, rangeEnd time.Time) error {
	dates := []sql.NullTime{sem.Start, sem.EndKBM, sem.Assessment}
	for _, d := range dates {
		if d.Valid {
			if d.Time.Before(rangeStart) || d.Time.After(rangeEnd) {
				return fmt.Errorf("%w: %s di luar rentang %s – %s",
					ErrSemesterOutOfRange,
					d.Time.Format("2006-01-02"),
					rangeStart.Format("2006-01-02"),
					rangeEnd.Format("2006-01-02"),
				)
			}
		}
	}
	return nil
}

// validateSemesterOrder memastikan urutan: Start <= EndKBM <= Assessment.
func validateSemesterOrder(sem parsedSemester) error {
	if sem.Start.Valid && sem.EndKBM.Valid {
		if sem.EndKBM.Time.Before(sem.Start.Time) {
			return ErrSemesterDateOrder
		}
	}
	if sem.EndKBM.Valid && sem.Assessment.Valid {
		if sem.Assessment.Time.Before(sem.EndKBM.Time) {
			return ErrSemesterDateOrder
		}
	}
	if sem.Start.Valid && sem.Assessment.Valid && !sem.EndKBM.Valid {
		if sem.Assessment.Time.Before(sem.Start.Time) {
			return ErrSemesterDateOrder
		}
	}
	return nil
}
