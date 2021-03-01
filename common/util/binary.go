package util

import (
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

var Binary serialization

type serialization struct {}


func (serialization)WriteBytes(w io.Writer, order binary.ByteOrder, data []byte, lengthSize int) error {
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

func (serialization)ReadBytes(r io.Reader, order binary.ByteOrder, lengthSize int) ([]byte, error) {
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
