package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
)

func main() {
	r := gin.Default()
	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		// 模拟数据库查询
		c.JSON(http.StatusOK, gin.H{"id": id, "name": "Alice"})
	})

	// user-service 启动时注册
	client, _ := api.NewClient(api.DefaultConfig())
	registration := &api.AgentServiceRegistration{
		ID:   "user-service-1",
		Name: "user-service",
		Port: 8081,
	}
	client.Agent().ServiceRegister(registration)

	r.Run(":8081")
}
