package localtime

import (
	"errors"
	"time"
)

func MarshalRFC3339(t time.Time) ([]byte, error) {
	return quoteJSONString([]byte(t.Local().Format(time.RFC3339))), nil
}

func UnmarshalRFC3339(data []byte) (*time.Time, error) {
	if string(data) == "null" {
		return nil, nil
	}

	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return nil, errors.New("Time.UnmarshalJSON: input is not a JSON string")
	}
	data = data[len(`"`) : len(data)-len(`"`)]

	timestamp, err := time.ParseInLocation(time.RFC3339, string(data), time.Local)
	if err != nil {
		// Try to parse RFC3339 without timezone
		var err2 error
		timestamp, err2 = time.ParseInLocation("2006-01-02T15:04:05", string(data), time.Local)
		if err2 != nil {
			return nil, errors.Join(err, err2)
		}
	}
	utcTimestamp := timestamp.In(time.UTC)
	return &utcTimestamp, nil
}

func MarshalRFC3339Date(t time.Time) ([]byte, error) {
	return quoteJSONString([]byte(t.Local().Format(time.DateOnly))), nil
}

func UnmarshalRFC3339Date(data []byte) (*time.Time, error) {
	if string(data) == "null" {
		return nil, nil
	}
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return nil, errors.New("Time.UnmarshalJSON: input is not a JSON string")
	}
	data = data[len(`"`) : len(data)-len(`"`)]

	timestamp, err := time.ParseInLocation(time.DateOnly, string(data), time.Local)
	if err != nil {
		return nil, err
	}
	utcTimestamp := timestamp.Local()
	return &utcTimestamp, nil
}

func MarshalNorwegianDateTime(t time.Time) ([]byte, error) {
	return quoteJSONString([]byte(t.Local().Format("02.01.2006 15:04:05"))), nil
}

func UnmarshalNorwegianDateTime(data []byte) (*time.Time, error) {
	if string(data) == "null" {
		return nil, nil
	}
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return nil, errors.New("Time.UnmarshalJSON: input is not a JSON string")
	}
	data = data[len(`"`) : len(data)-len(`"`)]

	timestamp, err := time.ParseInLocation("02.01.2006 15:04:05", string(data), time.Local)
	if err != nil {
		return nil, err
	}
	utcTimestamp := timestamp.In(time.UTC)
	return &utcTimestamp, nil
}

// Quotes the string if it'is not already quoted.
func quoteJSONString(s []byte) []byte {
	if len(s) == 0 {
		return []byte{}
	}
	if s[0] != '"' {
		s = append([]byte{'"'}, s...)
	}
	if s[len(s)-1] != '"' {
		s = append(s, '"')
	}
	return s
}

// Unqotes the string if it is quoted.
func unquoteJSONString(data []byte) []byte {
	if len(data) < 2 {
		return data
	}
	if data[0] == '"' {
		data = data[1:]
	}
	if data[len(data)-1] == '"' {
		data = data[:len(data)-1]
	}
	return data
}
