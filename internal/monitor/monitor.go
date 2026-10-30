package monitor

import (
	"fmt"

	"Capstone-Project---Kubernetes-Cluster-Health-Monitoring-System/internal/alerts"
	v1 "k8s.io/api/core/v1"
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

// CheckNodes checks node readiness and resource-pressure conditions.
func CheckNodes(nodes []v1.Node) {
	for _, node := range nodes {
		for _, condition := range node.Status.Conditions {
			switch condition.Type {
			case v1.NodeReady:
				if condition.Status != v1.ConditionTrue {
					alerts.SendAlert(fmt.Sprintf(
						"Node %s is not ready: %s",
						node.Name,
						condition.Message,
					))
				}
			case v1.NodeMemoryPressure, v1.NodeDiskPressure, v1.NodePIDPressure:
				if condition.Status == v1.ConditionTrue {
					alerts.SendAlert(fmt.Sprintf(
						"Node %s has resource pressure: %s",
						node.Name,
						condition.Type,
					))
				}
			}
		}
	}
}
