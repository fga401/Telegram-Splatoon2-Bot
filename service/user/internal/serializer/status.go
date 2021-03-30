package serializer

import (
	"bytes"
	"encoding/binary"

	"telegram-splatoon2-bot/service/language"
	"telegram-splatoon2-bot/service/user/database"
)

// ToStatus deserialize Status
func ToStatus(value []byte) database.Status {
	ret := database.Status{}
	buf := bytes.NewBuffer(value)
	_ = binary.Read(buf, binary.LittleEndian, &(ret.UserID))
	sessionToken, _ := ReadBytes(buf, binary.LittleEndian, 16)
	iksm, _ := ReadBytes(buf, binary.LittleEndian, 8)
	lastBattle, _ := ReadBytes(buf, binary.LittleEndian, 8)
	lastSalmon, _ := ReadBytes(buf, binary.LittleEndian, 8)
	lang, _ := ReadBytes(buf, binary.LittleEndian, 8)
	_ = binary.Read(buf, binary.LittleEndian, &(ret.Timezone))
	ret.SessionToken = string(sessionToken)
	ret.IKSM = string(iksm)
	ret.LastBattle = string(lastBattle)
	ret.LastSalmon = string(lastSalmon)
	ret.Language = language.Language(lang)
	return ret
}

// FromStatus serialize Status
func FromStatus(status database.Status) []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, status.UserID)
	_ = WriteBytes(buf, binary.LittleEndian, []byte(status.SessionToken), 16)
	_ = WriteBytes(buf, binary.LittleEndian, []byte(status.IKSM), 8)
	_ = WriteBytes(buf, binary.LittleEndian, []byte(status.LastBattle), 8)
	_ = WriteBytes(buf, binary.LittleEndian, []byte(status.LastSalmon), 8)
	_ = WriteBytes(buf, binary.LittleEndian, []byte(status.Language), 8)
	_ = binary.Write(buf, binary.LittleEndian, status.Timezone)
	return buf.Bytes()
}
