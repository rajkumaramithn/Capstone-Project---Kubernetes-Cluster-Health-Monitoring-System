package main

import (
	"fmt"
	"time"

	"Capstone-Project---Kubernetes-Cluster-Health-Monitoring-System/internal/k8sclient"
	"Capstone-Project---Kubernetes-Cluster-Health-Monitoring-System/internal/monitor"
)

func main() {
	fmt.Println("Kubernetes Health Monitor Starting...")

	client := k8sclient.NewClient()

	for {
		pods := client.GetPods("default")

		monitor.CheckPods(pods)

		time.Sleep(10 * time.Second)
	}
}