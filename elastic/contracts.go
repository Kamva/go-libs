package elastic

type Type interface {
	GetIndexName() string
	GetTypeName() string
	GetIndexMapping() string
	GetBody() interface{}
}
