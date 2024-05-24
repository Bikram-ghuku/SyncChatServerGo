package main

import (
	"github.com/Bikram-ghuku/SyncChatServerGo/services"
	"gorm.io/gorm"
)

var DB *gorm.DB
func main() {
  DB = services.InitDB();
}
