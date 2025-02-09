package models

import "strconv"

type Rule int64

const (
	RuleUnknown     Rule = 0
	RuleMuteLetters Rule = 1
)

func (r *Rule) Int64Ref() *int64 {
	return (*int64)(r)
}

func (r Rule) StringID() string {
	return strconv.FormatInt(int64(r), 10)
}

type RulesList []Rule

func (l RulesList) IsEmpty() bool {
	return len(l) == 0
}
