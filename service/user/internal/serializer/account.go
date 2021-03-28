package serializer

import (
	"bytes"
	"encoding/binary"

	"telegram-splatoon2-bot/service/user/database"
)

// ToAccounts deserialize Account
func ToAccounts(value []byte) []database.Account {
	buf := bytes.NewBuffer(value)
	var size int32
	_ = binary.Read(buf, binary.LittleEndian, &size)
	ret := make([]database.Account, size)
	for i := int32(0); i < size; i++ {
		account := database.Account{}
		_ = binary.Read(buf, binary.LittleEndian, &(account.UserID))
		sessionToken, _ := ReadBytes(buf, binary.LittleEndian, 16)
		tag, _ := ReadBytes(buf, binary.LittleEndian, 16)
		account.SessionToken = string(sessionToken)
		account.Tag = string(tag)
		ret[i] = account
	}
	return ret
}

// FromAccounts serialize Account
func FromAccounts(accounts []database.Account) []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, int32(len(accounts)))
	for _, account := range accounts {
		_ = binary.Write(buf, binary.LittleEndian, account.UserID)
		_ = WriteBytes(buf, binary.LittleEndian, []byte(account.SessionToken), 16)
		_ = WriteBytes(buf, binary.LittleEndian, []byte(account.Tag), 16)
	}
	return buf.Bytes()
}
