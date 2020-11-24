package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ActionDuration prometheus.ObserverVec
	TasksTotal     prometheus.Counter
	TasksActive    prometheus.Gauge
)

func init() {
	ActionDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "pbnj_action_duration_seconds",
		Help:    "Duration taken to complete an action.",
		Buckets: []float64{0.2, 0.5, 0.7, 1, 2, 5, 10, 15, 20, 25, 30},
	}, []string{"service", "action"})

	labelValues := []prometheus.Labels{
		{"service": "bmc", "action": "create_user"},
		{"service": "bmc", "action": "update_user"},
		{"service": "bmc", "action": "delete_user"},
		{"service": "machine", "action": "boot_device"},
		{"service": "machine", "action": "power"},
	}

	initObserverLabels(ActionDuration, labelValues)

	TasksTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pbnj_tasks_total",
		Help: "Total number of tasks executed.",
	})
	TasksActive = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "pbnj_tasks_active",
		Help: "Number of tasks currently active.",
	})
}

func initObserverLabels(m prometheus.ObserverVec, l []prometheus.Labels) {
	for _, labels := range l {
		m.With(labels)
	}
}
