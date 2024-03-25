package models

import (
	"github.com/gin-gonic/gin"
	"github.com/my/repo/bin"
	"net/http"
)

type Position struct {
	ID            int
	CheckoutPrice float64 `db:"checkout_price"`
	Amount        int
	OrderID       int `db:"order_id"`
	ProductID     int `db:"product_id"`
	Order         Order
	Product       Product
}

func GetOrderPositions(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		var positions []Position
		err := bin.DB.Select(&positions, "select * from order_position")
		bin.GlobalCheck(err)
		for i := range positions {
			err = bin.DB.Get(&positions[i].Order, "select * from \"Order\" where id = $1", positions[i].OrderID)
			bin.GlobalCheck(err)
			err = bin.DB.Get(&positions[i].Product, "select * from product where id = $1", positions[i].ProductID)
			bin.GlobalCheck(err)
		}
		bin.FinalCheck(c, positions)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func GetOrderPositionByID(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		var position Position
		err := bin.DB.Get(&position, "select * from order_position where id = $1", id)
		bin.GlobalCheck(err)
		err = bin.DB.Get(&position.Order, "select * from \"Order\" where id = $1", position.OrderID)
		bin.GlobalCheck(err)
		err = bin.DB.Get(&position.Product, "select * from product where id = $1", position.ProductID)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, position)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func PostOrderPosition(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		var position Position
		err := c.BindJSON(&position)
		bin.GlobalCheck(err)
		_, err = bin.DB.Exec("insert into order_position (checkout_price, amount, order_id, product_id) values ($1, $2, $3, $4)", position.CheckoutPrice, position.Amount, position.OrderID, position.ProductID)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, position)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func PutOrderPosition(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		var position Position
		err := c.BindJSON(&position)
		bin.GlobalCheck(err)
		_, err = bin.DB.Exec("update order_position set checkout_price = $2, amount = $3, order_id = $4, product_id = $5 where id = $1", id, position.CheckoutPrice, position.Amount, position.OrderID, position.ProductID)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, position)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func DeleteOrderPosition(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		_, err := bin.DB.Exec("delete from order_position where id = $1", id)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, nil)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}
