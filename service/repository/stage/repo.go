package stage

import (
	"image"
	"image/draw"
	"sort"
	"time"

	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"telegram-splatoon2-bot/common/log"
	"telegram-splatoon2-bot/common/util"
	imageSvc "telegram-splatoon2-bot/service/image"
	"telegram-splatoon2-bot/service/language"
	"telegram-splatoon2-bot/service/nintendo"
	"telegram-splatoon2-bot/service/repository"
	"telegram-splatoon2-bot/service/repository/internal/dump"
	"telegram-splatoon2-bot/service/user"
)

var (
	dumperKey = struct {
		Stage string
	}{"stage"}
)

// WrappedSchedule stores stage schedules and downloaded images.
type WrappedSchedule struct {
	ImageID  imageSvc.Identifier
	Schedule nintendo.StageSchedule
}

type content struct {
	RegularSchedules []WrappedSchedule
	GachiSchedules   []WrappedSchedule
	LeagueSchedules  []WrappedSchedule
}

// Repository fetches stage schedules.
type Repository interface {
	repository.Repository
	// Content returns stage schedules filtered by primaryFilter and secondaryFilters.
	// If result number > limit, the excess will be omitted.
	Content(primaryFilter PrimaryFilter, secondaryFilters []SecondaryFilter, limit int) []WrappedSchedule
}

type repoImpl struct {
	nintendoSvc nintendo.Service
	imageSvc    imageSvc.Service
	userSvc     user.Service
	dumper      dump.Dumper

	writerChan chan *content
	readerChan chan (chan *content)

	content *content
}

func (repo *repoImpl) Content(primaryFilter PrimaryFilter, secondaryFilters []SecondaryFilter, limit int) []WrappedSchedule {
	content := repo.getContent()
	if content == nil {
		return nil
	}
	return content.Filter(primaryFilter, secondaryFilters, limit)
}

// NewRepository return a Repository object.
func NewRepository(nintendoSvc nintendo.Service, userSvc user.Service, imageSvc imageSvc.Service, config Config) Repository {
	dumpConfig := dump.Config{}
	dumpConfig.AddTarget(dumperKey.Stage, config.Dumper.StageFile)
	dumper := dump.New(dumpConfig)
	ret := &repoImpl{
		nintendoSvc: nintendoSvc,
		userSvc:     userSvc,
		imageSvc:    imageSvc,
		dumper:      dumper,
		writerChan:  make(chan *content),
		readerChan:  make(chan chan *content),
	}
	ret.runUpdater()
	return ret
}

func (repo *repoImpl) getContent() *content {
	retChan := make(chan *content)
	defer close(retChan)
	repo.readerChan <- retChan
	return <-retChan
}

func (repo *repoImpl) NextUpdateTime() time.Time {
	if repo.getContent() == nil {
		return time.Now()
	}
	return util.Time.SplatoonNextUpdateTime(time.Now())
}

func (repo *repoImpl) Name() string {
	return "Stage Schedule"
}

func (repo *repoImpl) Update() error {
	var err error
	admins := repo.userSvc.Admins()
	if len(admins) == 0 {
		return errors.New("no admin")
	}
	for _, admin := range admins {
		err = repo.updateByUID(admin)
		if err == nil {
			return nil
		}
	}
	return errors.Wrap(err, "can't update stage schedules by admins")
}

func (repo *repoImpl) updateByUID(uid user.ID) error {
	status, err := repo.userSvc.GetStatus(uid)
	if err != nil {
		return errors.Wrap(err, "can't fetch admin status")
	}
	schedules, err := repo.nintendoSvc.GetStageSchedules(status.IKSM, status.Timezone, language.English)
	if errors.Is(err, &nintendo.ErrIKSMExpired{}) {
		status, err = repo.userSvc.UpdateStatusIKSM(uid)
		if err != nil {
			return errors.Wrap(err, "can't update IKSM when fetching salmon schedules")
		}
		schedules, err = repo.nintendoSvc.GetStageSchedules(status.IKSM, status.Timezone, language.English)
	}
	if err != nil {
		return errors.Wrap(err, "can't fetch stage schedules")
	}

	repo.sortSchedules(&schedules)
	repo.populateSchedules(&schedules)
	err = repo.updateDumper(schedules)
	if err != nil {
		// go ahead
		log.Warn("can't update stage dumping file", zap.Error(err))
	}

	wrappedSchedules, err := repo.wrapSchedules(&schedules)
	if err != nil {
		return errors.Wrap(err, "can't upload stage schedules images")
	}
	repo.writerChan <- wrappedSchedules
	return nil
}

func (repo *repoImpl) sortSchedules(stageSchedules *nintendo.StageSchedules) {
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

func (repo *repoImpl) populateSchedules(stageSchedules *nintendo.StageSchedules) {
	for i := range stageSchedules.Regular {
		stage := &stageSchedules.Regular[i]
		stage.StageA.Image = nintendo.Endpoint + stage.StageA.Image
		stage.StageB.Image = nintendo.Endpoint + stage.StageB.Image
	}
	for i := range stageSchedules.Gachi {
		stage := &stageSchedules.Gachi[i]
		stage.StageA.Image = nintendo.Endpoint + stage.StageA.Image
		stage.StageB.Image = nintendo.Endpoint + stage.StageB.Image
	}
	for i := range stageSchedules.League {
		stage := &stageSchedules.League[i]
		stage.StageA.Image = nintendo.Endpoint + stage.StageA.Image
		stage.StageB.Image = nintendo.Endpoint + stage.StageB.Image
	}
}

func (repo *repoImpl) getNewItems(stages []nintendo.StageSchedule, imageIDs []WrappedSchedule) []nintendo.StageSchedule {
	var lastUpdateTimestamp int64
	if len(imageIDs) > 0 {
		lastUpdateTimestamp = imageIDs[len(imageIDs)-1].Schedule.StartTime
	}
	for i := 0; i < len(stages); i++ {
		if lastUpdateTimestamp < stages[i].StartTime {
			return stages[i:]
		}
	}
	return make([]nintendo.StageSchedule, 0)
}

func (repo *repoImpl) wrapSchedules(schedules *nintendo.StageSchedules) (*content, error) {
	regularNewStages := schedules.Regular
	gachiNewStages := schedules.Gachi
	leagueNewStages := schedules.League
	c := repo.getContent()
	if c != nil {
		regularNewStages = repo.getNewItems(schedules.Regular, c.RegularSchedules)
		gachiNewStages = repo.getNewItems(schedules.Gachi, c.GachiSchedules)
		leagueNewStages = repo.getNewItems(schedules.League, c.LeagueSchedules)
	}
	regularNewStagesCount := len(regularNewStages)
	gachiNewStagesCount := len(gachiNewStages)
	leagueNewStagesCount := len(leagueNewStages)
	log.Info("found new stages",
		zap.Int("regular", regularNewStagesCount),
		zap.Int("gachi", gachiNewStagesCount),
		zap.Int("league", leagueNewStagesCount),
	)
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
		RegularSchedules: make([]WrappedSchedule, 0),
		GachiSchedules:   make([]WrappedSchedule, 0),
		LeagueSchedules:  make([]WrappedSchedule, 0),
	}

	if c != nil {
		if len(c.RegularSchedules) > 0 {
			newSchedules.RegularSchedules = c.RegularSchedules[regularNewStagesCount:]
		}
		if len(c.GachiSchedules) > 0 {
			newSchedules.GachiSchedules = c.GachiSchedules[gachiNewStagesCount:]
		}
		if len(c.LeagueSchedules) > 0 {
			newSchedules.LeagueSchedules = c.LeagueSchedules[leagueNewStagesCount:]
		}
	}

	offset = 0
	for i := 0; i < regularNewStagesCount; i++ {
		newSchedules.RegularSchedules = append(newSchedules.RegularSchedules, WrappedSchedule{
			ImageID:  ids[offset+i],
			Schedule: regularNewStages[i],
		})
	}
	offset += regularNewStagesCount
	for i := 0; i < gachiNewStagesCount; i++ {
		newSchedules.GachiSchedules = append(newSchedules.GachiSchedules, WrappedSchedule{
			ImageID:  ids[offset+i],
			Schedule: gachiNewStages[i],
		})
	}
	offset += gachiNewStagesCount
	for i := 0; i < leagueNewStagesCount; i++ {
		newSchedules.LeagueSchedules = append(newSchedules.LeagueSchedules, WrappedSchedule{
			ImageID:  ids[offset+i],
			Schedule: leagueNewStages[i],
		})
	}
	return newSchedules, nil
}

func (repo *repoImpl) updateDumper(schedules nintendo.StageSchedules) error {
	stagesPtr, err := repo.dumper.Get(dumperKey.Stage, &stageCollection{})
	if err != nil {
		return errors.Wrap(err, "can't get stage dumper object")
	}
	stages := *stagesPtr.(*stageCollection)
	for _, stage := range schedules.Regular {
		stages[stage.StageA.Name] = stage.StageA
		stages[stage.StageB.Name] = stage.StageB
	}
	for _, stage := range schedules.Gachi {
		stages[stage.StageA.Name] = stage.StageA
		stages[stage.StageB.Name] = stage.StageB
	}
	for _, stage := range schedules.League {
		stages[stage.StageA.Name] = stage.StageA
		stages[stage.StageB.Name] = stage.StageB
	}
	err = repo.dumper.Save(dumperKey.Stage, &stages)
	if err != nil {
		return errors.Wrap(err, "can't save stage dumper object")
	}
	return nil
}

func (repo *repoImpl) runUpdater() {
	go func() {
		for {
			select {
			case content := <-repo.writerChan:
				repo.content = content
			case rc := <-repo.readerChan:
				rc <- repo.content
			}
		}
	}()
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func drawImage(imgA, imgB image.Image) image.Image {
	width := minInt(imgA.Bounds().Dx(), imgB.Bounds().Dx())
	halfHeight := minInt(imgA.Bounds().Dy(), imgB.Bounds().Dy())
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
