package jwt

import (
	"fmt"
	"strings"
	"time"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/cache"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/config"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/exception"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

// getSecretKey returns the JWT secret key from config
func getSecretKey() []byte {
	return []byte(config.GetSecurity().SecretKey)
}

var AllAccess = []string{"read", "delete", "create", "update", "list"}

type Claims struct {
	Username   string   `json:"username"`
	Role       string   `json:"role"`
	Permission []string `json:"permission"`
	SchoolCode string   `json:"school_code"`
	Id         uint     `json:"id"`
}

// GetSchoolCode returns the SchoolCode from Claims
// Implements interface for cache middleware to extract school code
func (c Claims) GetSchoolCode() string {
	return c.SchoolCode
}

func GenerateJWT(claim Claims) (string, error) {
	claims := jwt.MapClaims{
		"username":    claim.Username,
		"role":        claim.Role,
		"permission":  claim.Permission,
		"school_code": claim.SchoolCode,
		"exp":         time.Now().Add(time.Hour * 24).Unix(),
		"id":          claim.Id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(getSecretKey())
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.UnauthorizedError(c, "Please input your token")
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
			return getSecretKey(), nil
		})

		if err != nil || !token.Valid {
			response.UnauthorizedError(c, "Your token invalid")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.UnauthorizedError(c, "Your token invalid")
			c.Abort()
			return
		}

		id, ok := claims["id"]
		if !ok {
			id = "0"
		}

		userClaims := Claims{
			Username:   claims["username"].(string),
			Role:       claims["role"].(string),
			SchoolCode: claims["school_code"].(string),
			Permission: toStringSlice(claims["permission"]),
			Id:         uint(id.(float64)),
		}

		c.Set("claims", userClaims)
		c.Next()
	}
}

func RequirePermission(requiredRole []string, requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			response.ForbiddenError(c, "You are not authorized")
			c.Abort()
			return
		}

		userClaims, ok := claims.(Claims)
		if !ok {
			response.ForbiddenError(c, "You are not authorized")
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
			response.ForbiddenError(c, "You are not authorized")
			c.Abort()
			return
		}

		for _, perm := range permissions {
			if strings.EqualFold(perm, requiredPermission) {
				c.Next()
				return
			}
		}

		response.ForbiddenError(c, "You don't have access, please contact administrator")
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
	_ = cache.SetJSON(key, user, exp)
}

func ExtractDetailUser(key string, detail interface{}) {
	err := cache.GetJSON(key, detail)
	if err != nil {
		panic(exception.NewBadRequestExceptionStruct(response.Unauthorized, err.Error()))
	}
}

func GetID(key string) float64 {
	var data map[string]interface{}
	err := cache.GetJSON(key, &data)
	if data == nil {
		log.Error().Msg("No data found")
		panic(exception.NewBadRequestExceptionStruct(response.Unauthorized, "Data user not found"))
	}
	if err != nil {
		log.Error().Msg("No data found")
		panic(exception.NewBadRequestExceptionStruct(response.Unauthorized, err.Error()))
	}
	if val, ok := data["ID"]; ok {
		if floatVal, ok := val.(float64); ok {
			return floatVal
		} else {
			// handle jika bukan float64
			log.Error().Msg("ID bukan float64")
		}
	}
	panic(exception.NewBadRequestExceptionStruct(response.Unauthorized, "Data user not found"))
}

func GetIDClaims(c *gin.Context) float64 {
	claims := GetDataClaims(c)
	return GetID(claims.Username)
}

// GetSchoolCode returns the SchoolCode from JWT claims in context
// This is a helper function to easily get school code for multi-tenant filtering
func GetSchoolCode(c *gin.Context) string {
	claims := GetDataClaims(c)
	return claims.SchoolCode
}

// SchoolScopeMiddleware validates that the user has a valid school context
// Use this middleware after AuthMiddleware for endpoints that require school scope
func SchoolScopeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := GetDataClaims(c)

		if claims.SchoolCode == "" {
			response.UnauthorizedError(c, "Invalid school context. Please login again.")
			c.Abort()
			return
		}

		// Set school_code ke context untuk easy access
		c.Set("school_code", claims.SchoolCode)
		c.Next()
	}
}
