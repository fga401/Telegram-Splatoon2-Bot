package nintendo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"telegram-splatoon2-bot/common/log"
	"telegram-splatoon2-bot/service/language"
	"telegram-splatoon2-bot/service/timezone"
)

var (
	svc  Service
	iksm = "82b55de7e27a253e70154c5b0b9888a9d14f1da3"
)

func TestMain(m *testing.M) {
	log.InitLogger("debug")
	svc = New(Config{
		Timeout:    0,
		RetryTimes: 1,
	})
	os.Exit(m.Run())
}

func TestGetBattleResults(t *testing.T) {
	results, err := svc.GetAllBattleResults(
		iksm,
		timezone.UTCPlus8,
		language.English,
	)
	require.Nil(t, err)
	latest, err := svc.GetLatestBattleResults(
		results.Results[2].Metadata().BattleNumber,
		10,
		iksm,
		timezone.UTCPlus8,
		language.English,
	)
	require.Nil(t, err)
	require.Len(t, latest, 2)
	detail, err := svc.GetDetailedBattleResults(
		results.Results[0].Metadata().BattleNumber,
		iksm,
		timezone.UTCPlus8,
		language.English,
	)
	require.Nil(t, err)
	require.Equal(t, detail.Metadata().BattleNumber, results.Results[0].Metadata().BattleNumber)
	require.Len(t, detail.MyTeamPlayerResults(), 3)
	require.Len(t, detail.OtherTeamPlayerResults(), 4)
}

func TestGetSalmonResults(t *testing.T) {
	results, err := svc.GetAllSalmonResults(
		iksm,
		timezone.UTCPlus8,
		language.English,
	)
	require.Nil(t, err)
	latest, err := svc.GetLatestSalmonResults(
		results.Results[2].JobID,
		10,
		iksm,
		timezone.UTCPlus8,
		language.English,
	)
	require.Nil(t, err)
	require.Len(t, latest, 2)
	detail, err := svc.GetDetailedSalmonResults(
		results.Results[0].JobID,
		iksm,
		timezone.UTCPlus8,
		language.English,
	)
	require.Nil(t, err)
	require.Equal(t, detail.JobID, results.Results[0].JobID)
	require.Len(t, detail.OtherResults, 3)
}

func TestGetSchedules(t *testing.T) {
	salmon, err := svc.GetSalmonSchedules(
		iksm,
		timezone.UTCPlus8,
		language.English,
	)
	require.Nil(t, err)
	require.Len(t, salmon.Schedules, 5)
	require.Len(t, salmon.Details, 2)
	stage, err := svc.GetStageSchedules(
		iksm,
		timezone.UTCPlus8,
		language.English,
	)
	require.Nil(t, err)
	require.Len(t, stage.Regular, 12)
	require.Len(t, stage.Gachi, 12)
	require.Len(t, stage.League, 12)
}
