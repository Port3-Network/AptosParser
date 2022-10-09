package models

import "strconv"

func ParseInt64(val interface{}) int64 {
	switch v := val.(type) {
	case int8:
		return int64(v)
	case uint8:
		return int64(v)
	case int16:
		return int64(v)
	case uint16:
		return int64(v)
	case int32:
		return int64(v)
	case uint32:
		return int64(v)
	case int64:
		return v
	case uint64:
		return int64(v)
	case int:
		return int64(v)
	case uint:
		return int64(v)
	case float32:
		return int64(v)
	case float64:
		return int64(v)
	case string:
		num, _ := strconv.ParseInt(v, 0, 64)
		return num
	}

	return 0
}
