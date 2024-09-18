package db

import (
	"crud_user/model"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func FuncDB() {
	ConnectDB()
	MigrateDB()
}

func ConnectDB() {
	dsn := "root:@/belajar"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error connecting database", err.Error())
	} else {
		fmt.Println("Connection to the mysql")
	}

	DB = db

}

func MigrateDB(){
	err := DB.AutoMigrate(&model.User{})

	if err != nil{
		log.Fatal("error migration table user", err.Error())
	} else{
		fmt.Println("success migration table")
	}
}
