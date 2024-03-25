package models

import (
	"github.com/gin-gonic/gin"
	"github.com/my/repo/bin"
	"net/http"
)

type Characteristic struct {
	ID       int
	Name     string
	Type     string
	Relation int
	Products []Product
	IsExist  bool `db:"is_exist"`
}

func GetCharacteristics(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		var characteristics []Characteristic
		err := bin.DB.Select(&characteristics, "select * from characteristic")
		bin.GlobalCheck(err)
		for i := range characteristics {
			err = bin.DB.Select(&characteristics[i].Products, "select product.id, product.name, product.amount, product.price, product.discount, product.image_link from product inner join set on product.id = set.product_id where characteristic_id = $1", characteristics[i].ID)
			bin.GlobalCheck(err)
		}
		bin.FinalCheck(c, characteristics)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func GetCharacteristicByID(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		var characteristic Characteristic
		err := bin.DB.Get(&characteristic, "select * from characteristic where id = $1", id)
		bin.GlobalCheck(err)
		err = bin.DB.Select(&characteristic.Products, "select product.id, product.name, product.amount, product.price, product.discount, product.image_link from product inner join set on product.id = set.product_id where characteristic_id = $1", characteristic.ID)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, characteristic)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func PostCharacteristic(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		var characteristic Characteristic
		err := c.BindJSON(&characteristic)
		bin.GlobalCheck(err)
		_, err = bin.DB.Exec("insert into characteristic (name, type, relation) values ($1, $2, $3)", characteristic.Name, characteristic.Type, characteristic.Relation)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, characteristic)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func PutCharacteristic(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		var characteristic Characteristic
		err := c.BindJSON(&characteristic)
		bin.GlobalCheck(err)
		_, err = bin.DB.Exec("update characteristic set name = $2, type = $3, relation = $4, is_exist = $5 where id = $1", id, characteristic.Name, characteristic.Type, characteristic.Relation, characteristic.IsExist)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, characteristic)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func DeleteCharacteristic(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		var characteristic Characteristic
		err := bin.DB.Get(&characteristic, "select * from characteristic where id = $1", id)
		bin.GlobalCheck(err)
		if characteristic.Relation == 0 {
			_, err = bin.DB.Exec("delete from characteristic where relation = $1", characteristic.ID)
			bin.GlobalCheck(err)
		}
		_, err = bin.DB.Exec("delete from characteristic where id = $1", id)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, nil)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}
