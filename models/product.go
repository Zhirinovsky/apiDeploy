package models

import (
	"github.com/gin-gonic/gin"
	"github.com/my/repo/bin"
	"net/http"
)

type Product struct {
	ID         int
	Name       string
	Price      float64
	Amount     int
	Discount   int
	ImageLink  string `db:"image_link"`
	CategoryID string `db:"category_id"`
	Category   Category
	Sets       []Set
	IsExist    bool `db:"is_exist"`
}

func GetProducts(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		var products []Product
		err := bin.DB.Select(&products, "select * from Product order by Name ASC")
		bin.GlobalCheck(err)
		for i, product := range products {
			err = bin.DB.Select(&products[i].Sets, "select set.id, set.value, product_id, characteristic_id from set where product_id = $1", products[i].ID)
			bin.GlobalCheck(err)
			for j := range products[i].Sets {
				err = bin.DB.Get(&products[i].Sets[j].Characteristic, "select characteristic.id, characteristic.name, characteristic.Type, characteristic.Relation from characteristic where id = $1", products[i].Sets[j].CharacteristicID)
				bin.GlobalCheck(err)
				products[i].Sets[j].Product = product
			}
			err = bin.DB.Get(&products[i].Category, "select * from category where id = $1", products[i].CategoryID)
			bin.GlobalCheck(err)
		}
		bin.FinalCheck(c, products)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func GetProductByID(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		var product Product
		err := bin.DB.Get(&product, "select * from Product where id = $1", id)
		bin.GlobalCheck(err)
		err = bin.DB.Select(&product.Sets, "select set.id, set.value, product_id, characteristic_id from set where product_id = $1", product.ID)
		bin.GlobalCheck(err)
		for j := range product.Sets {
			err = bin.DB.Get(&product.Sets[j].Characteristic, "select characteristic.id, characteristic.name, characteristic.Type, characteristic.Relation from characteristic where id = $1", product.Sets[j].CharacteristicID)
			bin.GlobalCheck(err)
			//product.Sets[j].Product = product
		}
		err = bin.DB.Get(&product.Category, "select * from category where id = $1", product.CategoryID)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, product)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func PostProduct(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		var product Product
		err := c.BindJSON(&product)
		bin.GlobalCheck(err)
		_, err = bin.DB.Exec("insert into product (name, price, amount, discount, image_link, category_id) values ($1, $2, $3, $4, $5, $6)", product.Name, product.Price, product.Amount, product.Discount, product.ImageLink, product.CategoryID)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, product)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func PutProduct(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		var product Product
		err := c.BindJSON(&product)
		bin.GlobalCheck(err)
		_, err = bin.DB.Exec("update product set name = $2, price = $3, amount = $4, discount = $5, image_link = $6, category_id = $7, is_exist = $8 where id = $1", id, product.Name, product.Price, product.Amount, product.Discount, product.ImageLink, product.CategoryID, product.IsExist)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, product)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func DeleteProduct(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		_, err := bin.DB.Exec("delete from product where id = $1", id)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, nil)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}
