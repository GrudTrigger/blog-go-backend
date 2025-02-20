package main

import (
	"backend/blog/internal/posts"
	"backend/blog/internal/user"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.Migrator().DropTable("users")

	if err != nil {
		log.Fatalf("Ошибка при удалении таблицы users: %v", err)
	}

	db.AutoMigrate(&user.User{}, &posts.Post{})
}