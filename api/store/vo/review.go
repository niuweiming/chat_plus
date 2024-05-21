package vo

type Review struct {
	UserId  string `json:"userid" binding:"required"`
	Content string `json:"content" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Review  uint   `json:"review" binging:"required"`
}

type Mailbox struct {
	Id           int    `gorm:"column:id;AUTO_INCREMENT;primary_key"`
	Userid       string `gorm:"column:userid;NOT NULL"`
	Question     string `gorm:"column:question;NOT NULL"`
	Botsid       string `gorm:"column:botsid;NOT NULL"`
	ReplyId      int    `gorm:"column:reply_id;default:0"`
	ReplyContent string `gorm:"column:reply_content;default:该问题暂未解决"`
}
