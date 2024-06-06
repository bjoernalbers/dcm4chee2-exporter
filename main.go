// dcm4chee2-exporter
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const metricsPath = "/metrics"

var (
	version  = "unset" // version gets set via build flags
	homepage = "https://github.com/bjoernalbers/dcm4chee2-exporter"

	// All metrics collected from JMX server.
	dcm4chee2Up                           = prometheus.NewDesc("dcm4chee2_up", "Availability of Dcm4chee 2.", []string{}, nil)
	moveScuQueueMessageCountDesc          = prometheus.NewDesc("dcm4chee2_movescu_queue_message_count", "The number of messages in the queue.", []string{}, nil)
	moveScuQueueDeliveringCountDesc       = prometheus.NewDesc("dcm4chee2_movescu_queue_delivering_count", "The number of messages currently being delivered.", []string{}, nil)
	moveScuQueueScheduledMessageCountDesc = prometheus.NewDesc("dcm4chee2_movescu_queue_scheduled_message_count", "The number of scheduled messages in the queue.", []string{}, nil)
	scrapeDurationSeconds                 = prometheus.NewDesc("dcm4chee2_scrape_duration_seconds", "Duration of backend response in seconds.", []string{}, nil)
)

func init() {
	flag.Usage = usage
}

func usage() {
	header := fmt.Sprintf(`Dcm4chee 2 Prometheus Exporter (version %s)

This exporter exposes metrics from (deprecated) dcm4chee 2 for Prometheus.

Usage:
`, version)
	footer := fmt.Sprintf(`
Homepage: %s
`, homepage)
	fmt.Fprintf(flag.CommandLine.Output(), header)
	flag.PrintDefaults()
	fmt.Fprintf(flag.CommandLine.Output(), footer)
}

// jmxServer provides access to a JMX server.
type jmxServer struct {
	Script   string
	Username string
	Password string
}

// Fetch returns metrics from the JMX server.
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

// collector collects metrics for Prometheus.
type collector struct {
	jmx jmxServer
}

// Describe metrics.
func (c collector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

// Collect metrics.
func (c collector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	output, err := c.jmx.Fetch()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(dcm4chee2Up, prometheus.GaugeValue, float64(0))
		return
	}
	ch <- prometheus.MustNewConstMetric(scrapeDurationSeconds, prometheus.GaugeValue, float64(time.Since(start).Seconds()))
	ch <- prometheus.MustNewConstMetric(dcm4chee2Up, prometheus.GaugeValue, float64(1))
	metrics := Translate(output)
	if m, ok := metrics["MessageCount"]; ok {
		ch <- prometheus.MustNewConstMetric(moveScuQueueMessageCountDesc, prometheus.GaugeValue, float64(m))
	}
	if m, ok := metrics["DeliveringCount"]; ok {
		ch <- prometheus.MustNewConstMetric(moveScuQueueDeliveringCountDesc, prometheus.GaugeValue, float64(m))
	}
	if m, ok := metrics["ScheduledMessageCount"]; ok {
		ch <- prometheus.MustNewConstMetric(moveScuQueueScheduledMessageCountDesc, prometheus.GaugeValue, float64(m))
	}
}

// Translate metrics from JMX server into map
func Translate(in []byte) map[string]int {
	metrics := map[string]int{}
	scanner := bufio.NewScanner(bytes.NewReader(in))
	for scanner.Scan() {
		key, value, found := strings.Cut(scanner.Text(), "=")
		if !found {
			continue
		}
		metric, err := strconv.Atoi(value)
		if err != nil {
			continue
		}
		metrics[key] = metric
	}
	return metrics
}

func main() {
	jmx := jmxServer{}
	flag.StringVar(&jmx.Script, "s", "twiddle.sh", "Path to JBoss twiddle script")
	flag.StringVar(&jmx.Username, "u", "admin", "Username of JBoss admin")
	flag.StringVar(&jmx.Password, "p", "admin", "Password of JBoss admin")
	address := flag.String("a", ":9404", "Address to listen on")
	flag.Parse()
	prometheus.MustRegister(&collector{jmx})
	http.Handle(metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
	<head>
		<title>Dcm4chee 2 Prometheus Exporter</title>
	</head>
	<body>
		<h1>Dcm4chee 2 Prometheus Exporter</h1>
		<p><a href='` + metricsPath + `'>Metrics</a></p>
	</body>
</html>
`))
	})
	log.Fatal(http.ListenAndServe(*address, nil))
}
