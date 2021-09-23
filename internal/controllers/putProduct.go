package controllers

import (
	"database/sql"
	"net/http"

	"github.com/MarioSimou/crash-course-gin-gonic/internal"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type putProduct struct {
	Name        string  `json:"name" binding:"required_without_all=Price Description"`
	Price       float64 `json:"price" binding:"omitempty,gt=0"`
	Description string  `json:"description" binding:"omitempty,max=250"`
}

func PutProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var binding guidBinding
		var payload putProduct
		var ctx = c.Request.Context()

		if e := c.ShouldBindUri(&binding); e != nil {
			var res = internal.NewHTTPResponse(http.StatusBadRequest, e)
			c.JSON(http.StatusBadRequest, res)
			return
		}

		if e := c.ShouldBindJSON(&payload); e != nil {
			var res = internal.NewHTTPResponse(http.StatusBadRequest, e)
			c.JSON(http.StatusBadRequest, res)
			return
		}

		var row = db.QueryRowContext(ctx, "SELECT name,price,description FROM products WHERE guid=?", binding.GUID)
		var currentProduct Product

		if e := row.Scan(&currentProduct.Name, &currentProduct.Price, &currentProduct.Description); e != nil {
			if e == sql.ErrNoRows {
				var res = internal.NewHTTPResponse(http.StatusNotFound, e)
				c.JSON(http.StatusNotFound, res)
				return
			}

			var res = internal.NewHTTPResponse(http.StatusInternalServerError, e)
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		var option = copier.Option{
			IgnoreEmpty: true,
			DeepCopy:    true,
		}

		if e := copier.CopyWithOption(&currentProduct, &payload, option); e != nil {
			var res = internal.NewHTTPResponse(http.StatusInternalServerError, e)
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		if _, e := db.ExecContext(ctx, "UPDATE products SET name=?,price=?,description=? WHERE guid=?", currentProduct.Name, currentProduct.Price, currentProduct.Description, binding.GUID); e != nil {
			var res = internal.NewHTTPResponse(http.StatusInternalServerError, e)
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		var updatedRow = db.QueryRowContext(ctx, "SELECT guid,name,price,description,createdAt FROM products WHERE guid=?", binding.GUID)
		var product Product

		if e := updatedRow.Scan(&product.GUID, &product.Name, &product.Price, &product.Description, &product.CreatedAt); e != nil {
			var res = internal.NewHTTPResponse(http.StatusInternalServerError, e)
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		var res = internal.NewHTTPResponse(http.StatusOK, product)
		c.JSON(http.StatusOK, res)
	}
}
