package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
)

func main() {
	r := gin.Default()
	r.GET("/orders/:id", func(c *gin.Context) {
		id := c.Param("id")
		// 调用用户服务（后续替换为服务发现）
		c.JSON(http.StatusOK, gin.H{"order_id": id, "user": "Alice"})
	})
	// order-service 发现用户服务地址
	client, _ := api.NewClient(api.DefaultConfig())
	registration := &api.AgentServiceRegistration{
		ID:   "order-service-1",
		Name: "order-service",
		Port: 8082,
	}
	client.Agent().ServiceRegister(registration)

	r.Run(":8082")
}
