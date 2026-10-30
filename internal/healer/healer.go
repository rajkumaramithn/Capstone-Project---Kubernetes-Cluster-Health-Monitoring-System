package healer

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
)

type podHealer interface {
	DeleteFailedControllerPod(
		ctx context.Context,
		pod v1.Pod,
	) (bool, error)
}

// HealFailedPods requests recovery for failed, controller-managed pods.
func HealFailedPods(client podHealer, pods []v1.Pod) {
	for _, pod := range pods {
		if pod.Status.Phase != v1.PodFailed {
			continue
		}

		deleted, err := client.DeleteFailedControllerPod(
			context.Background(),
			pod,
		)
		if err != nil {
			fmt.Printf(
				"HEAL ERROR: Pod %s/%s: %v\n",
				pod.Namespace,
				pod.Name,
				err,
			)
			continue
		}

		if deleted {
			fmt.Printf(
				"HEAL: Deleted failed controller-managed pod %s/%s; controller may recreate it\n",
				pod.Namespace,
				pod.Name,
			)
		}
	}
}
