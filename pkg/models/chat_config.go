package models

type ChatConfig struct {
	ID                 int64 `gorm:"column:id;primaryKey;autoIncrement"`
	TgID               int64 `gorm:"column:tg_id;index:tg_id,unique"`
	CreatedAt          int64 `gorm:"column:created_at;default:0"`
	UpdatedAt          int64 `gorm:"column:updated_at;default:0"`
	EnabledMuteLetters bool  `gorm:"column:enabled_mute_letters;default:0"`
}

func (ChatConfig) TableName() string {
	return "chat_config"
}

func (ChatConfig) RuleColumn(rule Rule) string {
	switch rule {
	case RuleMuteLetters:
		return "enabled_mute_letters"
	default:
		return ""
	}
}

func (c *ChatConfig) RulesList() RulesList {
	res := make(RulesList, 0)

	if c.EnabledMuteLetters {
		res = append(res, RuleMuteLetters)
	}

	return res
}
