package models

type RuntimeConfig struct {
	Enabled  []ConfigID
	Antispam AntispamConfig
}

type AntispamConfig struct {
	Debug bool
}

func (c *RuntimeConfig) UpdateFromChatConfig(config *ChatConfig) {
	c.Enabled = config.EnabledList()
	c.Antispam = config.AntispamConfig()
}
