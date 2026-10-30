package monitor

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
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

func TestCheckNodesDoesNotPanicOnHealthyNode(t *testing.T) {
	nodes := []v1.Node{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "healthy-node",
			},
			Status: v1.NodeStatus{
				Conditions: []v1.NodeCondition{
					{
						Type:   v1.NodeReady,
						Status: v1.ConditionTrue,
					},
					{
						Type:   v1.NodeMemoryPressure,
						Status: v1.ConditionFalse,
					},
					{
						Type:   v1.NodeDiskPressure,
						Status: v1.ConditionFalse,
					},
					{
						Type:   v1.NodePIDPressure,
						Status: v1.ConditionFalse,
					},
				},
			},
		},
	}

	CheckNodes(nodes)
}

func TestCheckNodesDoesNotPanicOnUnhealthyNode(t *testing.T) {
	nodes := []v1.Node{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "unhealthy-node",
			},
			Status: v1.NodeStatus{
				Conditions: []v1.NodeCondition{
					{
						Type:    v1.NodeReady,
						Status:  v1.ConditionFalse,
						Message: "Kubelet is not responding",
					},
					{
						Type:   v1.NodeMemoryPressure,
						Status: v1.ConditionTrue,
					},
				},
			},
		},
	}

	CheckNodes(nodes)
}

func TestUpdateNodeMetricsForHealthyNode(t *testing.T) {
	nodes := []v1.Node{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "test-healthy-node"},
			Status: v1.NodeStatus{
				Conditions: []v1.NodeCondition{
					{Type: v1.NodeReady, Status: v1.ConditionTrue},
					{Type: v1.NodeMemoryPressure, Status: v1.ConditionFalse},
					{Type: v1.NodeDiskPressure, Status: v1.ConditionFalse},
					{Type: v1.NodePIDPressure, Status: v1.ConditionFalse},
				},
			},
		},
	}

	UpdateNodeMetrics(nodes)

	if got := testutil.ToFloat64(
		nodeCount.WithLabelValues("test-healthy-node", "Ready"),
	); got != 1 {
		t.Errorf("Ready node metric = %v, want 1", got)
	}

	if got := testutil.ToFloat64(
		nodePressure.WithLabelValues("test-healthy-node", string(v1.NodeMemoryPressure)),
	); got != 0 {
		t.Errorf("Memory pressure metric = %v, want 0", got)
	}
}

func TestUpdateNodeMetricsForUnhealthyNode(t *testing.T) {
	nodes := []v1.Node{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "test-unhealthy-node"},
			Status: v1.NodeStatus{
				Conditions: []v1.NodeCondition{
					{Type: v1.NodeReady, Status: v1.ConditionFalse},
					{Type: v1.NodeMemoryPressure, Status: v1.ConditionTrue},
				},
			},
		},
	}

	UpdateNodeMetrics(nodes)

	if got := testutil.ToFloat64(
		nodeCount.WithLabelValues("test-unhealthy-node", "NotReady"),
	); got != 1 {
		t.Errorf("NotReady node metric = %v, want 1", got)
	}

	if got := testutil.ToFloat64(
		nodePressure.WithLabelValues("test-unhealthy-node", string(v1.NodeMemoryPressure)),
	); got != 1 {
		t.Errorf("Memory pressure metric = %v, want 1", got)
	}
}
