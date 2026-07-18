package academic_year

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	
	"backend-lingualoop/internal/modules/class"
	"backend-lingualoop/internal/modules/student"
)

type Service interface {
	GetAll(ctx context.Context) ([]*AcademicYearResponse, error)
	GetByID(ctx context.Context, id string) (*AcademicYearResponse, error)
	Create(ctx context.Context, req AcademicYearRequest) (*AcademicYearResponse, error)
	Update(ctx context.Context, id string, req AcademicYearRequest) (*AcademicYearResponse, error)
	Activate(ctx context.Context, id string) (*AcademicYearResponse, error)
	UpdateSemesterStatus(ctx context.Context, id string, req SemesterStatusRequest) (*AcademicYearResponse, error)
	CloseSemester(ctx context.Context, id string, req CloseSemesterRequest) (*AcademicYearResponse, error)
	Delete(ctx context.Context, id string) error
	FinalizePromotion(ctx context.Context, id string, req FinalizePromotionRequest) error
}

type service struct {
	repo        Repository
	classRepo   class.Repository
	studentRepo student.Repository
	db          *sql.DB
}

func NewService(repo Repository, classRepo class.Repository, studentRepo student.Repository, db *sql.DB) Service {
	return &service{
		repo:        repo,
		classRepo:   classRepo,
		studentRepo: studentRepo,
		db:          db,
	}
}

// --- Helper Functions ---

func parseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, errors.New("date is empty")
	}
	return time.Parse("2006-01-02", dateStr)
}

func parseNullDate(dateStr string) (sql.NullTime, error) {
	if dateStr == "" {
		return sql.NullTime{Valid: false}, nil
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return sql.NullTime{Valid: false}, ErrInvalidDateFormat
	}
	return sql.NullTime{Time: t, Valid: true}, nil
}

func formatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

func formatNullDate(nt sql.NullTime) string {
	if !nt.Valid {
		return ""
	}
	return nt.Time.Format("2006-01-02")
}

// buildEntity mengkonstruksi AcademicYear entity dari request dan hasil parsing yang sudah tervalidasi.
func buildEntity(req AcademicYearRequest, dates *parsedDates) *AcademicYear {
	now := time.Now()
	return &AcademicYear{
		ID:                     uuid.New().String(),
		Year:                   req.Year,
		StartDate:              dates.StartDate,
		EndDate:                dates.EndDate,
		Status:                 StatusDraft,
		SemGanjilStartDate:     dates.Ganjil.Start,
		SemGanjilEndKBM:        dates.Ganjil.EndKBM,
		SemGanjilEndAssessment: dates.Ganjil.Assessment,
		SemGanjilStatus:        SemStatusNotActive,
		SemGenapStartDate:      dates.Genap.Start,
		SemGenapEndKBM:         dates.Genap.EndKBM,
		SemGenapEndAssessment:  dates.Genap.Assessment,
		SemGenapStatus:         SemStatusNotActive,
		CreatedAt:              now,
		UpdatedAt:              now,
	}
}

// applyDatesToEntity menerapkan data tanggal hasil validasi ke entity yang sudah ada (untuk Update).
func applyDatesToEntity(existing *AcademicYear, req AcademicYearRequest, dates *parsedDates) {
	existing.Year = req.Year
	existing.StartDate = dates.StartDate
	existing.EndDate = dates.EndDate
	existing.SemGanjilStartDate = dates.Ganjil.Start
	existing.SemGanjilEndKBM = dates.Ganjil.EndKBM
	existing.SemGanjilEndAssessment = dates.Ganjil.Assessment
	existing.SemGenapStartDate = dates.Genap.Start
	existing.SemGenapEndKBM = dates.Genap.EndKBM
	existing.SemGenapEndAssessment = dates.Genap.Assessment
	existing.UpdatedAt = time.Now()
}

func mapEntityToDTO(y *AcademicYear) *AcademicYearResponse {
	activeSemesterName := "-"
	if y.Status == StatusActive {
		if y.SemGanjilStatus == SemStatusActive || y.SemGanjilStatus == SemStatusAssessment || y.SemGanjilStatus == SemStatusReadyToClose {
			activeSemesterName = SemesterOddLabel
		} else if y.SemGenapStatus == SemStatusActive || y.SemGenapStatus == SemStatusAssessment || y.SemGenapStatus == SemStatusReadyToClose {
			activeSemesterName = SemesterEvenLabel
		}
	}

	return &AcademicYearResponse{
		ID:        y.ID,
		Year:      y.Year,
		StartDate: formatDate(y.StartDate),
		EndDate:   formatDate(y.EndDate),
		Status:    y.Status,
		IsActive:  y.Status == StatusActive,
		Semester:  activeSemesterName,
		SemesterGanjil: SemesterData{
			StartDate:     formatNullDate(y.SemGanjilStartDate),
			EndKBM:        formatNullDate(y.SemGanjilEndKBM),
			EndAssessment: formatNullDate(y.SemGanjilEndAssessment),
			Status:        y.SemGanjilStatus,
		},
		SemesterGenap: SemesterData{
			StartDate:     formatNullDate(y.SemGenapStartDate),
			EndKBM:        formatNullDate(y.SemGenapEndKBM),
			EndAssessment: formatNullDate(y.SemGenapEndAssessment),
			Status:        y.SemGenapStatus,
		},
		CreatedAt: y.CreatedAt,
		UpdatedAt: y.UpdatedAt,
	}
}

// --- Service Methods ---

func (s *service) GetAll(ctx context.Context) ([]*AcademicYearResponse, error) {
	years, err := s.repo.FindAll(ctx)
	if err != nil {
		slog.Error("Failed to query academic years", "error", err)
		return nil, ErrSystemFail
	}

	var responses []*AcademicYearResponse
	for _, y := range years {
		responses = append(responses, mapEntityToDTO(y))
	}
	if responses == nil {
		responses = []*AcademicYearResponse{}
	}
	return responses, nil
}

func (s *service) GetByID(ctx context.Context, id string) (*AcademicYearResponse, error) {
	y, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrAcademicYearNotFound) {
			return nil, ErrAcademicYearNotFound
		}
		slog.Error("Failed to query academic year by ID", "error", err, "id", id)
		return nil, ErrSystemFail
	}
	return mapEntityToDTO(y), nil
}

func (s *service) Create(ctx context.Context, req AcademicYearRequest) (*AcademicYearResponse, error) {
	// 1. Cek apakah ada draft yang menggantung
	draftExists, err := s.repo.CheckDraftExists(ctx)
	if err != nil {
		return nil, err
	}
	if draftExists {
		return nil, ErrDraftExists
	}

	// 2. Cek duplikat
	existingID, err := s.repo.GetIDByYear(ctx, req.Year)
	if err != nil && !errors.Is(err, ErrAcademicYearNotFound) {
		slog.Error("Failed to check unique year", "error", err)
		return nil, ErrSystemFail
	}
	if existingID != "" {
		return nil, ErrAcademicYearExists
	}

	// 2. Validasi + Parse (satu pintu)
	dates, err := validateAndParseRequest(req)
	if err != nil {
		return nil, err
	}

	// 3. Build entity
	ay := buildEntity(req, dates)

	// 4. Simpan
	if err := s.repo.Create(ctx, ay); err != nil {
		slog.Error("Failed to create academic year", "error", err)
		return nil, ErrSystemFail
	}

	return mapEntityToDTO(ay), nil
}

func (s *service) Update(ctx context.Context, id string, req AcademicYearRequest) (*AcademicYearResponse, error) {
	// 1. Ambil data existing
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrAcademicYearNotFound) {
			return nil, ErrAcademicYearNotFound
		}
		slog.Error("Failed to get year for update", "error", err, "id", id)
		return nil, ErrSystemFail
	}

	// 2. Cek duplikat jika year berubah
	if req.Year != existing.Year {
		existingID, err := s.repo.GetIDByYear(ctx, req.Year)
		if err != nil && !errors.Is(err, ErrAcademicYearNotFound) {
			slog.Error("Failed to check unique year on update", "error", err)
			return nil, ErrSystemFail
		}
		if existingID != "" && existingID != id {
			return nil, ErrAcademicYearExists
		}
	}

	// 3. Validasi + Parse (satu pintu)
	dates, err := validateAndParseRequest(req)
	if err != nil {
		return nil, err
	}

	// 4. Terapkan perubahan ke entity
	applyDatesToEntity(existing, req, dates)

	// 5. Simpan
	if err := s.repo.Update(ctx, existing); err != nil {
		slog.Error("Failed to update academic year", "error", err)
		return nil, ErrSystemFail
	}

	return mapEntityToDTO(existing), nil
}

func (s *service) Activate(ctx context.Context, id string) (*AcademicYearResponse, error) {
	err := s.repo.ActivateYear(ctx, id)
	if err != nil {
		if errors.Is(err, ErrAcademicYearNotFound) || errors.Is(err, ErrMultipleActiveYears) {
			return nil, err
		}
		slog.Error("Failed to activate year", "error", err, "id", id)
		return nil, ErrSystemFail
	}

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		slog.Error("Failed to get activated year", "error", err, "id", id)
		return nil, ErrSystemFail
	}

	return mapEntityToDTO(existing), nil
}

func (s *service) UpdateSemesterStatus(ctx context.Context, id string, req SemesterStatusRequest) (*AcademicYearResponse, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrAcademicYearNotFound) {
			return nil, ErrAcademicYearNotFound
		}
		slog.Error("Failed to get year for status update", "error", err, "id", id)
		return nil, ErrSystemFail
	}

	switch req.Status {
	case SemStatusNotActive, SemStatusActive, SemStatusAssessment, SemStatusReadyToClose, SemStatusLocked:
		// Valid
	default:
		return nil, ErrInvalidSemesterStatus
	}

	if req.Semester == SemesterOddKey {
		existing.SemGanjilStatus = req.Status
	} else if req.Semester == SemesterEvenKey {
		existing.SemGenapStatus = req.Status
	}
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, existing); err != nil {
		slog.Error("Failed to update semester status", "error", err)
		return nil, ErrSystemFail
	}

	return mapEntityToDTO(existing), nil
}

func (s *service) CloseSemester(ctx context.Context, id string, req CloseSemesterRequest) (*AcademicYearResponse, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrAcademicYearNotFound) {
			return nil, ErrAcademicYearNotFound
		}
		slog.Error("Failed to get year to close semester", "error", err, "id", id)
		return nil, ErrSystemFail
	}

	if req.Semester == SemesterOddKey {
		if existing.SemGanjilStatus == SemStatusNotActive || existing.SemGanjilStatus == SemStatusLocked {
			return nil, ErrSemesterNotActive
		}
		existing.SemGanjilStatus = SemStatusLocked
		existing.SemGenapStatus = SemStatusActive
	} else if req.Semester == SemesterEvenKey {
		if existing.SemGenapStatus == SemStatusNotActive || existing.SemGenapStatus == SemStatusLocked {
			return nil, ErrSemesterNotActive
		}
		existing.SemGenapStatus = SemStatusLocked
		existing.Status = StatusPendingPromotion // Menunggu Kenaikan
	}
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, existing); err != nil {
		slog.Error("Failed to close semester", "error", err)
		return nil, ErrSystemFail
	}

	return mapEntityToDTO(existing), nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrAcademicYearNotFound) {
			return ErrAcademicYearNotFound
		}
		slog.Error("Failed to get year for delete", "error", err, "id", id)
		return ErrSystemFail
	}

	if existing.Status == StatusActive {
		return ErrDeleteActiveYear
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		slog.Error("Failed to delete academic year", "error", err)
		return ErrSystemFail
	}

	return nil
}

func (s *service) FinalizePromotion(ctx context.Context, id string, req FinalizePromotionRequest) error {
	sourceYear, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return ErrAcademicYearNotFound
	}
	
	if sourceYear.Status != StatusPendingPromotion {
		return errors.New("academic year is not ready for promotion")
	}

	targetYear, err := s.repo.FindByID(ctx, req.TargetYearID)
	if err != nil {
		return errors.New("target academic year not found")
	}

	// Cek Idempotency / Job
	job, err := s.repo.GetPromotionJob(ctx, sourceYear.ID)
	if err != nil {
		return ErrSystemFail
	}
	if job != nil && (job.Status == JobStatusRunning || job.Status == JobStatusDone) {
		return errors.New("promotion already processed or is running")
	}

	now := time.Now()
	newJob := &PromotionJob{
		ID:              uuid.New().String(),
		AcademicYearID:  sourceYear.ID,
		Status:          JobStatusRunning,
		TotalStudents:   len(req.Promotions),
		StartedAt:       &now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.repo.CreatePromotionJob(ctx, newJob); err != nil {
		slog.Error("Failed to create promotion job", "error", err)
		return ErrSystemFail
	}

	// Error handler for job (update status to FAILED on defer panic/error)
	var finalErr error
	defer func() {
		if finalErr != nil {
			newJob.Status = JobStatusFailed
		} else {
			newJob.Status = JobStatusDone
		}
		newJob.FinishedAt = &now
		newJob.UpdatedAt = time.Now()
		_ = s.repo.UpdatePromotionJob(context.Background(), newJob)
	}()

	// 1. Begin Tx
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		finalErr = err
		return ErrSystemFail
	}
	defer tx.Rollback()

	// 2. Data Load (Classes & Students of source year)
	sourceClasses, err := s.classRepo.FindAllByAcademicYear(ctx, sourceYear.ID)
	if err != nil {
		finalErr = err
		return err
	}
	
	sourceStudents, err := s.studentRepo.FindAllByAcademicYear(ctx, sourceYear.ID)
	if err != nil {
		finalErr = err
		return err
	}

	// Bikin map untuk source classes supaya gampang nyari GradeLevel & ClassNumber nya
	classMap := make(map[string]*class.Class)
	for _, c := range sourceClasses {
		classMap[c.ID] = c
	}

	// 3. Class Generation (Cloning)
	var targetClasses []*class.Class
	newClassMapByGradeMajorNumber := make(map[string]string) // key: "grade_major_number", value: new_class_id

	// Untuk kelas target, kita clone saja dari kelas sumber, dengan id tahun ajaran target
	for _, c := range sourceClasses {
		// GradeLevel naik 1 (atau tetep jika misal mentok 12, tapi logic bisnis di sini kita clone semua)
		newGrade := c.GradeLevel + 1
		if newGrade > 12 {
			continue // Lulus, ga perlu dibikinin kelas di tahun depan untuk angkatan ini (kecuali kalau tinggal kelas)
		}
		
		levelID, levelName, err := s.classRepo.GetLevelByGrade(ctx, newGrade)
		if err != nil {
			// Kalau misal grade 13 ga ketemu, skip. Tapi udah di-handle newGrade > 12
			continue
		}

		majorCode := ""
		if c.MajorID != nil {
			mc, _ := s.classRepo.GetMajorCodeByID(ctx, *c.MajorID)
			majorCode = mc
		}

		className := fmt.Sprintf("%s %s %d", levelName, majorCode, c.ClassNumber)
		newClass := &class.Class{
			ID:                uuid.New().String(),
			AcademicYearID:    targetYear.ID,
			MajorID:           c.MajorID,
			LevelID:           levelID,
			GradeLevel:        newGrade,
			ClassNumber:       c.ClassNumber,
			ClassName:         className,
			Classroom:         nil, // Reset classroom
			Capacity:          c.Capacity,
			HomeroomTeacherID: nil, // Reset guru
			IsActive:          true,
			CreatedAt:         now,
			UpdatedAt:         now,
		}
		targetClasses = append(targetClasses, newClass)
		
		key := fmt.Sprintf("%d_%s_%d", newGrade, optionalStringVal(c.MajorID), c.ClassNumber)
		newClassMapByGradeMajorNumber[key] = newClass.ID
	}

	// Tambahkan class target untuk siswa yang "Tinggal Kelas" (mengulang grade_level yang sama)
	for _, c := range sourceClasses {
		key := fmt.Sprintf("%d_%s_%d", c.GradeLevel, optionalStringVal(c.MajorID), c.ClassNumber)
		if _, exists := newClassMapByGradeMajorNumber[key]; !exists {
			levelID, levelName, _ := s.classRepo.GetLevelByGrade(ctx, c.GradeLevel)
			majorCode := ""
			if c.MajorID != nil {
				mc, _ := s.classRepo.GetMajorCodeByID(ctx, *c.MajorID)
				majorCode = mc
			}
			
			className := fmt.Sprintf("%s %s %d", levelName, majorCode, c.ClassNumber)
			newClass := &class.Class{
				ID:                uuid.New().String(),
				AcademicYearID:    targetYear.ID,
				MajorID:           c.MajorID,
				LevelID:           levelID,
				GradeLevel:        c.GradeLevel,
				ClassNumber:       c.ClassNumber,
				ClassName:         className,
				Classroom:         nil,
				Capacity:          c.Capacity,
				HomeroomTeacherID: nil,
				IsActive:          true,
				CreatedAt:         now,
				UpdatedAt:         now,
			}
			targetClasses = append(targetClasses, newClass)
			newClassMapByGradeMajorNumber[key] = newClass.ID
		}
	}

	if err := s.classRepo.CreateBatchTx(ctx, tx, targetClasses); err != nil {
		finalErr = err
		return err
	}

	// Map request
	promoReqMap := make(map[string]string)
	for _, p := range req.Promotions {
		promoReqMap[p.StudentID] = p.Status
	}

	// 4. Student Movement
	var studentClasses []*student.StudentClass
	var studentPromotions []*StudentPromotionHistory
	var updatedStudents []*student.Student

	successCount := 0
	failedCount := 0

	for _, st := range sourceStudents {
		statusReq, ok := promoReqMap[st.ID]
		if !ok {
			failedCount++
			continue
		}

		oldClassID := st.CurrentClassID
		var oldClass *class.Class
		if oldClassID != nil {
			oldClass = classMap[*oldClassID]
		}

		var targetClassID *string
		newStudentStatus := st.Status
		var newLevelID *string = st.ClassLevel // default

		promoStatus := PromotionStatusFailed
		
		if statusReq == PromotionStatusPromoted {
			promoStatus = PromotionStatusPromoted
			successCount++
			
			if oldClass != nil {
				newGrade := oldClass.GradeLevel + 1
				key := fmt.Sprintf("%d_%s_%d", newGrade, optionalStringVal(oldClass.MajorID), oldClass.ClassNumber)
				if ncid, exists := newClassMapByGradeMajorNumber[key]; exists {
					targetClassID = &ncid
					
					// update student's level_id
					lid, _, _ := s.classRepo.GetLevelByGrade(ctx, newGrade)
					newLevelID = &lid
				}
			}
		} else if statusReq == PromotionStatusGraduated {
			promoStatus = PromotionStatusGraduated
			newStudentStatus = student.StudentGraduated
			successCount++
		} else if statusReq == PromotionStatusRetained {
			promoStatus = PromotionStatusRetained
			successCount++
			
			// Siswa tetap di grade lama, masuk ke kelas dengan grade yg sama
			if oldClass != nil {
				key := fmt.Sprintf("%d_%s_%d", oldClass.GradeLevel, optionalStringVal(oldClass.MajorID), oldClass.ClassNumber)
				if ncid, exists := newClassMapByGradeMajorNumber[key]; exists {
					targetClassID = &ncid
				}
			}
		} else {
			promoStatus = PromotionStatusFailed
			failedCount++
		}

		// Insert history
		studentPromotions = append(studentPromotions, &StudentPromotionHistory{
			ID:                 uuid.New().String(),
			StudentID:          st.ID,
			FromClassID:        oldClassID,
			ToClassID:          targetClassID,
			FromAcademicYearID: sourceYear.ID,
			ToAcademicYearID:   &targetYear.ID,
			Status:             promoStatus,
			CreatedAt:          now,
		})

		// Jika pindah kelas, tambahkan student_classes
		if targetClassID != nil {
			studentClasses = append(studentClasses, &student.StudentClass{
				ID:             uuid.New().String(),
				StudentID:      st.ID,
				ClassID:        *targetClassID,
				AcademicYearID: targetYear.ID,
				IsActive:       true,
				CreatedAt:      now,
			})
		}

		// Update student
		st.ClassLevel = newLevelID
		st.Status = newStudentStatus
		updatedStudents = append(updatedStudents, st)
	}

	// The Grand Batch Inserts
	if err := s.studentRepo.InsertStudentClassesBatchTx(ctx, tx, studentClasses); err != nil {
		finalErr = err
		return err
	}

	if err := s.repo.InsertPromotionsBatchTx(ctx, tx, studentPromotions); err != nil {
		finalErr = err
		return err
	}

	if err := s.studentRepo.UpdateLevelAndStatusBatchTx(ctx, tx, updatedStudents); err != nil {
		finalErr = err
		return err
	}

	// 5. Update Years
	if err := s.repo.UpdatePromotionCompletedAtTx(ctx, tx, sourceYear.ID); err != nil {
		finalErr = err
		return err
	}
	if err := s.repo.UpdateStatusTx(ctx, tx, sourceYear.ID, StatusFinished); err != nil {
		finalErr = err
		return err
	}

	newJob.SuccessStudents = successCount
	newJob.FailedStudents = failedCount

	if err := tx.Commit(); err != nil {
		finalErr = err
		return ErrSystemFail
	}

	return nil
}

func optionalStringVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
