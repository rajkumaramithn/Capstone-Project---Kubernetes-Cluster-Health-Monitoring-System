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

	// Allow a graceful shutdown when the process receives a termination signal.
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
		pods := client.GetPods("default")
		nodes := client.GetNodes()

		monitor.CheckPods(pods)
		monitor.CheckNodes(nodes)
		monitor.UpdateMetrics(pods)
		monitor.UpdateNodeMetrics(nodes)
		healer.HealFailedPods(client, pods)

		select {
		case <-ctx.Done():
			log.Println("Shutting down Kubernetes Health Monitor...")
			shutdownCtx, cancel := context.WithTimeout(
				context.Background(),
				5*time.Second,
			)
			defer cancel()

			if err := server.Shutdown(shutdownCtx); err != nil {
				log.Printf("Metrics server shutdown error: %v", err)
			}
			return
		case <-ticker.C:
		}
	}
}
