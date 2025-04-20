package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"taskflow/config"
	"taskflow/middlewares"
	"taskflow/models"
)

// RegisterUser User Registration
// @Summary      User Registration
// @Description  Create Users
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user  body      models.User  true  "User Information"
// @Success      200   {object}  models.User
// @Router       /register [post]
func RegisterUser(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// secret hashing
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "secret hashing failed"})
		return
	}
	input.Password = string(hashed)

	//Insert into the database
	result, err := config.DB.Exec("INSERT INTO users (username, password) VALUES (?, ?)",
		input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	id, _ := result.LastInsertId()
	input.ID = id

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful", "user": input})

}

// LoginUser Users login
// @Summary      User Login
// @Description  User logs in with username and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        loginInput  body      models.LoginInput  true  "Login Information"
// @Success      200  {object}  models.LoginResponse
// @Failure      401  {object}  models.ErrorResponse
// @Router       /login [post]
func LoginUser(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Request parameter fails"})
		return
	}
	// read from database
	var stored models.User
	err := config.DB.QueryRow(`
        SELECT id, username, password
        FROM users WHERE username = ?`, input.Username).Scan(&stored.ID, &stored.Username, &stored.Password)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "User is not exist"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(stored.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Username or Password is incorrect"})
		return
	}

	// Generate JWT Token
	tokenString, err := middlewares.GenerateToken(stored.ID, stored.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "GenerateToken Error"})
		return
	}
	// return successfully
	c.JSON(http.StatusOK, models.LoginResponse{
		Message: "Login successful",
		Token:   tokenString,
	})
}
