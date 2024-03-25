package models

import (
	"github.com/gin-gonic/gin"
	"github.com/my/repo/bin"
	"net/http"
)

type Role struct {
	ID      int
	Name    string
	IsExist bool `db:"is_exist"`
}

func GetRoles(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		var roles []Role
		err := bin.DB.Select(&roles, "select * from role")
		bin.GlobalCheck(err)
		bin.FinalCheck(c, roles)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func GetRoleByID(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		var role Role
		err := bin.DB.Get(&role, "select * from role where id = $1", id)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, role)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func PostRole(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		var role Role
		err := c.BindJSON(&role)
		bin.GlobalCheck(err)
		_, err = bin.DB.Exec("insert into role (name) values ($1)", role.Name)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, role)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func PutRole(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		var role Role
		err := c.BindJSON(&role)
		bin.GlobalCheck(err)
		_, err = bin.DB.Exec("update role set name = $2, is_exist = $3 where id = $1", id, role.Name, role.IsExist)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, role)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func DeleteRole(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		_, err := bin.DB.Exec("delete from role where id = $1", id)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, nil)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}
