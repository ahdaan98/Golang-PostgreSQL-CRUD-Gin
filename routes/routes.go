package routes

import (
	"github.com/ahdaan98/go-gorm-crud/controllers"
	"github.com/gin-gonic/gin"
)

func Routes(engine *gin.Engine) {
	books := engine.Group("/books")
	{
		books.POST("/", controllers.CreateBook)
		books.GET("/get",controllers.GetBookByID)
		books.GET("/", controllers.ListAllBook)
		books.PUT("/", controllers.UpdateBook)
		books.DELETE("/", controllers.DeleteBook)
	}
}