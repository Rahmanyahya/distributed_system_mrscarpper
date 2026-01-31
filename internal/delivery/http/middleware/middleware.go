package middleware

import (
	"distributed_system/internal/config"
	"distributed_system/internal/domain/admin"
	"distributed_system/pkg/crypto"
	"distributed_system/pkg/response"
	"strings"


	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func ValidationRegistrationAgent(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Unauthorized")
			c.Abort()
			return
		}

		token := strings.SplitN(authHeader, " ", 2)[1] 

		if err := bcrypt.CompareHashAndPassword([]byte(token), []byte(cfg.Security.AgentSecret)); err != nil {
			response.Unauthorized(c, "Unauthorized")
			c.Abort()
			return
		}

		c.Next()
	}
}

func InternalGetConfigVaidation(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Unauthorized")
			c.Abort()
			return
		}

		token := strings.SplitN(authHeader, " ", 2)[1] 

		isValid, uuid, err := crypto.Verify(token, cfg.Security.AgentSig); 
		if err != nil {
			response.Unauthorized(c, "Unauthorized")
			c.Abort()
			return
		}

		if !isValid {
			response.Unauthorized(c, "Unauthorized")
			c.Abort()
			return
		}

		c.Set("uuid", uuid)
		c.Next()
	}
}

func AdminValidation(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		token := strings.SplitN(authHeader, " ", 2)[1]
		
		payload, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.Security.JWTSecret), nil
		})
		if err != nil || !payload.Valid  {
			response.Unauthorized(c, "Unauthorized")
			c.Abort()
			return
		}

		claims, ok := payload.Claims.(*admin.Claims)
		if !ok {
			c.Next()
			return
		}

		if claims.Role != "admin" {
			response.Forbidden(c, "Forbidden")
			c.Abort()
			return
		}

		c.Next()
	}
}

func ValidationAgentWorker(cfg *config.WorkerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Unauthorized")
			c.Abort()
			return
		}

		token := strings.SplitN(authHeader, " ", 2)[1] 

		if cfg.Auth.InternalKey != token {
			response.Unauthorized(c, "Unauthorized")
			c.Abort()
			return
		}

		c.Next()
	}
}
