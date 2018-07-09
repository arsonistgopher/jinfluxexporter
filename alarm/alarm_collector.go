package alarm

import (
	// "net/rpc"

	"fmt"
	"regexp"

	"github.com/arsonistgopher/jkafkaexporter/rpc"

	"github.com/arsonistgopher/jkafkaexporter/collector"
)

const prefix = "junos_alarms_"

// Counter struct
type Counter struct {
	YellowCount float64
	RedCount    float64
}

type alarmCollector struct {
	filter *regexp.Regexp
}

// NewCollector creates a new collector
func NewCollector(alarmsFilter string) collector.RPCCollector {
	c := &alarmCollector{}

	if len(alarmsFilter) > 0 {
		c.filter = regexp.MustCompile(alarmsFilter)
	}

	return c
}

// Collect collects metrics from JunOS
func (c *alarmCollector) Collect(client *rpc.Client, ch chan<- string, label string) error {
	// func (c *alarmCollector) Collect(client *rpc.Client, ch chan<- prometheus.Metric, labelValues []string) error {
	counter, err := c.Counter(client)
	if err != nil {
		return err
	}

	jsonReturn := "{Node: %s, Status: {RedAlarm: %f, YellowAlarm: %f}}"
	ch <- fmt.Sprintf(jsonReturn, label, counter.RedCount, counter.YellowCount)

	return nil
}

func (c *alarmCollector) Counter(client *rpc.Client) (*Counter, error) {
	red := 0
	yellow := 0

	cmds := []string{
		"<get-system-alarm-information/>",
		"<get-alarm-information/>",
	}

	for _, cmd := range cmds {
		a := &AlarmRPC{}
		err := client.RunCommandAndParse(cmd, &a)
		if err != nil {
			return nil, err
		}

		for _, d := range a.Details {
			if c.shouldFilterAlarm(&d) {
				continue
			}

			if d.Class == "Major" {
				red++
			} else if d.Class == "Minor" {
				yellow++
			}
		}
	}

	return &Counter{RedCount: float64(red), YellowCount: float64(yellow)}, nil
}

func (c *alarmCollector) shouldFilterAlarm(a *AlarmDetails) bool {
	if c.filter == nil {
		return false
	}

	return c.filter.MatchString(a.Description) || c.filter.MatchString(a.Type)
}
