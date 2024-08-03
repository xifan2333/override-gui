package backend

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Manager struct {
	server *http.Server
	config config
}

func NewServerManager() *Manager {
	return &Manager{}
}

func (sm *Manager) Start() ResponseData {
	// 读取配置
	respData := ReadConfig()
	if respData.Status == "fail" {
		return ResponseData{
			Status: "fail",
			Msg:    "读取配置失败，请先创建配置文件",
		}
	}

	// 类型断言配置数据
	cfg, ok := respData.Data.(config)
	if !ok {
		return ResponseData{
			Status: "fail",
			Msg:    "配置解析失败",
		}
	}
	sm.config = cfg

	// 设置 Gin 为发布模式
	gin.SetMode(gin.ReleaseMode)

	// 初始化 Gin 路由
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true, // 允许所有源
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 初始化 Proxy 服务
	proxyService, err := NewProxyService(&cfg)
	if err != nil {
		return ResponseData{
			Status: "fail",
			Msg:    "初始化 Proxy 服务失败: " + err.Error(),
		}
	}
	proxyService.InitRoutes(router)

	// 创建 HTTP 服务器
	sm.server = &http.Server{
		Addr:    cfg.Bind,
		Handler: router,
	}

	// 创建一个传递错误的通道
	errChan := make(chan error)

	// 启动服务器的 Goroutine
	go func() {
		err := sm.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			// 如果启动有错误，发送到通道
			errChan <- err
		} else {
			// 正常关闭不视为错误
			errChan <- nil
		}
		close(errChan)
	}()

	// 这里假设服务器能正常启动
	return ResponseData{
		Status: "success",
		Msg:    "服务器已成功启动",
	}
}

func (sm *Manager) Stop() ResponseData {
	// 使用 context.Background() 立即停止服务器
	log.Println("正在尝试立即停止服务器...")

	// 停止服务器
	err := sm.server.Shutdown(context.Background())
	if err != nil {
		log.Println("停止服务器时出错:", err)
		return ResponseData{
			Status: "fail",
			Msg:    "服务器停止失败: " + err.Error(),
		}
	}

	log.Println("服务器已成功停止。")
	return ResponseData{
		Status: "success",
		Msg:    "服务器已成功停止",
	}
}
