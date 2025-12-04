package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"smart_elf_standalone/internal/model"
	"smart_elf_standalone/pkg/config"
	"strings"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	projSDK "github.com/larksuite/project-oapi-sdk-golang"
	"github.com/larksuite/project-oapi-sdk-golang/core"
	"github.com/larksuite/project-oapi-sdk-golang/service/field"
	"github.com/larksuite/project-oapi-sdk-golang/service/project"
	"github.com/larksuite/project-oapi-sdk-golang/service/workitem"
	"gorm.io/gorm"
)

// EventService 事件服务
type EventService struct {
	db            *gorm.DB
	configService *ConfigService
	feishuCfg     config.FeishuConfig
}

// NewEventService 创建事件服务实例
func NewEventService(db *gorm.DB, configService *ConfigService, feishuCfg config.FeishuConfig) *EventService {
	return &EventService{
		db:            db,
		configService: configService,
		feishuCfg:     feishuCfg,
	}
}

// parseMessageContent 解析消息内容
func (s *EventService) parseMessageContent(contentStr string) (*model.TextContent, error) {
	var content model.TextContent
	err := json.Unmarshal([]byte(contentStr), &content)
	if err != nil {
		return nil, err
	}
	return &content, nil
}

// getLarkSDKCli 获取飞书SDK客户端
func (s *EventService) getLarkSDKCli(config *model.AppConfig) (*lark.Client, error) {
	if config == nil || config.BotID == "" || config.BotSecret == "" {
		return nil, errors.New("invalid bot configuration")
	}

	// 创建飞书SDK客户端
	client := lark.NewClient(config.BotID, config.BotSecret, lark.WithOpenBaseUrl(s.feishuCfg.IMOpenAPIHost))
	return client, nil
}

func (s *EventService) HandleMessageEvent(req *model.LarkCallbackRequest) (err error) {
	ctx := context.Background()
	if req == nil || req.Event == nil || req.Event.Message == nil {
		return errors.New("invalid event request")
	}

	// 获取消息内容
	message := req.Event.Message
	content, err := s.parseMessageContent(message.Content)
	if err != nil {
		log.Printf("错误: 解析消息内容失败: %v", err)
		return err
	}

	// 获取发送者信息
	sender := req.Event.Sender
	senderID := ""
	if sender != nil && sender.SenderID != nil {
		senderID = sender.SenderID.OpenID
	}

	// 忽略机器人自己发送的消息
	if sender != nil && sender.IsBot {
		log.Printf("信息: 忽略机器人自己的消息: %s", message.MessageID)
		return nil
	}
	config, err := s.configService.GetConfigBySignature(req.Signature)
	if err != nil {
		log.Printf("错误: 验证签名失败: %v, signature=%s", err, req.Signature)
		return err
	}

	larkCli, _ := s.getLarkSDKCli(config)
	if larkCli == nil {
		return
	}
	reporterOpenID := senderID
	userResp, err := larkCli.Contact.User.Get(ctx, larkcontact.NewGetUserReqBuilder().
		UserIdType("open_id").UserId(reporterOpenID).Build())
	if err != nil {
		log.Printf("get lark user failed,err=%s", err.Error())
		return
	}
	if !userResp.Success() {
		log.Printf("get lark user failed,code=%s,msg=%s,requestID=%s", userResp.Code, userResp.Msg, userResp.RequestId())
		return
	}
	reporterDisplayName := userResp.Data.User.Name
	if reporterDisplayName == nil {
		log.Printf("reporterDisplayName is nil")
		return
	}

	type TextContent struct {
		Text string `json:"text"`
	}

	meegoCli, _ := s.GetFeishuProjectClient()

	reg := regexp.MustCompile(`@_user_[0-9]+`)
	contentText := reg.ReplaceAllString(content.Text, "")
	contentText = strings.TrimSpace(contentText)

	//创建工单工作项
	userKey := config.APIUserKey
	fields := make([]*field.FieldValuePair, 0, 1)
	fields = append(fields, &field.FieldValuePair{
		FieldValue: fmt.Sprintf("%s###%s", *reporterDisplayName, reporterOpenID),
		FieldKey:   config.CreatorFieldKey})
	wiReq := workitem.NewCreateWorkItemReqBuilder().WorkItemTypeKey(config.WorkItemTypeKey).
		ProjectKey(config.ProjectKey).Name(contentText).FieldValuePairs(fields).TemplateID(config.WorkItemTemplateID).Build()
	wiResp, err := meegoCli.WorkItem.CreateWorkItem(ctx, wiReq, core.WithUserKey(userKey))

	if err != nil {
		log.Printf("create workitem failed,err=%s, logid=%s", err.Error(), wiResp.Header.Get("x-tt-logid"))
		return
	}

	//开启了自动拉群功能
	if config.CreateGroupSwitch {
		go func() {
			titleCN := ("[工单]" + contentText)
			titleEN := ("[Ticket]" + contentText)
			reqCreateGroup := larkim.NewCreateChatReqBuilder().UserIdType("open_id").SetBotManager(true).
				Body(
					larkim.NewCreateChatReqBodyBuilder().
						Name("[工单]" + contentText).I18nNames(&larkim.I18nNames{
						ZhCn: &titleCN,
						EnUs: &titleEN,
						JaJp: &titleEN,
					}).
						OwnerId(reporterOpenID).
						BotIdList([]string{config.BotID}).
						Build()).
				Build()
			respGroup, err0 := larkCli.Im.Chat.Create(ctx, reqCreateGroup)
			//这里失败只打日志不影响后续流程
			if err0 != nil {
				log.Printf("create group failed,err=%s", err0.Error())
				return
			}
			if !respGroup.Success() {
				log.Printf("create group failed,code=%s,msg=%s,requestID=%s", respGroup.Code, respGroup.Msg, respGroup.RequestId())
				return
			}
			chatID := respGroup.Data.ChatId
			if chatID == nil {
				log.Printf("chat id is nil")
				return
			}
			upFields := make([]*field.FieldValuePair, 0, 1)
			upFields = append(upFields, &field.FieldValuePair{
				FieldKey:   "group_type",
				FieldValue: "bind"}, &field.FieldValuePair{
				FieldKey:   "group_id",
				FieldValue: *chatID})
			wiUpdateReq := workitem.NewUpdateWorkItemReqBuilder().WorkItemTypeKey(config.WorkItemTypeKey).
				ProjectKey(config.ProjectKey).UpdateFields(upFields).WorkItemID(wiResp.Data).Build()
			wiUpdateResp, err0 := meegoCli.WorkItem.UpdateWorkItem(ctx, wiUpdateReq, core.WithUserKey(userKey))

			if err0 != nil {
				log.Printf("update workitem group failed,err=%s, logid=%s", err0.Error(), wiUpdateResp.Header.Get("x-tt-logid"))
				return
			}
			if !wiUpdateResp.Success() {
				log.Printf("update workitem group failed,code=%s", wiUpdateResp.Error())
			}
		}()

	}

	//开启了创建后反馈工单链接功能时
	if config.ReplySwitch {
		newID := wiResp.Data
		go func() {
			//获取simplename，获取工作项apiname
			var simpleName string

			respProj, errP := meegoCli.Project.GetProjectDetail(ctx,
				project.NewGetProjectDetailReqBuilder().ProjectKeys([]string{config.ProjectKey}).UserKey(userKey).Build(),
				core.WithUserKey(userKey))
			if errP != nil {
				log.Printf("get project info failed,err=%s", errP.Error())
				return
			}
			if p, ok := respProj.Data[config.ProjectKey]; ok {
				simpleName = p.SimpleName
			}
			wiURL := fmt.Sprintf("%s/%s/%s/detail/%d", s.feishuCfg.ProjectAPIHost, simpleName, config.WorkItemAPIName, newID)
			cnContent := make([][]map[string]interface{}, 0, 2)
			cnContentLine1 := make([]map[string]interface{}, 0, 2)
			cnContentLine1 = append(cnContentLine1, map[string]interface{}{
				"tag":   "text",
				"text":  "工单内容: ",
				"style": []string{"bold"},
			})
			cnContentLine1 = append(cnContentLine1, map[string]interface{}{
				"tag":  "text",
				"text": content,
			})
			cnContentLine2 := make([]map[string]interface{}, 0, 2)
			cnContentLine2 = append(cnContentLine2, map[string]interface{}{
				"tag":   "text",
				"text":  "工单链接: ",
				"style": []string{"bold"},
			})
			cnContentLine2 = append(cnContentLine2, map[string]interface{}{
				"tag":  "a",
				"text": "查看详情",
				"href": wiURL,
			})
			cnContent = append(cnContent, cnContentLine1, cnContentLine2)

			enContent := make([][]map[string]interface{}, 0, 2)
			enContentLine1 := make([]map[string]interface{}, 0, 2)
			enContentLine1 = append(enContentLine1, map[string]interface{}{
				"tag":   "text",
				"text":  "Ticket Content: ",
				"style": []string{"bold"},
			})
			enContentLine1 = append(enContentLine1, map[string]interface{}{
				"tag":  "text",
				"text": content,
			})
			enContentLine2 := make([]map[string]interface{}, 0, 2)
			enContentLine2 = append(enContentLine2, map[string]interface{}{
				"tag":   "text",
				"text":  "Ticket Link: ",
				"style": []string{"bold"},
			})
			enContentLine2 = append(enContentLine2, map[string]interface{}{
				"tag":  "a",
				"text": "View Detail",
				"href": wiURL,
			})
			enContent = append(enContent, enContentLine1, enContentLine2)
			msg := map[string]interface{}{
				"zh_cn": map[string]interface{}{
					"title":   "🆕工单创建成功️",
					"content": cnContent,
				},
				"en_us": map[string]interface{}{
					"title":   "🆕Ticket created",
					"content": enContent,
				},
			}
			msgStr, _ := json.Marshal(msg)
			respIm, err1 := larkCli.Im.Message.Create(ctx,
				larkim.NewCreateMessageReqBuilder().
					ReceiveIdType("open_id").
					Body(
						larkim.NewCreateMessageReqBodyBuilder().
							ReceiveId(reporterOpenID).
							MsgType("post").
							Content(string(msgStr)).
							Build()).
					Build())
			if err1 != nil {
				log.Printf("send msg failed,err=%s", err1.Error())
				return
			}
			if !respIm.Success() {
				log.Printf("send msg failed,code=%s,msg=%s,requestID=%s", respIm.Code, respIm.Msg, respIm.RequestId())
				return
			}
		}()
	}

	return

}

func (s *EventService) GetFeishuProjectClient() (*projSDK.Client, error) {

	clientV2 := projSDK.NewClient(s.feishuCfg.PluginID, s.feishuCfg.PluginSecret,
		projSDK.WithOpenBaseUrl(s.feishuCfg.ProjectAPIHost), projSDK.WithAccessTokenType(core.AccessTokenTypePlugin))

	return clientV2, nil
}
