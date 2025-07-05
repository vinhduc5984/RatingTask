package helper

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"todo-list/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// UserClaims struct holds custom jwt claim
type UserClaims struct {
	jwt.StandardClaims
	UserID    string `json:"userId" bson:"userId"`
	AccountId string `json:"accountId" bson:"accountId"`
	DiffHour  string `json:"diffHour" bson:"diffHour"`
	DeviceId  string `json:"deviceId" bson:"deviceId"`
	Role      string `json:"role" bson:"role"`
	Ip        string `json:"ip" bson:"ip"`
	ExpireAt  int64  `json:"expireAt" bson:"expireAt"`
}

// hàm tạo token
func GenerateToken(userInfo models.User) (string, error) {

	token_lifespan, err := strconv.Atoi(os.Getenv("TOKEN_HOUR_LIFESPAN"))

	if err != nil {
		return "", err
	}
	var claims UserClaims

	claims = UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix(),
		},
		UserID:    fmt.Sprintf("%v", userInfo.Id.Hex()),
		AccountId: fmt.Sprintf("%v", userInfo.AccountId.Hex()),
		DiffHour:  fmt.Sprintf("%v", 7),
		DeviceId:  fmt.Sprintf("%v", userInfo.DeviceToken),
		Role:      fmt.Sprintf("%v", userInfo.Role),
		Ip:        fmt.Sprintf("%v", userInfo.Ip),
		ExpireAt:  time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

// valid token đúng định dạng
func TokenValid(c *gin.Context) error {
	tokenString := ExtractToken(c)
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return err
	}
	return nil
}

// lấy token từ request header api
func ExtractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func Verify(c *gin.Context) (*jwt.MapClaims, error) {

	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		errorValue, _ := err.(*jwt.ValidationError)
		switch errorValue.Errors {
		case jwt.ValidationErrorExpired:
			return nil, errors.New("SYS.MSG.VALIDATION_EXPIRED_ERROR")
		case jwt.ValidationErrorMalformed:
			return nil, errors.New("SYS.MSG.VALIDATION_MALFORMED_ERROR")
		}
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return &claims, nil
	}
	return nil, nil
}
