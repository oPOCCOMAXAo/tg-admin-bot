package texts

import (
	"slices"
	"strconv"
	"strings"
)

const (
	QueryParamDelimiter = " "
	QueryValueDelimiter = "="
	QuerySliceDelimiter = ","
)

type Query struct {
	Command string // Command is the first word in the query.
	Params  map[string][]string
}

func QueryCommand(command string) *Query {
	return &Query{
		Command: command,
		Params:  make(map[string][]string),
	}
}

func (q *Query) AddCommand(command string) *Query {
	q.Command = command

	return q
}

func (q *Query) AddParam(key string, value string) *Query {
	q.Params[key] = append(q.Params[key], value)

	return q
}

func (q *Query) GetInt64(key string) (int64, bool) {
	values, ok := q.Params[key]
	if !ok {
		return 0, false
	}

	if len(values) == 0 {
		return 0, true
	}

	res, err := strconv.ParseInt(values[0], 10, 64)
	if err != nil {
		return 0, false
	}

	return res, true
}

func (q *Query) GetInt64Slice(key string) ([]int64, bool) {
	values, ok := q.Params[key]
	if !ok {
		return nil, false
	}

	res := make([]int64, 0, len(values))

	for _, value := range values {
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, false
		}

		res = append(res, val)
	}

	return res, true
}

func (q *Query) GetInt64Into(key string, into *int64) bool {
	value, ok := q.GetInt64(key)
	if !ok {
		return false
	}

	*into = value

	return true
}

func (q *Query) GetString(key string) (string, bool) {
	values, ok := q.Params[key]
	if !ok {
		return "", false
	}

	if len(values) == 0 {
		return "", true
	}

	return values[0], true
}

func (q *Query) GetStringInto(key string, into *string) bool {
	value, ok := q.GetString(key)
	if !ok {
		return false
	}

	*into = value

	return true
}

func (q *Query) GetStringSlice(key string) ([]string, bool) {
	values, ok := q.Params[key]
	if !ok {
		return nil, false
	}

	return values, true
}

func (q *Query) Encode() string {
	parts := make([]string, 0, 1+len(q.Params))

	if q.Command != "" {
		parts = append(parts, JoinQueryKeyValues(q.Command, q.Params[q.Command]))
	}

	params := make([]string, 0, len(q.Params))

	for key, values := range q.Params {
		if key == q.Command {
			continue
		}

		if len(values) == 0 {
			params = append(params, key)

			continue
		}

		params = append(params, JoinQueryKeyValues(key, values))
	}

	if len(params) > 0 {
		slices.Sort(params)

		parts = append(parts, strings.Join(params, QueryParamDelimiter))
	}

	return strings.Join(parts, QueryParamDelimiter)
}

func DecodeQuery(query string) Query {
	paramsList := strings.Split(query, QueryParamDelimiter)

	res := Query{
		Params: make(map[string][]string, len(paramsList)),
	}

	res.Command, _ = SplitQueryKeyValues(paramsList[0])

	for _, param := range paramsList {
		if param == "" {
			continue
		}

		key, values := SplitQueryKeyValues(param)
		res.Params[key] = append(res.Params[key], values...)
	}

	return res
}

func JoinQueryKeyValues(key string, values []string) string {
	if len(values) == 0 {
		return key
	}

	return key + QueryValueDelimiter + strings.Join(values, QuerySliceDelimiter)
}

//nolint:nonamedreturns,mnd
func SplitQueryKeyValues(keyValues string) (key string, values []string) {
	parts := strings.SplitN(keyValues, QueryValueDelimiter, 2)
	if len(parts) >= 1 {
		key = parts[0]
	}

	if len(parts) >= 2 {
		values = strings.Split(parts[1], QuerySliceDelimiter)
	}

	return key, values
}
