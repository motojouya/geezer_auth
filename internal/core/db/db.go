package db

type DBAccess struct {
	Driver string `env:"DB_DRIVER,notEmpty"`
	Name   string `env:"DB_NAME,notEmpty"`
	Host   string `env:"DB_HOST,notEmpty"`
	Port   uint   `env:"DB_PORT,notEmpty"`
	User   string `env:"DB_USER,notEmpty"`
	Pass   string `env:"DB_PASSWORD,notEmpty"`
}

func NewDBAccess(
	driver string,
	name   string,
	host   string,
	port   uint,
	user   string,
	pass   string,
) DBAccess {
	return DBAccess{
		Driver: driver,
		Name:   name,
		Host:   host,
		Port:   port,
		User:   user,
		Pass:   pass,
	}
}
