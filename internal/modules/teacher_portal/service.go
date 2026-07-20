package teacher_portal

import (
	"context"
	"fmt"
	"time"

	"backend-lingualoop/pkg/dateutil"
)

type Service interface {
	GetMyClasses(ctx context.Context, userID string) ([]TeacherClassResponse, error)
	GetMySchedules(ctx context.Context, userID string) (*TeacherScheduleResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetMyClasses(ctx context.Context, userID string) ([]TeacherClassResponse, error) {
	teacherID, err := s.repo.GetTeacherIDByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	classes, err := s.repo.GetClassesByTeacherID(ctx, teacherID)
	if err != nil {
		return nil, err
	}
	
	var responses []TeacherClassResponse
	for _, cls := range classes {
		responses = append(responses, s.mapTeacherClass(cls))
	}

	return responses, nil
}

func (s *service) GetMySchedules(ctx context.Context, userID string) (*TeacherScheduleResponse, error) {
	teacherID, err := s.repo.GetTeacherIDByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	cfg, err := s.repo.GetScheduleConfig(ctx)
	if err != nil {
		return nil, err
	}

	schedules, err := s.repo.GetSchedulesByTeacherID(ctx, teacherID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	todayStr := dateutil.WeekdayID(now.Weekday())

	var todaySchedules []ScheduleItem
	if todaySchedules == nil {
		todaySchedules = []ScheduleItem{}
	}

	weekly := make(map[string]*DailySchedule)
	
	// Initialize 7 days
	days := []string{"Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu", "Minggu"}
	for _, day := range days {
		weekly[day] = &DailySchedule{Classes: []ScheduleItem{}}
	}

	for _, tsch := range schedules {
		item := s.mapSchedule(tsch, cfg, now, todayStr)
		
		if item.Day == todayStr {
			todaySchedules = append(todaySchedules, item)
		}
		
		if daily, ok := weekly[item.Day]; ok {
			daily.Classes = append(daily.Classes, item)
		}
	}

	return &TeacherScheduleResponse{
		TodaySchedule:  todaySchedules,
		WeeklySchedule: weekly,
	}, nil
}

func (s *service) mapTeacherClass(cls TeacherClass) TeacherClassResponse {
	return TeacherClassResponse{
		ID:       cls.ID,
		Name:     cls.Name,
		Major:    cls.Major,
		Enrolled: cls.Enrolled,
		Capacity: cls.Capacity,
		Subject:  cls.Subject,
		Room:     cls.Room,
	}
}

func (s *service) mapSchedule(tsch TeacherSchedule, cfg *ScheduleConfig, now time.Time, todayStr string) ScheduleItem {
	startPeriod, endPeriod := s.buildScheduleTime(cfg, tsch.Period)
	
	item := ScheduleItem{
		ID:        tsch.ID,
		Day:       tsch.Day,
		Period:    tsch.Period,
		ClassName: tsch.ClassName,
		Room:      tsch.Room,
		Subject:   tsch.Subject,
		Major:     tsch.Major,
		Time:      fmt.Sprintf("%02d:%02d - %02d:%02d", startPeriod.Hour(), startPeriod.Minute(), endPeriod.Hour(), endPeriod.Minute()),
	}

	item.Status = s.calculateStatus(item.Day, todayStr, now, startPeriod, endPeriod)
	return item
}

func (s *service) buildScheduleTime(cfg *ScheduleConfig, period int) (time.Time, time.Time) {
	// Format of cfg.StartTime is "HH:MM" e.g. "07:00"
	var hour, min int
	fmt.Sscanf(cfg.StartTime, "%d:%d", &hour, &min)
	
	startTime := time.Date(2000, 1, 1, hour, min, 0, 0, time.UTC)
	startPeriod := startTime.Add(time.Duration((period-1)*cfg.PeriodDuration) * time.Minute)
	endPeriod := startPeriod.Add(time.Duration(cfg.PeriodDuration) * time.Minute)
	return startPeriod, endPeriod
}

func (s *service) calculateStatus(scheduleDay, todayStr string, now time.Time, startPeriod, endPeriod time.Time) string {
	if scheduleDay != todayStr {
		return StatusUpcoming
	}

	nowTotalMins := now.Hour()*60 + now.Minute()
	startTotalMins := startPeriod.Hour()*60 + startPeriod.Minute()
	endTotalMins := endPeriod.Hour()*60 + endPeriod.Minute()

	if nowTotalMins > endTotalMins {
		return StatusFinished
	}
	if nowTotalMins >= startTotalMins && nowTotalMins <= endTotalMins {
		return StatusRunning
	}
	return StatusUpcoming
}
