package salmon

import (
	nintendo2 "telegram-splatoon2-bot/service/nintendo"
)

//easyjson:json
type stageCollection map[string]nintendo2.SalmonStage

//easyjson:json
type weaponCollection map[string]nintendo2.SalmonWeaponType
