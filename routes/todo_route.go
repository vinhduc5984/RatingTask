package routes

import (
	"todo-list/controllers"

	"github.com/gin-gonic/gin"
)

func ToDoRoute(router *gin.Engine) {
	router.POST("/todo", controllers.CreateToDo())
	router.GET("/get-one-ToDo/:todoId", controllers.GetOneToDo())
	router.PUT("/update-ToDo/:todoId", controllers.EditToDo())
	router.DELETE("/delete/:todoId", controllers.DeleteDoto())
	router.GET("/get-all-todo", controllers.GetAllToDos())
	router.POST("/upload", controllers.UploadFile())
}
