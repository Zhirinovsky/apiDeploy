package bin

import (
	_ "database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var DB *sqlx.DB
var Client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       1,
})

func ConnectDB() {
	data := GetConfigData()
	info := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", data["User"], data["Password"], data["Dbname"])
	var err error
	DB, err = sqlx.Connect("postgres", info)
	CheckErr(err)
}
