package bin

import (
	_ "database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"strconv"
)

var data = GetConfigData()
var DB *sqlx.DB
var Client = redis.NewClient(&redis.Options{})

func ConnectDB() {
	info := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", data["User"], data["Password"], data["Dbname"], data["Host"], data["Port"])
	db, err := strconv.Atoi(data["RedisDB"])
	CheckErr(err)
	DB, err = sqlx.Connect("postgres", info)
	CheckErr(err)
	var password string
	if data["RedisPassword"] == "-" {
		password = ""
	} else {
		password = data["RedisPassword"]
	}
	Client = redis.NewClient(&redis.Options{
		Addr:     data["RedisAddr"],
		Password: password,
		DB:       db,
	})
}
