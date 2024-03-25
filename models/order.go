package models

import (
	"github.com/gin-gonic/gin"
	"github.com/my/repo/bin"
	"net/http"
	"time"
)

type Order struct {
	ID        int
	Date      string
	Address   string
	StatusID  int `db:"status_id"`
	UserID    int `db:"user_id"`
	Status    Status
	User      User
	Positions []Position
}

func GetOrders(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		var orders []Order
		err := bin.DB.Select(&orders, "select * from \"Order\"")
		bin.GlobalCheck(err)
		for i := range orders {
			err = bin.DB.Get(&orders[i].Status, "select * from status where id = $1", orders[i].StatusID)
			bin.GlobalCheck(err)
			err = bin.DB.Get(&orders[i].User, "select * from \"User\" where id = $1", orders[i].UserID)
			bin.GlobalCheck(err)
			err = bin.DB.Select(&orders[i].Positions, "select * from Order_Position where order_id = $1", orders[i].ID)
			bin.GlobalCheck(err)
			for j := range orders[i].Positions {
				err = bin.DB.Get(&orders[i].Positions[j].Order, "select * from \"Order\" where id = $1", orders[i].Positions[j].OrderID)
				bin.GlobalCheck(err)
				err = bin.DB.Get(&orders[i].Positions[j].Product, "select * from product where id = $1", orders[i].Positions[j].ProductID)
				bin.GlobalCheck(err)
				err = bin.DB.Get(&orders[i].Positions[j].Product.Category, "select * from category where id = $1", orders[i].Positions[j].Product.CategoryID)
				bin.GlobalCheck(err)
			}
			date, _ := time.Parse(time.RFC3339, orders[i].Date)
			orders[i].Date = time.Time.Format(date, time.DateTime)[0 : len(time.Time.Format(date, time.DateTime))-3]
		}
		bin.FinalCheck(c, orders)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func GetOrderByID(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		var order Order
		err := bin.DB.Get(&order, "select * from \"Order\" where id = $1", id)
		bin.GlobalCheck(err)
		err = bin.DB.Get(&order.Status, "select * from status where id = $1", order.StatusID)
		bin.GlobalCheck(err)
		err = bin.DB.Get(&order.User, "select * from \"User\" where id = $1", order.UserID)
		bin.GlobalCheck(err)
		err = bin.DB.Select(&order.Positions, "select * from Order_Position where order_id = $1", order.ID)
		bin.GlobalCheck(err)
		for j := range order.Positions {
			err = bin.DB.Get(&order.Positions[j].Order, "select * from \"Order\" where id = $1", order.Positions[j].OrderID)
			bin.GlobalCheck(err)
			err = bin.DB.Get(&order.Positions[j].Product, "select * from product where id = $1", order.Positions[j].ProductID)
			bin.GlobalCheck(err)
			err = bin.DB.Get(&order.Positions[j].Product.Category, "select * from category where id = $1", order.Positions[j].Product.CategoryID)
			bin.GlobalCheck(err)
		}
		date, _ := time.Parse(time.RFC3339, order.Date)
		order.Date = time.Time.Format(date, time.DateTime)[0 : len(time.Time.Format(date, time.DateTime))-3]
		bin.FinalCheck(c, order)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func PostOrder(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		var order Order
		err := c.BindJSON(&order)
		bin.GlobalCheck(err)
		_, err = bin.DB.Exec("insert into \"Order\" (date, address, status_id, user_id) values ($1, $2, $3, $4)", order.Date, order.Address, order.StatusID, order.UserID)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, order)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func PutOrder(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		var order Order
		err := c.BindJSON(&order)
		bin.GlobalCheck(err)
		_, err = bin.DB.Exec("update \"Order\" set date = $2, address = $3, status_id = $4, user_id = $5 where id = $1", id, order.Date, order.Address, order.StatusID, order.UserID)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, order)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}

func DeleteOrder(c *gin.Context) {
	if bin.ValidateToken(c.GetHeader("Authorization")) {
		id := c.Param("id")
		var order Order
		err := bin.DB.Get(&order, "select * from \"Order\" where id = $1", id)
		bin.GlobalCheck(err)
		_, err = bin.DB.Exec("delete from order_position where order_id = $1", order.ID)
		bin.GlobalCheck(err)
		_, err = bin.DB.Exec("delete from \"Order\" where id = $1", id)
		bin.GlobalCheck(err)
		bin.FinalCheck(c, nil)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": bin.InvalidTokenMessage})
	}
}
