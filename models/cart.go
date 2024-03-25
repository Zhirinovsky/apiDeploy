package models

import (
	"github.com/gin-gonic/gin"
	"github.com/my/repo/bin"
	"net/http"
)

type Cart struct {
	ID         int
	Amount     int
	Visibility bool
	UserId     string `db:"user_id"`
	ProductId  string `db:"product_id"`
	User       User
	Product    Product
	IsExist    bool `db:"is_exist"`
}

func GetCarts(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		var carts []Cart
		err := bin.DB.Select(&carts, "select * from cart_position")
		bin.GlobalCheck(err)
		for i := range carts {
			err = bin.DB.Get(&carts[i].User, "select * from \"User\" where id = $1", carts[i].UserId)
			bin.GlobalCheck(err)
			err = bin.DB.Get(&carts[i].Product, "select * from Product where id = $1", carts[i].ProductId)
			bin.GlobalCheck(err)
		}
		bin.FinalCheck(c, carts)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}
