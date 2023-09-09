package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ActionDuration   prometheus.ObserverVec
	TasksTotal       prometheus.Counter
	TotalGauge       prometheus.Gauge
	TasksActive      prometheus.Gauge
	PerIDQueue       prometheus.GaugeVec
	IngestionQueue   prometheus.Gauge
	Ingested         prometheus.Gauge
	FIFOQueue        prometheus.Gauge
	NumWorkers       prometheus.Gauge
	NumPerIDEnqueued prometheus.Gauge
	WorkerMap        prometheus.Gauge
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
		Name: "pbnj_tasks_processes",
		Help: "Total number of tasks executed.",
	})
	TasksActive = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "pbnj_tasks_active",
		Help: "Number of tasks currently active.",
	})
	PerIDQueue = *promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pbnj_per_id_queue",
		Help: "Number of tasks in perID queue.",
	}, []string{"host"})
	IngestionQueue = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "pbnj_ingestion_queue",
		Help: "Number of tasks in ingestion queue.",
	})
	Ingested = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "pbnj_ingested",
		Help: "Number of tasks ingested.",
	})
	FIFOQueue = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "pbnj_fifo_queue",
		Help: "Number of tasks in FIFO queue.",
	})
	TotalGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "pbnj_total",
		Help: "Total number of tasks.",
	})
	NumWorkers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "pbnj_num_workers",
		Help: "Number of workers.",
	})
	NumPerIDEnqueued = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "pbnj_num_per_id_enqueued",
		Help: "Number of perID enqueued.",
	})
	WorkerMap = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "pbnj_worker_map_size",
		Help: "Worker map size.",
	})
}

func initObserverLabels(m prometheus.ObserverVec, l []prometheus.Labels) {
	for _, labels := range l {
		m.With(labels)
	}
}
