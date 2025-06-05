package routes

import (
	"github.com/gin-gonic/gin"
)

func GetAllRoute(router *gin.Engine) {
	ToDoRoute(router)
	AuthRoute(router)
}
