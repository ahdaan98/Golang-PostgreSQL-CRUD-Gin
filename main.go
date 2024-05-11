package main

import (
	"github.com/ahdaan98/go-gorm-crud/config"
	"github.com/ahdaan98/go-gorm-crud/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	cfg := config.LoadEnv()
	config.ConnectDB(cfg)
}

func main(){
	r := gin.Default()

	routes.Routes(r)

	r.Run()
}