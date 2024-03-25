package models

import (
	"github.com/gin-gonic/gin"
	"github.com/my/repo/bin"
	"net/http"
)

type Status struct {
	ID      int
	Status  string
	IsExist bool `db:"is_exist"`
}

func GetStatuses(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		var statuses []Status
		err := bin.DB.Select(&statuses, "select * from status")
		bin.GlobalCheck(err)
		bin.FinalCheck(c, statuses)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func GetStatusByID(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		var status Status
		err := bin.DB.Get(&status, "select * from status where id = $1", id)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, status)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func PostStatus(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		var status Status
		err := c.BindJSON(&status)
		bin.GlobalCheck(err)
		_, err = bin.DB.Exec("insert into status (status) values ($1)", status.Status)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, status)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func PutStatus(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		var status Status
		err := c.BindJSON(&status)
		bin.GlobalCheck(err)
		_, err = bin.DB.Exec("update status set status = $2, is_exist = $3 where id = $1", id, status.Status, status.IsExist)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, status)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func DeleteStatus(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		_, err := bin.DB.Exec("delete from status where id = $1", id)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, nil)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}
