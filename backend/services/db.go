package services

import (
	"fmt"
	"os"

	"github.com/Bikram-ghuku/SyncChatServerGo/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	var DB *gorm.DB
	godotenv.Load()
	dbhost := os.Getenv("POSTGRES_HOST")
	dbname := os.Getenv("POSTGRES_DBNAME")
	dbuser := os.Getenv("POSTGRES_USER")
	dbpswd := os.Getenv("POSTGRES_PASSWORD")
	dbport := os.Getenv("PGPORT")
	conn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbuser, dbpswd, dbhost, dbport, dbname)
	fmt.Println(conn)
	var db, err = gorm.Open(postgres.Open(conn), &gorm.Config{})
	if err != nil {
		panic("DB Connection error")
	}

	DB = db
	fmt.Println("DB Connection Successful")
	DB.AutoMigrate(&models.Users{}, &models.Chats{}, &models.Messages{})
	fmt.Println("Migration Done")
	return DB
}
