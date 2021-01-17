package stage

import (
	"image"
	"image/draw"
	"sort"
	"strings"

	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	log "telegram-splatoon2-bot/common/log"
	"telegram-splatoon2-bot/service/db"
	imageSvc "telegram-splatoon2-bot/service/image"
	nintendoSvc "telegram-splatoon2-bot/service/nintendo"
	"telegram-splatoon2-bot/service/repository"
	"telegram-splatoon2-bot/service/todo"
	"telegram-splatoon2-bot/service/user"
)

type UID = int64

type CompositeSchedule struct {
	FileID   imageSvc.Identifier
	Schedule *nintendoSvc.StageSchedule
}

type content struct {
	regularSchedules []CompositeSchedule
	gachiSchedules   []CompositeSchedule
	leagueSchedules  []CompositeSchedule
}

type repoImpl struct {
	imageSvc  imageSvc.Service
	userSvc   user.Service
	schedules *content
	dumper    *DumperImpl
}

func NewStageScheduleRepo(userSvc user.Service, imageSvc imageSvc.Service, config Config) (repository.Repository, error) {
	dumper := NewDumper(config.dumper)
	err := dumper.Load()
	if err != nil {
		return nil, errors.Wrap(err, "can't load stage dumping files")
	}
	return &repoImpl{
		userSvc: userSvc,
		dumper:  dumper,
	}, nil
}

func (repo *repoImpl) HasInit() bool {
	return repo.schedules != nil
}

func (repo *repoImpl) Name() string {
	return "StageScheduleRepo"
}

func (repo *repoImpl) Update() error {
	var err error
	for _, admin := range repo.userSvc.Admins() {
		err := repo.updateByUid(admin)
		if err == nil {
			return nil
		}
	}
	return errors.Wrap(err, "can't update stage schedules by admins")
}

func (repo *repoImpl) updateByUid(uid UID) error {
	wrapper := func(iksm string, timezone int, acceptLang string, _ ...interface{}) (interface{}, error) {
		return nintendoSvc.GetStageSchedules(iksm, timezone, acceptLang)
	}
	result, err := todo.FetchResourceWithUpdate(uid, wrapper)
	if err != nil {
		return errors.Wrap(err, "can't fetch runtime")
	}
	schedules := result.(*nintendoSvc.StageSchedules)

	repo.sortSchedules(schedules)
	repo.populateSchedules(schedules)
	err = repo.dumper.Update(schedules)
	if err != nil {
		log.Warn("can't update stage dumping file", zap.Error(err))
	}
	err = repo.dumper.Save()
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

func (repo *repoImpl) sortSchedules(stageSchedules *nintendoSvc.StageSchedules) {
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

func (repo *repoImpl) populateSchedules(stageSchedules *nintendoSvc.StageSchedules) {
	for _, stage := range stageSchedules.Regular {
		stage.StageA.Image = nintendoSvc.Endpoint + stage.StageA.Image
		stage.StageB.Image = nintendoSvc.Endpoint + stage.StageB.Image
	}
	for _, stage := range stageSchedules.Gachi {
		stage.StageA.Image = nintendoSvc.Endpoint + stage.StageA.Image
		stage.StageB.Image = nintendoSvc.Endpoint + stage.StageB.Image
	}
	for _, stage := range stageSchedules.League {
		stage.StageA.Image = nintendoSvc.Endpoint + stage.StageA.Image
		stage.StageB.Image = nintendoSvc.Endpoint + stage.StageB.Image
	}
}

func (repo *repoImpl) getNewItems(stages []*nintendoSvc.StageSchedule, imageIDs []CompositeSchedule) []*nintendoSvc.StageSchedule {
	var lastUpdateTimestamp int64
	if len(imageIDs) > 0 {
		lastUpdateTimestamp = imageIDs[len(imageIDs)-1].Schedule.StartTime
	}
	for i := 0; i < len(stages); i++ {
		if lastUpdateTimestamp < stages[i].StartTime {
			return stages[i:]
		}
	}
	return make([]*nintendoSvc.StageSchedule, 0)
}

func (repo *repoImpl) wrapSchedules(stageSchedules *nintendoSvc.StageSchedules) (*content, error) {
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
	imgs, err := repo.imageSvc.DownloadAll(urls)
	if err != nil {
		return nil, err
	}

	concatImgs := make([]image.Image, 0)
	offset := 0
	regularImgs := imgs[offset : regularNewStagesCount*2]
	offset += regularNewStagesCount * 2
	gachiImgs := imgs[offset : offset+gachiNewStagesCount*2]
	offset += gachiNewStagesCount * 2
	leagueImgs := imgs[offset : offset+leagueNewStagesCount*2]
	for i := 0; i < len(regularImgs); i += 2 {
		img := drawImage(regularImgs[i], regularImgs[i+1])
		concatImgs = append(concatImgs, img)
	}
	for i := 0; i < len(gachiImgs); i += 2 {
		img := drawImage(gachiImgs[i], gachiImgs[i+1])
		concatImgs = append(concatImgs, img)
	}
	for i := 0; i < len(leagueImgs); i += 2 {
		img := drawImage(leagueImgs[i], leagueImgs[i+1])
		concatImgs = append(concatImgs, img)
	}
	ids, err := repo.imageSvc.UploadAll(concatImgs)
	if err != nil {
		return nil, errors.Wrap(err, "can't upload images")
	}

	newSchedules := &content{
		regularSchedules: make([]CompositeSchedule, 0),
		gachiSchedules:   make([]CompositeSchedule, 0),
		leagueSchedules:  make([]CompositeSchedule, 0),
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
		newSchedules.regularSchedules = append(newSchedules.regularSchedules, CompositeSchedule{
			FileID:   ids[offset+i],
			Schedule: regularNewStages[i],
		})
	}
	offset += regularNewStagesCount
	for i := 0; i < gachiNewStagesCount; i++ {
		newSchedules.gachiSchedules = append(newSchedules.gachiSchedules, CompositeSchedule{
			FileID:   ids[offset+i],
			Schedule: gachiNewStages[i],
		})
	}
	offset += gachiNewStagesCount
	for i := 0; i < leagueNewStagesCount; i++ {
		newSchedules.leagueSchedules = append(newSchedules.leagueSchedules, CompositeSchedule{
			FileID:   ids[offset+i],
			Schedule: leagueNewStages[i],
		})
	}
	return newSchedules, nil
}

// QueryStageSchedules handle this  command: /stage ([rgl]+)? ((\d+)|([czrt]+)|(b\d+-\d+))+
func QueryStageSchedules(update *botApi.Update) error {
	user := update.Message.From
	schedules := stageScheduleRepo.schedules
	if schedules == nil {
		return errors.Errorf("no cached schedules")
	}
	runtime, err := todo.FetchRuntime(int64(user.ID))
	if err != nil {
		return errors.Wrap(err, "can't fetch runtime")
	}

	var args = []string{"lgr", "1"} // default filter
	argsText := update.Message.CommandArguments()
	if argsText != "" {
		args = strings.Split(argsText, " ")
	}
	if len(StageScheduleHelper.primaryFilterRegExp.FindStringSubmatch(args[0])) == 0 {
		primaryFilterArg := "lgr"
		idx := StageScheduleHelper.firstIndexOfSecondaryFilterParam(args[0])
		if idx > 0 {
			primaryFilterArg = args[0][:idx]
			args[0] = args[0][idx:]
		}
		args = append([]string{primaryFilterArg}, args...) // add primary filter
	}
	primaryFilter, err := NewPrimaryFilter(args[0])
	if err != nil {
		msg := StageScheduleHelper.newFilterErrorMessage(update.Message.Chat.ID, runtime, user)
		_ = service.sendWithRetry(service.bot, msg)
		return err
	}
	secondaryFilters := make([]SecondaryFilter, 0)
	for _, arg := range args[1:] {
		f, err := NewSecondaryFilter(arg, runtime.Timezone)
		if err != nil {
			msg := StageScheduleHelper.newFilterErrorMessage(update.Message.Chat.ID, runtime, user)
			_ = service.sendWithRetry(service.bot, msg)
			return err
		}
		secondaryFilters = append(secondaryFilters, f)
	}
	if len(secondaryFilters) == 0 {
		secondaryFilters = append(secondaryFilters, NewNextNSecondaryFilter("2"))
	}

	stages := primaryFilter.Filter(schedules, secondaryFilters, service.proposedStageNumber)
	if len(stages) >= service.proposedStageNumber {
		msg := StageScheduleHelper.newNumberWarningMessage(update.Message.Chat.ID, runtime, user)
		err = service.sendWithRetry(service.bot, msg)
		if err != nil {
			return err
		}
	}
	for i := len(stages) - 1; i >= 0; i-- {
		msg := StageScheduleHelper.formatStage(stages[i], update.Message.Chat.ID, runtime, user)
		err = service.sendWithRetry(service.bot, msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (stageScheduleHelper) formatStage(stage CompositeSchedule, chatID int64, runtime *db.Runtime, user *botApi.User) botApi.Chattable {
	msg := botApi.NewPhotoShare(chatID, stage.FileID)
	timeTemplate := service.getI18nText(runtime.Language, user, todo.NewI18nKey(todo.TimeTemplateTextKey))[0]
	startTime := todo.TimeHelper.getLocalTime(stage.Schedule.StartTime, runtime.Timezone).Format(timeTemplate)
	endTime := todo.TimeHelper.getLocalTime(stage.Schedule.EndTime, runtime.Timezone).Format(timeTemplate)
	texts := service.getI18nText(runtime.Language, user, todo.NewI18nKey(service.stageSchedulesImageCaptionTextKey,
		startTime, endTime,
		stage.Schedule.GameMode.Name, stage.Schedule.Rule.Name,
		stage.Schedule.StageB.Name, stage.Schedule.StageA.Name,
		strings.Replace(stage.Schedule.GameMode.Name, " ", `\_`, -1),
		strings.Replace(stage.Schedule.Rule.Name, " ", `\_`, -1),
	))
	msg.Caption = texts[0]
	msg.ParseMode = "Markdown"
	return msg
}

func (stageScheduleHelper) newFilterErrorMessage(chatID int64, runtime *db.Runtime, user *botApi.User) botApi.Chattable {
	texts := service.getI18nText(runtime.Language, user, todo.NewI18nKey(service.stageSchedulesFilterErrorTextKey))
	msg := botApi.NewMessage(chatID, texts[0])
	msg.ParseMode = "Markdown"
	return msg
}

func (stageScheduleHelper) newNumberWarningMessage(chatID int64, runtime *db.Runtime, user *botApi.User) botApi.Chattable {
	texts := service.getI18nText(runtime.Language, user, todo.NewI18nKey(service.stageSchedulesNumberWarningTextKey))
	msg := botApi.NewMessage(chatID, texts[0])
	msg.ParseMode = "Markdown"
	return msg
}

func MinInt(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func drawImage(imgA, imgB image.Image) image.Image {
	width := MinInt(imgA.Bounds().Dx(), imgB.Bounds().Dx())
	halfHeight := MinInt(imgA.Bounds().Dy(), imgB.Bounds().Dy())
	height := halfHeight * 2
	// resize
	imgA = resize.Resize(uint(width), uint(halfHeight), imgA, resize.Lanczos3)
	imgB = resize.Resize(uint(width), uint(halfHeight), imgB, resize.Lanczos3)
	// prepare canvas
	r := image.Rectangle{Min: image.Point{}, Max: image.Point{X: width, Y: height}}
	rgba := image.NewRGBA(r)
	//draw: imgB is in top half, imgA is in bottom half
	draw.Draw(rgba,
		image.Rectangle{Min: image.Point{}, Max: image.Point{X: width, Y: height / 2}},
		imgB, image.Point{}, draw.Src)
	draw.Draw(rgba,
		image.Rectangle{Min: image.Point{Y: height / 2}, Max: image.Point{X: width, Y: height}},
		imgA, image.Point{}, draw.Src)
	return rgba
}
