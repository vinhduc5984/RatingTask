package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"todo-list/configs"
	"todo-list/models"
	"todo-list/responses"

	common "todo-list/common"
	utils "todo-list/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var todoCollection *mongo.Collection = configs.GetCollection(configs.DB, "todos")
var validate = validator.New()

var redisTest  *redis.Client = configs.RDB

func CreateToDo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var todo models.ToDo
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&todo); err != nil {
			c.JSON(http.StatusBadRequest, responses.ToDoResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&todo); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.ToDoResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newToDo := models.ToDo{
			Id:          primitive.NewObjectID(),
			Title:       todo.Title,
			Description: todo.Description,
			StartDate:   todo.StartDate,
			EndDate:     todo.EndDate,
			Email: 		 todo.Email,
		}

		// check exactly email domain
		if ok, suggest := common.IsTypoDomain(newToDo.Email); ok {
				mess := fmt.Sprintf("Bạn có muốn dùng domain đúng là %s không?\n", suggest)
				c.JSON(http.StatusBadRequest, responses.ToDoResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": mess}})
				return
			} else {
				fmt.Println("Domain ok")
			}

		result, err := todoCollection.InsertOne(ctx, newToDo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ToDoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		// get COUNT_REG redis
		numCount := int32(0)
		val, err := redisTest.Get(ctx,"COUNT_REG").Result()
		if err != nil{
			fmt.Println("get redis fail: ",err.Error())
		}
		if  val != ""{
			numCount = utils.ConvertToInt32(val, 0)
		}

		// set number user register
		err1 := redisTest.Set(ctx, "COUNT_REG", numCount + 1, 0).Err()
		if err1!=nil{
			fmt.Println("set redis fail: ", err1.Error())
		}
		c.JSON(http.StatusCreated, responses.ToDoResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetOneToDo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		todoId := c.Param("todoId")
		var user models.ToDo
		defer cancel()

		// objId, _ := primitive.ObjectIDFromHex(todoId)

		err := todoCollection.FindOne(ctx, bson.M{"id": todoId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ToDoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.ToDoResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": user}})
	}
}

func EditToDo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		todoId := c.Param("todoId")
		var todo models.ToDo
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(todoId)

		//validate the request body
		if err := c.BindJSON(&todo); err != nil {
			c.JSON(http.StatusBadRequest, responses.ToDoResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&todo); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.ToDoResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{"title": todo.Title, "description": todo.Description, "startDate": todo.StartDate, "endDate": todo.EndDate}
		result, err := todoCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ToDoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		fmt.Println("data after update::: ", common.JsonPrettyAny(result))

		//get updated user details
		var updatedUser models.ToDo
		if result.MatchedCount == 1 {
			err := todoCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.ToDoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, responses.ToDoResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedUser}})
	}
}

func DeleteDoto() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		todoId := c.Param("todoId")
		defer cancel()

		// objId, _ := primitive.ObjectIDFromHex(todoId)

		result, err := todoCollection.DeleteOne(ctx, bson.M{"id": todoId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ToDoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				responses.ToDoResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "ToDo with specified ID not found!"}},
			)
			return
		}

		// // delete RDB with key
		// configs.DeleteRDBByKey(ctx,[]string{"COUNT_REG"})

		c.JSON(http.StatusOK,
			responses.ToDoResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "ToDo successfully deleted!"}},
		)
	}
}

func GetAllToDos() gin.HandlerFunc {

	

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var users []models.ToDo
		defer cancel()

		// get redis
		val, err := redisTest.Get(ctx, "COUNT_REG").Result()
		if err != nil{
			fmt.Println("get redis fail: ",err.Error())
			
		}
		fmt.Println("data count redis::: ",val)

		results, err := todoCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ToDoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleUser models.ToDo
			if err = results.Decode(&singleUser); err != nil {
				c.JSON(http.StatusInternalServerError, responses.ToDoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			users = append(users, singleUser)
		}

		c.JSON(http.StatusOK,
			responses.ToDoResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": users}},
		)
	}
}
