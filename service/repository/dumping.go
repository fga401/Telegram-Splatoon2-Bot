package repository

type Dumper interface {
	Update(src interface{}) error
	Load() error
	Save() error
}
