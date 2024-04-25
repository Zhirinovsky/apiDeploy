package models

import (
	"github.com/gin-gonic/gin"
	"github.com/my/repo/bin"
	"net/http"
)

type Cart struct {
	ID        int
	Amount    int
	ProductID int `db:"product_id"`
	UserID    int `db:"user_id"`
	Product   Product
	User      User
}

func GetCartPositions(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		var carts []Cart
		err := bin.DB.Select(&carts, "select * from Cart_Position")
		bin.GlobalCheck(err)
		for i := range carts {
			err = bin.DB.Get(&carts[i].Product, "select * from Product where id = $1", carts[i].ProductID)
			bin.GlobalCheck(err)
			err = bin.DB.Get(&carts[i].User, "select * from \"User\" where id = $1", carts[i].UserID)
			bin.GlobalCheck(err)
		}
		bin.FinalCheck(c, carts)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func GetCartPositionByID(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		var cart Cart
		err := bin.DB.Get(&cart, "select * from Cart_Position where id = $1", id)
		bin.GlobalCheck(err)
		err = bin.DB.Get(&cart.Product, "select * from Product where id = $1", cart.ProductID)
		bin.GlobalCheck(err)
		err = bin.DB.Get(&cart.User, "select * from \"User\" where id = $1", cart.UserID)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, cart)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func PostCartPosition(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		var cart Cart
		err := c.BindJSON(&cart)
		bin.GlobalCheck(err)
		_, err = bin.DB.Exec("insert into cart_position (amount, product_id, user_id) values ($1, $2, $3)", cart.Amount, cart.ProductID, cart.UserID)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, cart)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func PutCartPosition(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		var cart Cart
		err := c.BindJSON(&cart)
		bin.GlobalCheck(err)
		_, err = bin.DB.Exec("update cart_position set amount = $2, product_id = $3, user_id = $4 where id = $1", id, cart.Amount, cart.ProductID, cart.UserID)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, cart)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func DeleteCartPosition(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		_, err := bin.DB.Exec("delete from cart_position where id = $1", id)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, nil)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}
