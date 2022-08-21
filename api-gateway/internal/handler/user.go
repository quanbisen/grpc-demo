package handler

import (
	"api-gateway/common/res"
	"api-gateway/common/util"
	"api-gateway/internal/service"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserRegister(c *gin.Context) {
	var userReq service.UserRequest
	PanicIfUserError(c.Bind(&userReq))
	userService := c.Keys["user"].(service.UserServiceClient)
	userResp, err := userService.UserRegister(context.Background(), &userReq)
	PanicIfUserError(err)
	c.JSON(http.StatusOK, userResp)
}

func UserLogin(c *gin.Context) {
	var userReq service.UserRequest
	PanicIfUserError(c.Bind(&userReq))
	userService := c.Keys["user"].(service.UserServiceClient)
	userResp, err := userService.UserLogin(context.Background(), &userReq)
	PanicIfUserError(err)
	token, err := util.GenerateToken(uint(userResp.UserDetail.UserId))
	PanicIfUserError(err)
	response := res.TokenData{
		Token: token,
		User:  userResp.UserDetail,
	}
	c.JSON(http.StatusOK, response)
}
