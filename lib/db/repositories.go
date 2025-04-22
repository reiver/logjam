package db

type TDBFieldType string

const (
	String            TDBFieldType = "string"
	Integer           TDBFieldType = "integer"
	UnsignedInteger64 TDBFieldType = "uint64"
	Float             TDBFieldType = "float"
	Boolean           TDBFieldType = "boolean"
	DateTime          TDBFieldType = "datetime"
)

type Field struct {
	Name string
	Type TDBFieldType
}

type IDBService interface {
	Ping() error
	CreateTableIfNotExists(cname string, fields []Field) error
	GetById(table, id string) (Record, error)
	GetByFilter(table string, filters map[string]any) ([]Record, error)
	GetAll(table string) ([]Record, error)
	Update(table, id string, data map[string]any) error
	UpdateByFilter(table string, filter map[string]any, data map[string]any) error
	Delete(table, id string) error
}
