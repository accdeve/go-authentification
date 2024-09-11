package handler

import (
	"crud_user/auth"
	"crud_user/db"
	"crud_user/model"
	"crud_user/util"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(c *gin.Context) {
	var userInput model.User
	var userFromDB model.User

	if err := c.ShouldBind(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if util.ValidateEmail(userInput.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	query := "SELECT id, username, email, password FROM users WHERE email = ?"
	err := db.DB.QueryRow(query, userInput.Email).Scan(&userFromDB.ID, &userFromDB.Username, &userFromDB.Email, &userFromDB.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Email not registrasion"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(userInput.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect password"})
		return
	}

	tokenString, err := auth.GenerateJWT(userFromDB.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login Successfull", "token": tokenString})
}

func RegistrationHandler(c *gin.Context) {
	var user model.User
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if util.ValidateEmail(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	if user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email can't empty"})
		return
	} else if user.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username can't empty"})
		return
	} else if user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password can't empty"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	query := "INSERT INTO users (username, email, password) VALUES (?,?,?)"
	_, err = db.DB.Exec(query, user.Username, user.Email, string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user in database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success add user"})
}

func GetAllUserHandler(c *gin.Context) {
	rows, err := db.DB.Query("SELECT id, username, email, password FROM USERS")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to select users from database"})
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan user"})
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success get all user", "data": users})

}

func UpdatePasswordHandler(c *gin.Context) {
	var user model.User

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Query("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter id is required"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	query := "UPDATE users SET password = ? WHERE id = ?"
	_, err = db.DB.Exec(query, string(hashedPassword), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password Update sucesfully"})
}

func DeleteAllUserHandler(c *gin.Context) {
	email := c.Query("email")

	if email == ""{
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter email is required"})
		return
	}

	query := "DELETE FROM users WHERE id = ?"
	_, err := db.DB.Exec(query, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user from the database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
