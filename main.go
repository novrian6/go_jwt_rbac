package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

var users = map[uint]User{
	1: {ID: 1, Name: "Admin", Role: "admin"},
	2: {ID: 2, Name: "User", Role: "user"},
}

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "Token required"})
			c.Abort()
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret-key"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		for _, role := range allowedRoles {
			if claims.Role == role {
				c.Set("user_id", claims.UserID)
				c.Set("role", claims.Role)
				c.Next()
				return
			}
		}

		c.JSON(403, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}

func main() {
	r := gin.Default()

	// Public endpoint
	r.GET("/public", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Public access"})
	})

	// Admin only
	r.GET("/admin", RequireRole("admin"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Admin area"})
	})

	// Admin and user
	r.GET("/dashboard", RequireRole("admin", "user"), func(c *gin.Context) {
		role := c.GetString("role")
		c.JSON(200, gin.H{
			"message": "Dashboard access",
			"role":    role,
		})
	})

	// Simple login endpoint for testing: POST /login { "user_id": 1 }
	r.POST("/login", func(c *gin.Context) {
		var req struct {
			UserID uint `json:"user_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}
		user, ok := users[req.UserID]
		if !ok {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}

		claims := Claims{
			UserID: user.ID,
			Role:   user.Role,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte("secret-key"))
		if err != nil {
			c.JSON(500, gin.H{"error": "Could not generate token"})
			return
		}
		c.JSON(200, gin.H{"token": tokenString})
	})

	r.Run(":8080")
}
