package controller

import (
	"go-api/usecase"
	"go-api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetUsersByFilters(c *gin.Context)
	GetUserById(c *gin.Context)
}
type userController struct {
	useCase usecase.UserUseCase
}

func NewUserController(useCase usecase.UserUseCase) UserController {
	return &userController{useCase: useCase}
}

// Register recebe um input JSON através do gin.Context e tenta registrar o usuário.
func (uc *userController) Register(c *gin.Context) {
	var i struct {
		Name     string `json:"name" binding:"required"`
		Phone    string `json:"phone" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		RoleId   int    `json:"role_id" binding:"required"`
	}
	// Valida o input de dados
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid registration input"})
		return
	}
	// Tenta criar um hash da senha do usuário
	hashedPassword, err := utils.HashPassword(i.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	if err := uc.useCase.Register(i.Name, i.Phone, i.Email, hashedPassword, i.RoleId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login recebe um input JSON através do gin.Context e tenta realizar o login do usuário.
func (uc *userController) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	// Valida o input de dados
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login input"})
		return
	}
	// Tenta logar o usuário
	token, err := uc.useCase.Login(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, token)
}

func (uc *userController) GetUsersByFilters(c *gin.Context) {
	name := c.Query("name")
	email := c.Query("email")

	userAccountList, err := uc.useCase.GetUsersByFilters(name, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userAccountList)
}

func (uc *userController) GetUserById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	userAccount, err := uc.useCase.GetUserById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userAccount)
}
