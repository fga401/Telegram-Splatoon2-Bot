package timezone

// Timezone of user
type Timezone int32

// all available timezones
const (
	UTCMinus12  = Timezone(-720)
	UTCMinus11  = Timezone(-660)
	UTCMinus10  = Timezone(-600)
	UTCMinus930 = Timezone(-570)
	UTCMinus9   = Timezone(-540)
	UTCMinus8   = Timezone(-480)
	UTCMinus7   = Timezone(-420)
	UTCMinus6   = Timezone(-360)
	UTCMinus5   = Timezone(-300)
	UTCMinus4   = Timezone(-240)
	UTCMinus330 = Timezone(-210)
	UTCMinus3   = Timezone(-180)
	UTCMinus2   = Timezone(-120)
	UTCMinus1   = Timezone(-60)
	UTCPlus0    = Timezone(0)
	UTCPlus1    = Timezone(60)
	UTCPlus2    = Timezone(120)
	UTCPlus3    = Timezone(180)
	UTCPlus330  = Timezone(210)
	UTCPlus4    = Timezone(240)
	UTCPlus430  = Timezone(270)
	UTCPlus5    = Timezone(300)
	UTCPlus530  = Timezone(330)
	UTCPlus545  = Timezone(345)
	UTCPlus6    = Timezone(360)
	UTCPlus630  = Timezone(390)
	UTCPlus7    = Timezone(420)
	UTCPlus8    = Timezone(480)
	UTCPlus9    = Timezone(540)
	UTCPlus930  = Timezone(570)
	UTCPlus10   = Timezone(600)
	UTCPlus1030 = Timezone(630)
	UTCPlus11   = Timezone(660)
	UTCPlus12   = Timezone(720)
	UTCPlus1245 = Timezone(765)
	UTCPlus13   = Timezone(780)
	UTCPlus1345 = Timezone(825)
	UTCPlus14   = Timezone(840)
)

// All available timezones
var All = []Timezone{
	UTCMinus12,
	UTCMinus11,
	UTCMinus10,
	UTCMinus930,
	UTCMinus9,
	UTCMinus8,
	UTCMinus7,
	UTCMinus6,
	UTCMinus5,
	UTCMinus4,
	UTCMinus330,
	UTCMinus3,
	UTCMinus2,
	UTCMinus1,
	UTCPlus0,
	UTCPlus1,
	UTCPlus2,
	UTCPlus3,
	UTCPlus330,
	UTCPlus4,
	UTCPlus430,
	UTCPlus5,
	UTCPlus530,
	UTCPlus545,
	UTCPlus6,
	UTCPlus630,
	UTCPlus7,
	UTCPlus8,
	UTCPlus9,
	UTCPlus930,
	UTCPlus10,
	UTCPlus1030,
	UTCPlus11,
	UTCPlus12,
	UTCPlus1245,
	UTCPlus13,
	UTCPlus1345,
	UTCPlus14,
}

// ByMinute converts offset in minute to Timezone.
// If the returned timezone is not in available timezone, it will return UTC+8.
func ByMinute(minutes int) Timezone {
	switch minutes {
	case -720:
		return UTCMinus12
	case -660:
		return UTCMinus11
	case -600:
		return UTCMinus10
	case -570:
		return UTCMinus930
	case -540:
		return UTCMinus9
	case -480:
		return UTCMinus8
	case -420:
		return UTCMinus7
	case -360:
		return UTCMinus6
	case -300:
		return UTCMinus5
	case -240:
		return UTCMinus4
	case -210:
		return UTCMinus330
	case -180:
		return UTCMinus3
	case -120:
		return UTCMinus2
	case -60:
		return UTCMinus1
	case 0:
		return UTCPlus0
	case 60:
		return UTCPlus1
	case 120:
		return UTCPlus2
	case 180:
		return UTCPlus3
	case 210:
		return UTCPlus330
	case 240:
		return UTCPlus4
	case 270:
		return UTCPlus430
	case 300:
		return UTCPlus5
	case 330:
		return UTCPlus530
	case 345:
		return UTCPlus545
	case 360:
		return UTCPlus6
	case 390:
		return UTCPlus630
	case 420:
		return UTCPlus7
	case 480:
		return UTCPlus8
	case 540:
		return UTCPlus9
	case 570:
		return UTCPlus930
	case 600:
		return UTCPlus10
	case 630:
		return UTCPlus1030
	case 660:
		return UTCPlus11
	case 720:
		return UTCPlus12
	case 765:
		return UTCPlus1245
	case 780:
		return UTCPlus13
	case 840:
		return UTCPlus14
	default:
		return UTCPlus8
	}
}
