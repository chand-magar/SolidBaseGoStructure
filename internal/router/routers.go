package router

import (
	controllers "github.com/chand-magar/SolidBaseGoStructure/internal/controllers"
	repositories "github.com/chand-magar/SolidBaseGoStructure/internal/repositories"
	services "github.com/chand-magar/SolidBaseGoStructure/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AllRouter(db *gorm.DB) *gin.Engine {

	r := gin.Default()

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	users := r.Group("/v1/webmaster")
	{
		users.POST("/users", userController.Create)
		users.GET("/users", userController.GetAll)
		users.GET("/users/:id", userController.FindOne)
		users.PUT("/users/:id", userController.Update)
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to Solid Base Go Structure API"})
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Unable to find the specified API"})
	})

	return r
}
