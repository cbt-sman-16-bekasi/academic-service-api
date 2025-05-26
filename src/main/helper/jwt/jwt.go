package jwt

import (
	"fmt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/redisstore"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yon-module/yon-framework/exception"
	"github.com/yon-module/yon-framework/logger"
	"github.com/yon-module/yon-framework/server/response"
	"log"
	"os"
	"strings"
	"time"
)

var secretKey = []byte(os.Getenv("yon.security.secret_key"))
var AllAccess = []string{"read", "delete", "create", "update", "list"}

type Claims struct {
	Username   string   `json:"username"`
	Role       string   `json:"role"`
	Permission []string `json:"permission"`
	SchoolCode string   `json:"school_code"`
	ID         uint     `json:"id"`
}

func GenerateJWT(claim Claims) (string, error) {
	claims := jwt.MapClaims{
		"username":    claim.Username,
		"role":        claim.Role,
		"permission":  claim.Permission,
		"school_code": claim.SchoolCode,
		"exp":         time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ErrorResponse(response.Unauthorized, "Please input your token", nil).Json(c)
			c.Abort()
			return
		}

		// Format token: "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			response.ErrorResponse(response.Unauthorized, "Your token invalid", nil).Json(c)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.ErrorResponse(response.Unauthorized, "Your token invalid", nil).Json(c)
			c.Abort()
			return
		}

		userClaims := Claims{
			Username:   claims["username"].(string),
			Role:       claims["role"].(string),
			SchoolCode: claims["school_code"].(string),
			Permission: toStringSlice(claims["permission"]),
		}

		c.Set("claims", userClaims)
		c.Next()
	}
}

func RequirePermission(requiredRole []string, requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			response.ErrorResponse(response.UnprocessableEntity, "You are not authorized", nil).Json(c)
			c.Abort()
			return
		}

		userClaims, ok := claims.(Claims)
		if !ok {
			response.ErrorResponse(response.UnprocessableEntity, "You are not authorized", nil).Json(c)
			c.Abort()
			return
		}

		role := userClaims.Role
		permissions := userClaims.Permission

		hasMatch := false
		for _, rl := range requiredRole {
			if role == rl {
				hasMatch = true
				break
			}
		}
		if !hasMatch {
			response.ErrorResponse(response.UnprocessableEntity, "You are not authorized", nil).Json(c)
			c.Abort()
			return
		}

		for _, perm := range permissions {
			if strings.EqualFold(perm, requiredPermission) {
				c.Next()
				return
			}
		}

		response.ErrorResponse(response.UnprocessableEntity, "You don't have access, please contact administrator", nil).Json(c)
		c.Abort()
		return
	}
}

func toStringSlice(input interface{}) []string {
	if input == nil {
		return []string{}
	}
	slice, ok := input.([]interface{})
	if !ok {
		return []string{}
	}
	result := make([]string, len(slice))
	for i, v := range slice {
		str, ok := v.(string)
		if !ok {
			continue
		}
		result[i] = str
	}
	return result
}

func GetDataClaims(c *gin.Context) Claims {
	claims, exists := c.Get("claims")
	if !exists {
		return Claims{}
	}
	return claims.(Claims)
}

func SaveDetailUser(key string, user interface{}, exp time.Duration) {
	_ = redisstore.SetJSON(key, user, exp)
}

func ExtractDetailUser(key string, detail interface{}) {
	err := redisstore.GetJSON(key, detail)
	if err != nil {
		panic(exception.NewBadRequestExceptionStruct(response.Unauthorized, err.Error()))
	}
}

func GetID(key string) float64 {
	var data map[string]interface{}
	err := redisstore.GetJSON(key, &data)
	if data == nil {
		logger.Log.Error().Msg("No data found")
		panic(exception.NewBadRequestExceptionStruct(response.Unauthorized, "Data user not found"))
	}
	if err != nil {
		logger.Log.Error().Msg("No data found")
		panic(exception.NewBadRequestExceptionStruct(response.Unauthorized, err.Error()))
	}
	if val, ok := data["ID"]; ok {
		if floatVal, ok := val.(float64); ok {
			return floatVal
		} else {
			// handle jika bukan float64, misalnya error log
			log.Println("ID bukan float64")
		}
	}
	panic(exception.NewBadRequestExceptionStruct(response.Unauthorized, "Data user not found"))
}

func GetIDClaims(c *gin.Context) float64 {
	claims := GetDataClaims(c)
	return GetID(claims.Username)
}
