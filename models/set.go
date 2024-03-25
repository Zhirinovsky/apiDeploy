package models

import (
	"github.com/gin-gonic/gin"
	"github.com/my/repo/bin"
	"net/http"
)

type Set struct {
	ID               int
	Value            float64
	ProductID        int `db:"product_id"`
	CharacteristicID int `db:"characteristic_id"`
	Product          Product
	Characteristic   Characteristic
}

func GetSets(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		var sets []Set
		err := bin.DB.Select(&sets, "select * from Set")
		bin.GlobalCheck(err)
		for i, _ := range sets {
			err = bin.DB.Get(&sets[i].Product, "select * from product where id = $1", &sets[i].ProductID)
			bin.GlobalCheck(err)
			err = bin.DB.Get(&sets[i].Characteristic, "select * from characteristic where id = $1", &sets[i].CharacteristicID)
			bin.GlobalCheck(err)
		}
		bin.FinalCheck(c, sets)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func GetSetByID(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		var set Set
		err := bin.DB.Get(&set, "select * from Set where id = $1", id)
		bin.GlobalCheck(err)
		err = bin.DB.Get(&set.Product, "select * from product where id = $1", &set.ProductID)
		bin.GlobalCheck(err)
		err = bin.DB.Get(&set.Characteristic, "select * from characteristic where id = $1", &set.CharacteristicID)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, set)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func PostSet(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		var set Set
		err := c.BindJSON(&set)
		bin.GlobalCheck(err)
		_, err = bin.DB.Exec("insert into set (value, product_id, characteristic_id) values ($1, $2, $3)", set.Value, set.ProductID, set.CharacteristicID)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, set)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func PutSet(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		var set Set
		err := c.BindJSON(&set)
		bin.GlobalCheck(err)
		_, err = bin.DB.Exec("update set set value = $2, product_id = $3, characteristic_id = $4 where id = $1", id, set.Value, set.ProductID, set.CharacteristicID)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, set)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func DeleteSet(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		_, err := bin.DB.Exec("delete from set where id = $1", id)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, nil)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}
