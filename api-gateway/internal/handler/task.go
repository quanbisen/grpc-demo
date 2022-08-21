package handler

import (
	"api-gateway/internal/service"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateTask(c *gin.Context) {
	var taskReq service.TaskRequest
	PanicIfTaskError(c.Bind(&taskReq))
	taskService := c.Keys["task"].(service.TaskServiceClient)
	userResp, err := taskService.TaskCreate(context.Background(), &taskReq)
	PanicIfTaskError(err)
	c.JSON(http.StatusOK, userResp)
}
