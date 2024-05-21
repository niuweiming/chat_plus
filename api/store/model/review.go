package model

type Reviews struct {
	Id      int    `gorm:"column:ID;type:int(11);AUTO_INCREMENT;primary_key" json:"ID"`
	UserId  string `gorm:"column:UserId;type:char(255);NOT NULL" json:"UserId"`
	Content string `gorm:"column:Content;type:text;NOT NULL" json:"Content"`
	Name    string `gorm:"column:Name;type:varchar(255);NOT NULL" json:"Name"`
	ModelId int    `gorm:"column:ModelId;type:int(11);NOT NULL" json:"ModelId"`
	Review  uint   `gorm:"column:Review;type:tinyint(4) unsigned;NOT NULL" json:"Review"`
}

type Mailboxs struct {
	Id           int    `gorm:"column:id;AUTO_INCREMENT;primary_key"`
	Userid       string `gorm:"column:userid;NOT NULL"`
	Question     string `gorm:"column:question;NOT NULL"`
	Botsid       string `gorm:"column:botsid;NOT NULL"`
	ReplyId      int    `gorm:"column:reply_id;default:0"`
	ReplyContent string `gorm:"column:reply_content;default:该问题暂未解决"`
}
