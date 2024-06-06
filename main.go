// dcm4chee2-exporter
package main

import (
	"flag"
	"log"
	"net/http"
	"os/exec"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	moveScuQueueMessageCountDesc = prometheus.NewDesc("dcm4chee2_movescu_queue_message_count", "The number of messages in the queue.", []string{}, nil)
	dcm4chee2Up                  = prometheus.NewDesc("dcm4chee2_up", "Availability of Dcm4chee 2.", []string{}, nil)
)

type collector struct {
	jmx jmxServer
}

func (c collector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func (c collector) Collect(ch chan<- prometheus.Metric) {
	_, err := c.jmx.Fetch()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(dcm4chee2Up, prometheus.GaugeValue, float64(0))
		return
	}
	ch <- prometheus.MustNewConstMetric(dcm4chee2Up, prometheus.GaugeValue, float64(1))
	// TODO: Un-dummy this metric!
	ch <- prometheus.MustNewConstMetric(moveScuQueueMessageCountDesc, prometheus.GaugeValue, float64(42))
}

type jmxServer struct {
	Script   string
	Username string
	Password string
}

// Fetch metrics from JMX server
func (j jmxServer) Fetch() ([]byte, error) {
	return exec.Command(
		j.Script,
		"-u",
		j.Username,
		"-p",
		j.Password,
		"get",
		"dcm4chee.archive:name=MoveScu,service=Queue",
		"MessageCount",
		"DeliveringCount",
		"ScheduledMessageCount",
	).Output()
}

func main() {
	jmx := jmxServer{}
	flag.StringVar(&jmx.Script, "s", "twiddle.sh", "Path to JBoss twiddle script")
	flag.StringVar(&jmx.Username, "u", "admin", "Username of JBoss admin")
	flag.StringVar(&jmx.Password, "p", "admin", "Password of JBoss admin")
	flag.Parse()
	prometheus.MustRegister(&collector{jmx})
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9404", nil))
}
