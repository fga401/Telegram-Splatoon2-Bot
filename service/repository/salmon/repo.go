package salmon

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
		Weapon string
		Stage  string
	}{"weapon", "stage"}

	SchedulesIdx = struct {
		Further int
		Latest  int
	}{0, 1}

	weaponID = struct {
		Random  string
		Grizzco string
	}{"-1", "-2"}
)

type Content struct {
	Schedules nintendo.SalmonSchedules
	ImageIDs  []imageSvc.Identifier
}

type Repository interface {
	repository.Repository
	Content() *Content
}

type repoImpl struct {
	content *Content

	nintendoSvc nintendo.Service
	userSvc     user.Service
	imageSvc    imageSvc.Service
	dumper      dump.Dumper

	writerChan chan *Content
	readerChan chan (chan *Content)
}

func NewRepository(nintendoSvc nintendo.Service, userSvc user.Service, imageSvc imageSvc.Service, config Config) Repository {
	dumpConfig := dump.Config{}
	dumpConfig.AddTarget(dumperKey.Weapon, config.Dumper.WeaponFile)
	dumpConfig.AddTarget(dumperKey.Stage, config.Dumper.StageFile)
	dumper := dump.New(dumpConfig)
	ret:= &repoImpl{
		userSvc:     userSvc,
		imageSvc:    imageSvc,
		nintendoSvc: nintendoSvc,
		dumper:      dumper,
		writerChan:  make(chan *Content),
		readerChan:  make(chan chan *Content),
	}
	ret.runUpdater()
	return ret
}

func (repo *repoImpl) Content() *Content {
	retChan := make(chan *Content)
	defer close(retChan)
	repo.readerChan <- retChan
	return <-retChan
}

func (repo *repoImpl) NextUpdateTime() time.Time {
	if repo.Content() == nil {
		return time.Now()
	}
	return util.Time.SplatoonNextUpdateTime(time.Now())
}

func (repo *repoImpl) Name() string {
	return "Salmon Schedules"
}

func (repo *repoImpl) Update() error {
	var err error
	admins := repo.userSvc.Admins()
	if len(admins) == 0 {
		return errors.New("no admin")
	}
	for _, admin := range admins {
		err = repo.updateByUid(admin)
		if err == nil {
			return nil
		}
	}
	if err != nil {
		return errors.Wrap(err, "can't update salmon schedules by admin")
	}
	return nil
}

func (repo *repoImpl) updateByUid(uid user.ID) error {
	status, err := repo.userSvc.GetStatus(uid)
	if err != nil {
		return errors.Wrap(err, "can't fetch admin status")
	}
	schedules, err := repo.nintendoSvc.GetSalmonSchedules(status.IKSM, status.Timezone, language.English)
	if err != nil {
		return errors.Wrap(err, "can't fetch salmon schedules")
	}
	sortSchedules(&schedules)
	populateSchedules(&schedules)

	if repo.content != nil && !hasUpdated(repo.content.Schedules, schedules) {
		log.Info("no new salmon schedules. skip update.")
		return nil
	}

	err = repo.updateDumper(schedules)
	if err != nil {
		// go ahead
		log.Warn("can't dump salmon schedules", zap.Error(err))
	}

	ids, err := repo.uploadSchedulesImages(schedules)
	if err != nil {
		return errors.Wrap(err, "can't upload salmon schedules images")
	}
	repo.writerChan <- &Content{
		Schedules: schedules,
		ImageIDs:  ids,
	}
	return nil
}

func sortSchedules(schedules *nintendo.SalmonSchedules) {
	// sort by start time in descending order
	sort.Slice(schedules.Details, func(i, j int) bool {
		return schedules.Details[i].StartTime > schedules.Details[j].StartTime
	})
	sort.Slice(schedules.Schedules, func(i, j int) bool {
		return schedules.Schedules[i].StartTime > schedules.Schedules[j].StartTime
	})
}

// populateSchedules fills Weapon's fields by SpecialWeapon, if Weapon is nil
func populateSchedules(schedules *nintendo.SalmonSchedules) {
	for i := range schedules.Details {
		detail := &schedules.Details[i]
		for j := range detail.Weapons {
			weapon := &detail.Weapons[j]
			if weapon.Weapon != nil {
				weapon.Weapon.Image = nintendo.Endpoint + weapon.Weapon.Image
				weapon.Weapon.Thumbnail = nintendo.Endpoint + weapon.Weapon.Thumbnail
			}
			if weapon.SpecialWeapon != nil {
				weapon.Weapon = &nintendo.SalmonWeapon{
					ID:   weapon.ID,
					Name: weapon.SpecialWeapon.Name,
				}
				if weapon.Weapon.ID == weaponID.Random {
					weapon.Weapon.Image = "file://./service/resources/salmon_random_weapon_green.png"
					weapon.Weapon.Thumbnail = "file://./service/resources/salmon_random_weapon_green.png"
				} else {
					weapon.Weapon.Image = "file://./service/resources/salmon_random_weapon_yellow.png"
					weapon.Weapon.Thumbnail = "file://./service/resources/salmon_random_weapon_yellow.png"
				}
				weapon.SpecialWeapon = nil
			}
		}
		detail.Stage.Image = nintendo.Endpoint + detail.Stage.Image
	}
}

func (repo *repoImpl) updateDumper(schedules nintendo.SalmonSchedules) error {
	stagesPtr, err := repo.dumper.Get(dumperKey.Stage, &stageCollection{})
	if err != nil {
		return errors.Wrap(err, "can't get stage dumper object")
	}
	weaponsPtr, err := repo.dumper.Get(dumperKey.Weapon, &weaponCollection{})
	if err != nil {
		return errors.Wrap(err, "can't get weapon dumper object")
	}
	stages := *stagesPtr.(*stageCollection)
	weapons := *weaponsPtr.(*weaponCollection)
	for _, detail := range schedules.Details {
		stages[detail.Stage.Name] = detail.Stage
		for _, weapon := range detail.Weapons {
			if weapon.Weapon != nil {
				weapons[weapon.Weapon.Name] = weapon
			} else if weapon.SpecialWeapon != nil {
				weapons[weapon.SpecialWeapon.Name] = weapon
			}
		}
	}
	err = repo.dumper.Save(dumperKey.Stage, &stages)
	if err != nil {
		return errors.Wrap(err, "can't save stage dumper object")
	}
	err = repo.dumper.Save(dumperKey.Weapon, &weapons)
	if err != nil {
		return errors.Wrap(err, "can't save weapon dumper object")
	}
	return nil
}

func hasUpdated(oldSchedules, newSchedules nintendo.SalmonSchedules) bool {
	return oldSchedules.Details[0].StartTime != newSchedules.Details[0].StartTime
}

func (repo *repoImpl) uploadSchedulesImages(schedules nintendo.SalmonSchedules) ([]imageSvc.Identifier, error) {
	urls := make([]string, 0, 10)
	for _, detail := range schedules.Details {
		urls = append(urls, detail.Stage.Image)
		for _, weapon := range detail.Weapons {
			if weapon.Weapon != nil {
				urls = append(urls, weapon.Weapon.Image)
			} else if weapon.SpecialWeapon != nil {
				urls = append(urls, weapon.SpecialWeapon.Image)
			} else {
				return nil, errors.Errorf("no image found")
			}
		}
	}
	imgs, err := repo.imageSvc.DownloadAll(urls)
	if err != nil {
		return nil, err
	}
	ids, err := repo.imageSvc.UploadAll([]image.Image{drawImage(imgs[0:5]), drawImage(imgs[5:10])})
	if err != nil {
		return nil, errors.Wrap(err, "can't upload images")
	}
	return ids, nil
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

// drawImage parameter:
//   imgs[0]: stage; imgs[1:5]: weapons
func drawImage(imgs []image.Image) image.Image {
	stage := imgs[0]
	weapons := imgs[1:5]
	width := stage.Bounds().Dx()
	qtrWidth := width / 4
	height := stage.Bounds().Dy() + qtrWidth
	// resize
	for i, img := range weapons {
		weapons[i] = resize.Resize(uint(qtrWidth), uint(qtrWidth), img, resize.Lanczos3)
	}
	// prepare canvas
	r := image.Rectangle{Min: image.Point{}, Max: image.Point{X: width, Y: height}}
	rgba := image.NewRGBA(r)
	//draw
	draw.Draw(rgba,
		image.Rectangle{Min: image.Point{Y: qtrWidth}, Max: image.Point{X: width, Y: height}},
		stage, image.Point{}, draw.Src)
	for i, img := range weapons {
		draw.Draw(rgba,
			image.Rectangle{Min: image.Point{X: i * qtrWidth}, Max: image.Point{X: (i + 1) * qtrWidth, Y: qtrWidth}},
			img, image.Point{}, draw.Src)
	}
	return rgba
}
