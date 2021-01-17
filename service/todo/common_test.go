package todo

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	json "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"telegram-splatoon2-bot/botutil"
	log "telegram-splatoon2-bot/common/log"
	"telegram-splatoon2-bot/service/cache"
	"telegram-splatoon2-bot/service/db"
	"telegram-splatoon2-bot/service/image"
	nintendo2 "telegram-splatoon2-bot/service/nintendo"
	"telegram-splatoon2-bot/service/repository/stage"
	"telegram-splatoon2-bot/service/user"
)


func prepareTest() {
	viper.SetConfigName("dev")
	viper.SetConfigType("json")
	viper.AddConfigPath("../config/")
	viper.AddConfigPath("./config/")
	viper.ReadInConfig()

	err := viper.BindEnv("token")
	if err != nil {
		panic(errors.Wrap(err, "can't bind token env"))
	}
	err = viper.BindEnv("admin")
	if err != nil {
		panic(errors.Wrap(err, "can't bind admin env"))
	}
	err = viper.BindEnv("store_channel")
	if err != nil {
		panic(errors.Wrap(err, "can't bind store_channel env"))
	}

	log.InitLogger()
	image.InitImageClient()
	db.InitDatabaseInstance()
	cache.InitCache()
	Cache = cache.Cache
	AccountTable = db.AccountTable
	UserTable = db.UserTable
	RuntimeTable = db.RuntimeTable
	Transactions = db.Transactions

	botConfig := botutil.BotConfig{
		UserProxy: viper.GetBool("bot.useProxy"),
		ProxyUrl:  viper.GetString("bot.proxyUrl"),
		Token:     viper.GetString("token"),
		Debug:     viper.GetBool("bot.debug"),
	}

	myBot := botutil.NewBot(botConfig)
	// default value
	bot = myBot
	userMaxAccount = viper.GetInt("account.maxAccount")
	userAllowPolling = viper.GetBool("account.allowPolling")
	callbackQueryCachedSecond = viper.GetInt("service.callbackQueryCachedSecond")
	RetryTimes = viper.GetInt("service.retryTimes")
	defaultAdmin, err = strconv.ParseInt(viper.GetString("admin"), 10, 64)
	if err != nil {
		panic(errors.Wrap(err, "viper get admin failed"))
	}
	storeChannelID, err = strconv.ParseInt(viper.GetString("store_channel"), 10, 64)
	if err != nil {
		panic(errors.Wrap(err, "viper get store_channel failed"))
	}
	updateFailureRetryInterval = viper.GetDuration("service.updateFailureRetryInterval")
	updateDelayInSecond = viper.GetInt64("service.updateDelayInSecond")
	// markup
	initMarkup()
	//service
	user.loadUsers()
	log.Info("Preparation done.")
}

func TestGetSplatoonNextUpdateTime(t *testing.T) {
	ret := TimeHelper.getSplatoonNextUpdateTime(time.Now())
	fmt.Println(ret)
	ret = TimeHelper.getSplatoonNextUpdateTime(time.Now().Add(time.Hour))
	fmt.Println(ret)
	ret = TimeHelper.getSplatoonNextUpdateTime(TimeHelper.getSplatoonNextUpdateTime(time.Now()))
	fmt.Println(ret)
	assert.Nil(t,nil)
}

func TestUpdateStage(t *testing.T) {
	prepareTest()
	var jsonText = `{"gachi":[{"game_mode":{"key":"gachi","name":"Ranked Battle"},"id":4780952683924113405,"end_time":1603209600,"start_time":1603202400,"rule":{"name":"Clam Blitz","multiline_name":"Clam\nBlitz","key":"clam_blitz"},"stage_a":{"name":"Musselforge Fitness","image":"/images/stage/83acec875a5bb19418d7b87d5df4ba1e38ceac66.png","id":"1"},"stage_b":{"image":"/images/stage/5c030a505ee57c889d3e5268a4b10c1f1f37880a.png","name":"Inkblot Art Academy","id":"4"}},{"game_mode":{"key":"gachi","name":"Ranked Battle"},"id":4780952683924113408,"end_time":1603216800,"start_time":1603209600,"rule":{"name":"Rainmaker","multiline_name":"Rainmaker","key":"rainmaker"},"stage_a":{"id":"11","image":"/images/stage/758338859615898a59e93b84f7e1ca670f75e865.png","name":"Blackbelly Skatepark"},"stage_b":{"image":"/images/stage/65c99da154295109d6fe067005f194f681762f8c.png","name":"Walleye Warehouse","id":"14"}},{"end_time":1603224000,"start_time":1603216800,"rule":{"name":"Splat Zones","multiline_name":"Splat\nZones","key":"splat_zones"},"stage_a":{"id":"7","image":"/images/stage/0907fc7dc325836a94d385919fe01dc13848612a.png","name":"Port Mackerel"},"game_mode":{"name":"Ranked Battle","key":"gachi"},"id":4780952683924113411,"stage_b":{"image":"/images/stage/98baf21c0366ce6e03299e2326fe6d27a7582dce.png","name":"The Reef","id":"0"}},{"game_mode":{"key":"gachi","name":"Ranked Battle"},"id":4780952683924113414,"end_time":1603231200,"start_time":1603224000,"rule":{"key":"tower_control","name":"Tower Control","multiline_name":"Tower\nControl"},"stage_a":{"image":"/images/stage/a12e4bf9f871677a5f3735d421317fbbf09e1a78.png","name":"Kelp Dome","id":"10"},"stage_b":{"id":"2","name":"Starfish Mainstage","image":"/images/stage/187987856bf575c4155d021cb511034931d06d24.png"}},{"id":4780952683924113416,"game_mode":{"key":"gachi","name":"Ranked Battle"},"end_time":1603238400,"start_time":1603231200,"rule":{"multiline_name":"Clam\nBlitz","name":"Clam Blitz","key":"clam_blitz"},"stage_a":{"id":"8","name":"Moray Towers","image":"/images/stage/96fd8c0492331a30e60a217c94fd1d4c73a966cc.png"},"stage_b":{"image":"/images/stage/dcf332bdcc80f566f3ae59c1c3a29bc6312d0ba8.png","name":"Arowana Mall","id":"15"}},{"stage_b":{"name":"Wahoo World","image":"/images/stage/555c356487ac3edb0088c416e8045576c6b37fcc.png","id":"20"},"game_mode":{"key":"gachi","name":"Ranked Battle"},"id":4780952683924113419,"end_time":1603245600,"rule":{"key":"rainmaker","name":"Rainmaker","multiline_name":"Rainmaker"},"start_time":1603238400,"stage_a":{"name":"New Albacore Hotel","image":"/images/stage/98a7d7a4009fae9fb7479554535425a5a604e88e.png","id":"19"}},{"stage_b":{"id":"21","name":"Ancho-V Games","image":"/images/stage/1430e5ac7ae9396a126078eeab824a186b490b5a.png"},"id":4780952683924113423,"game_mode":{"key":"gachi","name":"Ranked Battle"},"stage_a":{"id":"13","image":"/images/stage/d9f0f6c330aaa3b975e572637b00c4c0b6b89f7d.png","name":"MakoMart"},"start_time":1603245600,"end_time":1603252800,"rule":{"multiline_name":"Tower\nControl","name":"Tower Control","key":"tower_control"}},{"start_time":1603252800,"end_time":1603260000,"rule":{"multiline_name":"Splat\nZones","name":"Splat Zones","key":"splat_zones"},"stage_a":{"id":"18","image":"/images/stage/8cab733d543efc9dd561bfcc9edac52594e62522.png","name":"Goby Arena"},"game_mode":{"key":"gachi","name":"Ranked Battle"},"id":4780952683924113426,"stage_b":{"id":"6","image":"/images/stage/070d7ee287fdf3c5df02411950c2a1ce5b238746.png","name":"Manta Maria"}},{"stage_b":{"id":"1","image":"/images/stage/83acec875a5bb19418d7b87d5df4ba1e38ceac66.png","name":"Musselforge Fitness"},"stage_a":{"id":"12","name":"Shellendorf Institute","image":"/images/stage/23259c80272f45cea2d5c9e60bc0cedb6ce29e46.png"},"end_time":1603267200,"start_time":1603260000,"rule":{"key":"clam_blitz","multiline_name":"Clam\nBlitz","name":"Clam Blitz"},"id":4780952683924113429,"game_mode":{"name":"Ranked Battle","key":"gachi"}},{"stage_b":{"image":"/images/stage/bc794e337900afd763f8a88359f83df5679ddf12.png","name":"Sturgeon Shipyard","id":"3"},"stage_a":{"image":"/images/stage/fc23fedca2dfbbd8707a14606d719a4004403d13.png","name":"Humpback Pump Track","id":"5"},"start_time":1603267200,"end_time":1603274400,"rule":{"multiline_name":"Tower\nControl","name":"Tower Control","key":"tower_control"},"game_mode":{"name":"Ranked Battle","key":"gachi"},"id":4780952683924113432},{"id":4780952683924113435,"game_mode":{"key":"gachi","name":"Ranked Battle"},"stage_a":{"id":"22","name":"Skipper Pavilion","image":"/images/stage/132327c819abf2bd44d0adc0f4a21aad9cc84bb2.png"},"end_time":1603281600,"rule":{"key":"rainmaker","name":"Rainmaker","multiline_name":"Rainmaker"},"start_time":1603274400,"stage_b":{"image":"/images/stage/187987856bf575c4155d021cb511034931d06d24.png","name":"Starfish Mainstage","id":"2"}},{"id":4780952683924113439,"game_mode":{"name":"Ranked Battle","key":"gachi"},"stage_a":{"name":"Wahoo World","image":"/images/stage/555c356487ac3edb0088c416e8045576c6b37fcc.png","id":"20"},"end_time":1603288800,"start_time":1603281600,"rule":{"key":"tower_control","name":"Tower Control","multiline_name":"Tower\nControl"},"stage_b":{"id":"10","name":"Kelp Dome","image":"/images/stage/a12e4bf9f871677a5f3735d421317fbbf09e1a78.png"}}],"league":[{"stage_b":{"name":"The Reef","image":"/images/stage/98baf21c0366ce6e03299e2326fe6d27a7582dce.png","id":"0"},"game_mode":{"name":"League Battle","key":"league"},"id":4780952683924113405,"stage_a":{"id":"9","image":"/images/stage/8c95053b3043e163cbfaaf1ec1e5f3eb770e5e07.png","name":"Snapper Canal"},"rule":{"key":"splat_zones","name":"Splat Zones","multiline_name":"Splat\nZones"},"end_time":1603209600,"start_time":1603202400},{"stage_b":{"id":"21","name":"Ancho-V Games","image":"/images/stage/1430e5ac7ae9396a126078eeab824a186b490b5a.png"},"end_time":1603216800,"rule":{"name":"Tower Control","multiline_name":"Tower\nControl","key":"tower_control"},"start_time":1603209600,"stage_a":{"id":"22","image":"/images/stage/132327c819abf2bd44d0adc0f4a21aad9cc84bb2.png","name":"Skipper Pavilion"},"id":4780952683924113408,"game_mode":{"key":"league","name":"League Battle"}},{"stage_b":{"id":"12","name":"Shellendorf Institute","image":"/images/stage/23259c80272f45cea2d5c9e60bc0cedb6ce29e46.png"},"end_time":1603224000,"rule":{"multiline_name":"Clam\nBlitz","name":"Clam Blitz","key":"clam_blitz"},"start_time":1603216800,"stage_a":{"id":"13","name":"MakoMart","image":"/images/stage/d9f0f6c330aaa3b975e572637b00c4c0b6b89f7d.png"},"id":4780952683924113411,"game_mode":{"key":"league","name":"League Battle"}},{"stage_b":{"id":"4","name":"Inkblot Art Academy","image":"/images/stage/5c030a505ee57c889d3e5268a4b10c1f1f37880a.png"},"id":4780952683924113414,"game_mode":{"name":"League Battle","key":"league"},"rule":{"multiline_name":"Rainmaker","name":"Rainmaker","key":"rainmaker"},"end_time":1603231200,"start_time":1603224000,"stage_a":{"name":"Humpback Pump Track","image":"/images/stage/fc23fedca2dfbbd8707a14606d719a4004403d13.png","id":"5"}},{"id":4780952683924113416,"game_mode":{"name":"League Battle","key":"league"},"stage_a":{"id":"10","image":"/images/stage/a12e4bf9f871677a5f3735d421317fbbf09e1a78.png","name":"Kelp Dome"},"start_time":1603231200,"end_time":1603238400,"rule":{"key":"splat_zones","name":"Splat Zones","multiline_name":"Splat\nZones"},"stage_b":{"id":"20","image":"/images/stage/555c356487ac3edb0088c416e8045576c6b37fcc.png","name":"Wahoo World"}},{"stage_b":{"id":"1","image":"/images/stage/83acec875a5bb19418d7b87d5df4ba1e38ceac66.png","name":"Musselforge Fitness"},"rule":{"multiline_name":"Tower\nControl","name":"Tower Control","key":"tower_control"},"end_time":1603245600,"start_time":1603238400,"stage_a":{"name":"Port Mackerel","image":"/images/stage/0907fc7dc325836a94d385919fe01dc13848612a.png","id":"7"},"id":4780952683924113419,"game_mode":{"key":"league","name":"League Battle"}},{"stage_b":{"id":"15","name":"Arowana Mall","image":"/images/stage/dcf332bdcc80f566f3ae59c1c3a29bc6312d0ba8.png"},"game_mode":{"name":"League Battle","key":"league"},"id":4780952683924113423,"stage_a":{"image":"/images/stage/98baf21c0366ce6e03299e2326fe6d27a7582dce.png","name":"The Reef","id":"0"},"end_time":1603252800,"start_time":1603245600,"rule":{"key":"clam_blitz","multiline_name":"Clam\nBlitz","name":"Clam Blitz"}},{"stage_b":{"image":"/images/stage/187987856bf575c4155d021cb511034931d06d24.png","name":"Starfish Mainstage","id":"2"},"rule":{"key":"rainmaker","name":"Rainmaker","multiline_name":"Rainmaker"},"end_time":1603260000,"start_time":1603252800,"stage_a":{"image":"/images/stage/bc794e337900afd763f8a88359f83df5679ddf12.png","name":"Sturgeon Shipyard","id":"3"},"game_mode":{"key":"league","name":"League Battle"},"id":4780952683924113426},{"id":4780952683924113429,"game_mode":{"key":"league","name":"League Battle"},"start_time":1603260000,"end_time":1603267200,"rule":{"key":"splat_zones","name":"Splat Zones","multiline_name":"Splat\nZones"},"stage_a":{"id":"16","name":"Camp Triggerfish","image":"/images/stage/e4c4800be9fff23112334b193abb0fdf36e05933.png"},"stage_b":{"id":"17","name":"Piranha Pit","image":"/images/stage/828e49a8414a4bbc0a5da3e61454ab148a9f4063.png"}},{"stage_b":{"id":"6","image":"/images/stage/070d7ee287fdf3c5df02411950c2a1ce5b238746.png","name":"Manta Maria"},"stage_a":{"id":"11","image":"/images/stage/758338859615898a59e93b84f7e1ca670f75e865.png","name":"Blackbelly Skatepark"},"end_time":1603274400,"start_time":1603267200,"rule":{"name":"Clam Blitz","multiline_name":"Clam\nBlitz","key":"clam_blitz"},"id":4780952683924113432,"game_mode":{"key":"league","name":"League Battle"}},{"start_time":1603274400,"end_time":1603281600,"rule":{"key":"tower_control","name":"Tower Control","multiline_name":"Tower\nControl"},"stage_a":{"id":"18","image":"/images/stage/8cab733d543efc9dd561bfcc9edac52594e62522.png","name":"Goby Arena"},"id":4780952683924113435,"game_mode":{"name":"League Battle","key":"league"},"stage_b":{"id":"14","name":"Walleye Warehouse","image":"/images/stage/65c99da154295109d6fe067005f194f681762f8c.png"}},{"stage_b":{"name":"New Albacore Hotel","image":"/images/stage/98a7d7a4009fae9fb7479554535425a5a604e88e.png","id":"19"},"game_mode":{"key":"league","name":"League Battle"},"id":4780952683924113439,"stage_a":{"name":"MakoMart","image":"/images/stage/d9f0f6c330aaa3b975e572637b00c4c0b6b89f7d.png","id":"13"},"end_time":1603288800,"start_time":1603281600,"rule":{"key":"splat_zones","multiline_name":"Splat\nZones","name":"Splat Zones"}}],"regular":[{"end_time":1603209600,"rule":{"multiline_name":"Turf\nWar","name":"Turf War","key":"turf_war"},"start_time":1603202400,"stage_a":{"id":"22","name":"Skipper Pavilion","image":"/images/stage/132327c819abf2bd44d0adc0f4a21aad9cc84bb2.png"},"id":4780952683924113405,"game_mode":{"key":"regular","name":"Regular Battle"},"stage_b":{"id":"18","name":"Goby Arena","image":"/images/stage/8cab733d543efc9dd561bfcc9edac52594e62522.png"}},{"game_mode":{"name":"Regular Battle","key":"regular"},"id":4780952683924113408,"rule":{"key":"turf_war","name":"Turf War","multiline_name":"Turf\nWar"},"end_time":1603216800,"start_time":1603209600,"stage_a":{"name":"Shellendorf Institute","image":"/images/stage/23259c80272f45cea2d5c9e60bc0cedb6ce29e46.png","id":"12"},"stage_b":{"id":"7","name":"Port Mackerel","image":"/images/stage/0907fc7dc325836a94d385919fe01dc13848612a.png"}},{"id":4780952683924113411,"game_mode":{"key":"regular","name":"Regular Battle"},"stage_a":{"id":"14","image":"/images/stage/65c99da154295109d6fe067005f194f681762f8c.png","name":"Walleye Warehouse"},"start_time":1603216800,"end_time":1603224000,"rule":{"multiline_name":"Turf\nWar","name":"Turf War","key":"turf_war"},"stage_b":{"id":"1","name":"Musselforge Fitness","image":"/images/stage/83acec875a5bb19418d7b87d5df4ba1e38ceac66.png"}},{"stage_b":{"name":"Ancho-V Games","image":"/images/stage/1430e5ac7ae9396a126078eeab824a186b490b5a.png","id":"21"},"id":4780952683924113414,"game_mode":{"key":"regular","name":"Regular Battle"},"stage_a":{"id":"11","name":"Blackbelly Skatepark","image":"/images/stage/758338859615898a59e93b84f7e1ca670f75e865.png"},"start_time":1603224000,"end_time":1603231200,"rule":{"key":"turf_war","name":"Turf War","multiline_name":"Turf\nWar"}},{"stage_b":{"id":"5","image":"/images/stage/fc23fedca2dfbbd8707a14606d719a4004403d13.png","name":"Humpback Pump Track"},"id":4780952683924113416,"game_mode":{"name":"Regular Battle","key":"regular"},"stage_a":{"image":"/images/stage/d9f0f6c330aaa3b975e572637b00c4c0b6b89f7d.png","name":"MakoMart","id":"13"},"start_time":1603231200,"end_time":1603238400,"rule":{"key":"turf_war","multiline_name":"Turf\nWar","name":"Turf War"}},{"rule":{"multiline_name":"Turf\nWar","name":"Turf War","key":"turf_war"},"end_time":1603245600,"start_time":1603238400,"stage_a":{"image":"/images/stage/5c030a505ee57c889d3e5268a4b10c1f1f37880a.png","name":"Inkblot Art Academy","id":"4"},"game_mode":{"name":"Regular Battle","key":"regular"},"id":4780952683924113419,"stage_b":{"image":"/images/stage/8c95053b3043e163cbfaaf1ec1e5f3eb770e5e07.png","name":"Snapper Canal","id":"9"}},{"stage_a":{"image":"/images/stage/828e49a8414a4bbc0a5da3e61454ab148a9f4063.png","name":"Piranha Pit","id":"17"},"start_time":1603245600,"end_time":1603252800,"rule":{"multiline_name":"Turf\nWar","name":"Turf War","key":"turf_war"},"id":4780952683924113423,"game_mode":{"key":"regular","name":"Regular Battle"},"stage_b":{"name":"Manta Maria","image":"/images/stage/070d7ee287fdf3c5df02411950c2a1ce5b238746.png","id":"6"}},{"stage_a":{"id":"0","image":"/images/stage/98baf21c0366ce6e03299e2326fe6d27a7582dce.png","name":"The Reef"},"rule":{"key":"turf_war","multiline_name":"Turf\nWar","name":"Turf War"},"end_time":1603260000,"start_time":1603252800,"game_mode":{"key":"regular","name":"Regular Battle"},"id":4780952683924113426,"stage_b":{"image":"/images/stage/a12e4bf9f871677a5f3735d421317fbbf09e1a78.png","name":"Kelp Dome","id":"10"}},{"id":4780952683924113429,"game_mode":{"name":"Regular Battle","key":"regular"},"stage_a":{"id":"2","image":"/images/stage/187987856bf575c4155d021cb511034931d06d24.png","name":"Starfish Mainstage"},"start_time":1603260000,"end_time":1603267200,"rule":{"key":"turf_war","name":"Turf War","multiline_name":"Turf\nWar"},"stage_b":{"image":"/images/stage/bc794e337900afd763f8a88359f83df5679ddf12.png","name":"Sturgeon Shipyard","id":"3"}},{"stage_b":{"name":"Arowana Mall","image":"/images/stage/dcf332bdcc80f566f3ae59c1c3a29bc6312d0ba8.png","id":"15"},"game_mode":{"key":"regular","name":"Regular Battle"},"id":4780952683924113432,"end_time":1603274400,"start_time":1603267200,"rule":{"key":"turf_war","name":"Turf War","multiline_name":"Turf\nWar"},"stage_a":{"image":"/images/stage/e4c4800be9fff23112334b193abb0fdf36e05933.png","name":"Camp Triggerfish","id":"16"}},{"id":4780952683924113435,"game_mode":{"key":"regular","name":"Regular Battle"},"stage_a":{"id":"20","image":"/images/stage/555c356487ac3edb0088c416e8045576c6b37fcc.png","name":"Wahoo World"},"end_time":1603281600,"rule":{"multiline_name":"Turf\nWar","name":"Turf War","key":"turf_war"},"start_time":1603274400,"stage_b":{"id":"7","name":"Port Mackerel","image":"/images/stage/0907fc7dc325836a94d385919fe01dc13848612a.png"}},{"id":4780952683924113439,"game_mode":{"key":"regular","name":"Regular Battle"},"stage_a":{"id":"8","name":"Moray Towers","image":"/images/stage/96fd8c0492331a30e60a217c94fd1d4c73a966cc.png"},"end_time":1603288800,"start_time":1603281600,"rule":{"key":"turf_war","name":"Turf War","multiline_name":"Turf\nWar"},"stage_b":{"name":"Skipper Pavilion","image":"/images/stage/132327c819abf2bd44d0adc0f4a21aad9cc84bb2.png","id":"22"}}]}`
	stageSchedulesOld := &nintendo2.StageSchedules{}
	stageSchedulesNew := &nintendo2.StageSchedules{}
	_ = json.Unmarshal([]byte(jsonText), stageSchedulesOld)
	_ = json.Unmarshal([]byte(jsonText), stageSchedulesNew)
	for _, stage := range stageSchedulesNew.Regular{
		stage.StartTime += 7200
		stage.EndTime += 7200
	}
	for _, stage := range stageSchedulesNew.Gachi{
		stage.StartTime += 7200
		stage.EndTime += 7200
	}
	for _, stage := range stageSchedulesNew.League{
		stage.StartTime += 7200
		stage.EndTime += 7200
	}
	stageScheduleRepo, _ := stage.NewStageScheduleRepo(user.admins)
	stageScheduleRepo.sortSchedules(stageSchedulesOld)
	stageScheduleRepo.populateFields(stageSchedulesOld)
	s, err := stageScheduleRepo.wrapSchedules(stageSchedulesOld)
	assert.Nil(t, err)
	stageScheduleRepo.schedules = s
	log.Info("stageSchedulesOld uploaded")
	stageScheduleRepo.sortSchedules(stageSchedulesNew)
	stageScheduleRepo.populateFields(stageSchedulesNew)
	s, err = stageScheduleRepo.wrapSchedules(stageSchedulesNew)
	log.Info("stageSchedulesNew uploaded")
	assert.Nil(t, err)
	stageScheduleRepo.schedules = s
}