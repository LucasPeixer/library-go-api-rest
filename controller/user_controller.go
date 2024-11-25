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
	GetUserLoans(c *gin.Context)
	ToggleUser(action string) gin.HandlerFunc
	DeleteUser(c *gin.Context)
	GetUserReservations(c *gin.Context)
	GetLoggedUserReservations(c *gin.Context)
	CancelUserReservation(c *gin.Context)
	CancelLoggedUserReservation(c *gin.Context)
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
		Cpf      string `json:"cpf" binding:"required"`
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

	if !utils.IsValidCPF(i.Cpf) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cpf input"})
		return
	}

	// Tenta criar um hash da senha do usuário
	hashedPassword, err := utils.HashPassword(i.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	userId, err := uc.useCase.Register(i.Name, i.Cpf, i.Phone, i.Email, hashedPassword, i.RoleId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user_id": userId})
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

func (uc *userController) GetUserLoans(c *gin.Context) {
	userIDValue, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	userID, err := strconv.Atoi(userIDValue.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	loans, err := uc.useCase.GetUserLoans(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, loans)
}

func (uc *userController) ToggleUser(action string) gin.HandlerFunc {
	if action != "activate" && action != "deactivate" {
		panic("Invalid action. Must be either 'activate' or 'deactivate'")
	}

	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
			return
		}

		var userErr error
		if action == "activate" {
			userErr = uc.useCase.ActivateUser(id)

		} else {
			userErr = uc.useCase.DeactivateUser(id)
		}

		if userErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": userErr.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User has been successfully " + action + "d"})
	}
}

func (uc *userController) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	err = uc.useCase.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User has been successfully deleted"})
}

func (uc *userController) GetUserReservations(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user Id"})
		return
	}

	reservations, err := uc.useCase.GetUserReservations(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reservations)
}

func (uc *userController) GetLoggedUserReservations(c *gin.Context) {
	userIdStr, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized user"})
		return
	}

	userId, err := strconv.Atoi(userIdStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user Id"})
		return
	}

	reservations, err := uc.useCase.GetUserReservations(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reservations)
}

func (uc *userController) CancelUserReservation(c *gin.Context) {
	adminIdStr, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized user"})
		return
	}

	adminId, err := strconv.Atoi(adminIdStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid admin user Id"})
		return
	}

	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user Id"})
		return
	}

	reservationId, err := strconv.Atoi(c.Param("reservation-id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reservation Id"})
		return
	}

	if err := uc.useCase.CancelUserReservation(userId, reservationId, &adminId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User reservation has been successfully canceled"})
}

func (uc *userController) CancelLoggedUserReservation(c *gin.Context) {
	userIdStr, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized user"})
		return
	}

	userId, err := strconv.Atoi(userIdStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user Id"})
		return
	}

	reservationId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reservation Id"})
		return
	}

	if err := uc.useCase.CancelUserReservation(userId, reservationId, nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User reservation has been successfully canceled"})
}
