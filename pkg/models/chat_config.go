package models

type ChatConfig struct {
	ID                 int64 `gorm:"column:id;primaryKey;autoIncrement"`
	TgID               int64 `gorm:"column:tg_id;index:tg_id,unique"`
	CreatedAt          int64 `gorm:"column:created_at;default:0"`
	UpdatedAt          int64 `gorm:"column:updated_at;default:0"`
	EnabledMuteLetters bool  `gorm:"column:enabled_mute_letters;default:0"`
	EnabledAntispam    bool  `gorm:"column:enabled_antispam;default:0"`
	AntispamDebug      bool  `gorm:"column:antispam_debug;default:0"`
}

func (ChatConfig) TableName() string {
	return "chat_config"
}

func (ChatConfig) Columns() map[ConfigID]ColumnConfig {
	return map[ConfigID]ColumnConfig{
		CfgEnabledMuteRuLetters: {Name: "enabled_mute_letters", IsBool: true},
		CfgEnabledAntispam:      {Name: "enabled_antispam", IsBool: true},
		CfgAntispamDebug:        {Name: "antispam_debug", IsBool: true},
	}
}

func (ChatConfig) ColumnByID(id ConfigID) string {
	switch id {
	case CfgEnabledMuteRuLetters:
		return "enabled_mute_letters"
	case CfgEnabledAntispam:
		return "enabled_antispam"
	case CfgAntispamDebug:
		return "antispam_debug"
	default:
		return ""
	}
}

func (c *ChatConfig) EnabledList() []ConfigID {
	res := make([]ConfigID, 0)

	if c.EnabledMuteLetters {
		res = append(res, CfgEnabledMuteRuLetters)
	}

	if c.EnabledAntispam {
		res = append(res, CfgEnabledAntispam)
	}

	return res
}

func (c *ChatConfig) AntispamConfig() AntispamConfig {
	return AntispamConfig{
		Debug: c.AntispamDebug,
	}
}
