package prometheus

import "github.com/prometheus/client_golang/prometheus"

func NewGauge(namespace, subsystem, help, name string) prometheus.Gauge {
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Help:      help,
		// Namespace 和 Subsystem 和 Name 都不能有 _ 以外的其它符号
		Name: name,
		//ConstLabels: map[string]string{
		//	"instance_id": b.InstanceId,
		//},
	})
	prometheus.MustRegister(gauge)
	gauge.Set(0)
	return gauge
}
