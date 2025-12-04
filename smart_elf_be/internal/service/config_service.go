package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"smart_elf_standalone/internal/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ConfigService 配置服务
type ConfigService struct {
	db *gorm.DB
}

// NewConfigService 创建配置服务实例
func NewConfigService(db *gorm.DB) *ConfigService {
	return &ConfigService{
		db: db,
	}
}

// UpdateConfig 更新配置
func (s *ConfigService) UpdateConfig(req *model.ConfigRequest) error {
    if req == nil || req.Config == nil {
        return errors.New("invalid config request")
    }

	var appConfig model.AppConfig
	result := s.db.Where("project_key = ?", req.ProjectKey).First(&appConfig)

	// 如果不存在，则创建新配置
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 生成新签名
			signature, err := s.generateSignature(req.ProjectKey)
			if err != nil {
				log.Printf("错误: 生成签名失败: %v", err)
				return err
			}

            appConfig = model.AppConfig{
                ProjectKey:           req.ProjectKey,
                BotID:                req.Config.Bot.BotID,
                BotSecret:            req.Config.Bot.BotSecret,
                BotVerificationToken: getStringValue(req.Config.Bot.VerificationToken, ""),
                Signature:            signature,
                TenantKey:            "", // 可以根据实际情况设置
                WorkItemTypeKey:      req.Config.WorkItemType,
                WorkItemAPIName:      req.Config.WorkItemAPIName,
                WorkItemTemplateID:   req.Config.WorkItemTemplateID,
                CreatorFieldKey:      req.Config.CreatorFieldKey,
                ReplySwitch:          req.Config.ReplySwitch,
                CreateGroupSwitch:    req.Config.CreateGroupSwitch,
                APIUserKey:           req.Config.APIUserKey,
            }

			if err := s.db.Create(&appConfig).Error; err != nil {
				log.Printf("错误: 创建配置失败: %v", err)
				return err
			}
			log.Printf("信息: 创建配置成功: project_key=%s", req.ProjectKey)
		} else {
			log.Printf("错误: 查询配置失败: %v", result.Error)
			return result.Error
		}
	} else {
		// 更新现有配置
        updates := map[string]interface{}{
            "bot_id":                 req.Config.Bot.BotID,
            "bot_secret":             req.Config.Bot.BotSecret,
            "bot_verification_token": getStringValue(req.Config.Bot.VerificationToken, ""),
            "work_item_type_key":     req.Config.WorkItemType,
            "work_item_api_name":     req.Config.WorkItemAPIName,
            "work_item_template_id":  req.Config.WorkItemTemplateID,
            "creator_field_key":      req.Config.CreatorFieldKey,
            "reply_switch":           req.Config.ReplySwitch,
            "create_group_switch":    req.Config.CreateGroupSwitch,
            "api_user_key":           req.Config.APIUserKey,
            "updated_at":             time.Now(),
        }

		if err := s.db.Model(&appConfig).Updates(updates).Error; err != nil {
			log.Printf("错误: 更新配置失败: %v", err)
			return err
		}
		log.Printf("信息: 更新配置成功: project_key=%s", req.ProjectKey)
	}

	return nil
}

// QueryConfig 查询配置
func (s *ConfigService) QueryConfig(projectKey string) (*model.ConfigResponse, error) {
	var appConfig model.AppConfig
	result := s.db.Where("project_key = ?", projectKey).First(&appConfig)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("警告: 配置不存在: project_key=%s", projectKey)
			return nil, errors.New("config not found")
		}
		log.Printf("错误: 查询配置失败: %v", result.Error)
		return nil, result.Error
	}

	// 构建响应
    response := &model.ConfigResponse{
        Config: &model.Config{
            Bot: model.BotInfo{
                BotID:             appConfig.BotID,
                BotSecret:         appConfig.BotSecret,
                VerificationToken: &appConfig.BotVerificationToken,
            },
            WorkItemType:       appConfig.WorkItemTypeKey,
            WorkItemAPIName:    appConfig.WorkItemAPIName,
            WorkItemTemplateID: appConfig.WorkItemTemplateID,
            CreatorFieldKey:    appConfig.CreatorFieldKey,
            ReplySwitch:        appConfig.ReplySwitch,
            CreateGroupSwitch:  appConfig.CreateGroupSwitch,
            APIUserKey:         appConfig.APIUserKey,
        },
    }

	return response, nil
}

// GetSignature 获取或生成签名
func (s *ConfigService) GetSignature(projectKey string) (string, error) {
	var appConfig model.AppConfig
	result := s.db.Where("project_key = ?", projectKey).First(&appConfig)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 如果配置不存在，创建一个默认配置并生成签名
			signature, err := s.generateSignature(projectKey)
			if err != nil {
				return "", err
			}

			appConfig = model.AppConfig{
				ProjectKey: projectKey,
				Signature:  signature,
				BotID:      "",
				BotSecret:  "",
			}

			if err := s.db.Create(&appConfig).Error; err != nil {
				return "", err
			}

			return signature, nil
		}
		return "", result.Error
	}

	// 如果签名为空，生成新签名
	if appConfig.Signature == "" {
		signature, err := s.generateSignature(projectKey)
		if err != nil {
			return "", err
		}

		appConfig.Signature = signature
		if err := s.db.Save(&appConfig).Error; err != nil {
			return "", err
		}

		return signature, nil
	}

	return appConfig.Signature, nil
}

// generateSignature 生成签名
func (s *ConfigService) generateSignature(projectKey string) (string, error) {
	uuidStr := uuid.New().String()
	timestamp := time.Now().UnixNano()
	rawStr := projectKey + uuidStr + fmt.Sprintf("%d", timestamp)

	hasher := md5.New()
	hasher.Write([]byte(rawStr))
	hashBytes := hasher.Sum(nil)
	signature := hex.EncodeToString(hashBytes)

	return signature, nil
}

// GetConfigByProjectKey 根据项目密钥获取配置
func (s *ConfigService) GetConfigByProjectKey(projectKey string) (*model.AppConfig, error) {
	var appConfig model.AppConfig
	result := s.db.Where("project_key = ?", projectKey).First(&appConfig)
	if result.Error != nil {
		return nil, result.Error
	}
	return &appConfig, nil
}

// GetConfigBySignature 根据签名获取配置
func (s *ConfigService) GetConfigBySignature(signature string) (*model.AppConfig, error) {
	var appConfig model.AppConfig
	result := s.db.Where("signature = ?", signature).First(&appConfig)
	if result.Error != nil {
		return nil, result.Error
	}
	return &appConfig, nil
}

// 辅助函数
// getStringValue 从指针获取字符串值，如果为nil则返回默认值
func getStringValue(val *string, defaultValue string) string {
	if val != nil {
		return *val
	}
	return defaultValue
}

// getBoolValue 从指针获取布尔值，如果为nil则返回默认值
func getBoolValue(val *bool, defaultValue bool) bool {
	if val != nil {
		return *val
	}
	return defaultValue
}

// getInt64Value 从指针获取int64值，如果为nil则返回默认值
func getInt64Value(val *int64, defaultValue int64) int64 {
	if val != nil {
		return *val
	}
	return defaultValue
}
