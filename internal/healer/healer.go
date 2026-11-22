package healer

import (
	"context"
	"log"

	"Capstone-Project---Kubernetes-Cluster-Health-Monitoring-System/internal/k8sclient"
	"Capstone-Project---Kubernetes-Cluster-Health-Monitoring-System/internal/monitor"

	v1 "k8s.io/api/core/v1"
)

type podHealer interface {
	DeleteFailedControllerPod(
		ctx context.Context,
		pod v1.Pod,
	) (bool, error)
}

// HealFailedPods requests recovery for failed, evicted or crash-looping
// controller-managed pods. The owning controller recreates each deleted pod.
func HealFailedPods(client podHealer, pods []v1.Pod) {
	for _, pod := range pods {
		reason := k8sclient.UnhealthyReason(&pod)
		if reason == "" {
			continue
		}

		deleted, err := client.DeleteFailedControllerPod(
			context.Background(),
			pod,
		)
		if err != nil {
			log.Printf("HEAL ERROR ns=%s pod=%s err=%v", pod.Namespace, pod.Name, err)
			continue
		}

		if deleted {
			log.Printf("HEAL action=delete_pod ns=%s pod=%s reason=%s", pod.Namespace, pod.Name, reason)
			monitor.RecordHealingAction(pod.Namespace, reason)
		}
	}
}
