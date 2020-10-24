package service

import (
	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"image"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	log "telegram-splatoon2-bot/logger"
	"telegram-splatoon2-bot/nintendo"
	"telegram-splatoon2-bot/service/db"
	"time"
)

var (
	stageScheduleRepo *StageScheduleRepo
)

type StageDumpling struct {
	stageFileName string
	stages        map[string]*nintendo.Stage
}

func NewStageDumpling() *StageDumpling {
	return &StageDumpling{
		stageFileName: viper.GetString("service.stage.stageFileName"),
		stages:        make(map[string]*nintendo.Stage),
	}
}

func (d *StageDumpling) Save() error {
	err := marshalToFile(d.stageFileName, d.stages)
	if err != nil {
		return errors.Wrap(err, "can't save stage")
	}
	return nil
}

func (d *StageDumpling) Load() error {
	if _, err := os.Stat(d.stageFileName); err == nil {
		if err := unmarshalFromFile(d.stageFileName, &d.stages); err != nil {
			return errors.Wrap(err, "can't load stage stage")
		}
	} else {
		log.Warn("can't open stage file", zap.Error(err))
	}
	return nil
}

func (d *StageDumpling) Update(src interface{}) error {
	stageSchedules, ok := src.(*nintendo.StageSchedules)
	if !ok {
		return errors.Errorf("unknown input type")
	}
	for _, stage := range stageSchedules.Regular {
		d.stages[stage.StageA.Name] = stage.StageA
		d.stages[stage.StageB.Name] = stage.StageB
	}
	for _, stage := range stageSchedules.Gachi {
		d.stages[stage.StageA.Name] = stage.StageA
		d.stages[stage.StageB.Name] = stage.StageB
	}
	for _, stage := range stageSchedules.League {
		d.stages[stage.StageA.Name] = stage.StageA
		d.stages[stage.StageB.Name] = stage.StageB
	}
	return nil
}

type StageScheduleWrapper struct {
	FileID   string
	Schedule *nintendo.StageSchedule
}

type StageSchedulesWrapper struct {
	regularSchedules []StageScheduleWrapper
	gachiSchedules   []StageScheduleWrapper
	leagueSchedules  []StageScheduleWrapper
}

type StageScheduleRepo struct {
	schedules    *StageSchedulesWrapper
	admins       *SyncUserSet
	stageDumping *StageDumpling
}

func NewStageScheduleRepo(admins *SyncUserSet) (*StageScheduleRepo, error) {
	stageDumping := NewStageDumpling()
	err := stageDumping.Load()
	if err != nil {
		return nil, errors.Wrap(err, "can't load stage dumping files")
	}
	return &StageScheduleRepo{
		admins:       admins,
		stageDumping: stageDumping,
	}, nil
}

func (repo *StageScheduleRepo) HasInit() bool {
	return repo.schedules != nil
}

func (repo *StageScheduleRepo) RepoName() string {
	return "StageScheduleRepo"
}

func (repo *StageScheduleRepo) Update() error {
	var err error
	repo.admins.Range(func(uid int64) (continued bool) {
		err = repo.updateByUid(uid)
		return err != nil
	})
	if err == nil {
		return nil
	}
	log.Warn("can't update stage schedules by admin", zap.Error(err))
	runtime, err := RuntimeTable.GetFirstRuntime()
	if err != nil {
		return errors.Wrap(err, "can't get first runtime object")
	}
	err = repo.updateByUid(runtime.Uid)
	if err != nil {
		return errors.Wrap(err, "can't update stage schedules by other user")
	}
	return nil
}

func (repo *StageScheduleRepo) updateByUid(uid int64) error {
	wrapper := func(iksm string, timezone int, acceptLang string, _ ...interface{}) (interface{}, error) {
		return nintendo.GetStageSchedules(iksm, timezone, acceptLang)
	}
	result, err := fetchResourceWithUpdate(uid, wrapper)
	if err != nil {
		return errors.Wrap(err, "can't fetch runtime")
	}
	schedules := result.(*nintendo.StageSchedules)

	repo.sortSchedules(schedules)
	repo.populateFields(schedules)
	err = repo.stageDumping.Update(schedules)
	if err != nil {
		log.Warn("can't update stage dumping file", zap.Error(err))
	}
	err = repo.stageDumping.Save()
	if err != nil {
		log.Warn("can't update stage dumping file", zap.Error(err))
	} else {
		log.Info("dumped stages to files")
	}

	wrappedSchedules, err := repo.wrapSchedules(schedules)
	if err != nil {
		return errors.Wrap(err, "can't upload stage schedules images")
	}
	repo.schedules = wrappedSchedules
	return nil
}

func (repo *StageScheduleRepo) sortSchedules(stageSchedules *nintendo.StageSchedules) {
	// sort by start time in ascendant order
	sort.Slice(stageSchedules.League, func(i, j int) bool {
		return stageSchedules.League[i].StartTime < stageSchedules.League[j].StartTime
	})
	sort.Slice(stageSchedules.Gachi, func(i, j int) bool {
		return stageSchedules.Gachi[i].StartTime < stageSchedules.Gachi[j].StartTime
	})
	sort.Slice(stageSchedules.Regular, func(i, j int) bool {
		return stageSchedules.Regular[i].StartTime < stageSchedules.Regular[j].StartTime
	})
}

func (repo *StageScheduleRepo) populateFields(stageSchedules *nintendo.StageSchedules) {
	for _, stage := range stageSchedules.Regular {
		stage.StageA.Image = nintendo.Host + stage.StageA.Image
		stage.StageB.Image = nintendo.Host + stage.StageB.Image
	}
	for _, stage := range stageSchedules.Gachi {
		stage.StageA.Image = nintendo.Host + stage.StageA.Image
		stage.StageB.Image = nintendo.Host + stage.StageB.Image
	}
	for _, stage := range stageSchedules.League {
		stage.StageA.Image = nintendo.Host + stage.StageA.Image
		stage.StageB.Image = nintendo.Host + stage.StageB.Image
	}
}

func (repo *StageScheduleRepo) getNewItems(stages []*nintendo.StageSchedule, imageIDs []StageScheduleWrapper) []*nintendo.StageSchedule {
	var lastUpdateTimestamp int64
	if len(imageIDs) > 0 {
		lastUpdateTimestamp = imageIDs[len(imageIDs)-1].Schedule.StartTime
	}
	for i := 0; i < len(stages); i++ {
		if lastUpdateTimestamp < stages[i].StartTime {
			return stages[i:]
		}
	}
	return make([]*nintendo.StageSchedule, 0)
}

func (repo *StageScheduleRepo) wrapSchedules(stageSchedules *nintendo.StageSchedules) (*StageSchedulesWrapper, error) {
	regularNewStages := stageSchedules.Regular
	gachiNewStages := stageSchedules.Gachi
	leagueNewStages := stageSchedules.League
	if repo.HasInit() {
		regularNewStages = repo.getNewItems(stageSchedules.Regular, repo.schedules.regularSchedules)
		gachiNewStages = repo.getNewItems(stageSchedules.Gachi, repo.schedules.gachiSchedules)
		leagueNewStages = repo.getNewItems(stageSchedules.League, repo.schedules.leagueSchedules)
	}
	regularNewStagesCount := len(regularNewStages)
	gachiNewStagesCount := len(gachiNewStages)
	leagueNewStagesCount := len(leagueNewStages)
	log.Info("find new stages",
		zap.Int("regular", regularNewStagesCount),
		zap.Int("gachi", gachiNewStagesCount),
		zap.Int("league", leagueNewStagesCount))
	urls := make([]string, 0)
	for _, stage := range regularNewStages {
		urls = append(urls, stage.StageA.Image, stage.StageB.Image)
	}
	for _, stage := range gachiNewStages {
		urls = append(urls, stage.StageA.Image, stage.StageB.Image)
	}
	for _, stage := range leagueNewStages {
		urls = append(urls, stage.StageA.Image, stage.StageB.Image)
	}
	imgs, err := downloadImages(urls)
	if err != nil {
		return nil, err
	}

	now := strconv.FormatInt(time.Now().Round(time.Hour).Unix(), 10)
	concatImgs := make([]image.Image, 0)
	concatImgNames := make([]string, 0)
	offset := 0
	regularImgs := imgs[offset : regularNewStagesCount*2]
	offset += regularNewStagesCount * 2
	gachiImgs := imgs[offset : offset+gachiNewStagesCount*2]
	offset += gachiNewStagesCount * 2
	leagueImgs := imgs[offset : offset+leagueNewStagesCount*2]
	for i := 0; i < len(regularImgs); i += 2 {
		img := concatStageScheduleImage(regularImgs[i], regularImgs[i+1])
		concatImgs = append(concatImgs, img)
		concatImgNames = append(concatImgNames, "regelar_"+strconv.FormatInt(stageSchedules.Regular[i/2].StartTime, 10)+"_"+now)
	}
	for i := 0; i < len(gachiImgs); i += 2 {
		img := concatStageScheduleImage(gachiImgs[i], gachiImgs[i+1])
		concatImgs = append(concatImgs, img)
		concatImgNames = append(concatImgNames, "gachi_"+strconv.FormatInt(stageSchedules.Gachi[i/2].StartTime, 10)+"_"+now)
	}
	for i := 0; i < len(leagueImgs); i += 2 {
		img := concatStageScheduleImage(leagueImgs[i], leagueImgs[i+1])
		concatImgs = append(concatImgs, img)
		concatImgNames = append(concatImgNames, "league_"+strconv.FormatInt(stageSchedules.League[i/2].StartTime, 10)+"_"+now)
	}
	ids, err := uploadImages(concatImgs, concatImgNames)
	if err != nil {
		return nil, errors.Wrap(err, "can't upload images")
	}

	newSchedules := &StageSchedulesWrapper{
		regularSchedules: make([]StageScheduleWrapper, 0),
		gachiSchedules:   make([]StageScheduleWrapper, 0),
		leagueSchedules:  make([]StageScheduleWrapper, 0),
	}

	if repo.HasInit() {
		if len(repo.schedules.regularSchedules) > 0 {
			newSchedules.regularSchedules = repo.schedules.regularSchedules[regularNewStagesCount:]
		}
		if len(repo.schedules.gachiSchedules) > 0 {
			newSchedules.gachiSchedules = repo.schedules.gachiSchedules[gachiNewStagesCount:]
		}
		if len(repo.schedules.leagueSchedules) > 0 {
			newSchedules.leagueSchedules = repo.schedules.leagueSchedules[leagueNewStagesCount:]
		}
	}

	offset = 0
	for i := 0; i < regularNewStagesCount; i++ {
		newSchedules.regularSchedules = append(newSchedules.regularSchedules, StageScheduleWrapper{
			FileID:   ids[offset+i],
			Schedule: regularNewStages[i],
		})
	}
	offset += regularNewStagesCount
	for i := 0; i < gachiNewStagesCount; i++ {
		newSchedules.gachiSchedules = append(newSchedules.gachiSchedules, StageScheduleWrapper{
			FileID:   ids[offset+i],
			Schedule: gachiNewStages[i],
		})
	}
	offset += gachiNewStagesCount
	for i := 0; i < leagueNewStagesCount; i++ {
		newSchedules.leagueSchedules = append(newSchedules.leagueSchedules, StageScheduleWrapper{
			FileID:   ids[offset+i],
			Schedule: leagueNewStages[i],
		})
	}
	return newSchedules, nil
}

type GameModeName string

const (
	GameModeRegular GameModeName = "r"
	GameModeGachi   GameModeName = "g"
	GameModeLeague  GameModeName = "l"
)

type PrimaryFilter struct {
	orderByName []GameModeName
}

func (filter PrimaryFilter) Filter(schedules *StageSchedulesWrapper, secondaryFilters []SecondaryFilter, proposedN int) []StageScheduleWrapper {
	order := make([][]StageScheduleWrapper, 0)
	for _, name := range filter.orderByName {
		switch name {
		case GameModeLeague:
			order = append(order, schedules.leagueSchedules)
		case GameModeGachi:
			order = append(order, schedules.gachiSchedules)
		case GameModeRegular:
			order = append(order, schedules.regularSchedules)
		}
	}
	ret := make([]StageScheduleWrapper, 0)
	for i := 0; i < 12; i++ {
		for j := range order {
			s := order[j][i]
			keep := true
			for _, f := range secondaryFilters {
				if f.Filter(s) == false {
					keep = false
					break
				}
			}
			if keep {
				ret = append(ret, s)
			}
		}
		if len(ret) >= proposedN {
			break
		}
	}
	return ret
}

func NewPrimaryFilter(text string) (PrimaryFilter, error) {
	existed := make(map[GameModeName]struct{})
	ret := PrimaryFilter{}
	for _, c := range text {
		switch c {
		case 'l':
			if _, found := existed[GameModeLeague]; !found {
				ret.orderByName = append(ret.orderByName, GameModeLeague)
				existed[GameModeLeague] = struct{}{}
			}
		case 'r':
			if _, found := existed[GameModeRegular]; !found {
				ret.orderByName = append(ret.orderByName, GameModeRegular)
				existed[GameModeRegular] = struct{}{}
			}
		case 'g':
			if _, found := existed[GameModeGachi]; !found {
				ret.orderByName = append(ret.orderByName, GameModeGachi)
				existed[GameModeGachi] = struct{}{}
			}
		default:
			return ret, errors.Errorf("unknown character")
		}
	}
	return ret, nil
}

type SecondaryFilter interface {
	Filter(stage StageScheduleWrapper) bool
}

type RuleSecondaryFilter struct {
	allowRules map[string]struct{}
}

func (filter RuleSecondaryFilter) Filter(stage StageScheduleWrapper) bool {
	_, found := filter.allowRules[stage.Schedule.Rule.Key]
	return found
}
func NewRuleSecondaryFilter(text string) RuleSecondaryFilter {
	filter := RuleSecondaryFilter{allowRules: make(map[string]struct{})}
	for _, c := range text {
		switch c {
		case 'z':
			filter.allowRules["splat_zones"] = struct{}{}
		case 't':
			filter.allowRules["tower_control"] = struct{}{}
		case 'c':
			filter.allowRules["clam_blitz"] = struct{}{}
		case 'r':
			filter.allowRules["rainmaker"] = struct{}{}
		}
	}
	return filter
}

type TimeSecondaryFilter struct {
	begin, end int64
}

func (filter TimeSecondaryFilter) Filter(stage StageScheduleWrapper) bool {
	if filter.begin <= stage.Schedule.StartTime && stage.Schedule.StartTime < filter.end {
		return true
	}
	if filter.begin < stage.Schedule.EndTime && stage.Schedule.EndTime <= filter.end {
		return true
	}
	return false
}
func NewBetweenHourSecondaryFilter(begin string, end string, offset int) TimeSecondaryFilter {
	expectedBeginHour, err := strconv.Atoi(begin)
	if err != nil {
		expectedBeginHour = time.Now().Hour()
	}
	expectedEndHour, err := strconv.Atoi(end)
	if err != nil {
		expectedEndHour = (expectedBeginHour + 23) % 24
	}

	now := time.Now()
	userCurrentHour := getLocalTime(now.Unix(), offset).Hour()
	beginHourOffset := expectedBeginHour - userCurrentHour
	endHourOffset := expectedEndHour - userCurrentHour
	if beginHourOffset < 0 {
		beginHourOffset += 24
	}
	if endHourOffset < 0 {
		endHourOffset += 24
	}
	if endHourOffset-beginHourOffset < 0 {
		endHourOffset += 24
	}
	beginTime := now.Truncate(time.Hour).Add(time.Hour * time.Duration(beginHourOffset)).Unix()
	endTime := now.Truncate(time.Hour).Add(time.Hour * time.Duration(endHourOffset)).Unix()

	return TimeSecondaryFilter{begin: beginTime, end: endTime}
}
func NewNextNSecondaryFilter(text string) TimeSecondaryFilter {
	n, err := strconv.Atoi(text)
	if err != nil {
		n = 1
	}
	now := time.Now()
	beginTime := getSplatoonNextUpdateTime(now).Add(time.Hour * time.Duration(-2)).Unix()
	endTime := getSplatoonNextUpdateTime(now).Add(time.Hour * time.Duration(n*2-2)).Unix()
	return TimeSecondaryFilter{begin: beginTime, end: endTime}
}

var primaryFilterRegExp = regexp.MustCompile(`^(?P<primary>[lgrLGR]+)$`)
var ruleSecondaryFilterRegExp = regexp.MustCompile(`^(?P<primary>[czrtCZRT]+)$`)
var nextNSecondFilterRegExp = regexp.MustCompile(`(?P<n>^\d+$)`)
var betweenHourSecondFilterRegExp = regexp.MustCompile(`^[bB](?P<begin>\d+)-(?P<end>\d+)$`)

func NewSecondaryFilter(text string, offset int) (SecondaryFilter, error) {
	text = strings.ToLower(text)
	if args := ruleSecondaryFilterRegExp.FindStringSubmatch(text); len(args) != 0 {
		return NewRuleSecondaryFilter(args[1]), nil
	}
	if args := nextNSecondFilterRegExp.FindStringSubmatch(text); len(args) != 0 {
		return NewNextNSecondaryFilter(args[1]), nil
	}
	if args := betweenHourSecondFilterRegExp.FindStringSubmatch(text); len(args) != 0 {
		return NewBetweenHourSecondaryFilter(args[1], args[2], offset), nil
	}
	return nil, errors.Errorf("unknown supported filter format")
}

// QueryStageSchedules handle this  command: /stage ([rgl]+)? ((\d+)|([czrt]+)|(b\d+-\d+))+
func QueryStageSchedules(update *botapi.Update) error {
	user := update.Message.From
	schedules := stageScheduleRepo.schedules
	if schedules == nil {
		return errors.Errorf("no cached schedules")
	}
	runtime, err := fetchRuntime(int64(user.ID))
	if err != nil {
		return errors.Wrap(err, "can't fetch runtime")
	}

	var args = []string{"lgr", "1"} // default filter
	argsText := update.Message.CommandArguments()
	if argsText != "" {
		args = strings.Split(argsText, " ")
	}
	if len(primaryFilterRegExp.FindStringSubmatch(args[0])) == 0 {
		primaryFilterArg := "lgr"
		idx := firstIndexOfSecondaryFilterParam(args[0])
		if idx > 0 {
			primaryFilterArg = args[0][:idx]
			args[0] = args[0][idx:]
		}
		args = append([]string{primaryFilterArg}, args...) // add primary filter
	}
	primaryFilter, err := NewPrimaryFilter(args[0])
	if err != nil {
		msg := newFilterErrorMessage(update.Message.Chat.ID, runtime, user)
		_ = sendWithRetry(bot, msg)
		return err
	}
	secondaryFilters := make([]SecondaryFilter, 0)
	for _, arg := range args[1:] {
		f, err := NewSecondaryFilter(arg, runtime.Timezone)
		if err != nil {
			msg := newFilterErrorMessage(update.Message.Chat.ID, runtime, user)
			_ = sendWithRetry(bot, msg)
			return err
		}
		secondaryFilters = append(secondaryFilters, f)
	}
	if len(secondaryFilters) == 0 {
		secondaryFilters = append(secondaryFilters, NewNextNSecondaryFilter("2"))
	}

	stages := primaryFilter.Filter(schedules, secondaryFilters, proposedStageNumber)
	if len(stages) >= proposedStageNumber {
		msg := newNumberWarningMessage(update.Message.Chat.ID, runtime, user)
		err = sendWithRetry(bot, msg)
		if err != nil {
			return err
		}
	}
	for i := len(stages) - 1; i >= 0; i-- {
		msg := formatStage(stages[i], update.Message.Chat.ID, runtime, user)
		err = sendWithRetry(bot, msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func firstIndexOfSecondaryFilterParam(text string) int {
	for i, c := range text {
		if c != 'l' && c != 'r' && c != 'g' {
			return i
		}
	}
	return len(text)
}

func formatStage(stage StageScheduleWrapper, chatID int64, runtime *db.Runtime, user *botapi.User) botapi.Chattable {
	msg := botapi.NewPhotoShare(chatID, stage.FileID)
	timeTemplate := getI18nText(runtime.Language, user, NewI18nKey(TimeTemplateTextKey))[0]
	startTime := getLocalTime(stage.Schedule.StartTime, runtime.Timezone).Format(timeTemplate)
	endTime := getLocalTime(stage.Schedule.EndTime, runtime.Timezone).Format(timeTemplate)
	texts := getI18nText(runtime.Language, user, NewI18nKey(stageSchedulesImageCaptionTextKey,
		startTime, endTime,
		stage.Schedule.GameMode.Name, stage.Schedule.Rule.Name,
		stage.Schedule.StageB.Name, stage.Schedule.StageA.Name,
		strings.Replace(stage.Schedule.GameMode.Name, " ", "_", -1),
		strings.Replace(stage.Schedule.Rule.Name, " ", "_", -1),
	))
	msg.Caption = texts[0]
	msg.ParseMode = "Markdown"
	return msg
}

func newFilterErrorMessage(chatID int64, runtime *db.Runtime, user *botapi.User) botapi.Chattable {
	texts := getI18nText(runtime.Language, user, NewI18nKey(stageSchedulesFilterErrorTextKey))
	msg := botapi.NewMessage(chatID, texts[0])
	msg.ParseMode = "Markdown"
	return msg
}

func newNumberWarningMessage(chatID int64, runtime *db.Runtime, user *botapi.User) botapi.Chattable {
	texts := getI18nText(runtime.Language, user, NewI18nKey(stageSchedulesNumberWarningTextKey))
	msg := botapi.NewMessage(chatID, texts[0])
	msg.ParseMode = "Markdown"
	return msg
}
