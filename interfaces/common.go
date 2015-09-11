package interfaces

//Config interface is describe methods for works with it
type Config interface {
	SetBasePath(path string)
	Load(env string) error
	Get(key string) (interface{}, error)
	GetAll() *map[string]interface{}
}

//Orm interface is describe methods for persistant data in DB
type Orm interface {
	Connect(cong Config) error
	Create(interface{}) error
	Find(interface{}, interface{}) error
	First(interface{}, interface{}) error
	Update(interface{}, map[string]interface{}) error
	Save(interface{}) error
	Delete(interface{}) error
}

//App interface is describe methods for simple application
type App interface {
	GetDB() Orm
	GetConfig() Config
}
