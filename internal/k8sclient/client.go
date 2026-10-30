package k8sclient

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	set kubernetes.Interface
}

func NewClient() *Client {
	config, err := rest.InClusterConfig()
	if err != nil {
		// Not running inside a cluster — fall back to local kubeconfig
		kubeconfig := clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return &Client{set: clientset}
}

func (c *Client) GetPods(namespace string) []v1.Pod {
	pods, err := c.set.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error fetching pods:", err)
		return nil
	}

	return pods.Items
}
func (c *Client) GetNodes() []v1.Node {
	nodes, err := c.set.CoreV1().Nodes().List(
		context.Background(),
		metav1.ListOptions{},
	)
	if err != nil {
		fmt.Println("Error fetching nodes:", err)
		return nil
	}

	return nodes.Items
}
func (c *Client) DeleteFailedControllerPod(ctx context.Context, pod v1.Pod) (bool, error) {
	// Re-read the pod to avoid acting on stale information.
	current, err := c.set.CoreV1().Pods(pod.Namespace).Get(
		ctx,
		pod.Name,
		metav1.GetOptions{},
	)
	if err != nil {
		return false, err
	}

	// Confirm that this is still the same pod.
	if current.UID != pod.UID {
		return false, nil
	}

	// Only heal pods that Kubernetes has marked Failed.
	if current.Status.Phase != v1.PodFailed {
		return false, nil
	}

	// Never delete a standalone pod.
	controllerManaged := false
	for _, owner := range current.OwnerReferences {
		if owner.Controller != nil && *owner.Controller {
			controllerManaged = true
			break
		}
	}
	if !controllerManaged {
		return false, nil
	}

	// Use a UID precondition to avoid deleting a replacement pod.
	uid := current.UID
	err = c.set.CoreV1().Pods(current.Namespace).Delete(
		ctx,
		current.Name,
		metav1.DeleteOptions{
			Preconditions: &metav1.Preconditions{UID: &uid},
		},
	)
	if err != nil {
		return false, err
	}

	return true, nil
}
