package configs

import (
	"github.com/joho/godotenv"
	"os"
)

type Configs struct {
	DSN    string
	Secret string
}

func LoadConfig() *Configs {
	err := godotenv.Load("../.env")
	if err != nil {
		panic("Ошибка при загрузке env файла")
	}
	return &Configs{
		DSN: os.Getenv("DSN"),
	}
}
