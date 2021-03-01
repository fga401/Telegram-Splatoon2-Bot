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

// PlayType JSON structure
type PlayType struct {
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
	PlayType      PlayType   `json:"play_type"`
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
	PlayResult       PlayResult `json:"play_result"`
	StartTime        int64      `json:"start_time"`
	MyTeamResult     TeamResult `json:"my_team_result"`
	OtherTeamResult  TeamResult `json:"other_team_result"`
	WeaponPaintPoint int32      `json:"weapon_paint_point"`
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
func (r RegularBattleResult) Type() BattleResultType {
	return BattleResultTypeEnum.Regular
}

// Type return BattleResultType of GachiBattleResult.
func (r GachiBattleResult) Type() BattleResultType {
	return BattleResultTypeEnum.Gachi
}

// Type return BattleResultType of LeagueBattleResult.
func (r LeagueBattleResult) Type() BattleResultType {
	return BattleResultTypeEnum.League
}

// Metadata return BattleResultMetadata of RegularBattleResult.
func (r RegularBattleResult) Metadata() BattleResultMetadata {
	return r.BattleResultMetadata
}

// Metadata return BattleResultMetadata of GachiBattleResult.
func (r GachiBattleResult) Metadata() BattleResultMetadata {
	return r.BattleResultMetadata
}

// Metadata return BattleResultMetadata of LeagueBattleResult.
func (r LeagueBattleResult) Metadata() BattleResultMetadata {
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

// PlayResult JSON structure
type PlayResult struct {
	GamePaintPoint float32 `json:"game_paint_point"`
	DeathCount     int32   `json:"death_count"`
	KillCount      int32   `json:"kill_count"`
	AssertCount    int32   `json:"assert_count"`
	SpecialCount   int32   `json:"special_count"`
	SortScore      int32   `json:"sort_score"`
	Player         Player  `json:"player"`
}

// RawBattleResult JSON structure
type RawBattleResult struct {
	BattleNumber string   `json:"battle_number"`
	Rule         Rule     `json:"rule"`
	Type         string   `json:"type"`
	Stage        Stage    `json:"stage"`
	GameMode     GameMode `json:"game_mode"`
	Rank
	PlayResult       PlayResult `json:"play_result"`
	StartTime        int64      `json:"start_time"`
	MyTeamResult     TeamResult `json:"my_team_result"`
	OtherTeamResult  TeamResult `json:"other_team_result"`
	WeaponPaintPoint int32      `json:"weapon_paint_point"`

	// Regular
	WinMeter            float32 `json:"win_meter"`
	MyTeamPercentage    float32 `json:"my_team_percentage"`
	OtherTeamPercentage float32 `json:"other_team_percentage"`

	// Gachi
	XPower             float32 `json:"x_power"`
	ElapsedTime        int32   `json:"elapsed_time"`
	MyTeamCount        int32   `json:"my_team_count"`
	OtherTeamCount     int32   `json:"other_team_count"`
	EstimateGachiPower float32 `json:"estimate_gachi_power"`
	EstimateXPower     float32 `json:"estimate_x_power"`
	Udemae             Udemae  `json:"udemae"`
	// todo: CrownPlayers

	// League
	MaxLeaguePoint               float32 `json:"max_league_point"`
	MyTeamEstimateLeaguePoint    float32 `json:"my_team_estimate_league_point"`
	OtherTeamEstimateLeaguePoint float32 `json:"other_team_estimate_league_point"`
	TagID                        string  `json:"tag_id"`
	LeaguePoint                  float32 `json:"league_point"`
	// Duplicated
	// MyTeamCount               int32   `json:"my_team_count"`
	// OtherTeamCount            int32   `json:"other_team_count"`
	// Udemae                    Udemae  `json:"udemae"`

	// Festival
	FesMode                 GameMode  `json:"fes_mode"`
	FesID                   int64     `json:"fes_id"`
	EventType               EventType `json:"event_type"`
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
	UniformBonus            float32   `json:"uniform_bonus"`
	Version                 int32     `json:"version"`
	FesGrade                FesGrade  `json:"fes_grade"`
	ContributionPoint       float32   `json:"contribution_point"`
	FesPoint                float32   `json:"fes_point"`
	ContributionPointTotal  float32   `json:"contribution_point_total"`
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
func (r FesBattleResult) Type() BattleResultType {
	return BattleResultTypeEnum.Festival
}

// Metadata return BattleResultMetadata of FesBattleResult.
func (r FesBattleResult) Metadata() BattleResultMetadata {
	return r.BattleResultMetadata
}
