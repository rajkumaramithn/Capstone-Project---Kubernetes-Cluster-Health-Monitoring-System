package monitor

import (
	"fmt"

	v1 "k8s.io/api/core/v1"

	"Capstone-Project---Kubernetes-Cluster-Health-Monitoring-System/internal/alerts"
)

// CheckPods checks pod and container health and reports unhealthy states.
func CheckPods(pods []v1.Pod) {
	for _, pod := range pods {
		switch pod.Status.Phase {
		case v1.PodFailed:
			alerts.SendAlert(fmt.Sprintf(
				"Pod %s/%s has failed",
				pod.Namespace, pod.Name,
			))
		case v1.PodUnknown:
			alerts.SendAlert(fmt.Sprintf(
				"Pod %s/%s has unknown status",
				pod.Namespace, pod.Name,
			))
		}

		for _, container := range pod.Status.ContainerStatuses {
			if container.RestartCount > 3 {
				alerts.SendAlert(fmt.Sprintf(
					"Container %s in pod %s/%s has restarted %d times",
					container.Name,
					pod.Namespace,
					pod.Name,
					container.RestartCount,
				))
			}

			if waiting := container.State.Waiting; waiting != nil {
				alerts.SendAlert(fmt.Sprintf(
					"Container %s in pod %s/%s is waiting: %s",
					container.Name,
					pod.Namespace,
					pod.Name,
					waiting.Reason,
				))
			}
		}
	}
}
