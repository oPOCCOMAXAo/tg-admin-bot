package models

type MessageInfo struct {
	ID            int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Time          int64  `gorm:"column:time;not null;index:search_idx"`
	ChatID        int64  `gorm:"column:chat_id;not null"`                  // tg chat id
	MessageID     int64  `gorm:"column:message_id;not null"`               // tg message id
	UserID        int64  `gorm:"column:user_id;not null;index:search_idx"` // tg user id
	IsProcessed   bool   `gorm:"column:is_processed;not null;default:0"`
	Score         uint16 `gorm:"column:score;not null;default:0"`
	HasRULetters  bool   `gorm:"column:has_ru_letters;not null;default:0"`
	HasCaps       bool   `gorm:"column:has_caps;not null;default:0"`
	HasShort      bool   `gorm:"column:has_short;not null;default:0"`
	HasLong       bool   `gorm:"column:has_long;not null;default:0"`
	CountLinks    uint8  `gorm:"column:count_links;not null;default:0"`
	CountEmbeds   uint8  `gorm:"column:count_embeds;not null;default:0"`
	CountMedias   uint8  `gorm:"column:count_media;not null;default:0"`
	CountMentions uint8  `gorm:"column:count_mentions;not null;default:0"`
	IsFast        bool   `gorm:"column:is_fast;not null;default:0"`
}

func (MessageInfo) TableName() string {
	return "message_info"
}
