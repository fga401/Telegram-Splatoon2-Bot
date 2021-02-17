package nintendo

// Common

// Endpoint of Nintendo API
var Endpoint = "https://app.splatoon2.nintendo.net"

type ScheduleTime struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}
type Rule struct {
	MultilineName string `json:"multiline_name"`
	Name          string `json:"name"`
	Key           string `json:"key"`
}
type Stage struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}
type GameMode struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}
type Udemae struct {
	Number          int32  `json:"number"`
	IsNumberReached bool   `json:"is_number_reached"`
	IsX             bool   `json:"is_x"`
	SPlusNumber     int32 `json:"s_plus_number"`
	Name            string `json:"name"`
}
type Rank struct {
	PlayerRank int32 `json:"player_rank"`
	StarRank   int32 `json:"star_rank"`
}
type PlayType struct {
	Species string `json:"species"`
	Style   string `json:"style"`
}
type GearSkill struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}
type GearSkills struct {
	Subs []GearSkill `json:"subs"`
	Main GearSkill   `json:"main"`
}
type Brand struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Image         string    `json:"image"`
	FrequentSkill GearSkill `json:"frequent_skill"`
}
type Gear struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Brand     Brand  `json:"brand"`
	Thumbnail string `json:"thumbnail"`
	Image     string `json:"image"`
	Rarity    int32  `json:"rarity"`
}
type WeaponSkill struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	ImageA string `json:"image_a"`
	ImageB string `json:"image_b"`
}
type Weapon struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Image     string      `json:"image"`
	Thumbnail string      `json:"thumbnail"`
	Special   WeaponSkill `json:"special"`
	Sub       WeaponSkill `json:"sub"`
}
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

// Battle

var BattleResultTypeEnum = struct {
	Regular  BattleResultType
	Gachi    BattleResultType
	League   BattleResultType
	Festival BattleResultType
}{"regular", "gachi", "league", "fes"}

type BattleResultType string
type BattleResult interface {
	Type() BattleResultType
	Metadata() BattleResultMetadata
}

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
type RegularBattleResult struct {
	BattleResultMetadata
	WinMeter            float32 `json:"win_meter"`
	MyTeamPercentage    float32 `json:"my_team_percentage"`
	OtherTeamPercentage float32 `json:"other_team_percentage"`
}
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
type LeagueBattleResult struct {
	BattleResultMetadata
	TagID                    string  `json:"tag_id"`
	Udemae                   Udemae `json:"udemae"`
	ElapsedTime              int32   `json:"elapsed_time"`
	MyTeamCount              int32   `json:"my_team_count"`
	OtherTeamCount           int32   `json:"other_team_count"`
	LeaguePoint              float32 `json:"league_point"`
	MaxLeaguePoint           float32 `json:"max_league_point"`
	MyEstimateLeaguePoint    float32 `json:"my_estimate_league_point"`
	OtherEstimateLeaguePoint float32 `json:"other_estimate_league_point"`
	EstimateGachiPower       float32 `json:"estimate_gachi_power"`
}

func (r RegularBattleResult) Type() BattleResultType {
	return BattleResultTypeEnum.Regular
}
func (r GachiBattleResult) Type() BattleResultType {
	return BattleResultTypeEnum.Gachi
}
func (r LeagueBattleResult) Type() BattleResultType {
	return BattleResultTypeEnum.League
}

func (r RegularBattleResult) Metadata() BattleResultMetadata {
	return r.BattleResultMetadata
}
func (r GachiBattleResult) Metadata() BattleResultMetadata {
	return r.BattleResultMetadata
}
func (r LeagueBattleResult) Metadata() BattleResultMetadata {
	return r.BattleResultMetadata
}

type BattleResults struct {
	ID      string         `json:"unique_id"`
	Summary BattleSummary `json:"summary"`
	Results []BattleResult `json:"results"`
}
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
type TeamResult struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}
type PlayResult struct {
	GamePaintPoint float32 `json:"game_paint_point"`
	DeathCount     int32   `json:"death_count"`
	KillCount      int32   `json:"kill_count"`
	AssertCount    int32   `json:"assert_count"`
	SpecialCount   int32   `json:"special_count"`
	SortScore      int32   `json:"sort_score"`
	Player         Player `json:"player"`
}

// Raw Battle

type RawBattleResult struct {
	BattleNumber string    `json:"battle_number"`
	Rule         Rule     `json:"rule"`
	Type         string    `json:"type"`
	Stage        Stage    `json:"stage"`
	GameMode     GameMode `json:"game_mode"`
	Rank
	PlayResult       PlayResult `json:"play_result"`
	StartTime        int64       `json:"start_time"`
	MyTeamResult     TeamResult `json:"my_team_result"`
	OtherTeamResult  TeamResult `json:"other_team_result"`
	WeaponPaintPoint int32       `json:"weapon_paint_point"`

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

// Salmon

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
type SalmonWeaponType struct {
	ID            string               `json:"id"`
	Weapon        *SalmonWeapon        `json:"weapon"`
	SpecialWeapon *SalmonSpecialWeapon `json:"coop_special_weapon"`
}
type SalmonStage struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}
type SalmonScheduleDetail struct {
	Weapons []SalmonWeaponType `json:"weapons"`
	Stage   SalmonStage        `json:"stage"`
	ScheduleTime
}
type SalmonSchedules struct {
	Details   []SalmonScheduleDetail `json:"details"`
	Schedules []ScheduleTime         `json:"schedules"`
}

// Stage

type StageSchedule struct {
	ID       int64     `json:"id"`
	Rule     Rule     `json:"rule"`
	GameMode GameMode `json:"game_mode"`
	StageA   Stage    `json:"stage_a"`
	StageB   Stage    `json:"stage_b"`
	ScheduleTime
}
type StageSchedules struct {
	League  []StageSchedule `json:"league"`
	Gachi   []StageSchedule `json:"gachi"`
	Regular []StageSchedule `json:"regular"`
}

// Festival

type Color struct {
	A      float64 `json:"a"`
	B      float64 `json:"b"`
	R      float64 `json:"r"`
	G      float64 `json:"g"`
	CSSRGB string  `json:"css_rgb"`
}
type Theme struct {
	Color Color `json:"color"`
	Key   string `json:"key"`
	Name  string `json:"name"`
}
type EventType struct {
	MultilineName string `json:"multiline_name"`
	ClassName     string `json:"class_name"`
	Key           string `json:"key"`
	Name          string `json:"name"`
}
type FesGrade struct {
	Name string `json:"name"`
	Rank string `json:"rank"`
}

type FesBattleResult struct {
	RegularBattleResult
	FesMode                 GameMode  `json:"fes_mode"`
	FesID                   int64      `json:"fes_id"`
	UniformBonus            float32    `json:"uniform_bonus"`
	EventType               EventType `json:"event_type"`
	FesGrade                FesGrade  `json:"fes_grade"`
	MyTeamFesTheme          Theme     `json:"my_team_fes_theme"`
	OtherTeamFesTheme       Theme     `json:"other_team_fes_theme"`
	MyEstimateFesPower      float32    `json:"my_estimate_fes_power"`
	OtherEstimateFesPower   float32    `json:"other_estimate_fes_power"`
	MyTeamAnotherName       string     `json:"my_team_another_name"`
	OtherTeamAnotherName    string     `json:"other_team_another_name"`
	MyTeamConsecutiveWin    int32      `json:"my_team_consecutive_win"`
	OtherTeamConsecutiveWin int32      `json:"other_team_consecutive_win"`
	FesPower                float32    `json:"fes_power"`
	MaxFesPower             float32    `json:"max_fes_power"`
	FesPoint                float32    `json:"fes_point"`
	ContributionPoint       float32    `json:"contribution_point"`
	ContributionPointTotal  float32    `json:"contribution_point_total"`
	Version                 int32      `json:"version"`
}

func (r FesBattleResult) Type() BattleResultType {
	return BattleResultTypeEnum.Festival
}

func (r FesBattleResult) Metadata() BattleResultMetadata {
	return r.BattleResultMetadata
}
