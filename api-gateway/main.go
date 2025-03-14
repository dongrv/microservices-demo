package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/api/users/:id", func(c *gin.Context) {
		proxyRequest(c, "http://user-service:8081")
	})
	r.GET("/api/orders/:id", func(c *gin.Context) {
		proxyRequest(c, "http://order-service:8082")
	})
	r.Run(":8080")
}

func proxyRequest(c *gin.Context, targetURL string) {
	// 1. 解析 targetURL
	target, err := url.Parse(targetURL)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid target URL: %v", err),
		})
		return
	}
	// 2. 创建反向代理
	proxy := httputil.NewSingleHostReverseProxy(target)

	// 3. 修改请求配置
	proxy.Director = func(r *http.Request) {
		// 保留原始请求头
		r.Header = c.Request.Header

		// 设置目标地址
		r.URL.Scheme = target.Scheme
		r.URL.Host = target.Host
		r.Host = target.Host // 重要：不分服务依赖Host头

		// 保留原始路径和查询参数
		r.URL.Path = c.Param("proxyPath")
		if c.Request.URL.RawQuery != "" || target.RawQuery != "" {
			r.URL.RawQuery = c.Request.URL.RawQuery + "&" + target.RawQuery
		}

		// 添加X-Forward头
		r.Header.Set("X-Forward-For", c.ClientIP())
		r.Header.Set("X-Forward-Host", c.Request.Host)
	}

	// 4. 自定义传输层配置（优化性能）
	proxy.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	// 5. 错误处理中间件
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("proxy err: %v", err),
		})
	}

	// 6. 创建带取消的上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 监听客户端断开
	go func() {
		<-ctx.Done()
		cancel()
	}()

	// 7. 执行代理请求
	proxy.ServeHTTP(c.Writer, c.Request.WithContext(ctx))
}
