package main

import (
	"ecommerce/apigateway/pkg/auth"
	"ecommerce/apigateway/pkg/config"
	"ecommerce/apigateway/pkg/order"
	"ecommerce/apigateway/pkg/product"

	"github.com/gin-gonic/gin"
)

func main() {
	c := config.NewConfig().InLocalConfig().InDockerComposeEnv()

	r := gin.Default()

	authSvc := *auth.RegisterRoutes(r, c)
	product.RegisterRoutes(r, c, &authSvc)
	order.RegisterRoutes(r, c, &authSvc)

	r.Run(c.Port)
}
