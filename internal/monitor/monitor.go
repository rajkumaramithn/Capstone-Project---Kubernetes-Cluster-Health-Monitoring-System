package monitor

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
)

func CheckPods(pods []v1.Pod) {
	for _, pod := range pods {
		for _, c := range pod.Status.ContainerStatuses {
			if c.RestartCount > 3 || c.State.Waiting != nil {
				fmt.Printf("🚨 ALERT: Pod %s is unhealthy\n", pod.Name)
			}
		}
	}
}
