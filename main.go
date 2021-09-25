package main

import (
	"bwastartupgolang/auth"
	"bwastartupgolang/handler"
	"bwastartupgolang/helper"
	"bwastartupgolang/user"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	authSerevice := auth.NewService()

	// userService.SaveAvatar(1, "images/1-profile.png")

	userHandler := handler.NewUserHandler(userService, authSerevice)

	router := gin.Default()
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddlerware(authSerevice, userService), userHandler.UploadAvatar)

	router.Run()

	// userInput := user.RegisterUserInput{}
	// userInput.Name = "Test simpan dari service"
	// userInput.Email = "Contoh@gmail.com"
	// userInput.Occupation = "Band"
	// userInput.Password = "password"

	// userService.RegisterUser(userInput)

	// user := user.User{
	// 	Name: "Test",
	// }

	// userRepository.Save(user)

	// fmt.Println("Connection to db is good")

	// var users []user.User

	// db.Find(&users)

	// router := gin.Default()
	// router.GET("/handler", handler)
	// router.Run()
}

// func handler(c *gin.Context) {
// 	dsn := "root:@tcp(127.0.0.1:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}

// 	var users []user.User
// 	db.Find(&users)

// 	c.JSON(http.StatusOK, users)
// }

//========how!
//input dari user
//handler = mapping input dari user -> struct input
//service = melakukan mapping dari struct input -> struct user
//repository
//db

//========middleware
//ambil nilai header authorization : Bearer tokentoken
//dari header authorization ambil nilai tokennya saja
//kita validasi token
//ambil user_id
//ambil user dari db berdasarkan user_id lewat service
//kita set context isinya user

func authMiddlerware(authService auth.Service, userService user.Service) gin.HandlerFunc {

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.ApiResponse("Unauthorize", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)

		if err != nil {
			response := helper.ApiResponse("Unauthorize", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.ApiResponse("Unauthorize", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.ApiResponse("Unauthorize", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}

}
