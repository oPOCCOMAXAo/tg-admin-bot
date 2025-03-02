package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ColumnConfig struct {
	Name   string
	IsBool bool
	IsInt  bool
}

func (c ColumnConfig) ValueInt(value int64) clause.Expr {
	if c.IsBool {
		if value > 0 {
			value = 1
		}
	}

	return gorm.Expr("?", value)
}
