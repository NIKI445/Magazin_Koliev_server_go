package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"RestApiGo/internal/database"
	"RestApiGo/internal/models"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Неверный формат запроса",
		})
		return
	}

	var user models.User
	query := `SELECT user_id, full_name, middle_name, email FROM users WHERE email = $1 AND password = $2`

	err := database.DB.QueryRow(query, req.Email, req.Password).Scan(
		&user.UserID, &user.FullName, &user.MiddleName, &user.Email,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Неверный логин или пароль",
			})
			return
		}
		log.Println("Login error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Ошибка выполнения запроса",
		})
		return
	}

	// Получаем корзину пользователя
	var cartID int
	cartQuery := `SELECT cart_id FROM carts WHERE user_id = $1`
	err = database.DB.QueryRow(cartQuery, user.UserID).Scan(&cartID)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Cart query error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Ошибка получения корзины",
		})
		return
	}

	response := models.LoginResponse{
		UserID:     user.UserID,
		FullName:   user.FullName,
		Email:      user.Email,
		CartID:     cartID,
		MiddleName: user.MiddleName,
	}
	log.Println(user.MiddleName)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// POST /auth/signup
func (h *AuthHandler) Signup(c *gin.Context) {
	var req models.SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Неверный формат запроса",
		})
		return
	}

	// Проверяем существование пользователя
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := database.DB.QueryRow(checkQuery, req.Email).Scan(&exists)
	if err != nil {
		log.Println("Check user error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Ошибка проверки пользователя",
		})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Пользователь уже существует",
		})
		return
	}

	// Создаем пользователя
	var userID int
	insertQuery := `INSERT INTO users (full_name, middle_name, email, password) 
	                VALUES ($1, $2, $3, $4) RETURNING user_id`
	err = database.DB.QueryRow(insertQuery, req.FullName, req.MiddleName, req.Email, req.Password).Scan(&userID)
	if err != nil {
		log.Println("Insert user error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Ошибка регистрации",
		})
		return
	}

	// Создаем корзину
	var cartID int
	cartQuery := `INSERT INTO carts (user_id) VALUES ($1) RETURNING cart_id`
	err = database.DB.QueryRow(cartQuery, userID).Scan(&cartID)
	if err != nil {
		log.Println("Create cart error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Ошибка создания корзины",
		})
		return
	}

	// Получаем данные пользователя
	var user models.User
	getQuery := `SELECT user_id, full_name, middle_name, email FROM users WHERE user_id = $1`
	err = database.DB.QueryRow(getQuery, userID).Scan(
		&user.UserID, &user.FullName, &user.MiddleName, &user.Email,
	)
	if err != nil {
		log.Println("Get user error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Ошибка получения данных",
		})
		return
	}

	response := models.LoginResponse{
		UserID:     user.UserID,
		FullName:   user.FullName,
		Email:      user.Email,
		CartID:     cartID,
		MiddleName: user.MiddleName,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}
