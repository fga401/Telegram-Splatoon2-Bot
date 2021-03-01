// Package serializer serializes/deserializes database struct to bytes.
// No error is returned since Read and Write of bytes.Buffer
// always return nil error.
package serializer

import (
	"bytes"
	"encoding/binary"

	"telegram-splatoon2-bot/service/user/database"
)

// ToID deserialize UserID
func ToID(key []byte) database.UserID {
	var id database.UserID
	buf := bytes.NewBuffer(key)
	_ = binary.Read(buf, binary.LittleEndian, &(id))
	return id
}

// FromID serialize UserID
func FromID(id database.UserID) []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, id)
	return buf.Bytes()
}
