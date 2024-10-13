package handler

import (
	"Auth_service/api/email"
	"Auth_service/pkg/models"
	"context"
	"github.com/badoux/checkmail"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// RegisterStudent godoc
// @Summary Registers user
// @Description Registers a new user`
// @Tags auth
// @Param user body models.RegisterRequest1 true "User data"
// @Success 200 {object} models.RegisterResponse
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /user/register [post]
func (h *Handler) RegisterStudent(c *gin.Context) {
	h.log.Info("RegisterStudent handler called.")

	var req models.RegisterRequest1
	if err := c.ShouldBind(&req); err != nil {
		h.log.Error("Invalid data provided", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := checkmail.ValidateFormat(req.Email)
	if err != nil {
		h.log.Error("Invalid email provided", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email provided"})
		return
	}

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	code, err := email.Email(req.Email)
	if err != nil {
		h.log.Error("Invalid email provided", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email provided: " + err.Error()})
		return
	}
	req1 := models.RegisterRequest{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Username:    req.Username,
		Email:       req.Email,
		Password:    req.Password,
	}
	req1.Code = code

	err = h.redis.SetRegister(ctx, req1)
	if err != nil {
		h.log.Error("Failed to register user", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}
	h.log.Info("Successfully saved to redis")

	c.JSON(http.StatusOK, gin.H{"info": "code sent to this email " + req.Email})
}

// AcceptCodeToRegister godoc
// @Summary Accept code to register
// @Description it accepts code to register
// @Tags auth
// @Param token body models.AcceptCode true "enough"
// @Success 200 {object} models.RegisterResponse
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /user/accept-code [post]
func (h *Handler) AcceptCodeToRegister(c *gin.Context) {
	var req models.AcceptCode
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("Invalid data provided", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	register, err := h.redis.GetRegister(ctx, req.Email)
	if err != nil {
		h.log.Error("Failed to get register from redis", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get register from redis; " + err.Error()})
		return
	}

	if register.Code != req.Code {
		h.log.Error("Invalid code", "code", req.Code)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid code"})
		return
	}

	response, err := h.auth.Register(ctx, register)
	if err != nil {
		h.log.Error("Failed to register student", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register student; " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// RegisterAdmin godoc
// @Summary Registers user
// @Description Registers a new user`
// @Tags auth
// @Success 200 {object} models.Message
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /admin/register [post]
func (h *Handler) RegisterAdmin(c *gin.Context) {
	h.log.Info("RegisterStudent handler called.")

	err := h.auth.RegisterAdmin(c)
	if err != nil {
		h.log.Error("Error registering ADMIN", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("Successfully registered user")
	c.JSON(http.StatusOK, models.Message{Message: "FOR SURE!"})
}

// Login godoc
// @Summary Login a user
// @Description Login a user
// @Tags auth
// @Param user body models.LoginRequest true "User data"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /user/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBind(&req); err != nil {
		h.log.Error("Invalid data provided", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	res, err := h.auth.Login(ctx, req)
	if err != nil {
		h.log.Error("Error logging in user", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("Successfully logged in user")
	c.JSON(http.StatusOK, res)
}

// ForgotPassword godoc
// @Summary Forgot Password
// @Description it sends code to your email address
// @Tags auth
// @Param token body models.ForgotPasswordRequest true "enough"
// @Success 200 {object} string "message"
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /user/forgot-password [post]
func (h Handler) ForgotPassword(c *gin.Context) {
	h.log.Info("ForgotPassword is working")
	var req models.ForgotPasswordRequest
	if err := c.BindJSON(&req); err != nil {
		h.log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.auth.GetUserByEmail(c, req.Email)
	if err != nil {
		h.log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not registered"})
		return
	}

	code, err := email.Email(req.Email)
	if err != nil {
		h.log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending email " + err.Error()})
		return
	}
	err = h.redis.SetCode(c, req.Email, code)
	if err != nil {
		h.log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error storing codes in Redis " + err.Error()})
		return
	}
	h.log.Info("ForgotPassword succeeded")
	c.JSON(200, gin.H{"message": "Password reset code sent to your email"})
}

// ResetPassword godoc
// @Summary Reset Password
// @Description it Reset your Password
// @Tags auth
// @Param token body models.ResetPassReq true "enough"
// @Success 200 {object} string "message"
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /user/reset-password [post]
func (h *Handler) ResetPassword(c *gin.Context) {
	h.log.Info("ResetPassword is working")
	var req models.ResetPassReq
	if err := c.BindJSON(&req); err != nil {
		h.log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	code, err := h.redis.GetCodes(c, req.Email)
	if err != nil {
		h.log.Error(err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid or expired code " + err.Error()})
		return
	}
	if code != req.Code {
		h.log.Error("Invalid code")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid code "})
		return
	}
	res, err := h.auth.GetUserByEmail(c, req.Email)
	if err != nil {
		h.log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user" + err.Error()})
		return
	}

	err = h.auth.UpdatePassword(c, &models.UpdatePasswordReq{Id: res.Id, Password: req.Password})
	if err != nil {
		h.log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating password" + err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Password reset successfully"})
}
