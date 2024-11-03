// api/user.go
package api

import (
	"github.com/gocql/gocql"
	"github.com/kataras/iris/v12"
	"net/http"
	"project_go/config"
	"project_go/models"
)

func CreateUser(ctx iris.Context) {
	var user models.User
	if err := ctx.ReadJSON(&user); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid JSON"})
		return
	}
	user.ID = gocql.TimeUUID()
	query := `INSERT INTO users (id, name, email) VALUES (?, ?, ?)`
	if err := config.Session.Query(query, user.ID, user.Name, user.Email).Exec(); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to create user"})
		return
	}
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(user)
}

func GetUser(ctx iris.Context) {
	id := ctx.Params().Get("id")
	uuid, err := gocql.ParseUUID(id)
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid user ID"})
		return
	}
	var user models.User
	query := `SELECT id, name, email FROM users WHERE id = ? LIMIT 1`
	if err := config.Session.Query(query, uuid).Consistency(gocql.One).Scan(&user.ID, &user.Name, &user.Email); err != nil {
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(iris.Map{"error": "User not found"})
		return
	}
	ctx.JSON(user)
}

func UpdateUser(ctx iris.Context) {
	id := ctx.Params().Get("id")
	uuid, err := gocql.ParseUUID(id)
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid user ID"})
		return
	}
	var user models.User
	if err := ctx.ReadJSON(&user); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid JSON"})
		return
	}
	query := `UPDATE users SET name = ?, email = ? WHERE id = ?`
	if err := config.Session.Query(query, user.Name, user.Email, uuid).Exec(); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to update user"})
		return
	}
	ctx.JSON(iris.Map{"message": "User updated successfully"})
}

func DeleteUser(ctx iris.Context) {
	id := ctx.Params().Get("id")
	uuid, err := gocql.ParseUUID(id)
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid user ID"})
		return
	}
	query := `DELETE FROM users WHERE id = ?`
	if err := config.Session.Query(query, uuid).Exec(); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to delete user"})
		return
	}
	ctx.JSON(iris.Map{"message": "User deleted successfully"})
}

func ListUsers(ctx iris.Context) {
	var users []models.User
	query := `SELECT id, name, email FROM users`
	iter := config.Session.Query(query).Iter()
	var user models.User
	for iter.Scan(&user.ID, &user.Name, &user.Email) {
		users = append(users, user)
	}
	if err := iter.Close(); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to fetch users"})
		return
	}
	ctx.JSON(users)
}
