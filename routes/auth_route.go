package routes

import (
	"todo-list/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoute(router *gin.Engine) {

	api := router.Group("/auth")
	{
		api.POST("/register", controllers.RegisterAccount())
		api.POST("/login", controllers.Login())
		api.POST("/change-password", controllers.EditToDo())

	}

}
