package models

type MessageInfo struct {
	ID   int64 `gorm:"column:id;primaryKey;autoIncrement"`
	Time int64 `gorm:"column:time;not null;index:search_idx_v3"`
	// tg chat id. It is a chat where the message was sent.
	ChatID int64 `gorm:"column:chat_id;not null;index:search_idx_v3"`
	// tg message id. It is a message id in the chat.
	MessageID int64 `gorm:"column:message_id;not null;index:search_idx_v3"`
	// tg user id. It is a user who sent the message.
	UserID int64 `gorm:"column:user_id;not null;default:0;index:search_idx_v3"`
	// tg sender chat id. It is a chat which sent the message, premium feature.
	SenderChatID  int64  `gorm:"column:sender_chat_id;not null;default:0;index:search_idx_v3"`
	GroupID       string `gorm:"column:group_id;not null;default:''"`
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
	IsGroupFirst  bool   `gorm:"column:is_group_first;not null;default:0"`
}

func (MessageInfo) TableName() string {
	return "message_info"
}

func (m *MessageInfo) IsAnonymousAdmin() bool {
	return m.ChatID == m.SenderChatID
}
