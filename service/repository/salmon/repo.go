package salmon

import (
	"image"
	"image/draw"
	"sort"
	"strconv"
	"time"

	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	log "telegram-splatoon2-bot/common/log"
	imageSvc "telegram-splatoon2-bot/service/image"
	nintendoSvc "telegram-splatoon2-bot/service/nintendo"
	"telegram-splatoon2-bot/service/todo"
	"telegram-splatoon2-bot/service/user"
)

type UID = int64

type content struct {
	schedules      *nintendoSvc.SalmonSchedules
	furtherImageID string
	laterImageID   string
}

type RepoImpl struct {
	//schedules      *nintendo.SalmonSchedules
	//furtherImageID string
	//laterImageID   string

	content *content

	userSvc       user.Service
	imageSvc      imageSvc.Service
	salmonDumping *DumperImpl
}

func NewRepo(userSvc user.Service, imageSvc imageSvc.Service, config Config) (*RepoImpl, error) {
	salmonDumping := NewDumper(config.dumper)
	err := salmonDumping.Load()
	if err != nil {
		return nil, errors.Wrap(err, "can't load salmon dumping files")
	}
	return &RepoImpl{
		content:       &content{},
		userSvc:       userSvc,
		imageSvc:      imageSvc,
		salmonDumping: salmonDumping,
	}, nil
}

func (repo *RepoImpl) HasInit() bool {
	return repo.content.schedules != nil
}

func (repo *RepoImpl) Name() string {
	return "SalmonRepo"
}

func (repo *RepoImpl) Update() error {
	admins := repo.userSvc.Admin.Snapshot()
	for _, admin := range admins {
		err := repo.updateByUid(admin)
		if err == nil {
			return nil
		}
	}

	log.Warn("can't update salmon schedules by admin")
	// todo: user responsibility?
	runtime, err := todo.RuntimeTable.GetFirstRuntime()
	if err != nil {
		return errors.Wrap(err, "can't get first runtime object")
	}
	err = repo.updateByUid(runtime.Uid)
	if err != nil {
		return errors.Wrap(err, "can't update salmon schedules by other user")
	}
	return nil
}

func (repo *RepoImpl) updateByUid(uid int64) error {
	wrapper := func(iksm string, timezone int, acceptLang string, _ ...interface{}) (interface{}, error) {
		return nintendoSvc.GetSalmonSchedules(iksm, timezone, acceptLang)
	}
	result, err := todo.FetchResourceWithUpdate(uid, wrapper)
	if err != nil {
		return errors.Wrap(err, "can't fetch runtime")
	}
	schedules := result.(*nintendoSvc.SalmonSchedules)
	sortSchedules(schedules)
	populateSchedules(schedules)

	currentSchedules := repo.content.schedules
	if currentSchedules != nil && currentSchedules.Details[0].StartTime == schedules.Details[0].StartTime {
		log.Info("no new salmon schedules. skip update.")
		return nil
	}

	err = repo.salmonDumping.Update(schedules)
	if err != nil {
		log.Warn("can't update salmon dumping file", zap.Error(err))
	}
	err = repo.salmonDumping.Save()
	if err != nil {
		log.Warn("can't update salmon dumping file", zap.Error(err))
	} else {
		log.Info("dumped salmon stages and weapons to files")
	}

	err = repo.uploadSchedulesImages(schedules)
	if err != nil {
		return errors.Wrap(err, "can't upload salmon schedules images")
	}
	repo.content.schedules = schedules
	return nil
}

func sortSchedules(salmonSchedules *nintendoSvc.SalmonSchedules) {
	// sort by start time in descending order
	sort.Slice(salmonSchedules.Details, func(i, j int) bool {
		return salmonSchedules.Details[i].StartTime > salmonSchedules.Details[j].StartTime
	})
	sort.Slice(salmonSchedules.Schedules, func(i, j int) bool {
		return salmonSchedules.Schedules[i].StartTime > salmonSchedules.Schedules[j].StartTime
	})
}

// populateSchedules fills Weapon's fields by SpecialWeapon, if Weapon is nil
func populateSchedules(salmonSchedules *nintendoSvc.SalmonSchedules) {
	for _, detail := range salmonSchedules.Details {
		for _, weapon := range detail.Weapons {
			if weapon.Weapon != nil {
				weapon.Weapon.Image = nintendoSvc.Endpoint + weapon.Weapon.Image
				weapon.Weapon.Thumbnail = nintendoSvc.Endpoint + weapon.Weapon.Thumbnail
			}
			if weapon.SpecialWeapon != nil {
				weapon.Weapon = &nintendoSvc.SalmonWeapon{
					ID:   weapon.ID,
					Name: weapon.SpecialWeapon.Name,
					// default images are too ugly
					//Image:     nintendo.Endpoint + weapon.SpecialWeapon.Image,
					//Thumbnail: nintendo.Endpoint + weapon.SpecialWeapon.Image,
				}
				if weapon.SpecialWeapon.Name == "Random" {
					weapon.Weapon.Image = "file://./service/resources/salmon_random_weapon_green.png"
					weapon.Weapon.Thumbnail = "file://./service/resources/salmon_random_weapon_green.png"
				} else {
					weapon.Weapon.Image = "file://./service/resources/salmon_random_weapon_yellow.png"
					weapon.Weapon.Thumbnail = "file://./service/resources/salmon_random_weapon_yellow.png"
				}
				weapon.SpecialWeapon = nil
			}
		}
		detail.Stage.Image = nintendoSvc.Endpoint + detail.Stage.Image
	}
}

func (repo *RepoImpl) uploadSchedulesImages(salmonSchedules *nintendoSvc.SalmonSchedules) error {
	urls := make([]string, 0, 10)
	for _, detail := range salmonSchedules.Details {
		urls = append(urls, detail.Stage.Image)
		for _, weapon := range detail.Weapons {
			if weapon.Weapon != nil {
				urls = append(urls, weapon.Weapon.Image)
			} else if weapon.SpecialWeapon != nil {
				urls = append(urls, weapon.SpecialWeapon.Image)
			} else {
				return errors.Errorf("no image found")
			}
		}
	}
	imgs, err := repo.imageSvc.DownloadAll(urls)
	if err != nil {
		return err
	}
	now := strconv.FormatInt(time.Now().Round(time.Hour).Unix(), 10)
	furtherImg := drawImage(imgs[0:5])
	laterImg := drawImage(imgs[5:10])
	ids, err := repo.imageSvc.UploadAll(
		[]image.Image{furtherImg, laterImg},
		[]string{"further_salmon_schedule_" + now, "later_salmon_schedule_" + now},
	)
	if err != nil {
		return errors.Wrap(err, "can't upload images")
	}
	repo.content.furtherImageID = ids[0]
	repo.content.laterImageID = ids[1]
	return nil
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

// todo: move
func QuerySalmonSchedules(update *botApi.Update) error {
	//user := update.Message.From
	//schedules := salmonScheduleRepo.schedules
	//if schedules == nil {
	//	return errors.Errorf("no cached schedules")
	//}
	//runtime, err := service.FetchRuntime(int64(user.ID))
	//if err != nil {
	//	return errors.Wrap(err, "can't fetch runtime")
	//}
	//now := time.Now().Unix()
	//startTime := schedules.Details[1].StartTime
	//endTime := schedules.Details[1].EndTime
	//var textKey string
	//var remainingTime time.Duration
	//if now > startTime {
	//	textKey = service.salmonSchedulesOpenTextKey
	//	remainingTime = time.Until(time.Unix(endTime, 0))
	//} else {
	//	textKey = service.salmonSchedulesSoonTextKey
	//	remainingTime = time.Until(time.Unix(startTime, 0))
	//}
	//remainingTime = remainingTime.Round(time.Minute)
	//hour := remainingTime / time.Hour
	//remainingTime -= hour * time.Hour
	//minute := remainingTime / time.Minute
	//
	//timeTemplate := service.getI18nText(runtime.Language, user, service.NewI18nKey(service.TimeTemplateTextKey))[0]
	//keys := []service.I18nKeys{
	//	service.NewI18nKey(service.salmonSchedulesFutureTextKey),
	//	service.NewI18nKey(service.salmonSchedulesNextTextKey),
	//	service.NewI18nKey(textKey, hour, minute),
	//}
	//future := schedules.Schedules[:len(schedules.Schedules)-2]
	//for _, s := range future {
	//	startTime := service.TimeHelper.getLocalTime(s.StartTime, runtime.Timezone).Format(timeTemplate)
	//	endTime := service.TimeHelper.getLocalTime(s.EndTime, runtime.Timezone).Format(timeTemplate)
	//	keys = append(keys, service.NewI18nKey(service.salmonSchedulesScheduleTextKey, startTime, endTime))
	//}
	//for _, s := range schedules.Details {
	//	startTime := service.TimeHelper.getLocalTime(s.StartTime, runtime.Timezone).Format(timeTemplate)
	//	endTime := service.TimeHelper.getLocalTime(s.EndTime, runtime.Timezone).Format(timeTemplate)
	//	keys = append(keys, service.NewI18nKey(service.salmonSchedulesDetailTextKey,
	//		startTime, endTime, s.Stage.Name,
	//		s.Weapons[0].Weapon.Name,
	//		s.Weapons[1].Weapon.Name,
	//		s.Weapons[2].Weapon.Name,
	//		s.Weapons[3].Weapon.Name))
	//}
	//texts := service.getI18nText(runtime.Language, user, keys...)
	//futureText := strings.Join(texts[3:len(future)+3], "") + texts[0]
	//futureMsg := botApi.NewMessage(update.Message.Chat.ID, futureText)
	//futureMsg.ParseMode = "Markdown"
	//furtherText := texts[3+len(future)] + texts[1]
	//furtherMsg := botApi.NewPhotoShare(update.Message.Chat.ID, salmonScheduleRepo.furtherImageID)
	//furtherMsg.Caption = furtherText
	//furtherMsg.ParseMode = "Markdown"
	//laterText := texts[4+len(future)] + texts[2]
	//laterMsg := botApi.NewPhotoShare(update.Message.Chat.ID, salmonScheduleRepo.laterImageID)
	//laterMsg.Caption = laterText
	//laterMsg.ParseMode = "Markdown"
	//err = service.sendWithRetry(service.bot, futureMsg)
	//if err != nil {
	//	return err
	//}
	//err = service.sendWithRetry(service.bot, furtherMsg)
	//if err != nil {
	//	return err
	//}
	//err = service.sendWithRetry(service.bot, laterMsg)
	//if err != nil {
	//	return err
	//}
	return nil
}
