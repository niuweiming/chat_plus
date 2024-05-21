package vo

type Review struct {
	UserId  string `json:"userid" binding:"required"`
	Content string `json:"content" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Review  uint   `json:"review" binging:"required"`
}

type Mailbox struct {
	Id           int    `json:"id"`
	Userid       string `json:"userid"`
	Question     string `json:"question"`
	Botsid       string `json:"botsid"`
	ReplyId      int    `json:"reply_id"`
	ReplyContent string `json:"reply_content"`
}
