package database

import (
	"telegram-splatoon2-bot/driver/database"
	"telegram-splatoon2-bot/service/language"
	"telegram-splatoon2-bot/service/timezone"
)

func init() {
	registerStatements([]database.Declaration{
		{
			Token:    tokenEnum.Status.SelectByUid,
			Stmt:     "SELECT * FROM status WHERE uid=?;",
			Named:    false,
			Prepared: true,
		},
		{
			Token:    tokenEnum.Status.UpdateLanguage,
			Stmt:     "UPDATE status SET language=? WHERE uid=?;",
			Named:    false,
			Prepared: false,
		},
		{
			Token:    tokenEnum.Status.UpdateTimezone,
			Stmt:     "UPDATE status SET timezone=? WHERE uid=?;",
			Named:    false,
			Prepared: false,
		},
		{
			Token:    tokenEnum.Status.UpdateIKSM,
			Stmt:     "UPDATE status SET iksm=? WHERE uid=?;",
			Named:    false,
			Prepared: true,
		},
	})
}

func (svc *serviceImpl) SelectStatus(uid UserID) (Status, error) {
	ret := Status{}
	err := svc.db.Get(tokenEnum.Status.SelectByUid, &ret, uid)
	return ret, err
}

func (svc *serviceImpl) UpdateStatusIKSM(uid UserID, iksm string) error {
	return svc.db.Exec(tokenEnum.Status.UpdateIKSM, iksm, uid)
}

func (svc *serviceImpl) UpdateStatusTimezone(uid UserID, timezone timezone.Timezone) error {
	return svc.db.Exec(tokenEnum.Status.UpdateTimezone, timezone, uid)
}

func (svc *serviceImpl) UpdateStatusLanguage(uid UserID, language language.Language) error {
	return svc.db.Exec(tokenEnum.Status.UpdateLanguage, language, uid)
}