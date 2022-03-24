// main.go
package main

import (
    "log"
    "os"

    "github.com/joho/godotenv"



    "github.com/khelechy/rielzapi/api/controllers"
)

// @title Rielz API
// @version 1.0
// @description This is a simple real estate management service
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email soberkoder@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:5000
// @BasePath /
func main() {
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    app := controllers.App{}
    app.Initialize(
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PASSWORD"))

    app.RunServer()
}