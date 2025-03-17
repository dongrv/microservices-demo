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
	// 添加健康检查端点
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// order-service 发现用户服务地址
	config := api.DefaultConfig()
	config.Address = "consul:8500"
	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}
	registration := &api.AgentServiceRegistration{
		ID:   "order-service-1",
		Name: "order-service",
		Port: 8082,
		Check: &api.AgentServiceCheck{
			HTTP:     "http://order-service:8082/health", // 健康检查地址
			Interval: "5s",                               // 健康检查间隔
			Timeout:  "2s",                               // 健康检查超时
		},
	}
	if err = client.Agent().ServiceRegister(registration); err != nil {
		panic(err)
	}

	r.Run(":8082")
}
