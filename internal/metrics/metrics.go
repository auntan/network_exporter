package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Config struct {
	RTTHistogramBuckets []float64
}

var (
	RTTProbe  *prometheus.HistogramVec
	RTTHost   *prometheus.HistogramVec
	SentProbe *prometheus.CounterVec
	SentHost  *prometheus.CounterVec
	RecvProbe *prometheus.CounterVec
	RecvHost  *prometheus.CounterVec
)

func Initialize(conf *Config) {
	RTTProbe = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "network_exporter",
		Name:      "rtt_probe",
		Help:      "round trip time by probes",
		Buckets:   conf.RTTHistogramBuckets,
	}, []string{"probe", "host_deployment", "host", "target_deployment", "target_host", "target_address"})

	RTTHost = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "network_exporter",
		Name:      "rtt_host",
		Help:      "round trip time by hosts",
		Buckets:   conf.RTTHistogramBuckets,
	}, []string{"host", "target_host", "target_address"})

	SentProbe = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "network_exporter",
		Name:      "sent_probe",
		Help:      "total packets sent by probes",
	}, []string{"probe", "host_deployment", "host", "target_deployment", "target_host", "target_address"})

	SentHost = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "network_exporter",
		Name:      "sent_host",
		Help:      "total packets sent by hosts",
	}, []string{"host", "target_host", "target_address"})

	RecvProbe = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "network_exporter",
		Name:      "recv_probe",
		Help:      "total packets received by probes",
	}, []string{"probe", "host_deployment", "host", "target_deployment", "target_host", "target_address"})

	RecvHost = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "network_exporter",
		Name:      "recv_host",
		Help:      "total packets received by hosts",
	}, []string{"host", "target_host", "target_address"})
}
