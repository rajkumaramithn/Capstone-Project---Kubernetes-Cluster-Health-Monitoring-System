package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Capstone-Project---Kubernetes-Cluster-Health-Monitoring-System/internal/healer"
	"Capstone-Project---Kubernetes-Cluster-Health-Monitoring-System/internal/k8sclient"
	"Capstone-Project---Kubernetes-Cluster-Health-Monitoring-System/internal/monitor"
)

func main() {
	fmt.Println("Kubernetes Health Monitor Starting...")

	// Start the Prometheus metrics endpoint.
	mux := http.NewServeMux()
	mux.Handle("/metrics", monitor.MetricsHandler())

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Println("Metrics endpoint listening on :8080/metrics")
		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Printf("Metrics server error: %v", err)
		}
	}()

	// Allow graceful shutdown.
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	client := k8sclient.NewClient()
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		// Retrieve Kubernetes pods and nodes.
		pods := client.GetPods("default")
		nodes := client.GetNodes()

		// Retrieve node CPU and memory usage.
		nodeMetrics, err := client.GetNodeMetricsData()
		if err != nil {
			log.Printf("Error fetching node metrics: %v", err)
		} else {
			// Calculate resource utilization and update Prometheus metrics.
			monitor.UpdateNodeResourceMetrics(nodes, nodeMetrics)
		}

		// Run existing health monitoring.
		monitor.CheckPods(pods)
		monitor.CheckNodes(nodes)

		// Update Prometheus monitoring metrics.
		monitor.UpdateMetrics(pods)
		monitor.UpdateNodeMetrics(nodes)

		// Run existing pod self-healing mechanism.
		healer.HealFailedPods(client, pods)

		select {
		case <-ctx.Done():
			log.Println("Shutting down Kubernetes Health Monitor...")

			shutdownCtx, cancel := context.WithTimeout(
				context.Background(),
				5*time.Second,
			)

			if err := server.Shutdown(shutdownCtx); err != nil {
				log.Printf("Metrics server shutdown error: %v", err)
			}

			cancel()
			return

		case <-ticker.C:
		}
	}
}
