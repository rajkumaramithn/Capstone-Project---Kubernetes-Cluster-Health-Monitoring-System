package monitor

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	v1 "k8s.io/api/core/v1"
)

var (
	podCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kubernetes_health_monitor_pods",
			Help: "Number of observed pods by namespace and phase.",
		},
		[]string{"namespace", "phase"},
	)

	containerRestarts = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kubernetes_health_monitor_container_restarts",
			Help: "Container restart counts reported by Kubernetes.",
		},
		[]string{"namespace", "pod", "container"},
	)

	unhealthyPods = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "kubernetes_health_monitor_unhealthy_pods",
			Help: "Number of observed pods considered unhealthy.",
		},
	)
)

func init() {
	prometheus.MustRegister(podCount)
	prometheus.MustRegister(containerRestarts)
	prometheus.MustRegister(unhealthyPods)
}

// UpdateMetrics refreshes metrics using the latest pod information.
func UpdateMetrics(pods []v1.Pod) {
	podCount.Reset()
	containerRestarts.Reset()

	unhealthy := 0

	for _, pod := range pods {
		phase := string(pod.Status.Phase)
		if phase == "" {
			phase = "Unknown"
		}

		podCount.WithLabelValues(pod.Namespace, phase).Inc()

		isUnhealthy := pod.Status.Phase == v1.PodFailed ||
			pod.Status.Phase == v1.PodUnknown

		for _, container := range pod.Status.ContainerStatuses {
			containerRestarts.WithLabelValues(
				pod.Namespace,
				pod.Name,
				container.Name,
			).Set(float64(container.RestartCount))

			if container.RestartCount > 3 || container.State.Waiting != nil {
				isUnhealthy = true
			}
		}

		if isUnhealthy {
			unhealthy++
		}
	}

	unhealthyPods.Set(float64(unhealthy))
}

// MetricsHandler returns the HTTP handler used by Prometheus to scrape metrics.
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
