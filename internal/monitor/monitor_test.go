package monitor

import (
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCheckPodsDoesNotPanicOnHealthyPod(t *testing.T) {
	pods := []v1.Pod{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "healthy-pod",
				Namespace: "default",
			},
			Status: v1.PodStatus{
				Phase: v1.PodRunning,
				ContainerStatuses: []v1.ContainerStatus{
					{
						Name:         "app",
						RestartCount: 0,
						State: v1.ContainerState{
							Running: &v1.ContainerStateRunning{},
						},
					},
				},
			},
		},
	}

	CheckPods(pods)
}

func TestCheckPodsDoesNotPanicOnUnhealthyPods(t *testing.T) {
	pods := []v1.Pod{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "failed-pod",
				Namespace: "default",
			},
			Status: v1.PodStatus{
				Phase: v1.PodFailed,
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "restarting-pod",
				Namespace: "default",
			},
			Status: v1.PodStatus{
				Phase: v1.PodRunning,
				ContainerStatuses: []v1.ContainerStatus{
					{
						Name:         "app",
						RestartCount: 4,
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "waiting-pod",
				Namespace: "default",
			},
			Status: v1.PodStatus{
				Phase: v1.PodPending,
				ContainerStatuses: []v1.ContainerStatus{
					{
						Name: "app",
						State: v1.ContainerState{
							Waiting: &v1.ContainerStateWaiting{
								Reason: "CrashLoopBackOff",
							},
						},
					},
				},
			},
		},
	}

	CheckPods(pods)
}
