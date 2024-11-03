// main.go
package main

import (
    "project_go/api"
    "project_go/config"
    "github.com/kataras/iris/v12"
)

func main() {
    config.InitCassandra()
    defer config.CloseSession()

    app := iris.New()
    app.Post("/users", api.CreateUser)
    app.Get("/users/{id}", api.GetUser)
    app.Put("/users/{id}", api.UpdateUser)
    app.Delete("/users/{id}", api.DeleteUser)
    app.Get("/users", api.ListUsers)

    app.Listen(":" + "8080")
}
