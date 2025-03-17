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
	// 添加健康检查端点
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// user-service 启动时注册
	config := api.DefaultConfig()
	config.Address = "consul:8500"
	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}
	registration := &api.AgentServiceRegistration{
		ID:   "user-service-1",
		Name: "user-service",
		Port: 8081,
		Check: &api.AgentServiceCheck{
			HTTP:     "http://user-service:8081/health", // 健康检查地址
			Interval: "5s",                              // 健康检查间隔
			Timeout:  "2s",                              // 健康检查超时
		},
	}
	if err = client.Agent().ServiceRegister(registration); err != nil {
		panic(err)
	}

	r.Run(":8081")
}
