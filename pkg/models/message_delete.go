package models

type MessageDelete struct {
	ID        int64 `gorm:"column:id;primaryKey;autoIncrement;not null"`
	ChatID    int64 `gorm:"column:chat_id;not null"`
	MessageID int64 `gorm:"column:message_id;not null"`
	ExecuteAt int64 `gorm:"column:execute_at;not null;default:0;index:idx_execute_at"`
}

func (MessageDelete) TableName() string {
	return "message_delete"
}
