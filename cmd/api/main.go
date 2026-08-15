package main

import "github.com/HemlockPham7/user-service/internal/infrastructure"

// @title User Service API
// @version 1.0.0
// @description This is the API documentation for the User service in the Bookmark Management system.
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Init api
	app := infrastructure.CreateAPI()

	// Run api
	err := app.Start()
	if err != nil {
		panic(err)
	}
}
