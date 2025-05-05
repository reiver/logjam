package db

type TDBFieldType string

const (
	StringType            TDBFieldType = "text"     // single‑line text
	EmailType             TDBFieldType = "email"    // email address
	URLType               TDBFieldType = "url"      // URL
	TextType              TDBFieldType = "text"     // multi‑line/plain text
	EditorType            TDBFieldType = "editor"   // rich‑text HTML editor
	IntegerType           TDBFieldType = "integer"  // int/float stored as integer
	UnsignedInteger64Type TDBFieldType = "uint64"   // unsigned 64‑bit integer
	FloatType             TDBFieldType = "float"    // floating‑point number
	BooleanType           TDBFieldType = "bool"     // boolean
	DateType              TDBFieldType = "date"     // datetime string
	AutodateType          TDBFieldType = "autodate" // auto‑set timestamp on create/update
	SelectType            TDBFieldType = "select"   // single/multi select from options
	FileType              TDBFieldType = "file"     // file upload(s)
	RelationType          TDBFieldType = "relation" // record relation(s)
	JSONType              TDBFieldType = "json"     // arbitrary JSON
	GeoPointType          TDBFieldType = "geoPoint" // {lon,lat} coords
)

type Field struct {
	Name   string
	Type   TDBFieldType
	Unique bool
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
	Insert(table string, data map[string]any) (id string, err error)
}
