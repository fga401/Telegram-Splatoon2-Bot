package nintendo

// Common

// Endpoint of Nintendo API
var Endpoint = "https://app.splatoon2.nintendo.net"

// ScheduleTime JSON structure
type ScheduleTime struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

// Rule JSON structure
type Rule struct {
	MultilineName string `json:"multiline_name"`
	Name          string `json:"name"`
	Key           string `json:"key"`
}

// Stage JSON structure
type Stage struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

// GameMode JSON structure
type GameMode struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

// Udemae JSON structure
type Udemae struct {
	Number          int32  `json:"number"`
	IsNumberReached bool   `json:"is_number_reached"`
	IsX             bool   `json:"is_x"`
	SPlusNumber     int32  `json:"s_plus_number"`
	Name            string `json:"name"`
}

// Rank JSON structure
type Rank struct {
	PlayerRank int32 `json:"player_rank"`
	StarRank   int32 `json:"star_rank"`
}

// PlayerType JSON structure
type PlayerType struct {
	Species string `json:"species"`
	Style   string `json:"style"`
}

// GearSkill JSON structure
type GearSkill struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

// GearSkills JSON structure
type GearSkills struct {
	Subs []GearSkill `json:"subs"`
	Main GearSkill   `json:"main"`
}

// Brand JSON structure
type Brand struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Image         string    `json:"image"`
	FrequentSkill GearSkill `json:"frequent_skill"`
}

// Gear JSON structure
type Gear struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Brand     Brand  `json:"brand"`
	Thumbnail string `json:"thumbnail"`
	Image     string `json:"image"`
	Rarity    int32  `json:"rarity"`
}

// WeaponSkill JSON structure
type WeaponSkill struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	ImageA string `json:"image_a"`
	ImageB string `json:"image_b"`
}

// Weapon JSON structure
type Weapon struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Image     string      `json:"image"`
	Thumbnail string      `json:"thumbnail"`
	Special   WeaponSkill `json:"special"`
	Sub       WeaponSkill `json:"sub"`
}

// Player JSON structure
type Player struct {
	PrincipalID string `json:"principal_id"`
	Nickname    string `json:"nickname"`
	Rank
	Udemae        Udemae     `json:"udemae"`
	PlayType      PlayerType `json:"player_type"`
	Weapon        Weapon     `json:"weapon"`
	Head          Gear       `json:"head"`
	HeadSkills    GearSkills `json:"head_skills"`
	Clothes       Gear       `json:"clothes"`
	ClothesSkills GearSkills `json:"clothes_skills"`
	Shoes         Gear       `json:"shoes"`
	ShoesSkills   GearSkills `json:"shoes_skills"`
}

// BattleResultTypeEnum lists all available types of battle.
var BattleResultTypeEnum = struct {
	Regular  BattleResultType
	Gachi    BattleResultType
	League   BattleResultType
	Festival BattleResultType
}{"regular", "gachi", "league", "fes"}

// BattleResultType is what its name says.
type BattleResultType string

// BattleResult is what its name says.
type BattleResult interface {
	Type() BattleResultType
	Metadata() BattleResultMetadata
}

// BattleResultMetadata JSON structure
type BattleResultMetadata struct {
	BattleNumber string   `json:"battle_number"`
	Rule         Rule     `json:"rule"`
	Type         string   `json:"type"`
	Stage        Stage    `json:"stage"`
	GameMode     GameMode `json:"game_mode"`
	Rank
	PlayerResult     PlayerResult `json:"player_result"`
	StartTime        int64        `json:"start_time"`
	MyTeamResult     TeamResult   `json:"my_team_result"`
	OtherTeamResult  TeamResult   `json:"other_team_result"`
	WeaponPaintPoint int32        `json:"weapon_paint_point"`
}

// RegularBattleResult JSON structure
type RegularBattleResult struct {
	BattleResultMetadata
	WinMeter            float32 `json:"win_meter"`
	MyTeamPercentage    float32 `json:"my_team_percentage"`
	OtherTeamPercentage float32 `json:"other_team_percentage"`
}

// GachiBattleResult JSON structure
type GachiBattleResult struct {
	BattleResultMetadata
	Udemae             Udemae  `json:"udemae"`
	XPower             float32 `json:"x_power"`
	ElapsedTime        int32   `json:"elapsed_time"`
	MyTeamCount        int32   `json:"my_team_count"`
	OtherTeamCount     int32   `json:"other_team_count"`
	EstimateGachiPower float32 `json:"estimate_gachi_power"`
	EstimateXPower     float32 `json:"estimate_x_power"`
	// todo: CrownPlayers
}

// LeagueBattleResult JSON structure
type LeagueBattleResult struct {
	BattleResultMetadata
	TagID                    string  `json:"tag_id"`
	Udemae                   Udemae  `json:"udemae"`
	ElapsedTime              int32   `json:"elapsed_time"`
	MyTeamCount              int32   `json:"my_team_count"`
	OtherTeamCount           int32   `json:"other_team_count"`
	LeaguePoint              float32 `json:"league_point"`
	MaxLeaguePoint           float32 `json:"max_league_point"`
	MyEstimateLeaguePoint    float32 `json:"my_estimate_league_point"`
	OtherEstimateLeaguePoint float32 `json:"other_estimate_league_point"`
	EstimateGachiPower       float32 `json:"estimate_gachi_power"`
}

// Type return BattleResultType of RegularBattleResult.
func (r *RegularBattleResult) Type() BattleResultType {
	return BattleResultTypeEnum.Regular
}

// Type return BattleResultType of GachiBattleResult.
func (r *GachiBattleResult) Type() BattleResultType {
	return BattleResultTypeEnum.Gachi
}

// Type return BattleResultType of LeagueBattleResult.
func (r *LeagueBattleResult) Type() BattleResultType {
	return BattleResultTypeEnum.League
}

// Metadata return BattleResultMetadata of RegularBattleResult.
func (r *RegularBattleResult) Metadata() BattleResultMetadata {
	return r.BattleResultMetadata
}

// Metadata return BattleResultMetadata of GachiBattleResult.
func (r *GachiBattleResult) Metadata() BattleResultMetadata {
	return r.BattleResultMetadata
}

// Metadata return BattleResultMetadata of LeagueBattleResult.
func (r *LeagueBattleResult) Metadata() BattleResultMetadata {
	return r.BattleResultMetadata
}

// BattleResults JSON structure
type BattleResults struct {
	ID      string         `json:"unique_id"`
	Summary BattleSummary  `json:"summary"`
	Results []BattleResult `json:"results"`
}

// BattleSummary JSON structure
type BattleSummary struct {
	KillCountAverage    float32 `json:"kill_count_average"`
	AssistCountAverage  float32 `json:"assist_count_average"`
	DeathCountAverage   float32 `json:"death_count_average"`
	SpecialCountAverage float32 `json:"special_count_average"`
	VictoryRate         float32 `json:"victory_rate"`
	VictoryCount        int32   `json:"victory_count"`
	DefeatCount         int32   `json:"defeat_count"`
	Count               int32   `json:"count"`
}

// TeamResult JSON structure
type TeamResult struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

// PlayerResult JSON structure
type PlayerResult struct {
	GamePaintPoint float32 `json:"game_paint_point"`
	DeathCount     int32   `json:"death_count"`
	KillCount      int32   `json:"kill_count"`
	AssistCount    int32   `json:"assist_count"`
	SpecialCount   int32   `json:"special_count"`
	SortScore      int32   `json:"sort_score"`
	Player         Player  `json:"player"`
}

// TeamPlayerResults JSON structure
type TeamPlayerResults struct {
	MyTeamMembers    []PlayerResult `json:"my_team_members"`
	OtherTeamMembers []PlayerResult `json:"other_team_members"`
}

// MyTeamPlayerResults return my team player results
func (t *TeamPlayerResults) MyTeamPlayerResults() []PlayerResult {
	return t.MyTeamMembers
}

// OtherTeamPlayerResults return my team player results
func (t *TeamPlayerResults) OtherTeamPlayerResults() []PlayerResult {
	return t.OtherTeamMembers
}

// DetailedBattleResult interface of detailed BattleResult.
type DetailedBattleResult interface {
	BattleResult
	MyTeamPlayerResults() []PlayerResult
	OtherTeamPlayerResults() []PlayerResult
}

// DetailedRegularBattleResult JSON structure
type DetailedRegularBattleResult struct {
	RegularBattleResult
	TeamPlayerResults
}

// DetailedGachiBattleResult JSON structure
type DetailedGachiBattleResult struct {
	GachiBattleResult
	TeamPlayerResults
}

// DetailedLeagueBattleResult JSON structure
type DetailedLeagueBattleResult struct {
	LeagueBattleResult
	TeamPlayerResults
}

// SalmonSummary JSON struct
type SalmonSummary struct {
	Card       SalmonCard     `json:"card"`
	Results    []SalmonResult `json:"results"`
	RewardGear Gear           `json:"reward_gear"`
}

// SalmonCard JSON struct
type SalmonCard struct {
	KumaPoint        int32 `json:"kuma_point"`
	IkuraTotal       int64 `json:"ikura_total"`
	JobNumber        int32 `json:"job_number"`
	GoldenIkuraTotal int32 `json:"golden_ikura_total"`
	HelpTotal        int32 `json:"help_total"`
	KumaPointTotal   int64 `json:"kuma_point_total"`
}

// SalmonStats JSON struct
type SalmonStats struct {
	ScheduleTime
	Grade                SalmonGrade          `json:"grade"`
	GradePoint           int32                `json:"grade_point"`
	JobNum               int32                `json:"job_num"`
	ClearNum             int32                `json:"clear_num"`
	FailureCounts        []int32              `json:"failure_counts"` // failure in the i-th wave
	HelpTotal            int32                `json:"help_total"`
	DeadTotal            int32                `json:"dead_total"`
	KumaPointTotal       int64                `json:"kuma_point_total"`
	MyIkuraTotal         int32                `json:"my_ikura_total"`
	MyGoldenIkuraTotal   int32                `json:"my_golden_ikura_total"`
	TeamIkuraTotal       int32                `json:"team_ikura_total"`
	TeamGoldenIkuraTotal int32                `json:"team_golden_ikura_total"`
	Schedule             SalmonScheduleDetail `json:"schedule"`
}

// SalmonGrade JSON struct
type SalmonGrade struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// SalmonResult JSON structure
type SalmonResult struct {
	ScheduleTime
	JobID           int32                `json:"job_id"`
	PlayTime        int64                `json:"play_time"`
	Schedule        SalmonScheduleDetail `json:"schedule"`
	JobRate         int32                `json:"job_rate"`
	DangerRate      float32              `json:"danger_rate"`
	JobScore        int32                `json:"job_score"`
	KumaPoint       int32                `json:"kuma_point"`
	Grade           SalmonDetailedGrade  `json:"grade"`
	GradePoint      int32                `json:"grade_point"`
	GradePointDelta int32                `json:"grade_point_delta"`
	JobResult       SalmonJobResult      `json:"job_result"`
	BossCounts      BossCounts           `json:"boss_counts"`
	WaveDetails     []SalmonWaveDetail   `json:"wave_details"`
	MyResult        SalmonPlayerResult   `json:"my_result"`
	PlayType        PlayerType           `json:"play_type"`
}

// SalmonDetailedGrade JSON structure
type SalmonDetailedGrade struct {
	SalmonGrade
	ShortName string `json:"short_name"`
	LongName  string `json:"long_name"`
}

// SalmonJobResult JSON structure
type SalmonJobResult struct {
	IsClear       bool   `json:"is_clear"`
	FailureWave   int32  `json:"failure_wave"`
	FailureReason string `json:"failure_reason"`
}

// SalmonBoss JSON structure
type SalmonBoss struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

// BossCounts JSON structure
type BossCounts map[string]SalmonBossCount

// SalmonBossCount JSON structure
type SalmonBossCount struct {
	Boss  SalmonBoss `json:"boss"`
	Count int32      `json:"count"`
}

// SalmonWaveDetail JSON structure
type SalmonWaveDetail struct {
	QuotaNum          int32            `json:"quota_num"`
	IkuraNum          int32            `json:"ikura_num"`
	GoldenIkuraNum    int32            `json:"golden_ikura_num"`
	GoldenIkuraPopNum int32            `json:"golden_ikura_pop_num"`
	WaterLevel        SalmonWaterLevel `json:"water_level"`
	EventType         SalmonEventType  `json:"event_type"`
}

// SalmonWaterLevel JSON structure
type SalmonWaterLevel struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

// SalmonEventType JSON structure
type SalmonEventType struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

// SalmonPlayerResult JSON structure
type SalmonPlayerResult struct {
	PID            string             `json:"pid"`
	Name           string             `json:"name"`
	PlayerType     PlayerType         `json:"player_type"`
	WeaponList     []SalmonWeaponType `json:"weapon_list"`
	Special        WeaponSkill        `json:"special"`
	SpecialCount   []int32            `json:"special_count"`
	BossKillCounts BossCounts         `json:"boss_kill_counts"`
	HelpCount      int32              `json:"help_count"`
	DeadCount      int32              `json:"dead_count"`
	IkuraNum       int32              `json:"ikura_num"`
	GoldenIkuraNum int32              `json:"golden_ikura_num"`
}

// SalmonDetailedResult JSON structure
type SalmonDetailedResult struct {
	SalmonResult
	OtherResults []SalmonPlayerResult `json:"other_results"`
}

// SalmonSpecialWeapon JSON structure
type SalmonSpecialWeapon struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

// SalmonWeapon JSON structure
type SalmonWeapon struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Image     string `json:"image"`
	Thumbnail string `json:"thumbnail"`
}

// SalmonWeaponType JSON structure
type SalmonWeaponType struct {
	ID            string               `json:"id"`
	Weapon        *SalmonWeapon        `json:"weapon"`
	SpecialWeapon *SalmonSpecialWeapon `json:"coop_special_weapon"`
}

// SalmonStage JSON structure
type SalmonStage struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

// SalmonScheduleDetail JSON structure
type SalmonScheduleDetail struct {
	Weapons []SalmonWeaponType `json:"weapons"`
	Stage   SalmonStage        `json:"stage"`
	ScheduleTime
}

// SalmonSchedules JSON structure
type SalmonSchedules struct {
	Details   []SalmonScheduleDetail `json:"details"`
	Schedules []ScheduleTime         `json:"schedules"`
}

// StageSchedule JSON structure
type StageSchedule struct {
	ID       int64    `json:"id"`
	Rule     Rule     `json:"rule"`
	GameMode GameMode `json:"game_mode"`
	StageA   Stage    `json:"stage_a"`
	StageB   Stage    `json:"stage_b"`
	ScheduleTime
}

// StageSchedules JSON structure
type StageSchedules struct {
	League  []StageSchedule `json:"league"`
	Gachi   []StageSchedule `json:"gachi"`
	Regular []StageSchedule `json:"regular"`
}

// Color JSON structure
type Color struct {
	A      float64 `json:"a"`
	B      float64 `json:"b"`
	R      float64 `json:"r"`
	G      float64 `json:"g"`
	CSSRGB string  `json:"css_rgb"`
}

// Theme JSON structure
type Theme struct {
	Color Color  `json:"color"`
	Key   string `json:"key"`
	Name  string `json:"name"`
}

// EventType JSON structure
type EventType struct {
	MultilineName string `json:"multiline_name"`
	ClassName     string `json:"class_name"`
	Key           string `json:"key"`
	Name          string `json:"name"`
}

// FesGrade JSON structure
type FesGrade struct {
	Name string `json:"name"`
	Rank string `json:"rank"`
}

// FesBattleResult JSON structure
type FesBattleResult struct {
	RegularBattleResult
	FesMode                 GameMode  `json:"fes_mode"`
	FesID                   int64     `json:"fes_id"`
	UniformBonus            float32   `json:"uniform_bonus"`
	EventType               EventType `json:"event_type"`
	FesGrade                FesGrade  `json:"fes_grade"`
	MyTeamFesTheme          Theme     `json:"my_team_fes_theme"`
	OtherTeamFesTheme       Theme     `json:"other_team_fes_theme"`
	MyEstimateFesPower      float32   `json:"my_estimate_fes_power"`
	OtherEstimateFesPower   float32   `json:"other_estimate_fes_power"`
	MyTeamAnotherName       string    `json:"my_team_another_name"`
	OtherTeamAnotherName    string    `json:"other_team_another_name"`
	MyTeamConsecutiveWin    int32     `json:"my_team_consecutive_win"`
	OtherTeamConsecutiveWin int32     `json:"other_team_consecutive_win"`
	FesPower                float32   `json:"fes_power"`
	MaxFesPower             float32   `json:"max_fes_power"`
	FesPoint                float32   `json:"fes_point"`
	ContributionPoint       float32   `json:"contribution_point"`
	ContributionPointTotal  float32   `json:"contribution_point_total"`
	Version                 int32     `json:"version"`
}

// Type return BattleResultType of FesBattleResult.
func (r *FesBattleResult) Type() BattleResultType {
	return BattleResultTypeEnum.Festival
}

// Metadata return BattleResultMetadata of FesBattleResult.
func (r *FesBattleResult) Metadata() BattleResultMetadata {
	return r.BattleResultMetadata
}

// DetailedFesBattleResult JSON structure
type DetailedFesBattleResult struct {
	FesBattleResult
	TeamPlayerResults
}
