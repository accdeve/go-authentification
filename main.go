package main

import (
	"crud_user/db"
	"crud_user/router"
)

func main() {
	db.FuncDB()

	r := router.SetupRouter()

	r.Run(":8080")
}
