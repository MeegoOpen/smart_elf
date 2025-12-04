package handler

import (
	"log"
	"net/http"
	"smart_elf_standalone/internal"
	"smart_elf_standalone/internal/model"
	"time"

	"github.com/gin-gonic/gin" // 确保 go.mod 中已添加依赖: go get -u github.com/gin-gonic/gin
)

// Handler HTTP处理器
type Handler struct {
	smartElf *internal.SmartElf
}

// NewHandler 创建处理器实例
func NewHandler(smartElf *internal.SmartElf) *Handler {
	return &Handler{
		smartElf: smartElf,
	}
}

// HandleLarkEvent 处理飞书事件回调
func (h *Handler) HandleLarkEvent(c *gin.Context) {
	var req model.LarkCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("错误: 绑定请求参数失败: %v", err)
		Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// 调用SmartElf处理事件
	resp, err := h.smartElf.HandleLarkEvent(&req)
	if err != nil {
		log.Printf("错误: 处理飞书事件失败: %v", err)
		Error(c, http.StatusInternalServerError, "Failed to handle event")
		return
	}

	Success(c, resp)
}

// GetSignature 获取插件签名
func (h *Handler) GetSignature(c *gin.Context) {
	var req model.SignatureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("错误: 绑定请求参数失败: %v", err)
		Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	signature, err := h.smartElf.GetSignature(req.ProjectKey)
	if err != nil {
		log.Printf("错误: 获取签名失败: project_key=%s, err=%v", req.ProjectKey, err)
		Error(c, http.StatusInternalServerError, "Failed to get signature")
		return
	}

	Success(c, model.SignatureResponse{
		Signature: signature,
	})
}

// UpdateConfig 更新插件配置
func (h *Handler) UpdateConfig(c *gin.Context) {
	var req model.ConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("错误: 绑定请求参数失败: %v", err)
		Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.smartElf.UpdateConfig(&req)
	if err != nil {
		log.Printf("错误: 更新配置失败: project_key=%s, err=%v", req.ProjectKey, err)
		Error(c, http.StatusInternalServerError, "Failed to update config")
		return
	}

	Success(c, gin.H{"message": "Config updated successfully"})
}

// QueryConfig 查询插件配置
func (h *Handler) QueryConfig(c *gin.Context) {
	projectKey := c.Query("project_key")
	if projectKey == "" {
		log.Printf("错误: 缺少project_key参数")
		Error(c, http.StatusBadRequest, "Missing project_key parameter")
		return
	}

	config, err := h.smartElf.QueryConfig(projectKey)
	if err != nil {
		log.Printf("错误: 查询配置失败: project_key=%s, err=%v", projectKey, err)
		Error(c, http.StatusInternalServerError, "Failed to query config")
		return
	}

	Success(c, config)
}

// HealthCheck 健康检查
func (h *Handler) HealthCheck(c *gin.Context) {
	Success(c, gin.H{
		"status":    "ok",
		"service":   "smart-elf-plugin",
		"timestamp": time.Now().Unix(),
	})
}
