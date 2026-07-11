package dashboard

import "time"

type DashboardResponse struct {
	Stats            []StatCard `json:"stats"`
	LoginStatistics  ChartData  `json:"loginStatistics"` // Akan dipakai untuk Demografi Gender
	TaskStatistics   ChartData  `json:"taskStatistics"`  // Akan dipakai untuk Distribusi Kelas
	RecentActivities []Activity `json:"recentActivities"`
	Shortcuts        []Shortcut `json:"shortcuts"`
}

type StatCard struct {
	Title      string `json:"title"`
	Value      string `json:"value"`
	Growth     string `json:"growth"`
	IconName   string `json:"iconName"`
	ColorClass string `json:"colorClass"`
}

type ChartData struct {
	Categories []string      `json:"categories"`
	Series     []ChartSeries `json:"series"`
}

type ChartSeries struct {
	Name string `json:"name"`
	Data []int  `json:"data"`
}

type Activity struct {
	ID         int    `json:"id"`
	User       string `json:"user"`
	Action     string `json:"action"`
	Target     string `json:"target"`
	Details    string `json:"details"`
	Time       string `json:"time"`
	IconName   string `json:"iconName"`
	ColorClass string `json:"colorClass"`
}

type Shortcut struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Path        string `json:"path"`
	IconName    string `json:"iconName"`
	ColorClass  string `json:"colorClass"`
}

// Struct internal untuk hasil kueri
type GenderStat struct {
	Gender string
	Count  int
}

type LevelStat struct {
	Level string
	Count int
}

type RecentRegistration struct {
	ID        string
	FullName  string
	Role      string // "Teacher" atau "Student"
	CreatedAt time.Time
}
