package config

type Server struct {
	Port string `env:"PORT,notEmpty"`
}

func NewServer(port string) Server {
	return Server{Port: port}
}

func (Server Server) GetEchoPort() string {
	return ":" + Server.Port
}
