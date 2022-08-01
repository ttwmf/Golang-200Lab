package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
)

type Restaurant struct {
	Id      int    `json:"id" gorm:"id"`
	Name    string `json:"name" gorm:"name"`
	Address string `json:"address" gorm:"address"`
}

func (Restaurant) TableName() string {
	return "restaurants"
}

type RestaurantUpdate struct {
	Name    *string `json:"name" gorm:"name"`
	Address *string `json:"address" gorm:"address"`
}

func (RestaurantUpdate) TableName() string {
	return Restaurant{}.TableName()
}

type RestaurantCreate struct {
	Id      int    `json:"id" gorm:"id"`
	Name    string `json:"name" gorm:"name"`
	Address string `json:"address" gorm:"address"`
}

func (RestaurantCreate) TableName() string {
	return Restaurant{}.TableName()
}
func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/golang_200lab?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln("Cannot connect to MySQL:", err)
	}

	log.Println("Connected:", db)

	router := gin.Default()
	v1 := router.Group("/v1")
	{
		restaurants := v1.Group("/restaurants")
		{
			restaurants.POST("", CreateRestaurant(db))
		}
	}
	router.Run(":3098")
}

func (res *RestaurantCreate) Validate() error {
	res.Name = strings.TrimSpace(res.Name)

	if len(res.Name) == 0 {
		return errors.New("Restaurant name can't be blank")
	}
	return nil
}

func CreateRestaurant(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data RestaurantCreate

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		if err := data.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		if err := db.Create(&data).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": data.Id})
	}
}
