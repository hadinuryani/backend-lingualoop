package dashboard

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
)

type Service interface {
	GetDashboardData(ctx context.Context) (*DashboardResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) GetDashboardData(ctx context.Context) (*DashboardResponse, error) {
	// 1. Ambil Aggregate Stats
	totalTeachers, err := s.repo.GetTotalTeachers(ctx)
	if err != nil {
		slog.Error("Failed to get total teachers", "error", err)
	}

	totalStudents, err := s.repo.GetTotalStudents(ctx)
	if err != nil {
		slog.Error("Failed to get total students", "error", err)
	}

	totalClasses, err := s.repo.GetTotalClasses(ctx)
	if err != nil {
		slog.Error("Failed to get total classes", "error", err)
	}

	totalMajors, err := s.repo.GetTotalMajors(ctx)
	if err != nil {
		slog.Error("Failed to get total majors", "error", err)
	}

	stats := []StatCard{
		{
			Title:      "Total Guru",
			Value:      strconv.Itoa(totalTeachers),
			Growth:     "Berdasarkan data master",
			IconName:   "GraduationCap",
			ColorClass: "bg-indigo-500/10 text-indigo-500 border-indigo-500/10",
		},
		{
			Title:      "Total Siswa",
			Value:      strconv.Itoa(totalStudents),
			Growth:     "Berdasarkan data master",
			IconName:   "Users",
			ColorClass: "bg-emerald-500/10 text-emerald-500 border-emerald-500/10",
		},
		{
			Title:      "Total Kelas Aktif",
			Value:      strconv.Itoa(totalClasses),
			Growth:     "Rombongan belajar terdaftar",
			IconName:   "BookOpen",
			ColorClass: "bg-pink-500/10 text-pink-500 border-pink-500/10",
		},
		{
			Title:      "Total Jurusan",
			Value:      strconv.Itoa(totalMajors),
			Growth:     "Program studi yang tersedia",
			IconName:   "Layers",
			ColorClass: "bg-amber-500/10 text-amber-500 border-amber-500/10",
		},
	}

	// 2. Ambil Demografi Gender Siswa (Chart 1)
	genderStats, err := s.repo.GetGenderDemographics(ctx)
	var genderCategories []string
	var genderData []int
	if err == nil {
		for _, stat := range genderStats {
			genderCategories = append(genderCategories, stat.Gender)
			genderData = append(genderData, stat.Count)
		}
	} else {
		slog.Error("Failed to get gender stats", "error", err)
	}

	loginStatistics := ChartData{
		Categories: genderCategories,
		Series: []ChartSeries{
			{
				Name: "Jumlah Siswa",
				Data: genderData,
			},
		},
	}

	// 3. Ambil Distribusi Level Kelas (Chart 2)
	levelStats, err := s.repo.GetClassLevelDistribution(ctx)
	var levelCategories []string
	var levelData []int
	if err == nil {
		for _, stat := range levelStats {
			levelCategories = append(levelCategories, stat.Level)
			levelData = append(levelData, stat.Count)
		}
	} else {
		slog.Error("Failed to get level stats", "error", err)
	}

	taskStatistics := ChartData{
		Categories: levelCategories,
		Series: []ChartSeries{
			{
				Name: "Jumlah Kelas",
				Data: levelData,
			},
		},
	}

	// 4. Ambil Recent Activities (Registrasi Terbaru)
	registrations, err := s.repo.GetRecentRegistrations(ctx)
	var recentActivities []Activity
	if err == nil {
		for i, reg := range registrations {
			action := "Pendaftaran Siswa Baru"
			iconName := "UserPlus"
			colorClass := "bg-emerald-500/10 text-emerald-500 border-emerald-500/20"

			if reg.Role == "Teacher" {
				action = "Pendaftaran Guru Baru"
				iconName = "GraduationCap"
				colorClass = "bg-indigo-500/10 text-indigo-500 border-indigo-500/20"
			}

			recentActivities = append(recentActivities, Activity{
				ID:         i + 1,
				User:       "Sistem",
				Action:     action,
				Target:     reg.FullName,
				Details:    fmt.Sprintf("ID: %s", reg.ID),
				Time:       reg.CreatedAt.Format("02 Jan 2006 15:04"),
				IconName:   iconName,
				ColorClass: colorClass,
			})
		}
	} else {
		slog.Error("Failed to get recent registrations", "error", err)
	}

	if recentActivities == nil {
		recentActivities = []Activity{} // ensure it's not null in JSON
	}

	// 5. Shortcuts (Static)
	shortcuts := []Shortcut{
		{
			Title:       "Tambah Guru",
			Description: "Pendaftaran & data tenaga pengajar",
			Path:        "/admin/teachers",
			IconName:    "GraduationCap",
			ColorClass:  "bg-slate-50 hover:bg-slate-100 dark:bg-zinc-800/50 dark:hover:bg-zinc-800 text-slate-700 dark:text-slate-300 border-slate-200 dark:border-zinc-700",
		},
		{
			Title:       "Tambah Siswa",
			Description: "Pendaftaran & data siswa baru",
			Path:        "/admin/students",
			IconName:    "UserPlus",
			ColorClass:  "bg-slate-50 hover:bg-slate-100 dark:bg-zinc-800/50 dark:hover:bg-zinc-800 text-slate-700 dark:text-slate-300 border-slate-200 dark:border-zinc-700",
		},
		{
			Title:       "Tambah Kelas",
			Description: "Buat & jadwalkan rombongan belajar",
			Path:        "/admin/classes",
			IconName:    "School",
			ColorClass:  "bg-slate-50 hover:bg-slate-100 dark:bg-zinc-800/50 dark:hover:bg-zinc-800 text-slate-700 dark:text-slate-300 border-slate-200 dark:border-zinc-700",
		},
		{
			Title:       "Manajemen Jurusan",
			Description: "Kelola program studi & kurikulum",
			Path:        "/admin/majors",
			IconName:    "Layers",
			ColorClass:  "bg-slate-50 hover:bg-slate-100 dark:bg-zinc-800/50 dark:hover:bg-zinc-800 text-slate-700 dark:text-slate-300 border-slate-200 dark:border-zinc-700",
		},
	}

	response := &DashboardResponse{
		Stats:            stats,
		LoginStatistics:  loginStatistics,
		TaskStatistics:   taskStatistics,
		RecentActivities: recentActivities,
		Shortcuts:        shortcuts,
	}

	return response, nil
}
