package main

import (
	"fmt"
	"os"
  "gorm.io/driver/postgres"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
  "github.com/Bikram-ghuku/SyncChatServerGo/models"
)

var DB *gorm.DB
func main() {
  godotenv.Load()
  dbhost := os.Getenv("POSTGRES_HOST");
  dbname := os.Getenv("POSTGRES_DBNAME");
  dbuser := os.Getenv("POSTGRES_USER");
  dbpswd := os.Getenv("POSTGRES_PASSWORD");
  conn := fmt.Sprintf("postgres://%s:%s@%s/%s", dbuser, dbpswd, dbhost, dbname);
  fmt.Println(conn)
  var db, err = gorm.Open(postgres.Open(conn), &gorm.Config{});
  if err != nil{
    panic("DB Connection error")
  }
  
  DB = db;
  fmt.Println("DB Connection Successful")
  DB.AutoMigrate(&models.Users{}, &models.Chats{})
  fmt.Println("Migration Done")
}
