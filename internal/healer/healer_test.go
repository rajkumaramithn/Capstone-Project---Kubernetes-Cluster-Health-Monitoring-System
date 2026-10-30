package healer

import (
	"context"
	"errors"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type fakePodHealer struct {
	calls   int
	deleted bool
	err     error
	lastPod v1.Pod
}

func (f *fakePodHealer) DeleteFailedControllerPod(
	ctx context.Context,
	pod v1.Pod,
) (bool, error) {
	f.calls++
	f.lastPod = pod
	return f.deleted, f.err
}

func TestHealFailedPodsRequestsRecoveryForFailedPod(t *testing.T) {
	client := &fakePodHealer{deleted: true}
	pod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "failed-pod",
			Namespace: "default",
		},
		Status: v1.PodStatus{Phase: v1.PodFailed},
	}

	HealFailedPods(client, []v1.Pod{pod})

	if client.calls != 1 {
		t.Fatalf("expected 1 recovery request, got %d", client.calls)
	}
	if client.lastPod.Name != "failed-pod" {
		t.Fatalf("unexpected pod: %s", client.lastPod.Name)
	}
}

func TestHealFailedPodsSkipsHealthyPod(t *testing.T) {
	client := &fakePodHealer{}
	pod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "healthy-pod"},
		Status:     v1.PodStatus{Phase: v1.PodRunning},
	}

	HealFailedPods(client, []v1.Pod{pod})

	if client.calls != 0 {
		t.Fatalf("expected no recovery requests, got %d", client.calls)
	}
}

func TestHealFailedPodsContinuesAfterError(t *testing.T) {
	client := &fakePodHealer{
		err: errors.New("simulated deletion error"),
	}
	pods := []v1.Pod{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "failed-one"},
			Status:     v1.PodStatus{Phase: v1.PodFailed},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "failed-two"},
			Status:     v1.PodStatus{Phase: v1.PodFailed},
		},
	}

	HealFailedPods(client, pods)

	if client.calls != 2 {
		t.Fatalf("expected recovery to be attempted for both pods, got %d calls", client.calls)
	}
}
