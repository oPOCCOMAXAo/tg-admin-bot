package domain

import (
	"sync"

	"github.com/opoccomaxao/tg-admin-bot/pkg/models"
)

type RulesCache struct {
	mu   sync.RWMutex
	data map[int64][]models.Rule
}

func NewRulesCache() *RulesCache {
	return &RulesCache{
		data: make(map[int64][]models.Rule),
	}
}

func (c *RulesCache) GetRules(chatID int64) []models.Rule {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.data[chatID]
}

func (c *RulesCache) SetRules(chatID int64, rules []models.Rule) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(rules) == 0 {
		delete(c.data, chatID)

		return
	}

	c.data[chatID] = rules
}
