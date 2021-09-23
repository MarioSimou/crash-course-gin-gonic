package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/MarioSimou/crash-course-gin-gonic/internal"
	"github.com/gin-gonic/gin"
)

type postProduct struct {
	Name        string  `json:"name" binding:"required"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Description string  `json:"description" binding:"omitempty,max=250"`
}

type Product struct {
	GUID        string  `json:"guid"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	CreatedAt   string  `json:"createdAt"`
}

func PostProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload postProduct
		var ctx = c.Request.Context()

		if e := c.ShouldBindJSON(&payload); e != nil {
			var res = internal.NewHTTPResponse(http.StatusBadRequest, e)
			c.JSON(http.StatusBadRequest, res)
			return
		}

		var guid = uuid.New().String()
		var createdAt = time.Now().Format(time.RFC3339)
		if _, e := db.ExecContext(ctx, "INSERT INTO products(guid, name, price, description, createdAt) VALUES(?,?,?,?,?)", guid, payload.Name, payload.Price, payload.Description, createdAt); e != nil {
			var res = internal.NewHTTPResponse(http.StatusInternalServerError, e)
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		var product Product
		var row = db.QueryRow("SELECT guid,name,price,description,createdAt FROM products WHERE guid=?", guid)

		if e := row.Scan(&product.GUID, &product.Name, &product.Price, &product.Description, &product.CreatedAt); e != nil {
			var res = internal.NewHTTPResponse(http.StatusInternalServerError, e)
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		var res = internal.NewHTTPResponse(http.StatusCreated, product)

		c.Writer.Header().Add("Location", fmt.Sprintf("/products/%s", guid))
		c.JSON(http.StatusCreated, res)
	}
}
