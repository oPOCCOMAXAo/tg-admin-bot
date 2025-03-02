package models

import "strconv"

type ConfigID int64

const (
	CfgUnknown              ConfigID = 0
	CfgEnabledMuteRuLetters ConfigID = 1
	CfgEnabledAntispam      ConfigID = 2
	CfgAntispamDebug        ConfigID = 3
)

func (id *ConfigID) Int64Ref() *int64 {
	return (*int64)(id)
}

func (id ConfigID) StringID() string {
	return strconv.FormatInt(int64(id), 10)
}
