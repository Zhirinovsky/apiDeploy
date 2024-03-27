package bin

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthForm struct {
	Email    string
	Password string
	Role     string
}

var InvalidTokenMessage = "Токен недействителен"

func Register(c *gin.Context) {
	var user AuthForm
	err := c.BindJSON(&user)
	GlobalCheck(err)
	var hashedPassword []byte
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	CheckErr(err)
	var roleId int
	err = DB.QueryRow("select id from role where name = $1", "Клиент").Scan(&roleId)
	CheckErr(err)
	_, err = DB.Exec("insert into \"User\" (email, password, role_id) values ($1, $2, $3)", user.Email, hashedPassword, roleId)
	GlobalCheck(err)
	if GlobalErr != nil {
		GlobalErr = nil
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Данный почтовый адрес уже зарегестрирован"})
	} else {
		c.IndentedJSON(http.StatusCreated, user)
	}
}

func Login(c *gin.Context) {
	var user AuthForm
	err := c.BindJSON(&user)
	GlobalCheck(err)
	var users []AuthForm
	err = DB.Select(&users, "select email, password from \"User\"")
	GlobalCheck(err)
	check := false
	for _, userDB := range users {
		passwordCheck := bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(user.Password))
		if userDB.Email == user.Email && passwordCheck == nil {
			check = true
			if GlobalErr != nil {
				message := GlobalErr.Error()
				GlobalErr = nil
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": message})
			} else {
				newToken := jwt.New(jwt.SigningMethodHS256)
				randText := make([]byte, 10)
				_, err = rand.Read(randText)
				GlobalCheck(err)
				tokenString, err := newToken.SignedString([]byte(GetConfigData()["Signing Key"] + base64.StdEncoding.EncodeToString(randText)))
				CheckErr(err)
				err = Client.HSet(context.Background(), "sessions:"+tokenString, "login", user.Email).Err()
				CheckErr(err)
				Client.Expire(context.Background(), "sessions:"+tokenString, 2*time.Hour)
				err = DB.Get(&user.Role, "select Role.Name from Role inner join \"User\" on Role.id = role_id where email = $1", user.Email)
				CheckErr(err)
				//c.SetCookie("access_token", access_token, 2*60*60, "/", "localhost", false, true)
				c.JSON(http.StatusOK, gin.H{"status": true, "token": tokenString, "role": user.Role})
			}
		}
	}
	if !check {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Неверные данные авторизации"})
	}
}

func Refresh(c *gin.Context) {
	accessToken := c.GetHeader("Authorization")
	if ValidateToken(accessToken) {
		user := Client.HGetAll(context.Background(), "sessions:"+accessToken).Val()
		Client.Del(context.Background(), "sessions:"+accessToken)
		newToken := jwt.New(jwt.SigningMethodHS256)
		tokenString, err := newToken.SignedString([]byte(GetConfigData()["Signing Key"]))
		CheckErr(err)
		Client.HSet(context.Background(), "sessions:"+tokenString, "login", user["login"])
		Client.Expire(context.Background(), "sessions:"+tokenString, 2*time.Hour)
		//c.SetCookie("access_token", accessToken, 2*60*60, "/", "localhost", false, true)
		c.IndentedJSON(http.StatusOK, gin.H{"status": true, "token": tokenString})
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": InvalidTokenMessage})
	}
}

func Logout(c *gin.Context) {
	accessToken := c.GetHeader("Authorization")
	if ValidateToken(accessToken) {
		Client.Del(context.Background(), "sessions:"+accessToken)
		//c.SetCookie("access_token", "", -1, "/", "localhost", false, true)
		c.IndentedJSON(http.StatusOK, gin.H{"status": true})
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": InvalidTokenMessage})
	}
}

func ValidateToken(accessToken string) bool {
	user := Client.HGetAll(context.Background(), "sessions:"+accessToken).Val()
	if len(user) == 0 {
		return false
	} else {
		return true
	}
}

func CheckToken(c *gin.Context) {
	accessToken := c.GetHeader("Authorization")
	if ValidateToken(accessToken) {
		c.IndentedJSON(http.StatusOK, gin.H{"status": true})
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": InvalidTokenMessage})
	}
}
