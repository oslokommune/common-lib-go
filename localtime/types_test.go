package localtime_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/oslokommune/common-lib-go/localtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestDateTime(t *testing.T) {
	suite.Run(t, &DateTimeSuite{})
}

func TestDate(t *testing.T) {
	suite.Run(t, &DateSuite{})
}

func TestNorwegianDateTime(t *testing.T) {
	suite.Run(t, &NorwegianDateTimeSuite{})
}

func TestStructFieldSupport(t *testing.T) {
	suite.Run(t, &StructFieldSupportSuite{})
}

type DateTimeSuite struct {
	suite.Suite
}

func (suite *DateTimeSuite) SetupTest() {
	localtime.SetLocalTimeZone("Europe/Oslo")
}

func (suite *DateTimeSuite) TestUnmarshal_ParsesTimestampIntoUTC_WhenTimeZoneSpecified() {
	t := suite.T()

	timestamp := []byte(`"2020-02-01T13:34:56+02:00"`)

	var datetime localtime.DateTime
	err := json.Unmarshal(timestamp, &datetime)
	assert.NoError(t, err)
	assert.Equal(t, 2020, datetime.Year())
	assert.Equal(t, time.Month(2), datetime.Month())
	assert.Equal(t, 1, datetime.Day())

	assert.Equal(t, time.UTC, datetime.Location())
	assert.Equal(t, 11, datetime.Hour())
	assert.Equal(t, 34, datetime.Minute())
	assert.Equal(t, 56, datetime.Second())
}

func (suite *DateTimeSuite) TestUnmarshal_ParsesLocalTimestampIntoUTC_WhenTimeZoneNotSpecified() {
	t := suite.T()

	timestamp := []byte(`"2020-02-01T12:34:56"`)

	var datetime localtime.DateTime
	err := json.Unmarshal(timestamp, &datetime)
	assert.NoError(t, err)
	assert.Equal(t, 2020, datetime.Year())
	assert.Equal(t, time.Month(2), datetime.Month())
	assert.Equal(t, 1, datetime.Day())

	assert.Equal(t, time.UTC, datetime.Location())
	assert.Equal(t, 11, datetime.Hour())
	assert.Equal(t, 34, datetime.Minute())
	assert.Equal(t, 56, datetime.Second())
}

func (suite *DateTimeSuite) TestMarshal_FormatsTimeIntoLocalTimezone() {
	t := suite.T()

	datetime := localtime.DateTime{time.Date(2020, time.February, 1, 11, 34, 56, 0, time.UTC)}

	marshalled, err := json.Marshal(datetime)

	assert.NoError(t, err)
	assert.Equal(t, []byte(`"2020-02-01T12:34:56+01:00"`), marshalled)
}

func (suite *DateTimeSuite) TestMarshalUnmarshal_AreInverse() {
	t := suite.T()

	timestamp := []byte(`"2020-02-01T12:34:56+01:00"`)
	var datetime localtime.DateTime
	err := json.Unmarshal(timestamp, &datetime)
	assert.NoError(t, err)
	marshalled, err := json.Marshal(datetime)
	assert.NoError(t, err)
	assert.Equal(t, timestamp, marshalled)
}

func (suite *DateTimeSuite) TestParseDateTime_ParsesZonedTimestampIntoUTC() {
	t := suite.T()

	datetime, err := localtime.ParseDateTime("2020-02-01T13:34:56+02:00")
	assert.NoError(t, err)
	assert.Equal(t, 2020, datetime.Year())
	assert.Equal(t, time.Month(2), datetime.Month())
	assert.Equal(t, 1, datetime.Day())

	assert.Equal(t, time.UTC, datetime.Location())
	assert.Equal(t, 11, datetime.Hour())
	assert.Equal(t, 34, datetime.Minute())
	assert.Equal(t, 56, datetime.Second())
}

func (suite *DateTimeSuite) TestParseDateTime_ParsesLocalTimestampIntoUTC() {
	t := suite.T()

	datetime, err := localtime.ParseDateTime("2020-02-01T12:34:56")
	assert.NoError(t, err)
	assert.Equal(t, 2020, datetime.Year())
	assert.Equal(t, time.Month(2), datetime.Month())
	assert.Equal(t, 1, datetime.Day())

	assert.Equal(t, time.UTC, datetime.Location())
	assert.Equal(t, 11, datetime.Hour())
	assert.Equal(t, 34, datetime.Minute())
	assert.Equal(t, 56, datetime.Second())
}

func (suite *DateTimeSuite) TestString_ReturnsLocalTimestamp() {
	t := suite.T()

	datetime := localtime.DateTime{time.Date(2020, time.February, 1, 11, 34, 56, 0, time.UTC)}

	timestamp := datetime.String()

	assert.Equal(t, "2020-02-01T12:34:56+01:00", timestamp)
}

type DateSuite struct {
	suite.Suite
}

func (suite *DateSuite) SetupTest() {
	localtime.SetLocalTimeZone("Europe/Oslo")
}

func (suite *DateSuite) TestUnmarshal_ReturnsCorrectDate() {
	t := suite.T()

	data := []byte(`"2020-02-01"`)

	var date localtime.Date
	err := json.Unmarshal(data, &date)
	assert.NoError(t, err)
	assert.Equal(t, 2020, date.Year())
	assert.Equal(t, time.Month(2), date.Month())
	assert.Equal(t, 1, date.Day())
}

// This behaviour is convenient - if we converted from GMT+1 to UTC in the code it'd be confusing as the date would change
func (suite *DateSuite) TestUnmarshal_ParsesDateIntoLocalTimezone() {
	t := suite.T()

	data := []byte(`"2020-02-01"`)

	var date localtime.Date
	err := json.Unmarshal(data, &date)
	assert.NoError(t, err)
	assert.Equal(t, time.Local, date.Location())
}

func (suite *DateSuite) TestMarshal_ReturnsCorrectDate() {
	t := suite.T()

	date := localtime.Date{time.Date(2020, 2, 1, 0, 0, 0, 0, time.Local)}
	timestamp, err := json.Marshal(date)
	assert.NoError(t, err)
	assert.Equal(t, []byte("\"2020-02-01\""), timestamp)
}

func (suite *DateSuite) TestMarshalUnmarshal_AreInverse() {
	t := suite.T()

	data := []byte(`"2020-02-01"`)

	var date localtime.Date
	err := json.Unmarshal(data, &date)
	assert.NoError(t, err)

	marshalled, err := json.Marshal(date)

	assert.NoError(t, err)
	assert.Equal(t, data, marshalled)
}

func (suite *DateSuite) TestParseDate_ParsesDateIntoLocalDate() {
	t := suite.T()

	datetime, err := localtime.ParseDate("2020-02-01")
	assert.NoError(t, err)
	assert.Equal(t, 2020, datetime.Year())
	assert.Equal(t, time.Month(2), datetime.Month())
	assert.Equal(t, 1, datetime.Day())

	assert.Equal(t, time.Local, datetime.Location())
}

func (suite *DateSuite) TestString_ReturnsDate() {
	t := suite.T()

	datetime := localtime.Date{time.Date(2020, time.February, 1, 11, 34, 56, 0, time.UTC)}

	datestamp := datetime.String()

	assert.Equal(t, "2020-02-01", datestamp)
}

type NorwegianDateTimeSuite struct {
	suite.Suite
}

func (suite *NorwegianDateTimeSuite) SetupTest() {
	localtime.SetLocalTimeZone("Europe/Oslo")
}

func (suite *NorwegianDateTimeSuite) TestUnmarshal_ParsesLocalTimestampIntoUTC() {
	t := suite.T()
	timestamp := []byte(`"01.02.2020 12:34:56"`)

	var datetime localtime.NorwegianDateTime
	err := json.Unmarshal(timestamp, &datetime)
	assert.NoError(t, err)
	assert.Equal(t, 2020, datetime.Year())
	assert.Equal(t, time.Month(2), datetime.Month())
	assert.Equal(t, 1, datetime.Day())

	assert.Equal(t, time.UTC, datetime.Location())
	assert.Equal(t, 11, datetime.Hour())
	assert.Equal(t, 34, datetime.Minute())
	assert.Equal(t, 56, datetime.Second())
}

func (suite *NorwegianDateTimeSuite) TestMarshal_FormatsTimeInLocalTimesone() {
	t := suite.T()

	datetime := localtime.NorwegianDateTime{time.Date(2020, time.February, 1, 11, 34, 56, 0, time.UTC)}
	timestamp, err := json.Marshal(datetime)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`"01.02.2020 12:34:56"`), timestamp)
}

func (suite *NorwegianDateTimeSuite) TestMarshalUnmarshal_AreInverse() {
	t := suite.T()
	timestamp := []byte(`"01.02.2020 12:34:56"`)

	var datetime localtime.NorwegianDateTime
	err := json.Unmarshal(timestamp, &datetime)
	assert.NoError(t, err)
	assert.Equal(t, time.UTC, datetime.Location())

	marshalled, err := json.Marshal(datetime)
	assert.NoError(t, err)
	assert.Equal(t, timestamp, marshalled)
}

func (suite *NorwegianDateTimeSuite) TestParseDateTime_ParsesLocalTimestampIntoUTC() {
	t := suite.T()

	datetime, err := localtime.ParseNorwegianDateTime("01.02.2020 12:34:56")
	assert.NoError(t, err)
	assert.Equal(t, 2020, datetime.Year())
	assert.Equal(t, time.Month(2), datetime.Month())
	assert.Equal(t, 1, datetime.Day())

	assert.Equal(t, time.UTC, datetime.Location())
	assert.Equal(t, 11, datetime.Hour())
	assert.Equal(t, 34, datetime.Minute())
	assert.Equal(t, 56, datetime.Second())
}

func (suite *NorwegianDateTimeSuite) TestString_ReturnsLocalTimestamp() {
	t := suite.T()

	datetime := localtime.NorwegianDateTime{time.Date(2020, time.February, 1, 11, 34, 56, 0, time.UTC)}

	timestamp := datetime.String()

	assert.Equal(t, "01.02.2020 12:34:56", timestamp)
}

type StructFieldSupportSuite struct {
	suite.Suite
}

func (suite *StructFieldSupportSuite) SetupTest() {
	localtime.SetLocalTimeZone("Europe/Oslo")
}

func (suite *StructFieldSupportSuite) TestUnmarshal_StructWithLocalDateTypes() {
	t := suite.T()

	bytes := []byte(`{"date": "2020-02-01", "datetime": "2020-02-01T13:34:56+02:00", "norwegianDateTime": "01.02.2020 12:34:56"}`)

	var data struct {
		Date              localtime.Date              `json:"date" binding:"required"`
		DateTime          localtime.DateTime          `json:"datetime" binding:"required"`
		NorwegianDateTime localtime.NorwegianDateTime `json:"norwegianDateTime" binding:"required"`
	}

	err := json.Unmarshal(bytes, &data)
	assert.NoError(t, err)
	assert.Equal(t, 2020, data.Date.Year())
	assert.Equal(t, time.Month(2), data.Date.Month())
	assert.Equal(t, 1, data.Date.Day())
	assert.Equal(t, time.Local, data.Date.Location())

	assert.Equal(t, 2020, data.DateTime.Year())
	assert.Equal(t, time.Month(2), data.DateTime.Month())
	assert.Equal(t, 1, data.DateTime.Day())
	assert.Equal(t, time.UTC, data.DateTime.Location())
	assert.Equal(t, 11, data.DateTime.Hour())
	assert.Equal(t, 34, data.DateTime.Minute())
	assert.Equal(t, 56, data.DateTime.Second())

	assert.Equal(t, 2020, data.NorwegianDateTime.Year())
	assert.Equal(t, time.Month(2), data.NorwegianDateTime.Month())
	assert.Equal(t, 1, data.NorwegianDateTime.Day())
	assert.Equal(t, time.UTC, data.NorwegianDateTime.Location())
	assert.Equal(t, 11, data.NorwegianDateTime.Hour())
	assert.Equal(t, 34, data.NorwegianDateTime.Minute())
	assert.Equal(t, 56, data.NorwegianDateTime.Second())

}

func (suite *StructFieldSupportSuite) TestMarshal_StructWithLocalDateTypes() {
	t := suite.T()

	data := struct {
		Date              localtime.Date              `json:"date" binding:"required"`
		DateTime          localtime.DateTime          `json:"datetime" binding:"required"`
		NorwegianDateTime localtime.NorwegianDateTime `json:"norwegianDateTime" binding:"required"`
	}{
		Date:              localtime.Date{time.Date(2020, time.February, 1, 0, 0, 0, 0, time.Local)},
		DateTime:          localtime.DateTime{time.Date(2020, time.February, 1, 11, 34, 56, 0, time.UTC)},
		NorwegianDateTime: localtime.NorwegianDateTime{time.Date(2020, time.February, 1, 11, 34, 56, 0, time.UTC)},
	}

	expected := []byte(`{"date":"2020-02-01","datetime":"2020-02-01T12:34:56+01:00","norwegianDateTime":"01.02.2020 12:34:56"}`)

	bytes, err := json.Marshal(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, bytes)
}
