package estype

import "strings"

// DateFormat represents the date format for DateProperty.
// https://www.elastic.co/docs/reference/elasticsearch/mapping-reference/mapping-date-format
type DateFormat string

// String returns the string representation of the DateFormat.
func (df DateFormat) String() string {
	return string(df)
}

const (
	// DateFormatDefault is the default format combining strict date optional time and epoch millis.
	// Format: strict_date_optional_time||epoch_millis
	DateFormatDefault DateFormat = "strict_date_optional_time||epoch_millis"

	// DateFormatEpochMillis formats number of milliseconds since the epoch.
	// Format: milliseconds since epoch (e.g., 1609459200000)
	// Note: Subject to Java Long.MIN_VALUE and Long.MAX_VALUE limits.
	DateFormatEpochMillis DateFormat = "epoch_millis"

	// DateFormatEpochSecond formats number of seconds since the epoch.
	// Format: seconds since epoch (e.g., 1609459200)
	// Note: Subject to Java Long.MIN_VALUE and Long.MAX_VALUE / 1000 limits.
	DateFormatEpochSecond DateFormat = "epoch_second"

	// DateFormatDateOptionalTime is a generic ISO datetime parser with optional time.
	// Format: yyyy-MM-dd'T'HH:mm:ss.SSSZ or yyyy-MM-dd
	// Note: Lenient parsing may parse numbers as years (e.g., 292278994).
	DateFormatDateOptionalTime DateFormat = "date_optional_time"

	// DateFormatStrictDateOptionalTime is a strict ISO datetime parser with optional time.
	// Format: yyyy-MM-dd'T'HH:mm:ss.SSSZ or yyyy-MM-dd
	DateFormatStrictDateOptionalTime DateFormat = "strict_date_optional_time"

	// DateFormatStrictDateOptionalTimeNanos is a strict ISO datetime parser with nanosecond resolution.
	// Format: yyyy-MM-dd'T'HH:mm:ss.SSSSSSZ or yyyy-MM-dd
	DateFormatStrictDateOptionalTimeNanos DateFormat = "strict_date_optional_time_nanos"

	// DateFormatBasicDate formats a full date without separators.
	// Format: yyyyMMdd
	DateFormatBasicDate DateFormat = "basic_date"

	// DateFormatBasicDateTime combines basic date and time with separators.
	// Format: yyyyMMdd'T'HHmmss.SSSZ
	DateFormatBasicDateTime DateFormat = "basic_date_time"

	// DateFormatBasicDateTimeNoMillis combines basic date and time without millis.
	// Format: yyyyMMdd'T'HHmmssZ
	DateFormatBasicDateTimeNoMillis DateFormat = "basic_date_time_no_millis"

	// DateFormatBasicOrdinalDate formats a full ordinal date.
	// Format: yyyyDDD (four digit year and three digit dayOfYear)
	DateFormatBasicOrdinalDate DateFormat = "basic_ordinal_date"

	// DateFormatBasicOrdinalDateTime formats a full ordinal date and time.
	// Format: yyyyDDD'T'HHmmss.SSSZ
	DateFormatBasicOrdinalDateTime DateFormat = "basic_ordinal_date_time"

	// DateFormatBasicOrdinalDateTimeNoMillis formats a full ordinal date and time without millis.
	// Format: yyyyDDD'T'HHmmssZ
	DateFormatBasicOrdinalDateTimeNoMillis DateFormat = "basic_ordinal_date_time_no_millis"

	// DateFormatBasicTime formats time with millis and timezone.
	// Format: HHmmss.SSSZ
	DateFormatBasicTime DateFormat = "basic_time"

	// DateFormatBasicTimeNoMillis formats time without millis.
	// Format: HHmmssZ
	DateFormatBasicTimeNoMillis DateFormat = "basic_time_no_millis"

	// DateFormatBasicTTime formats time prefixed with T.
	// Format: 'T'HHmmss.SSSZ
	DateFormatBasicTTime DateFormat = "basic_t_time"

	// DateFormatBasicTTimeNoMillis formats time prefixed with T without millis.
	// Format: 'T'HHmmssZ
	DateFormatBasicTTimeNoMillis DateFormat = "basic_t_time_no_millis"

	// DateFormatBasicWeekDate formats a full weekyear date.
	// Format: xxxx'W'wwe
	DateFormatBasicWeekDate DateFormat = "basic_week_date"

	// DateFormatStrictBasicWeekDate formats a full weekyear date (strict).
	// Format: xxxx'W'wwe
	DateFormatStrictBasicWeekDate DateFormat = "strict_basic_week_date"

	// DateFormatBasicWeekDateTime combines weekyear date and time.
	// Format: xxxx'W'wwe'T'HHmmss.SSSZ
	DateFormatBasicWeekDateTime DateFormat = "basic_week_date_time"

	// DateFormatStrictBasicWeekDateTime combines weekyear date and time (strict).
	// Format: xxxx'W'wwe'T'HHmmss.SSSZ
	DateFormatStrictBasicWeekDateTime DateFormat = "strict_basic_week_date_time"

	// DateFormatBasicWeekDateTimeNoMillis combines weekyear date and time without millis.
	// Format: xxxx'W'wwe'T'HHmmssZ
	DateFormatBasicWeekDateTimeNoMillis DateFormat = "basic_week_date_time_no_millis"

	// DateFormatStrictBasicWeekDateTimeNoMillis combines weekyear date and time without millis (strict).
	// Format: xxxx'W'wwe'T'HHmmssZ
	DateFormatStrictBasicWeekDateTimeNoMillis DateFormat = "strict_basic_week_date_time_no_millis"

	// DateFormatDate formats a full date.
	// Format: yyyy-MM-dd
	DateFormatDate DateFormat = "date"

	// DateFormatStrictDate formats a full date (strict).
	// Format: yyyy-MM-dd
	DateFormatStrictDate DateFormat = "strict_date"

	// DateFormatDateHour combines date and hour.
	// Format: yyyy-MM-dd'T'HH
	DateFormatDateHour DateFormat = "date_hour"

	// DateFormatStrictDateHour combines date and hour (strict).
	// Format: yyyy-MM-dd'T'HH
	DateFormatStrictDateHour DateFormat = "strict_date_hour"

	// DateFormatDateHourMinute combines date, hour and minute.
	// Format: yyyy-MM-dd'T'HH:mm
	DateFormatDateHourMinute DateFormat = "date_hour_minute"

	// DateFormatStrictDateHourMinute combines date, hour and minute (strict).
	// Format: yyyy-MM-dd'T'HH:mm
	DateFormatStrictDateHourMinute DateFormat = "strict_date_hour_minute"

	// DateFormatDateHourMinuteSecond combines date, hour, minute and second.
	// Format: yyyy-MM-dd'T'HH:mm:ss
	DateFormatDateHourMinuteSecond DateFormat = "date_hour_minute_second"

	// DateFormatStrictDateHourMinuteSecond combines date, hour, minute and second (strict).
	// Format: yyyy-MM-dd'T'HH:mm:ss
	DateFormatStrictDateHourMinuteSecond DateFormat = "strict_date_hour_minute_second"

	// DateFormatDateHourMinuteSecondFraction combines date, hour, minute, second and fraction.
	// Format: yyyy-MM-dd'T'HH:mm:ss.SSS
	DateFormatDateHourMinuteSecondFraction DateFormat = "date_hour_minute_second_fraction"

	// DateFormatStrictDateHourMinuteSecondFraction combines date, hour, minute, second and fraction (strict).
	// Format: yyyy-MM-dd'T'HH:mm:ss.SSS
	DateFormatStrictDateHourMinuteSecondFraction DateFormat = "strict_date_hour_minute_second_fraction"

	// DateFormatDateHourMinuteSecondMillis combines date, hour, minute, second and millis.
	// Format: yyyy-MM-dd'T'HH:mm:ss.SSS
	DateFormatDateHourMinuteSecondMillis DateFormat = "date_hour_minute_second_millis"

	// DateFormatStrictDateHourMinuteSecondMillis combines date, hour, minute, second and millis (strict).
	// Format: yyyy-MM-dd'T'HH:mm:ss.SSS
	DateFormatStrictDateHourMinuteSecondMillis DateFormat = "strict_date_hour_minute_second_millis"

	// DateFormatDateTime combines date and time with timezone.
	// Format: yyyy-MM-dd'T'HH:mm:ss.SSSZ
	DateFormatDateTime DateFormat = "date_time"

	// DateFormatStrictDateTime combines date and time with timezone (strict).
	// Format: yyyy-MM-dd'T'HH:mm:ss.SSSZ
	DateFormatStrictDateTime DateFormat = "strict_date_time"

	// DateFormatDateTimeNoMillis combines date and time without millis.
	// Format: yyyy-MM-dd'T'HH:mm:ssZ
	DateFormatDateTimeNoMillis DateFormat = "date_time_no_millis"

	// DateFormatStrictDateTimeNoMillis combines date and time without millis (strict).
	// Format: yyyy-MM-dd'T'HH:mm:ssZ
	DateFormatStrictDateTimeNoMillis DateFormat = "strict_date_time_no_millis"

	// DateFormatHour formats hour of day.
	// Format: HH
	DateFormatHour DateFormat = "hour"

	// DateFormatStrictHour formats hour of day (strict).
	// Format: HH
	DateFormatStrictHour DateFormat = "strict_hour"

	// DateFormatHourMinute formats hour and minute.
	// Format: HH:mm
	DateFormatHourMinute DateFormat = "hour_minute"

	// DateFormatStrictHourMinute formats hour and minute (strict).
	// Format: HH:mm
	DateFormatStrictHourMinute DateFormat = "strict_hour_minute"

	// DateFormatHourMinuteSecond formats hour, minute and second.
	// Format: HH:mm:ss
	DateFormatHourMinuteSecond DateFormat = "hour_minute_second"

	// DateFormatStrictHourMinuteSecond formats hour, minute and second (strict).
	// Format: HH:mm:ss
	DateFormatStrictHourMinuteSecond DateFormat = "strict_hour_minute_second"

	// DateFormatHourMinuteSecondFraction formats hour, minute, second and fraction.
	// Format: HH:mm:ss.SSS
	DateFormatHourMinuteSecondFraction DateFormat = "hour_minute_second_fraction"

	// DateFormatStrictHourMinuteSecondFraction formats hour, minute, second and fraction (strict).
	// Format: HH:mm:ss.SSS
	DateFormatStrictHourMinuteSecondFraction DateFormat = "strict_hour_minute_second_fraction"

	// DateFormatHourMinuteSecondMillis formats hour, minute, second and millis.
	// Format: HH:mm:ss.SSS
	DateFormatHourMinuteSecondMillis DateFormat = "hour_minute_second_millis"

	// DateFormatStrictHourMinuteSecondMillis formats hour, minute, second and millis (strict).
	// Format: HH:mm:ss.SSS
	DateFormatStrictHourMinuteSecondMillis DateFormat = "strict_hour_minute_second_millis"

	// DateFormatOrdinalDate formats a full ordinal date.
	// Format: yyyy-DDD
	DateFormatOrdinalDate DateFormat = "ordinal_date"

	// DateFormatStrictOrdinalDate formats a full ordinal date (strict).
	// Format: yyyy-DDD
	DateFormatStrictOrdinalDate DateFormat = "strict_ordinal_date"

	// DateFormatOrdinalDateTime formats ordinal date and time.
	// Format: yyyy-DDD'T'HH:mm:ss.SSSZ
	DateFormatOrdinalDateTime DateFormat = "ordinal_date_time"

	// DateFormatStrictOrdinalDateTime formats ordinal date and time (strict).
	// Format: yyyy-DDD'T'HH:mm:ss.SSSZ
	DateFormatStrictOrdinalDateTime DateFormat = "strict_ordinal_date_time"

	// DateFormatOrdinalDateTimeNoMillis formats ordinal date and time without millis.
	// Format: yyyy-DDD'T'HH:mm:ssZ
	DateFormatOrdinalDateTimeNoMillis DateFormat = "ordinal_date_time_no_millis"

	// DateFormatStrictOrdinalDateTimeNoMillis formats ordinal date and time without millis (strict).
	// Format: yyyy-DDD'T'HH:mm:ssZ
	DateFormatStrictOrdinalDateTimeNoMillis DateFormat = "strict_ordinal_date_time_no_millis"

	// DateFormatTime formats time with millis and timezone.
	// Format: HH:mm:ss.SSSZ
	DateFormatTime DateFormat = "time"

	// DateFormatStrictTime formats time with millis and timezone (strict).
	// Format: HH:mm:ss.SSSZ
	DateFormatStrictTime DateFormat = "strict_time"

	// DateFormatTimeNoMillis formats time without millis.
	// Format: HH:mm:ssZ
	DateFormatTimeNoMillis DateFormat = "time_no_millis"

	// DateFormatStrictTimeNoMillis formats time without millis (strict).
	// Format: HH:mm:ssZ
	DateFormatStrictTimeNoMillis DateFormat = "strict_time_no_millis"

	// DateFormatTTime formats time prefixed with T.
	// Format: 'T'HH:mm:ss.SSSZ
	DateFormatTTime DateFormat = "t_time"

	// DateFormatStrictTTime formats time prefixed with T (strict).
	// Format: 'T'HH:mm:ss.SSSZ
	DateFormatStrictTTime DateFormat = "strict_t_time"

	// DateFormatTTimeNoMillis formats time prefixed with T without millis.
	// Format: 'T'HH:mm:ssZ
	DateFormatTTimeNoMillis DateFormat = "t_time_no_millis"

	// DateFormatStrictTTimeNoMillis formats time prefixed with T without millis (strict).
	// Format: 'T'HH:mm:ssZ
	DateFormatStrictTTimeNoMillis DateFormat = "strict_t_time_no_millis"

	// DateFormatWeekDate formats a full week date (ISO week-date).
	// Format: YYYY-'W'ww-e
	DateFormatWeekDate DateFormat = "week_date"

	// DateFormatStrictWeekDate formats a full week date (strict, ISO week-date).
	// Format: YYYY-'W'ww-e
	DateFormatStrictWeekDate DateFormat = "strict_week_date"

	// DateFormatWeekDateTime combines week date and time (ISO week-date).
	// Format: YYYY-'W'ww-e'T'HH:mm:ss.SSSZ
	DateFormatWeekDateTime DateFormat = "week_date_time"

	// DateFormatStrictWeekDateTime combines week date and time (strict, ISO week-date).
	// Format: YYYY-'W'ww-e'T'HH:mm:ss.SSSZ
	DateFormatStrictWeekDateTime DateFormat = "strict_week_date_time"

	// DateFormatWeekDateTimeNoMillis combines week date and time without millis (ISO week-date).
	// Format: YYYY-'W'ww-e'T'HH:mm:ssZ
	DateFormatWeekDateTimeNoMillis DateFormat = "week_date_time_no_millis"

	// DateFormatStrictWeekDateTimeNoMillis combines week date and time without millis (strict, ISO week-date).
	// Format: YYYY-'W'ww-e'T'HH:mm:ssZ
	DateFormatStrictWeekDateTimeNoMillis DateFormat = "strict_week_date_time_no_millis"

	// DateFormatWeekyear formats a four digit weekyear (ISO week-date).
	// Format: YYYY
	DateFormatWeekyear DateFormat = "weekyear"

	// DateFormatStrictWeekyear formats a four digit weekyear (strict, ISO week-date).
	// Format: YYYY
	DateFormatStrictWeekyear DateFormat = "strict_weekyear"

	// DateFormatWeekyearWeek formats weekyear and week (ISO week-date).
	// Format: YYYY-'W'ww
	DateFormatWeekyearWeek DateFormat = "weekyear_week"

	// DateFormatStrictWeekyearWeek formats weekyear and week (strict, ISO week-date).
	// Format: YYYY-'W'ww
	DateFormatStrictWeekyearWeek DateFormat = "strict_weekyear_week"

	// DateFormatWeekyearWeekDay formats weekyear, week and day (ISO week-date).
	// Format: YYYY-'W'ww-e
	DateFormatWeekyearWeekDay DateFormat = "weekyear_week_day"

	// DateFormatStrictWeekyearWeekDay formats weekyear, week and day (strict, ISO week-date).
	// Format: YYYY-'W'ww-e
	DateFormatStrictWeekyearWeekDay DateFormat = "strict_weekyear_week_day"

	// DateFormatYear formats a four digit year.
	// Format: yyyy
	DateFormatYear DateFormat = "year"

	// DateFormatStrictYear formats a four digit year (strict).
	// Format: yyyy
	DateFormatStrictYear DateFormat = "strict_year"

	// DateFormatYearMonth formats year and month.
	// Format: yyyy-MM
	DateFormatYearMonth DateFormat = "year_month"

	// DateFormatStrictYearMonth formats year and month (strict).
	// Format: yyyy-MM
	DateFormatStrictYearMonth DateFormat = "strict_year_month"

	// DateFormatYearMonthDay formats year, month and day.
	// Format: yyyy-MM-dd
	DateFormatYearMonthDay DateFormat = "year_month_day"

	// DateFormatStrictYearMonthDay formats year, month and day (strict).
	// Format: yyyy-MM-dd
	DateFormatStrictYearMonthDay DateFormat = "strict_year_month_day"
)

// JoinDateFormats joins multiple DateFormat values with "||" separator
// for use in Elasticsearch date field format configuration.
func JoinDateFormats(formats ...DateFormat) string {
	var result strings.Builder
	for i, f := range formats {
		if i > 0 {
			result.WriteString("||")
		}
		result.WriteString(string(f))
	}
	return result.String()
}
