package models

type Restriction struct {
	ID           int64 `gorm:"column:id;primaryKey;autoIncrement"`
	ExecuteAt    int64 `gorm:"column:execute_at;not null;index:search_idx"`
	ChatID       int64 `gorm:"column:chat_id;not null;index:search_idx"`
	UserID       int64 `gorm:"column:user_id;not null;default:0;index:search_idx"`
	SenderChatID int64 `gorm:"column:sender_chat_id;not null;default:0;index:search_idx"`
	IsBan        bool  `gorm:"column:is_ban;not null;default:0"`
	IsUnban      bool  `gorm:"column:is_unban;not null;default:0"`
	IsMute       bool  `gorm:"column:is_mute;not null;default:0"`
	Duration     int64 `gorm:"column:duration;not null;default:0"` // In seconds.
}

func (Restriction) TableName() string {
	return "restrictions"
}
