package localtime

import (
	"encoding/json"
	"time"
)

// -------------------------------------------------
// Type definitions

// time.Time wrapper that formats and parses dates in RFC3339 datetime format, with or without timezone specified.
// In the code, the time is stored in UTC, but when it is marshalled to/from JSON, it is converted to the local timezone.
//
// Before marshalling/unmarshalling, remember to set the local timezone with SetLocalTimeZone().
type DateTime struct {
	time.Time
}

// time.Time wrapper that formats and parses dates in RFC3339 date format.
// In the code and when marshalling to/from JSON, it's stored in the local timezone.
//
// Before marshalling/unmarshalling, remember to set the local timezone with SetLocalTimeZone().
type Date struct {
	time.Time
}

// time.Time wrapper that formats and parses dates in dd.mm.yyyy hh:mm:ss format.
// In the code, the time is stored in UTC, but when it is marshalled to/from JSON, it is converted to the local timezone.
//
// Before marshalling/unmarshalling, remember to set the local timezone with SetLocalTimeZone().
type NorwegianDateTime struct {
	time.Time
}

// -------------------------------------------------
// Interface implementations

var _ json.Marshaler = (*DateTime)(nil)
var _ json.Unmarshaler = (*DateTime)(nil)
var _ json.Marshaler = (*Date)(nil)
var _ json.Unmarshaler = (*Date)(nil)
var _ json.Marshaler = (*NorwegianDateTime)(nil)
var _ json.Unmarshaler = (*NorwegianDateTime)(nil)

// -------------------------------------------------
// DateTime

func (t DateTime) MarshalJSON() ([]byte, error) {
	return MarshalRFC3339(t.Time)
}
func (t *DateTime) UnmarshalJSON(data []byte) error {
	timestamp, err := UnmarshalRFC3339(data)
	if err != nil {
		return err
	}
	*t = DateTime{*timestamp}
	return err
}

// Parses timestamp in RFC3339 datetime format (yyyy-mm-ddThh:mm:ssZ). If timezone is omitted, it's interpreted as a local time.
func ParseDateTime(timestamp string) (*DateTime, error) {
	jsonTimestamp := quoteJSONString([]byte(timestamp))

	var datetime DateTime
	err := json.Unmarshal(jsonTimestamp, &datetime)
	if err != nil {
		return nil, err
	}
	return &datetime, nil
}

func (t *DateTime) String() string {
	bytes, err := t.MarshalJSON()
	if err != nil {
		return ""
	}
	return string(unquoteJSONString(bytes))
}

// -------------------------------------------------
// Date

func (t Date) MarshalJSON() ([]byte, error) {
	return MarshalRFC3339Date(t.Time)
}
func (t *Date) UnmarshalJSON(data []byte) error {
	timestamp, err := UnmarshalRFC3339Date(data)
	if err != nil {
		return err
	}
	*t = Date{*timestamp}
	return err
}

func (t *Date) String() string {
	bytes, err := t.MarshalJSON()
	if err != nil {
		return ""
	}
	return string(unquoteJSONString(bytes))
}

// Parses timestamp in RFC3339 date format (yyyy-mm-dd), in local timezone.
func ParseDate(timestamp string) (*Date, error) {
	jsonTimestamp := quoteJSONString([]byte(timestamp))

	var date Date
	err := json.Unmarshal(jsonTimestamp, &date)
	if err != nil {
		return nil, err
	}
	return &date, nil
}

// -------------------------------------------------
// NorwegianDateTime

func (t NorwegianDateTime) MarshalJSON() ([]byte, error) {
	return MarshalNorwegianDateTime(t.Time)
}

func (t *NorwegianDateTime) UnmarshalJSON(data []byte) error {
	timestamp, err := UnmarshalNorwegianDateTime(data)
	if err != nil {
		return err
	}
	*t = NorwegianDateTime{*timestamp}
	return err
}

// Parses timestamp in dd.mm.yyyy hh:mm:ss format, interpreted as local time.
func ParseNorwegianDateTime(timestamp string) (*NorwegianDateTime, error) {
	jsonTimestamp := quoteJSONString([]byte(timestamp))

	var datetime NorwegianDateTime
	err := json.Unmarshal(jsonTimestamp, &datetime)
	if err != nil {
		return nil, err
	}
	return &datetime, nil
}

func (t *NorwegianDateTime) String() string {
	bytes, err := t.MarshalJSON()
	if err != nil {
		return ""
	}
	return string(unquoteJSONString(bytes))
}
