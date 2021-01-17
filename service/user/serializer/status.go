package serializer

import (
	"bytes"
	"encoding/binary"

	"telegram-splatoon2-bot/common/util"
	"telegram-splatoon2-bot/service/language"
	"telegram-splatoon2-bot/service/user/database"
)

func ToStatus(value []byte) database.Status {
	ret := database.Status{}
	buf := bytes.NewBuffer(value)
	_ = binary.Read(buf, binary.LittleEndian, &(ret.UserID))
	sessionToken, _ := util.Binary.ReadBytes(buf, binary.LittleEndian, 16)
	iksm, _ := util.Binary.ReadBytes(buf, binary.LittleEndian, 8)
	lang, _ := util.Binary.ReadBytes(buf, binary.LittleEndian, 8)
	_ = binary.Read(buf, binary.LittleEndian, &(ret.Timezone))
	ret.SessionToken = string(sessionToken)
	ret.IKSM = string(iksm)
	ret.Language = language.Language(lang)
	return ret
}

func FromStatus(status database.Status) []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, status.UserID)
	_ = util.Binary.WriteBytes(buf, binary.LittleEndian, []byte(status.SessionToken), 16)
	_ = util.Binary.WriteBytes(buf, binary.LittleEndian, []byte(status.IKSM), 8)
	_ = util.Binary.WriteBytes(buf, binary.LittleEndian, []byte(status.Language), 8)
	_ = binary.Write(buf, binary.LittleEndian, status.Timezone)
	return buf.Bytes()
}

