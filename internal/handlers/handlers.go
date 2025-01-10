package handlers

import (
	"net/http"
	"wonder-interview/internal/config"
	"wonder-interview/internal/models"
	"wonder-interview/internal/utils"

	"github.com/gin-gonic/gin"
)

// generate swagger doc
// @Summary Login
// @Description Login with account and password
// @Produce json
// @Param user body models.User true "User account and password"
// @Success 200 {string} string "JWT token"
// @Failure 400 {object} utils.ErrorResponse "Invalid data"
// @Failure 403 {object} utils.ErrorResponse "Invalid password"
// @Failure 500 {object} utils.ErrorResponse "Failed to generate token"
// @Router /login [post]
func LoginHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		code := http.StatusBadRequest
		c.JSON(code, utils.NewErrorResponse(code, "Invalid data", nil))
		return
	}

	originPwd := user.Password
	mockUser := models.FindUserByAccount(user.Account)
	if mockUser == nil {
		code := http.StatusBadRequest
		c.JSON(code, utils.NewErrorResponse(code, "Invalid account", nil))
		return
	}
	if err := utils.CheckPassword(mockUser.PasswordHash, originPwd); err != nil {
		code := http.StatusForbidden
		c.JSON(code, utils.NewErrorResponse(code, "Invalid password", nil))
		return
	}

	if token, err := utils.GenerateJWT(mockUser.ID, config.SECRET_KEY); err != nil {
		code := http.StatusInternalServerError
		c.JSON(code, utils.NewErrorResponse(code, "Failed to generate token", nil))
		return
	} else {
		c.JSON(http.StatusOK, token)
	}
}
