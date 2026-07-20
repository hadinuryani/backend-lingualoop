package teacher_portal

type TeacherClass struct {
	ID            string
	Name          string
	Major         string
	Enrolled      int
	Capacity      int
	Subject       string
	Room          string
}

type TeacherSchedule struct {
	ID        string
	Day       string
	Period    int
	ClassName string
	Room      string
	Subject   string
	Major     string
}

type ScheduleConfig struct {
	PeriodsPerDay  int
	PeriodDuration int
	StartTime      string // format "HH:MM"
}
