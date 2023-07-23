package config

type Config struct {
	DBHost string
}

var AppConfig = Config{

	DBHost: "popop",
}
