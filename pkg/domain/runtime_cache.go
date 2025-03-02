package domain

import (
	"sync"

	"github.com/opoccomaxao/tg-admin-bot/pkg/models"
)

type RuntimeCache struct {
	mu   sync.RWMutex
	data map[int64]*models.RuntimeConfig
}

func NewRuntimeCache() *RuntimeCache {
	return &RuntimeCache{
		data: make(map[int64]*models.RuntimeConfig),
	}
}

func (c *RuntimeCache) GetConfig(chatID int64) *models.RuntimeConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.data[chatID]
}

func (c *RuntimeCache) SetFromChatConfig(chatID int64, config *models.ChatConfig) {
	c.mu.Lock()
	defer c.mu.Unlock()

	cfg := c.data[chatID]
	if cfg == nil {
		cfg = &models.RuntimeConfig{}
		c.data[chatID] = cfg
	}

	cfg.UpdateFromChatConfig(config)
}
