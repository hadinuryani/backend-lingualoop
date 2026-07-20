package teacher_portal

const (
	StatusUpcoming = "Mendatang"
	StatusRunning  = "Sedang Berlangsung"
	StatusFinished = "Selesai"
)

type TeacherClassResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Major    string `json:"major"`
	Enrolled int    `json:"enrolled"`
	Capacity int    `json:"capacity"`
	Subject  string `json:"subject"`
	Room     string `json:"room"`
}

type ScheduleItem struct {
	ID        string `json:"id"`
	Day       string `json:"day"`
	Period    int    `json:"period"`
	Time      string `json:"time,omitempty"`
	ClassName string `json:"className"`
	Room      string `json:"room"`
	Subject   string `json:"subject"`
	Major     string `json:"major"`
	Status    string `json:"status,omitempty"`
}

type DailySchedule struct {
	IsHoliday   bool           `json:"isHoliday"`
	HolidayName string         `json:"holidayName,omitempty"`
	Classes     []ScheduleItem `json:"classes"`
}

type TeacherScheduleResponse struct {
	TodaySchedule  []ScheduleItem            `json:"todaySchedule"`
	WeeklySchedule map[string]*DailySchedule `json:"weeklySchedule"`
}
