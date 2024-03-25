package models

import (
	"github.com/gin-gonic/gin"
	"github.com/my/repo/bin"
	"net/http"
)

type Category struct {
	ID       int
	Name     string
	Relation int
	IsExist  bool `db:"is_exist"`
}

func GetCategories(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		var categories []Category
		err := bin.DB.Select(&categories, "select * from category")
		bin.GlobalCheck(err)
		bin.FinalCheck(c, categories)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func GetCategoryByID(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		var category Category
		err := bin.DB.Get(&category, "select * from Category where id = $1", id)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, category)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func PostCategory(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		var category Category
		err := c.BindJSON(&category)
		bin.GlobalCheck(err)
		_, err = bin.DB.Exec("insert into category (name, relation) values ($1, $2)", category.Name, category.Relation)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, category)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func PutCategory(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		var category Category
		err := c.BindJSON(&category)
		bin.GlobalCheck(err)
		_, err = bin.DB.Exec("update category set name = $2, relation = $3, is_exist = $4 where id = $1", id, category.Name, category.Relation, category.IsExist)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, category)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func DeleteCategory(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		var category Category
		err := bin.DB.Get(&category, "select * from category where id = $1", id)
		bin.GlobalCheck(err)
		if category.Relation == 0 {
			_, err = bin.DB.Exec("delete from category where relation = $1", category.ID)
			bin.GlobalCheck(err)
		}
		_, err = bin.DB.Exec("delete from category where id = $1", id)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, nil)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}
