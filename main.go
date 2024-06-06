// dcm4chee2-exporter
package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var moveScuQueueMessageCountDesc = prometheus.NewDesc("dcm4chee2_movescu_queue_message_count", "The number of messages in the queue", []string{}, nil)

type collector struct{}

func (c collector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func (c collector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(moveScuQueueMessageCountDesc, prometheus.GaugeValue, float64(42))
}

func main() {
	prometheus.MustRegister(&collector{})
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9404", nil))
}
