package enum

import (
	"reflect"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"telegram-splatoon2-bot/common/log"
)

// Enum is the basic type of all enumeration type.
type Enum int64

var enumType = reflect.TypeOf(Enum(0))

// Assign automatically assigns different value to all fields in a structure whose all fields are Enum type.
// inStructPtr should be a pointer of structure.
// Assign returns assigned inStructPtr.
func Assign(inStructPtr interface{}) interface{} {
	rType := reflect.TypeOf(inStructPtr)
	rVal := reflect.ValueOf(inStructPtr)
	if rType.Kind() == reflect.Ptr {
		rType = rType.Elem()
		rVal = rVal.Elem()
	} else {
		log.Panic("inStructPtr must be ptr to struct")
	}

	for i := 0; i < rType.NumField(); i++ {
		t := rType.Field(i)
		f := rVal.Field(i)
		if f.Kind() == reflect.Struct {
			Assign(f.Addr().Interface())
		} else {
			if f.Type().ConvertibleTo(enumType) {
				f.Set(reflect.ValueOf(gen.next()).Convert(f.Type()))
			} else {
				log.Panic("can't assign enum value", zap.Error(errors.Errorf(t.Name+" type is not Enum")))
			}
		}
	}
	return inStructPtr
}
