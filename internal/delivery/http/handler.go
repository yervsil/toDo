package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/yervsil/toDo-microservice/internal/service"
	"github.com/yervsil/toDo-microservice/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
    "github.com/swaggo/files" // swagger embed files
    _ "github.com/yervsil/toDo-microservice/docs"
)

type Handler struct {
	service *service.Service
	logger *logger.Logger
}

func NewHandler(services *service.Service, logger *logger.Logger) *Handler{
	return &Handler{
		service: services,
		logger: logger,
	}
}


func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	api := router.Group("/api")
	{
		v1 := api.Group("/todo-list")
	{
		v1.POST("/tasks", h.createTask)
		v1.PUT("/tasks/:id", h.updateTask)
		v1.DELETE("/tasks/:id", h.deleteTask)
		v1.PATCH("/tasks/:id/done", h.statusUpdate)
		v1.GET("/tasks", h.getTasks)
	}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func parseIdFromPath(c *gin.Context, param string) (primitive.ObjectID, error) {
	idParam := c.Param(param)
	if idParam == "" {
		return primitive.ObjectID{}, errors.New("empty id param")
	}

	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return primitive.ObjectID{}, errors.New("invalid id param")
	}

	return id, nil
}
