package bin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var GlobalErr error

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GlobalCheck(err error) {
	if err != nil && GlobalErr == nil {
		GlobalErr = err
	}
}

func FinalCheck(c *gin.Context, result any) {
	if GlobalErr != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": false, "message": GlobalErr.Error()})
	} else {
		c.IndentedJSON(http.StatusOK, result)
	}
	GlobalErr = nil
}
