package handlers

import (
	"log"
	"net/http"
	"strconv"

	"RestApiGo/internal/database"
	"RestApiGo/internal/models"

	"github.com/gin-gonic/gin"
)

type RequestHandler struct{}

func NewRequestHandler() *RequestHandler {
	return &RequestHandler{}
}

// GET /api/products - получение всех товаров
func (h *RequestHandler) GetProducts(c *gin.Context) {
	query := `SELECT product_id, name, price, description, image, category, stock FROM products`

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Println("Ошибка получения товаров:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Ошибка выполнения запроса",
		})
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ProductID, &p.Name, &p.Price, &p.Description, &p.Image, &p.Category, &p.Stock)
		if err != nil {
			log.Println("Ошибка сканирования товара:", err)
			continue
		}
		products = append(products, p)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    products,
	})
}

// POST /api/product - добавление товара в корзину
func (h *RequestHandler) AddToCart(c *gin.Context) {
	var req models.AddToCartRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Неверный формат запроса",
		})
		return
	}

	if req.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Не указан userID",
		})
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Ошибка сервера",
		})
		return
	}
	defer tx.Rollback()

	// Создаём корзину если не существует (без created_at/updated_at)
	_, err = tx.Exec(`
		INSERT INTO carts (cart_id, user_id) 
		VALUES ($1, $2)
		ON CONFLICT (cart_id) DO UPDATE 
		SET user_id = EXCLUDED.user_id
	`, req.CartID, req.UserID)

	if err != nil {
		log.Println("Ошибка создания/обновления корзины:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Ошибка подготовки корзины",
		})
		return
	}

	// UPSERT для cart_items (без created_at/updated_at)
	query := `
		INSERT INTO cart_items (cart_id, product_id, quantity) 
		VALUES ($1, $2, $3)
		ON CONFLICT (cart_id, product_id) 
		DO UPDATE SET 
			quantity = cart_items.quantity + $3
		RETURNING cart_item_id, quantity
	`

	var cartItemID int
	var newQuantity int
	err = tx.QueryRow(query, req.CartID, req.ProductID, req.Quantity).Scan(&cartItemID, &newQuantity)
	if err != nil {
		log.Println("Ошибка UPSERT:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Ошибка добавления товара",
		})
		return
	}

	// Получаем данные товара
	var cartItem models.EnrichedCartItem
	getItemQuery := `
		SELECT 
			ci.cart_item_id,
			ci.cart_id,
			ci.product_id,
			ci.quantity,
			p.name,
			p.price,
			p.description,
			p.image,
			p.category,
			(ci.quantity * p.price) as total_price
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.product_id
		WHERE ci.cart_item_id = $1
	`

	err = tx.QueryRow(getItemQuery, cartItemID).Scan(
		&cartItem.CartItemID,
		&cartItem.CartID,
		&cartItem.ProductID,
		&cartItem.Quantity,
		&cartItem.Name,
		&cartItem.Price,
		&cartItem.Description,
		&cartItem.Image,
		&cartItem.Category,
		&cartItem.TotalPrice,
	)

	if err != nil {
		log.Println("Ошибка получения товара:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Ошибка получения данных товара",
		})
		return
	}

	if err = tx.Commit(); err != nil {
		log.Println("Ошибка коммита:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Ошибка сохранения",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Товар добавлен в корзину",
		"data":    cartItem,
	})
}



// GET /api/cart/:id - получение корзины
func (h *RequestHandler) GetCart(c *gin.Context) {
	idParam := c.Param("id")
	cartID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Неверный ID корзины",
		})
		return
	}

	cartItems, err := h.getCartItems(cartID)
	if err != nil {
		log.Println("Ошибка получения корзины:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Ошибка получения корзины",
		})
		return
	}

	// Если cartItems == nil, отправляем пустой массив
	if cartItems == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    []interface{}{}, // пустой массив
		})
		return
	}

	log.Println("Корзина получена:", cartItems)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    cartItems,
	})
}

// Вспомогательная функция для получения товаров в корзине
func (h *RequestHandler) getCartItems(cartID int) ([]models.EnrichedCartItem, error) {
	query := `
        SELECT 
            ci.cart_item_id,
            ci.cart_id,
            ci.product_id,
            ci.quantity,
            p.name,
            p.price,
            p.description,
            p.image,
            p.category,
            (ci.quantity * p.price) as total_price
        FROM cart_items ci
        JOIN products p ON ci.product_id = p.product_id
        WHERE ci.cart_id = $1
    `

	rows, err := database.DB.Query(query, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.EnrichedCartItem
	for rows.Next() {
		var item models.EnrichedCartItem
		err := rows.Scan(
			&item.CartItemID,
			&item.CartID,
			&item.ProductID,
			&item.Quantity,
			&item.Name,
			&item.Price,
			&item.Description,
			&item.Image,
			&item.Category,
			&item.TotalPrice,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// POST /api/cart/count/:id - обновление количества товара в корзине
func (h *RequestHandler) UpdateCartCount(c *gin.Context) {
	idParam := c.Param("id")
	cartItemID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Неверный ID товара в корзине",
		})
		return
	}

	var req models.UpdateCartCountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Неверный формат запроса",
		})
		return
	}

	if req.Quantity <= 0 {
		// Если количество 0 или меньше, удаляем товар
		deleteQuery := `DELETE FROM cart_items WHERE cart_item_id = $1`
		_, err = database.DB.Exec(deleteQuery, cartItemID)
		if err != nil {
			log.Println("Delete cart item error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Ошибка удаления товара",
			})
			return
		}
	} else {
		// Обновляем количество
		updateQuery := `UPDATE cart_items SET quantity = $1 WHERE cart_item_id = $2`
		_, err = database.DB.Exec(updateQuery, req.Quantity, cartItemID)
		if err != nil {
			log.Println("Update cart count error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Ошибка обновления количества",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    req,
	})
}

// DELETE /api/cart/:id - удаление конкретного товара из корзины
func (h *RequestHandler) DeleteCartItem(c *gin.Context) {
	idParam := c.Param("id")
	cartItemID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Неверный ID товара в корзине",
		})
		return
	}

	// Проверяем существование товара
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM cart_items WHERE cart_item_id = $1)`
	err = database.DB.QueryRow(checkQuery, cartItemID).Scan(&exists)
	if err != nil {
		log.Println("Check cart item error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Ошибка проверки товара",
		})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Запись не найдена",
		})
		return
	}

	// Удаляем товар
	deleteQuery := `DELETE FROM cart_items WHERE cart_item_id = $1`
	result, err := database.DB.Exec(deleteQuery, cartItemID)
	if err != nil {
		log.Println("Delete cart item error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Ошибка удаления товара",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"message":      "Товар успешно удален",
		"affectedRows": rowsAffected,
	})
}

// DELETE /api/cartAll/:id - очистка всей корзины
func (h *RequestHandler) ClearCart(c *gin.Context) {
	idParam := c.Param("id")
	cartID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Неверный ID корзины",
		})
		return
	}

	// Проверяем наличие товаров в корзине
	var count int
	checkQuery := `SELECT COUNT(*) FROM cart_items WHERE cart_id = $1`
	err = database.DB.QueryRow(checkQuery, cartID).Scan(&count)
	if err != nil {
		log.Println("Check cart error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Ошибка проверки корзины",
		})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Корзина пуста",
		})
		return
	}

	// Очищаем корзину
	deleteQuery := `DELETE FROM cart_items WHERE cart_id = $1`
	result, err := database.DB.Exec(deleteQuery, cartID)
	if err != nil {
		log.Println("Clear cart error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Ошибка очистки корзины",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"message":      "Корзина успешно очищена",
		"affectedRows": rowsAffected,
	})
}
