package dump

import (
	"os"

	"github.com/mailru/easyjson"
	"github.com/pkg/errors"
)

// Dumper saves splatoon2 item keys to files for translation.
type Dumper interface {
	// Get returns the newest dumping object.
	// Key should be added at Config. File is determined by the pair of key and file in Config.
	Get(key string, obj easyjson.MarshalerUnmarshaler) (easyjson.MarshalerUnmarshaler, error)
	// Load loads the dumping object from file.
	// Key should be added at Config. File is determined by the pair of key and file in Config.
	Load(key string, obj easyjson.MarshalerUnmarshaler) error
	// Save saves the dumping object to file.
	// Key should be added at Config. File is determined by the pair of key and file in Config.
	Save(key string, obj easyjson.MarshalerUnmarshaler) error
}

type dumperImpl struct {
	targets map[string]string
	objs    map[string]easyjson.MarshalerUnmarshaler
}

// New returns a new Dumper.
func New(config Config) Dumper {
	ret := &dumperImpl{
		targets: make(map[string]string),
		objs:    make(map[string]easyjson.MarshalerUnmarshaler),
	}
	for k, v := range config.targets {
		ret.targets[k] = v
	}
	return ret
}

func (d *dumperImpl) Get(key string, obj easyjson.MarshalerUnmarshaler) (easyjson.MarshalerUnmarshaler, error) {
	if v, ok := d.objs[key]; ok {
		return v, nil
	}
	err := d.Load(key, obj)
	return obj, err
}

func (d *dumperImpl) Load(key string, obj easyjson.MarshalerUnmarshaler) error {
	fileName := d.targets[key]
	file, err := os.Open(fileName)
	if err != nil {
		return errors.Wrap(err, "can't open dumping file: "+fileName)
	}
	err = easyjson.UnmarshalFromReader(file, obj)
	if err != nil {
		return errors.Wrap(err, "can't unmarshal from file: "+fileName)
	}
	d.objs[key] = obj
	return nil
}

func (d *dumperImpl) Save(key string, obj easyjson.MarshalerUnmarshaler) error {
	fileName := d.targets[key]
	file, err := os.Create(fileName)
	if err != nil {
		return errors.Wrap(err, "can't open dumping file: "+fileName)
	}
	_, err = easyjson.MarshalToWriter(obj, file)
	if err != nil {
		return errors.Wrap(err, "can't marshal to file: "+fileName)
	}
	d.objs[key] = obj
	return nil
}
