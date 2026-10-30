package k8sclient

import (
	"context"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestDeleteFailedControllerPod(t *testing.T) {
	controller := true

	tests := []struct {
		name         string
		pod          v1.Pod
		shouldDelete bool
	}{
		{
			name: "failed controller-managed pod",
			pod: v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "failed-pod", Namespace: "default", UID: "pod-1",
					OwnerReferences: []metav1.OwnerReference{
						{APIVersion: "apps/v1", Kind: "ReplicaSet", Name: "test-rs", Controller: &controller},
					},
				},
				Status: v1.PodStatus{Phase: v1.PodFailed},
			},
			shouldDelete: true,
		},
		{
			name: "failed standalone pod",
			pod: v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "standalone-pod", Namespace: "default", UID: "pod-2",
				},
				Status: v1.PodStatus{Phase: v1.PodFailed},
			},
			shouldDelete: false,
		},
		{
			name: "running controller-managed pod",
			pod: v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "running-pod", Namespace: "default", UID: "pod-3",
					OwnerReferences: []metav1.OwnerReference{
						{APIVersion: "apps/v1", Kind: "ReplicaSet", Name: "test-rs", Controller: &controller},
					},
				},
				Status: v1.PodStatus{Phase: v1.PodRunning},
			},
			shouldDelete: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientset := fake.NewSimpleClientset(&tt.pod)
			client := &Client{set: clientset}

			deleted, err := client.DeleteFailedControllerPod(context.Background(), tt.pod)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if deleted != tt.shouldDelete {
				t.Fatalf("deleted = %v, want %v", deleted, tt.shouldDelete)
			}

			deleteActions := 0
			for _, action := range clientset.Actions() {
				if action.GetVerb() == "delete" && action.GetResource().Resource == "pods" {
					deleteActions++
				}
			}

			if tt.shouldDelete && deleteActions != 1 {
				t.Errorf("delete actions = %d, want 1", deleteActions)
			}
			if !tt.shouldDelete && deleteActions != 0 {
				t.Errorf("delete actions = %d, want 0", deleteActions)
			}
		})
	}
}
