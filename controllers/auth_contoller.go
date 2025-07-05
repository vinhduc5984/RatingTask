package controllers

import (
	"context"
	"fmt"
	"os"

	"net/http"
	"time"
	"todo-list/configs"
	"todo-list/controllers/helper"
	"todo-list/models"
	"todo-list/responses"

	utils "todo-list/utils"

	"crypto/subtle"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "user")
var accountCollection *mongo.Collection = configs.GetCollection(configs.DB, "account")
var validateUser = validator.New()
var validateAccount = validator.New()

func RegisterAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var account models.Account
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&account); err != nil {
			c.JSON(http.StatusBadRequest, responses.AuthResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := models.ValidateAccount(account); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.AuthResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr}})
			return
		}

		// 1. Check Username duy nhất
		var existing models.Account
		err := accountCollection.FindOne(ctx, bson.M{"username": account.Username}).Decode(&existing)
		if err == nil {
			c.JSON(http.StatusConflict, responses.AuthResponse{Status: 409, Message: "Username đã tồn tại", Data: nil})
			return
		}
		if err != mongo.ErrNoDocuments {
			c.JSON(http.StatusInternalServerError, responses.AuthResponse{Status: 500, Message: "DB Error", Data: map[string]interface{}{"error": err.Error()}})
			return
		}

		// encoding password (SHA1)
		// get config key from .env
		// get PORT config
		pKey := os.Getenv("ENDCODE_PASSWORD_KEY")
		encodePassword := utils.EncodeSHA1Password(account.Password, pKey)

		newAccount := models.Account{
			Id:           primitive.NewObjectID(),
			Username:     account.Username,
			Password:     encodePassword,
			Role:         account.Role,
			Token:        "",
			RefreshToken: "",
			ExpiryDate:   utils.NULL_DATE,
			Disabled:     0,
			DeletedBy:    primitive.NilObjectID,
			DeletedAt:    utils.NULL_DATE,
			UpdatedBy:    primitive.NilObjectID,
			UpdatedAt:    utils.NULL_DATE,
			CreatedBy:    primitive.NilObjectID,
			CreatedAt:    utils.NULL_DATE,
		}

		result, err := accountCollection.InsertOne(ctx, newAccount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.AuthResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		// insert default user when create account
		objId, ok := result.InsertedID.(primitive.ObjectID)
		if !ok {
			c.JSON(http.StatusConflict, responses.AuthResponse{Status: 409, Message: "Lỗi", Data: nil})
			return
		}

		if result.InsertedID != primitive.NilObjectID && result.InsertedID != "" {
			var accountDB models.Account
			err := accountCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&accountDB)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.ToDoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			newUser := models.User{
				Id:        primitive.NewObjectID(),
				AccountId: accountDB.Id,
				Role:      accountDB.Role,
				DeletedBy: accountDB.DeletedBy,
				DeletedAt: accountDB.DeletedAt,
				UpdatedBy: accountDB.UpdatedBy,
				UpdatedAt: accountDB.UpdatedAt,
				CreatedBy: accountDB.CreatedBy,
				CreatedAt: accountDB.CreatedAt,
			}

			_, err1 := userCollection.InsertOne(ctx, newUser)
			if err1 != nil {
				c.JSON(http.StatusConflict, responses.AuthResponse{Status: 409, Message: "Lỗi tạo user", Data: nil})
				return
			}
		}

		c.JSON(http.StatusCreated, responses.AuthResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var account models.AccountLogin
		var accountDB models.Account
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&account); err != nil {
			c.JSON(http.StatusBadRequest, responses.AuthResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validateAccount.Struct(account); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.AuthResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		userId, _ := primitive.ObjectIDFromHex("60d5f483b3c2c66b28f4f123")
		accountId, _ := primitive.ObjectIDFromHex("60d5f483b3c2c66b28f4f124")
		text := models.User{Id: userId, AccountId: accountId, DeviceToken: "Device-123", Role: "ADMIN", Ip: "142.168.1.1"}

		token, err1 := helper.GenerateToken(text)
		if err1 != nil {
			fmt.Println("loi roi: ", err1.Error())
		}

		fmt.Println("data token ne::: ", token)

		err := accountCollection.FindOne(ctx, map[string]interface{}{"username": account.Username}).Decode(&accountDB)
		if err != nil {
			fmt.Println("login error: ", err.Error())
			c.JSON(http.StatusBadRequest, responses.ToDoResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Tài khoản không tồn tại"}})
			return
		}

		if accountDB.Username == "" {
			c.JSON(http.StatusBadRequest, responses.ToDoResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Tài khoản không tồn tại"}})
			return
		} else {
			//process enter password
			// encoding password (SHA1)

			pKey := os.Getenv("ENDCODE_PASSWORD_KEY")
			encodePassword := utils.EncodeSHA1Password(account.Password, pKey)
			if accountDB.Password != "" && len(accountDB.Password) > 0 && subtle.ConstantTimeCompare([]byte(accountDB.Password), []byte(encodePassword)) == 1 {
				accountDB.HidePassWord()
				c.JSON(http.StatusOK, responses.AuthResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": accountDB}})
				// handle create token here
				return
			}
		}
		c.JSON(http.StatusBadRequest, responses.ToDoResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Sai mật khẩu"}})
	}
}
