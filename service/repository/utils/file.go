package utils

import (
	"github.com/mailru/easyjson"
	"github.com/pkg/errors"
	"os"
)

func MarshalToFile(fileName string, obj easyjson.Marshaler) error {
	file, err := os.Create(fileName)
	if err != nil {
		return errors.Wrap(err, "can't open dumping file: "+fileName)
	}
	_, err = easyjson.MarshalToWriter(obj, file)
	if err != nil {
		return errors.Wrap(err, "can't marshal to file: "+fileName)
	}
	return nil
}

func UnmarshalFromFile(fileName string, obj easyjson.Unmarshaler) error {
	file, err := os.Open(fileName)
	if err != nil {
		return errors.Wrap(err, "can't open dumping file: "+fileName)
	}
	err = easyjson.UnmarshalFromReader(file, obj)
	if err != nil {
		return errors.Wrap(err, "can't unmarshal from file: "+fileName)
	}
	return nil
}
