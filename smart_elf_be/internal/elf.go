package internal

import (
	"errors"
	"log"
	"smart_elf_standalone/internal/model"
	"smart_elf_standalone/internal/service"

	"gorm.io/gorm"
)

// SmartElf 智能插件的核心结构体
type SmartElf struct {
	ConfigService *service.ConfigService
	EventService  *service.EventService
}

// NewSmartElf 创建新的SmartElf实例
func NewSmartElf(
	configService *service.ConfigService,
	eventService *service.EventService,
) *SmartElf {
	return &SmartElf{
		ConfigService: configService,
		EventService:  eventService,
	}
}

// HandleLarkEvent 处理飞书事件回调
func (e *SmartElf) HandleLarkEvent(req *model.LarkCallbackRequest) (*model.LarkCallbackResponse, error) {
	// 处理URL验证
	if req.Type == "url_verification" {
		log.Printf("信息: 处理URL验证请求: type=url_verification")
		return &model.LarkCallbackResponse{
			Challenge: req.Challenge,
		}, nil
	}

	// 处理消息事件
	if req.Header.EventType == "im.message.receive_v1" {
		log.Printf("信息: 处理消息接收事件: event_type=im.message.receive_v1")

		// 交给EventService处理具体的事件逻辑
		err := e.EventService.HandleMessageEvent(req)
		if err != nil {
			log.Printf("错误: 处理消息事件失败: %v", err)
			return nil, err
		}

		// 返回空响应（飞书回调不需要额外响应内容）
		return &model.LarkCallbackResponse{}, nil
	}

	log.Printf("警告: 收到未处理的事件类型: %s", req.Type)
	return &model.LarkCallbackResponse{}, nil
}

// GetSignature 获取插件签名
func (e *SmartElf) GetSignature(projectKey string) (string, error) {
	signature, err := e.ConfigService.GetSignature(projectKey)
	if err != nil {
		log.Printf("错误: 获取签名失败: %v, project_key=%s", err, projectKey)
		return "", err
	}
	return signature, nil
}

// UpdateConfig 更新插件配置
func (e *SmartElf) UpdateConfig(req *model.ConfigRequest) error {
	err := e.ConfigService.UpdateConfig(req)
	if err != nil {
		log.Printf("错误: 更新配置失败: %v, project_key=%s", err, req.ProjectKey)
		return err
	}
	return nil
}

// QueryConfig 查询插件配置
func (e *SmartElf) QueryConfig(projectKey string) (*model.ConfigResponse, error) {
	config, err := e.ConfigService.QueryConfig(projectKey)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("错误: 查询配置失败: %v, project_key=%s", err, projectKey)
		return nil, err
	}
	return config, nil
}
