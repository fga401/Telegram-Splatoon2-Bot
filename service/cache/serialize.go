package cache

import (
	"bytes"
	"encoding/binary"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"io"
	"strconv"
	"telegram-splatoon2-bot/service/db"
)

func userToStringKey(user *tgbotapi.User) string {
	key := strconv.FormatInt(int64(user.ID), 10)
	key = "PK" + key
	return key
}

func userToBytesKey(user *tgbotapi.User) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 8))
	err := binary.Write(buf, binary.LittleEndian, int64(user.ID))
	if err != nil {
		return nil, errors.Wrap(err, "can't encode user id")
	}
	key := buf.Bytes()
	return key, nil
}

func serializeBytes(w io.Writer, order binary.ByteOrder, data []byte, lengthSize int) error {
	var err error
	switch lengthSize {
	case 8:
		err = binary.Write(w, order, int8(len(data)))
	case 16:
		err = binary.Write(w, order, int16(len(data)))
	case 32:
		err = binary.Write(w, order, int32(len(data)))
	case 64:
		err = binary.Write(w, order, int64(len(data)))
	default:
		return errors.Errorf("known length size")
	}
	if err != nil {
		return errors.Wrap(err, "can't encode length")
	}
	_, err = w.Write(data)
	if err != nil {
		return errors.Wrap(err, "can't write data")
	}
	return nil
}

func deserializeBytes(r io.Reader, order binary.ByteOrder, lengthSize int) ([]byte, error) {
	var err error
	length := 0
	switch lengthSize {
	case 8:
		var length8 int8
		err = binary.Read(r, order, &length8)
		length = int(length8)
	case 16:
		var length16 int16
		err = binary.Read(r, order, &length16)
		length = int(length16)
	case 32:
		var length32 int32
		err = binary.Read(r, order, &length32)
		length = int(length32)
	case 64:
		var length64 int64
		err = binary.Read(r, order, &length64)
		length = int(length64)
	default:
		return nil, errors.Errorf("known length size")
	}
	if err != nil {
		return nil, errors.Wrap(err, "can't decode length")
	}
	ret := make([]byte, length)
	_, err = r.Read(ret)
	if err != nil {
		return nil, errors.Wrap(err, "can't read data")
	}
	return ret, nil
}

func serializeRuntime(runtime *db.Runtime) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := binary.Write(buf, binary.LittleEndian, runtime.Uid)
	if err != nil {
		return nil, errors.Wrap(err, "can't encode user id")
	}
	err = serializeBytes(buf, binary.LittleEndian, []byte(runtime.SessionToken), 16)
	if err != nil {
		return nil, errors.Wrap(err, "can't encode session token")
	}
	err = serializeBytes(buf, binary.LittleEndian, []byte(runtime.IKSM), 8)
	if err != nil {
		return nil, errors.Wrap(err, "can't encode iksm")
	}
	err = serializeBytes(buf, binary.LittleEndian, []byte(runtime.Language), 8)
	if err != nil {
		return nil, errors.Wrap(err, "can't encode language")
	}
	err = binary.Write(buf, binary.LittleEndian, int16(runtime.Timezone))
	if err != nil {
		return nil, errors.Wrap(err, "can't encode user timezone")
	}
	return buf.Bytes(), nil
}

func deserializeRuntime(data []byte) (*db.Runtime, error) {
	runtime := &db.Runtime{}
	buf := bytes.NewBuffer(data)
	err := binary.Read(buf, binary.LittleEndian, &(runtime.Uid))
	if err != nil {
		return nil, errors.Wrap(err, "can't decode user id")
	}

	sessionToken, err := deserializeBytes(buf, binary.LittleEndian, 16)
	if err != nil {
		return nil, errors.Wrap(err, "can't decode session token")
	}
	runtime.SessionToken = string(sessionToken)

	iksm, err := deserializeBytes(buf, binary.LittleEndian, 8)
	if err != nil {
		return nil, errors.Wrap(err, "can't decode iksm")
	}
	runtime.IKSM = iksm

	lang, err := deserializeBytes(buf, binary.LittleEndian, 8)
	if err != nil {
		return nil, errors.Wrap(err, "can't decode language")
	}
	runtime.Language = string(lang)

	var timezone int16
	err = binary.Read(buf, binary.LittleEndian, &(timezone))
	if err != nil {
		return nil, errors.Wrap(err, "can't decode user timezone")
	}
	runtime.Timezone = int(timezone)
	return runtime, nil
}
