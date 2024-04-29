package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/my/repo/bin"
	"github.com/my/repo/models"
	"os"
)

func main() {
	bin.ConnectDB()
	//f, _ := os.Create("info.log")cart
	//gin.DisableConsoleColor()
	//gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	router := gin.Default()
	//Авторизация
	router.POST("/registration", bin.Register)
	router.POST("/login", bin.Login)
	router.GET("/refresh", bin.Refresh)
	router.GET("/logout", bin.Logout)
	router.GET("/check", bin.CheckToken)
	//Пользователи
	router.GET("/users", models.GetUsers)
	router.GET("/users/:id", models.GetUserByID)
	router.POST("/users", models.PostUsers)
	router.PUT("/users/:id", models.PutUsers)
	router.DELETE("/users/:id", models.DeleteUsers)
	router.GET("/getcurrentuser", models.GetCurrentUser)
	//Роли
	router.GET("/roles", models.GetRoles)
	router.GET("/roles/:id", models.GetRoleByID)
	router.POST("/roles", models.PostRole)
	router.PUT("/roles/:id", models.PutRole)
	router.DELETE("/roles/:id", models.DeleteRole)
	//Скидочные карты
	router.GET("/cards", models.GetCards)
	router.GET("/cards/:id", models.GetCardByID)
	router.POST("/cards", models.PostCards)
	router.PUT("/cards/:id", models.PutCards)
	router.DELETE("/cards/:id", models.DeleteCards)
	//Товары
	router.GET("/products", models.GetProducts)
	router.GET("/products/:id", models.GetProductByID)
	router.POST("/products", models.PostProduct)
	router.PUT("/products/:id", models.PutProduct)
	router.DELETE("/products/:id", models.DeleteProduct)
	//Характеристики товаров
	router.GET("/sets", models.GetSets)
	router.GET("/sets/:id", models.GetSetByID)
	router.POST("/sets", models.PostSet)
	router.PUT("/sets/:id", models.PutSet)
	router.DELETE("/sets/:id", models.DeleteSet)
	//Характеристики
	router.GET("/characteristics", models.GetCharacteristics)
	router.GET("/characteristics/:id", models.GetCharacteristicByID)
	router.POST("/characteristics", models.PostCharacteristic)
	router.PUT("/characteristics/:id", models.PutCharacteristic)
	router.DELETE("/characteristics/:id", models.DeleteCharacteristic)
	//Категории
	router.GET("/categories", models.GetCategories)
	router.GET("/categories/:id", models.GetCategoryByID)
	router.POST("/categories", models.PostCategory)
	router.PUT("/categories/:id", models.PutCategory)
	router.DELETE("/categories/:id", models.DeleteCategory)
	//Статусы
	router.GET("/statuses", models.GetStatuses)
	router.GET("/statuses/:id", models.GetStatusByID)
	router.POST("/statuses", models.PostStatus)
	router.PUT("/statuses/:id", models.PutStatus)
	router.DELETE("/statuses/:id", models.DeleteStatus)
	//Заказы
	router.GET("/orders", models.GetOrders)
	router.GET("/orders/:id", models.GetOrderByID)
	router.POST("/orders", models.PostOrder)
	router.PUT("/orders/:id", models.PutOrder)
	router.DELETE("/orders/:id", models.DeleteOrder)
	//Позиции заказов
	router.GET("/positions", models.GetOrderPositions)
	router.GET("/positions/:id", models.GetOrderPositionByID)
	router.POST("/positions", models.PostOrderPosition)
	router.PUT("/positions/:id", models.PutOrderPosition)
	router.DELETE("/positions/:id", models.DeleteOrderPosition)
	//Позиции корзин
	router.GET("/carts", models.GetCartPositions)
	router.GET("/carts/:id", models.GetCartPositionByID)
	router.POST("/carts", models.PostCartPosition)
	router.PUT("/carts/:id", models.PutCartPosition)
	router.DELETE("/carts/:id", models.DeleteCartPosition)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8079"
	}
	err := router.Run(":" + port)
	bin.CheckErr(err)
}
