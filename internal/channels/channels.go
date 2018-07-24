package channels

// Response is a channel used for moving data between the Kafka code and collectors
type Response struct {
	Data  string
	Topic string
}
