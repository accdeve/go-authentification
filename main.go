package main

import (
	"crud_user/db"
	"crud_user/router"
)

func main() {
	db.ConnectDB()

	db.MigrateDB()

	r := router.SetupRouter()

	r.Run(":8080")	
}