package chatimpl

import (
	"bytes"
	"chatplus/core"
	"chatplus/handler"
	"chatplus/service/oss"
	"chatplus/store/model"
	"chatplus/store/vo"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
	"gorm.io/gorm"
	"io"
	"net/http"
)

type ChatBotHandler struct {
	handler.BaseHandler
	redis         *redis.Client
	uploadManager *oss.UploaderManager
}

// 调用api返回的数据结构返回的(是通过botid找到KbIDs)
type Response struct {
	Code int `json:"code"`
	Data []struct {
		BotID          string   `json:"bot_id"`
		BotName        string   `json:"bot_name"`
		Description    string   `json:"description"`
		HeadImage      string   `json:"head_image"`
		KbIDs          []string `json:"kb_ids"`
		KbNames        []string `json:"kb_names"`
		Model          string   `json:"model"`
		PromptSetting  string   `json:"prompt_setting"`
		UpdateTime     string   `json:"update_time"`
		UserID         string   `json:"user_id"`
		WelcomeMessage string   `json:"welcome_message"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// 前端提供的json
type PassData struct {
	Sessionid string `json:"sessionid"`
	Chatid    string `json:"chatid"`
	UserID    string `json:"user_id"`
	Bot_id    string `json:"bot_id"`
	Question  string `json:"question"`
}

// RequestData 定义请求体数据结构
type RequestData struct {
	UserID   string   `json:"user_id"`
	KbIDs    []string `json:"kb_ids"`
	Question string   `json:"question"`
}

func NewChatBotHandler(app *core.AppServer, db *gorm.DB, redis *redis.Client, manager *oss.UploaderManager) *ChatBotHandler {
	return &ChatBotHandler{
		BaseHandler:   handler.BaseHandler{App: app, DB: db},
		redis:         redis,
		uploadManager: manager,
	}
}

func (h *ChatBotHandler) ChatBotsHandler(c *gin.Context) {
	var passData PassData
	if err := c.ShouldBindJSON(&passData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求体"})
		return
	}
	logger.Error(passData.Question)
	passData = PassData{
		UserID:   "zzp",
		Bot_id:   passData.Bot_id,
		Question: passData.Question,
	}
	kdIds := GetKbid(passData.UserID, passData.Bot_id)

	var requestData RequestData
	requestData = RequestData{
		UserID:   "zzp",
		KbIDs:    kdIds,
		Question: passData.Question,
	}
	fmt.Println(kdIds)
	logger.Error("我的问题是什么", requestData.Question)
	url := "http://122.51.6.203:8777/api/local_doc_qa/local_doc_chat"
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	reqBody, err := json.Marshal(requestData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "序列化请求体失败"})
		return
	}

	httpClient := &http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建请求失败"})
		return
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送请求失败"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取响应体失败"})
		return
	}

	var responseBody map[string]interface{}
	if err := json.Unmarshal(body, &responseBody); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解析响应体失败"})
		return
	}
	logger.Error(responseBody)
	c.JSON(http.StatusOK, responseBody)
}

// 通过机器id获取数据库id
func GetKbid(userid string, botid string) []string {
	var passData PassData
	passData = PassData{
		UserID: userid,
		Bot_id: botid,
	}

	url := "http://122.51.6.203:8777/api/local_doc_qa/get_bot_info"
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	reqBody, err := json.Marshal(passData)
	if err != nil {
		fmt.Println(err)
		fmt.Println("x序列化失败")
		return nil
	}

	httpClient := &http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println(err)
		fmt.Println("创建请求失败")
		return nil
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		fmt.Println("发送请求失败")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		fmt.Println("读取响应体失败")
		return nil
	}

	var responseBody Response
	if err := json.Unmarshal(body, &responseBody); err != nil {
		fmt.Println(err)
		fmt.Println("解析响应体失败")
		return nil
	}
	return responseBody.Data[0].KbIDs
}

// 不满意发送给专家
func (h *ChatBotHandler) Recordmail(c *gin.Context) {
	logger.Error("执行Recordmail函数")
	recordmail := &vo.Mailbox{}
	if err := c.ShouldBindJSON(recordmail); err != nil {
		logger.Error("用户满意度表参数绑定失败!")
		logger.Error(err)
	}
	recordmails := model.Mailboxs{
		Userid:   recordmail.Userid,
		Question: recordmail.Question,
		Botsid:   recordmail.Botsid,
	}
	if err := h.DB.Create(&recordmails).Error; err != nil {
		logger.Error("用户满意度表插入数据库失败!")
		logger.Error(err)
	}
}

// 查看个人信箱
func (h *ChatBotHandler) View(c *gin.Context) {
	logger.Error("查看个人信箱")
	userid := c.Query("userid")
	result := &[]vo.Mailbox{}
	if err := h.DB.Where("userid = ?", userid).Find(&result).Error; err != nil {
		logger.Error("查询失败", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, result)
}
