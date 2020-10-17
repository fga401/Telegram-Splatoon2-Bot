package nintendo

type ScheduleTime struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

type SalmonSpecialWeapon struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}
type SalmonWeapon struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Image     string `json:"image"`
	Thumbnail string `json:"thumbnail"`
}
type SalmonWeaponWrapper struct {
	ID            string               `json:"id"`
	Weapon        *SalmonWeapon        `json:"weapon"`
	SpecialWeapon *SalmonSpecialWeapon `json:"coop_special_weapon"`
}
type SalmonStage struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}
type SalmonScheduleDetail struct {
	Weapons []*SalmonWeaponWrapper `json:"weapons"`
	Stage   *SalmonStage           `json:"stage"`
	*ScheduleTime
}
type SalmonSchedules struct {
	Details   []*SalmonScheduleDetail `json:"details"`
	Schedules []*ScheduleTime         `json:"schedules"`
}
