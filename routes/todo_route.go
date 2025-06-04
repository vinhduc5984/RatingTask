package routes

import (
	"todo-list/controllers"

	"github.com/gin-gonic/gin"
)

func ToDoRoute(router *gin.Engine) {

	api := router.Group("/todo")
	{
		api.POST("/create", controllers.CreateToDo())
		api.GET("/get-one/:todoId", controllers.GetOneToDo())
		api.PUT("/update/:todoId", controllers.EditToDo())
		api.DELETE("/delete/:todoId", controllers.DeleteDoto())
		api.GET("/get-all", controllers.GetAllToDos())
		api.POST("/upload", controllers.UploadFile())
	}
}
