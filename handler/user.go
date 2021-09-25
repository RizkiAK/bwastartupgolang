package handler

import (
	"bwastartupgolang/auth"
	"bwastartupgolang/helper"
	"bwastartupgolang/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService  user.Service
	authSerevice auth.Service
}

func NewUserHandler(userService user.Service, authSerevice auth.Service) *userHandler {
	return &userHandler{userService, authSerevice}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	//tangkap inputan user
	//map input dari user ke struct RegisterUser
	//struct di atas passing sebagai  parameter service

	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.ApiResponse("Register account failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		response := helper.ApiResponse("Register account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	//
	token, err := h.authSerevice.GenerateToken(newUser.ID)
	if err != nil {
		response := helper.ApiResponse("Register account failed", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	//

	formatter := user.FormatUser(newUser, token)

	response := helper.ApiResponse("Account has been registered", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	//user memasukkan input (email & password)
	//input ditangkap handler
	//mapping dari input user ke input struct
	//input struct parsing ke service
	//service mencari dengan bantuan repository user dengan emailx
	//mencocokan password

	var input user.LoginInput

	err := c.ShouldBind(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.ApiResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedinUser, err := h.userService.Login(input)

	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ApiResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	token, err := h.authSerevice.GenerateToken(loggedinUser.ID)
	if err != nil {
		response := helper.ApiResponse("Login failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(loggedinUser, token)

	response := helper.ApiResponse("Succesfully loggedin", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) CheckEmailAvailability(c *gin.Context) {
	//ada input email dari user
	//input email di mapping ke struct input
	//struct input di parsing ke service
	//service akan memanggil repository - email sudah ada atau belum
	//repository - db

	var input user.CheckEmailInput
	err := c.ShouldBind(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.ApiResponse("Email checking field", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"errors": "Server Error"}
		response := helper.ApiResponse("Email checking field", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H{
		"is_available": isEmailAvailable,
	}

	metaMessage := "Email has been registered"

	if isEmailAvailable {
		metaMessage = "Email is available"
	}

	response := helper.ApiResponse(metaMessage, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	// c.SaveUploadedFile(file, )
	//input dari user
	//simpan gambar di folder "/images"
	//di service panggil repo
	//JWT (sementara hardcode, seakan2 user yg login ID = 1)
	//repo ambil data user yang ID = 1
	//repo update data user simpan lokasi file

	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.ApiResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	//harusnya dari JWT, tapi sekarang manual
	// userID := 1

	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	//format lama
	//images/namafile.png

	//format baru
	//images/ID(1-)-namafile.png

	// path := "images/" + file.Filename
	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.ApiResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.ApiResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.ApiResponse("Avatar successfully uploaded", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
