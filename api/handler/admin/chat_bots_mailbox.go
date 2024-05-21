package admin

import (
	"chatplus/core"
	"chatplus/handler"
	"chatplus/service/oss"
	"chatplus/store/vo"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"net/http"
)

type ChatBotmailbox struct {
	handler.BaseHandler
	redis         *redis.Client
	uploadManager *oss.UploaderManager
}

func NewChatBotmailbox(app *core.AppServer, db *gorm.DB, redis *redis.Client, manager *oss.UploaderManager) *ChatBotmailbox {
	return &ChatBotmailbox{
		BaseHandler:   handler.BaseHandler{App: app, DB: db},
		redis:         redis,
		uploadManager: manager,
	}
}

// 展示每一个问题
func (h *ChatBotmailbox) List(c *gin.Context) {
	botid := c.Query("botid")
	result := &[]vo.Mailbox{}
	err := h.DB.Where("bot_id = ? AND reply_id = ?", botid, 0).Find(&result).Error
	if err != nil {
		logger.Error("查询出错", err)
	}
	c.JSON(http.StatusOK, result)
}

// 对问题进行作答
func (h *ChatBotmailbox) Upload(c *gin.Context) {
	reply := &vo.Mailbox{}
	if err := c.ShouldBindJSON(&reply); err != nil {
		logger.Error("作答参数绑定失败", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	res := &vo.Mailbox{}
	err := h.DB.Where("useid = ? AND question = ?", reply.Userid, reply.Question).First(res).Error
	if err != nil {
		logger.Error("查询对应记录出错", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	res.ReplyId = reply.ReplyId
	res.Question = reply.Question
	if err := h.DB.Save(res).Error; err != nil {
		logger.Error("更新失败", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, res)
}
