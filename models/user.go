package models

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/my/repo/bin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type User struct {
	ID         int
	Email      string
	Password   string
	Cost       string
	Phone      string
	LastName   string `db:"last_name"`
	Name       string
	MiddleName string `db:"middle_name"`
	Gender     string
	RoleID     int `db:"role_id"`
	Role       Role
	Card       Card
	IsExist    bool `db:"is_exist"`
}

func GetUsers(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		var users []User
		err := bin.DB.Select(&users, "select * from \"User\"")
		bin.GlobalCheck(err)
		for i := range users {
			err = bin.DB.Get(&users[i].Role, "select Role.Id, Role.Name from Role where id = $1", users[i].RoleID)
			bin.GlobalCheck(err)
			err = bin.DB.Get(&users[i].Card, "select * from discount_card where id = $1", users[i].ID)
			if err != nil {
				if err.Error() != "sql: no rows in result set" {
					bin.GlobalCheck(err)
				}
			}
			date, _ := time.Parse(time.RFC3339, users[i].Card.Date)
			users[i].Card.Date = time.Time.Format(date, time.DateTime)[0 : len(time.Time.Format(date, time.DateTime))-9]
		}
		bin.FinalCheck(c, users)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func GetUserByID(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		var user User
		err := bin.DB.Get(&user, "select * from \"User\" where id = $1", id)
		bin.GlobalCheck(err)
		err = bin.DB.Get(&user.Role, "select Role.Id, Role.Name from Role where id = $1", user.RoleID)
		bin.GlobalCheck(err)
		err = bin.DB.Get(&user.Card, "select * from discount_card where id = $1", user.ID)
		bin.GlobalCheck(err)
		date, _ := time.Parse(time.RFC3339, user.Card.Date)
		user.Card.Date = time.Time.Format(date, time.DateTime)[0 : len(time.Time.Format(date, time.DateTime))-9]
		bin.FinalCheck(c, user)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func PostUsers(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		var user User
		err := c.BindJSON(&user)
		bin.GlobalCheck(err)
		var hashedPassword []byte
		hashedPassword, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		bin.GlobalCheck(err)
		_, err = bin.DB.Exec("insert into \"User\" (email, password, phone, last_name, name, middle_name, gender, role_id) values ($1, $2, $3, $4, $5, $6, $7, $8)", user.Email, hashedPassword, user.Phone, user.LastName, user.Name, user.MiddleName, user.Gender, user.RoleID)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, user)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func PutUsers(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		var user User
		err := c.BindJSON(&user)
		bin.GlobalCheck(err)
		_, err = bin.DB.Exec("update \"User\" set email = $2, password = $3, phone = $4, last_name = $5, name = $6, middle_name = $7, gender = $8, role_id = $9, is_exist = $10 where id = $1", id, user.Email, user.Password, user.Phone, user.LastName, user.Name, user.MiddleName, user.Gender, user.RoleID, user.IsExist)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, user)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func DeleteUsers(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		var userData User
		err := bin.DB.Get(&userData, "select * from \"User\" where id = $1;", id)
		bin.GlobalCheck(err)
		_, err = bin.DB.Exec("delete from discount_card where id = $1", id)
		bin.GlobalCheck(err)
		//_, err = bin.DB.Exec("delete from cart_position where user_id = $1", id)
		//bin.GlobalCheck(err)
		_, err = bin.DB.Exec("delete from \"User\" where id = $1", id)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, nil)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func GetCurrentUser(c *gin.Context) {
	accessToken := c.GetHeader("Authorization")
	user := bin.Client.HGetAll(context.Background(), "sessions:"+accessToken).Val()
	if len(user) != 0 {
		var userData User
		err := bin.DB.Get(&userData, "select * from \"User\" where email = $1;", user["login"])
		bin.GlobalCheck(err)
		err = bin.DB.Get(&userData.Role, "select Role.Id, Role.Name from Role inner join \"User\" on Role.id = role_id where email = $1", userData.Email)
		bin.GlobalCheck(err)
		err = bin.DB.Get(&userData.Card, "select * from discount_card where id = $1", userData.ID)
		if err == nil {
			date, _ := time.Parse(time.RFC3339, userData.Card.Date)
			userData.Card.Date = time.Time.Format(date, time.DateTime)[0 : len(time.Time.Format(date, time.DateTime))-9]
		}
		bin.FinalCheck(c, userData)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}
