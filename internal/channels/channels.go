package channels

import "time"

// InfluxDBMeasurement is a channel used for moving InfluxDB
type InfluxDBMeasurement struct {
	Measurement string
	TagSet      map[string]string
	FieldSet    map[string]interface{}
	TimeStamp   time.Time
}
