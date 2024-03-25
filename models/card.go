package models

import (
	"github.com/gin-gonic/gin"
	"github.com/my/repo/bin"
	"net/http"
)

type Card struct {
	ID       int
	Date     string
	Discount int
}

func GetCards(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		var cards []Card
		err := bin.DB.Select(&cards, "select * from discount_card")
		bin.GlobalCheck(err)
		bin.FinalCheck(c, cards)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func GetCardByID(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		var card Card
		err := bin.DB.Get(&card, "select * from discount_card where id = $1", id)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, card)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func PostCards(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		var card Card
		err := c.BindJSON(&card)
		bin.GlobalCheck(err)
		_, err = bin.DB.Exec("insert into discount_card (id, date, discount) values ($1, $2, $3)", card.ID, card.Date, card.Discount)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, card)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func PutCards(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		var card Card
		err := c.BindJSON(&card)
		bin.GlobalCheck(err)
		_, err = bin.DB.Exec("update discount_card set date = $2, discount = $3 where id = $1", id, card.Date, card.Discount)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, card)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func DeleteCards(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		_, err := bin.DB.Exec("delete from discount_card where id = $1", id)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, nil)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}
