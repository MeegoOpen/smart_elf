package model

import (
	"gorm.io/gorm"
)

// AppConfig 应用配置模型
type AppConfig struct {
	gorm.Model
	BotID                string `gorm:"column:bot_id" json:"bot_id"`
	BotSecret            string `gorm:"column:bot_secret" json:"bot_secret"`
	BotVerificationToken string `gorm:"column:bot_verification_token" json:"bot_verification_token"`
	ProjectKey           string `gorm:"column:project_key" json:"project_key"`
	TenantKey            string `gorm:"column:tenant_key" json:"tenant_key"`
	WorkItemTypeKey      string `gorm:"column:work_item_type_key" json:"work_item_type_key"`
	WorkItemAPIName      string `gorm:"column:work_item_api_name" json:"work_item_api_name"`
	WorkItemTemplateID   int64  `gorm:"column:work_item_template_id" json:"work_item_template_id"`
	CreatorFieldKey      string `gorm:"column:creator_field_key" json:"creator_field_key"`
	ReplySwitch          bool   `gorm:"column:reply_switch" json:"reply_switch"`
	CreateGroupSwitch    bool   `gorm:"column:create_group_switch" json:"create_group_switch"`
	Signature            string `gorm:"column:signature" json:"signature"`
	APIUserKey           string `gorm:"column:api_user_key" json:"api_user_key"`
}

// TableName 指定表名
func (a AppConfig) TableName() string {
	return "smart_elf"
}

// BotInfo 机器人信息
type BotInfo struct {
	BotID             string  `json:"bot_id" binding:"required"`
	BotSecret         string  `json:"bot_secret" binding:"required"`
	VerificationToken *string `json:"verification_token"`
}

// ConfigRequest 配置请求结构
type ConfigRequest struct {
	ProjectKey string  `json:"project_key" binding:"required"`
	Config     *Config `json:"config" binding:"required"`
}

// Config 配置信息
type Config struct {
	Bot                BotInfo `json:"bot_info" binding:"required"`
	WorkItemType       string  `json:"work_item_type_key"`
	WorkItemAPIName    string  `json:"work_item_api_name"`
	WorkItemTemplateID int64   `json:"work_item_template_id"`
	CreatorFieldKey    string  `json:"creator_field_key"`
	ReplySwitch        bool    `json:"reply_switch"`
	CreateGroupSwitch  bool    `json:"create_group_switch"`
	APIUserKey         string  `json:"api_user_key"`
}

// ConfigResponse 配置响应结构
type ConfigResponse struct {
	Config *Config `json:"config"`
}

// SignatureRequest 签名请求
type SignatureRequest struct {
	ProjectKey string `json:"project_key" binding:"required"`
}

// SignatureResponse 签名响应
type SignatureResponse struct {
	Signature string `json:"signature"`
}

// LarkCallbackRequest 飞书回调请求
type LarkCallbackRequest struct {
	Challenge string              `json:"challenge"`
	Token     string              `json:"token"`
	Type      string              `json:"type"`
	Header    *LarkCallbackHeader `json:"header"`
	Event     *LarkCallbackEvent  `json:"event"`
	Signature string              `json:"signature"`
}

// LarkCallbackHeader 飞书回调头部
type LarkCallbackHeader struct {
	Token      string `json:"token"`
	EventType  string `json:"event_type"`
	CreateTime string `json:"create_time"`
	AppID      string `json:"app_id"`
}

// LarkCallbackEvent 飞书回调事件
type LarkCallbackEvent struct {
	Message *LarkMessage `json:"message"`
	Sender  *LarkSender  `json:"sender"`
}

// LarkMessage 飞书消息
type LarkMessage struct {
	MessageID   string `json:"message_id"`
	RootID      string `json:"root_id"`
	ParentID    string `json:"parent_id"`
	MsgType     string `json:"msg_type"`
	Content     string `json:"content"`
	CreateTime  string `json:"create_time"`
	UpdatedTime string `json:"updated_time"`
}

// LarkSender 飞书发送者
type LarkSender struct {
	SenderID   *LarkSenderID `json:"sender_id"`
	SenderType string        `json:"sender_type"`
	IsBot      bool          `json:"is_bot"`
}

// LarkSenderID 飞书发送者ID
type LarkSenderID struct {
	OpenID  string `json:"open_id"`
	UserID  string `json:"user_id"`
	UnionID string `json:"union_id"`
}

// LarkCallbackResponse 飞书回调响应
type LarkCallbackResponse struct {
	Challenge string `json:"challenge,omitempty"`
}

// TextContent 文本内容
type TextContent struct {
	Text string `json:"text"`
}
