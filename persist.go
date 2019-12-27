package ddd

type Persist interface {
	Load(id string, object interface{}) error
	Save(id string, object interface{}) error
}

