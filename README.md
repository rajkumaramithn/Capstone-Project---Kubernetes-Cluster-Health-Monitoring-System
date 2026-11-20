# Capstone-Project---Kubernetes-Cluster-Health-Monitoring-System
Capstone Project by Hero Vired  - Authors : rajkumaramithn , dhruvkapoor604-collab &amp; hitesh-1984-pardeshi


# Kubernetes Cluster Health Checker and Auto-Healing

An automated health monitoring and self-healing tool for Kubernetes clusters. It detects common issues — failed pods, unresponsive nodes, resource pressure — and takes corrective action automatically, reducing the manual overhead for small DevOps teams managing large-scale clusters.

## Problem statement

Manually monitoring and troubleshooting Kubernetes clusters is time-consuming and error-prone, especially for small teams managing large deployments. This project automates detection and recovery of common failure modes to improve availability and reduce operator burden.

## Goals

1. Automated health monitoring of node health, pod status, and resource utilization.
2. Self-healing actions: pod restarts, rescheduling, and scaling to balance workloads.
3. Real-time alerts and notifications for issues requiring manual intervention.
4. A web dashboard showing live health status, historical data, and auto-healing logs.

## Tech stack

| Layer | Tool | Purpose |
|---|---|---|
| Core service | Go | Kubernetes API interaction, performance-critical logic |
| Scripting | Python | Data processing, glue scripts |
| Cluster control | Kubernetes API (client-go) | Read/write access to nodes, pods, deployments |
| Metrics | Prometheus | Cluster and application metric collection |
| Visualization | Grafana | Real-time and historical dashboards |
| Alerting | Alertmanager | Threshold-based alert routing |
| Notifications | Slack API | Delivering alerts to the DevOps team |
| Packaging | Docker | Containerizing the service for deployment |

## Architecture

<img src="assets/tech_architecture.png" alt="Project Architectural Image" width="100%" height="100%">


## Kubernetes Cluster (Nodes & Pods)

The environment being monitored by the system. Nodes represent the machines running the cluster and report health conditions such as readiness and resource pressure, while Pods represent the workloads scheduled on them and report status such as running, pending, or failed. This layer is the foundation the entire architecture is built around.

## Health Checker (Go + Python, Docker)

The core application of the project, responsible for watching cluster state and triggering self-healing actions when issues are detected. Go handles direct interaction with the Kubernetes API for performance, while Python supports auxiliary scripting, and Docker packages it for deployment. It connects directly to the cluster via the Kubernetes API and is the only component that actively modifies cluster state.

## Prometheus

The metrics collection layer that scrapes health and resource data from the cluster and the Health Checker at regular intervals, storing it as time-series data. It sits at the center of the observability stack, feeding data to both the visualization and alerting components.

## Grafana

The visualization layer that queries Prometheus to render real-time and historical dashboards of cluster health and auto-healing activity. It connects only to Prometheus and provides the primary interface for reviewing system status.

## Alertmanager

The alert routing layer that evaluates rules against Prometheus data and manages deduplication, grouping, and routing of alerts. It receives triggered alerts from Prometheus and forwards them to the appropriate notification channel based on severity.

## Slack

The notification endpoint that delivers alerts routed through Alertmanager. It connects to the end of the pipeline, surfacing critical issues that may require manual attention.


## Project Execution

## Sprint 1 - Project Setup and Kubernetes Cluster Access

Project File Structure Definition 

<img src="assets/folder_structure.png" alt="Project Folder Architecture Image" width="100%" height="100%">

## Step 1 - Initializing the Go Project

<b>On the Command Line Terminal Run :</b>

go mod init Capstone-Project---Kubernetes-Cluster-Health-Monitoring-System
go get k8s.io/client-go@latest
go get k8s.io/apimachinery@latest
go get k8s.io/api@latest
go get github.com/gin-gonic/gin


<b>File Structure Responsibilities</b>

📁 cmd/healer/main.go
</br><i>Entry point (starts everything)</i>

📁 internal/k8sclient
</br><i>Talks to Kubernetes (list pods, get status)</i>

📁 internal/monitor
</br><i>Checks health rules (CrashLoopBackOff detection)</i>

📁 internal/alerts
</br><i>Prints/logs alerts</i>

📁 internal/healer
</br><i>auto-fix logic</i>


## Step 2 - Creating and Setting up the Go Files 

</br><b>Create MAIN ENTRY (cmd/healer/main.go)</b>

Refer to: cmd\healer\main.go

Starts a Kubernetes health monitoring loop that continuously fetches pods from the “default” namespace and checks their status every 10 seconds.
It connects to the cluster using k8sclient, then passes pod data to monitor.CheckPods() for health evaluation and detection of issues.


</br><b>Kubernetes Client (internal/k8sclient)</b>

Refer to: internal\k8sclient\client.go

Creates a Kubernetes client in Go that connects either to a local kubeconfig (Minikube) or in-cluster config, then fetches all pods from a given namespace.
GetPods() uses the Kubernetes API to list pods and returns them so your monitor can check their health status.


</br><b>Monitor Logic</b>

Refer to: internal\monitor\monitor.go

This code:

Loops through all Kubernetes pods and checks each container’s status to detect unhealthy conditions like high restart count or waiting state.
If a problem is found, it prints an alert indicating that the pod is unhealthy (basic health monitoring logic).


</br><b>Alerts Module (simple for Sprint 1)</b>

Refer to: internal\alerts\alerts.go

Creates an alerts package with a function that prints an alert message to the console.
SendAlert() is used to simulate sending notifications when a problem is detected in the Kubernetes cluster.


## Step 3 - Creating the RBAC YAML file for Kubernetes Access

Refer to: deploy\manifests\rbac.yaml


It creates a ServiceAccount + RBAC rules so your app can list and watch pods in the cluster.

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: health-monitor-role
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: health-monitor-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: health-monitor-role
subjects:
- kind: ServiceAccount
  name: health-monitor-sa
  namespace: default

---

## Sprint 2 - Health Monitoring Module (Node & Pod Checks)

In Sprint 2 the basic pod loop from Sprint 1 becomes a real health monitoring module. The Health Checker now watches both nodes and pods, measures node CPU and memory usage, exposes the results as Prometheus metrics, runs inside the cluster as a container, and raises threshold-based alerts. Prometheus is installed to collect the metrics and Grafana is connected to visualize them.

<b>Cluster used :</b> Minikube (with the metrics-server addon)

<b>Sprint 2 Overview</b>

📁 internal/k8sclient </br><i>Now also lists nodes and reads node CPU and memory usage (Metrics API)</i>

📁 internal/monitor </br><i>Checks pod and node health and exposes Prometheus metrics</i>

📁 cmd/healer </br><i>Runs the metrics endpoint and the monitoring loop</i>

📁 deploy/manifests and deploy/prometheus </br><i>Run the monitor, Prometheus and Grafana inside the cluster</i>

---

## Step 1 - Enabling the Metrics API on Minikube

Node CPU and memory usage are read from the Kubernetes Metrics API, which is provided by the metrics-server addon.

<b>On the Command Line Terminal Run :</b>

minikube start --cpus 2 --memory 2200
minikube addons enable metrics-server
kubectl get nodes
kubectl top nodes

kubectl top nodes can take a minute or two to return values after the addon is enabled.

<img src="assets/Sprint2/Sprint2_Kube_Cluster_Ready_1.png" alt="Cluster ready and kubectl top nodes output" width="100%" height="100%">

<img src="assets/Sprint2/Sprint2_Kube_Cluster_Ready_2.png" alt="Cluster ready and kubectl top nodes output" width="100%" height="100%">

---

## Step 2 - Updating the Kubernetes Client and Adding Node Usage

</br><b>Kubernetes Client (internal/k8sclient)</b>

Refer to: internal\k8sclient\client.go

The client now holds a second client for the Metrics API next to the normal Kubernetes client. A new method, GetNodes(), lists the cluster nodes, and GetNodeMetricsData() reads the current CPU and memory usage of every node so the monitor can calculate utilization.

go
type Client struct {
	set           kubernetes.Interface
	metricsClient metricsclient.Interface
}

// GetNodeMetricsData retrieves node metrics for utilization calculations.
func (c *Client) GetNodeMetricsData() ([]metricsv1beta1.NodeMetrics, error) {
	metrics, err := c.metricsClient.MetricsV1beta1().NodeMetricses().List(
		context.Background(),
		metav1.ListOptions{},
	)
	if err != nil {
		return nil, err
	}

	return metrics.Items, nil
}

The new dependency (k8s.io/metrics) is added to go.mod by running:

go mod tidy

---

## Step 3 - Extending the Monitor Logic

</br><b>Monitor Logic</b>

Refer to: internal\monitor\monitor.go

This code:

CheckPods() now reports pods in the Failed or Unknown phase, containers that restarted more than 3 times, and containers stuck in a waiting state such as ImagePullBackOff. The new CheckNodes() reports nodes that are not Ready and nodes with memory, disk or PID pressure. Every finding is printed through alerts.SendAlert().

go
// CheckNodes checks node readiness and resource-pressure conditions.
func CheckNodes(nodes []v1.Node) {
	for _, node := range nodes {
		for _, condition := range node.Status.Conditions {
			switch condition.Type {
			case v1.NodeReady:
				if condition.Status != v1.ConditionTrue {
					alerts.SendAlert(fmt.Sprintf("Node %s is not ready: %s", node.Name, condition.Message))
				}
			case v1.NodeMemoryPressure, v1.NodeDiskPressure, v1.NodePIDPressure:
				if condition.Status == v1.ConditionTrue {
					alerts.SendAlert(fmt.Sprintf("Node %s has resource pressure: %s", node.Name, condition.Type))
				}
			}
		}
	}
}

---

## Step 4 - Exposing Prometheus Metrics

</br><b>Metrics Module</b>

Refer to: internal\monitor\metrics.go

This file registers the Prometheus metrics and serves them over HTTP. UpdateMetrics(), UpdateNodeMetrics() and UpdateNodeResourceMetrics() refresh the values on every loop.

Metric	Labels	Meaning
kubernetes_health_monitor_pods	namespace, phase	Number of pods per phase
kubernetes_health_monitor_container_restarts	namespace, pod, container	Restart count per container
kubernetes_health_monitor_unhealthy_pods	none	Pods that are failed, unknown, waiting or restarting too often
kubernetes_health_monitor_nodes	node, status	Node readiness (Ready / NotReady / Unknown)
kubernetes_health_monitor_node_pressure	node, condition	1 when memory, disk or PID pressure is active
kubernetes_health_monitor_node_cpu_usage_percent	node	Node CPU utilization
kubernetes_health_monitor_node_memory_usage_percent	node	Node memory utilization
kubernetes_health_monitor_node_scaling_recommendation	node, resource	1 when utilization is 85% or more (used by the scaling work in Sprint 4)

CPU and memory usage are converted into a percentage of the node's capacity:

go
cpuPercent = float64(nodeMetric.Usage.Cpu().MilliValue()) /
	float64(cpuCapacity) * 100

memoryPercent = float64(nodeMetric.Usage.Memory().Value()) /
	float64(memoryCapacity) * 100

The function also prints a WARNING when a node reaches 70% utilization and a SCALE RECOMMENDATION at 85%.


---

## Step 5 - Updating the Main Entry

</br><b>Create MAIN ENTRY (cmd/healer/main.go)</b>

Refer to: cmd\healer\main.go

The main file now does two jobs. It serves the /metrics endpoint on port 8080 in the background, and it runs the monitoring loop every 10 seconds. Each loop fetches pods, nodes and node usage, runs the checks, and refreshes the metrics. It also shuts down cleanly when the container is stopped.

go
for {
	pods := client.GetPods("default")
	nodes := client.GetNodes()

	nodeMetrics, err := client.GetNodeMetricsData()
	if err != nil {
		log.Printf("Error fetching node metrics: %v", err)
	} else {
		monitor.UpdateNodeResourceMetrics(nodes, nodeMetrics)
	}

	monitor.CheckPods(pods)
	monitor.CheckNodes(nodes)
	monitor.UpdateMetrics(pods)
	monitor.UpdateNodeMetrics(nodes)

	// wait 10 seconds, or stop on a shutdown signal
}

Automated pod recovery is not part of this sprint, so the loop does not call the healer package yet. It is added in Sprint 3.


---

## Step 6 - Building and Testing the Code

<b>On the Command Line Terminal Run :</b>

go mod tidy
go build ./...
go test ./...

All packages build without errors and the tests pass.

<img src="assets/Sprint2/Sprint2_The_Go_Test.png" alt="Go tests passing" width="100%" height="100%">

---

## Step 7 - Updating the RBAC YAML file

Refer to: deploy\manifests\rbac.yaml

The ServiceAccount from Sprint 1 now also needs to read nodes and the node metrics API. All permissions are read-only (get, list, watch), following the principle of least privilege.

yaml
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "list", "watch"]
- apiGroups: [""]
  resources: ["nodes"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["metrics.k8s.io"]
  resources: ["nodes"]
  verbs: ["get", "list", "watch"]

<b>Apply and verify :</b>

kubectl apply -f deploy/manifests/rbac.yaml
kubectl auth can-i list nodes --as=system:serviceaccount:default:health-monitor-sa

The second command prints yes. (The warning that nodes are not namespace scoped is expected, because nodes are cluster-wide.)


---

## Step 8 - Containerizing and Deploying the Health Monitor

Refer to: Dockerfile and deploy\manifests\health-monitor.yaml

The Dockerfile builds the Go program in a build stage and copies only the binary into a small Alpine image that runs as a non-root user and exposes port 8080.

The Deployment runs the monitor with the health-monitor-sa service account, uses /metrics for the readiness and liveness probes, and sets small CPU and memory requests and limits so the monitor uses very little of the cluster (a cost optimization). A Service named health-monitor exposes port 8080 so Prometheus can scrape it.

yaml
resources:
  requests:
    cpu: 50m
    memory: 64Mi
  limits:
    cpu: 100m
    memory: 128Mi

<b>Build the image, load it into Minikube and deploy :</b>

docker build -t health-monitor:latest .
minikube image load health-monitor:latest
kubectl apply -f deploy/manifests/health-monitor.yaml
kubectl apply -f deploy/prometheus/prometheus-deployment.yaml
kubectl apply -f deploy/prometheus/grafana-deployment.yaml
kubectl get pods

<img src="assets/Sprint2/Sprint2_Docker_Build_Status_1.png" alt="Docker image build" width="100%" height="100%">

<img src="assets/Sprint2/Sprint2_Docker_Build_Status_2.png" alt="Docker image build" width="100%" height="100%">

All three pods (health-monitor, prometheus, grafana) reach the Running state. The first start can take a few minutes because the Prometheus and Grafana images have to be downloaded.

<img src="assets/Sprint2/Sprint2_Pods_Running_1.png" alt="Monitor, Prometheus and Grafana pods running" width="100%" height="100%">

<img src="assets/Sprint2/Sprint2_Pods_Running_2.png" alt="Monitor, Prometheus and Grafana pods running" width="100%" height="100%">

<img src="assets/Sprint2/Sprint2_Pods_Running_3.png" alt="Monitor, Prometheus and Grafana pods running" width="100%" height="100%">

<img src="assets/Sprint2/Sprint2_Pods_Running_4.png" alt="Monitor, Prometheus and Grafana pods running" width="100%" height="100%">

<b>Check the monitor logs :</b>

kubectl logs deploy/health-monitor --tail=5

<img src="assets/Sprint2/Sprint2_Monitor_Logs.png" alt="Health monitor logs" width="100%" height="100%">

<b>Check the metrics endpoint :</b>

kubectl port-forward svc/health-monitor 18080:8080

Open http://localhost:18080/metrics and search for kubernetes_health_monitor. (Local port 18080 is used because Windows blocked port 8080. Any free local port works.)

<img src="assets/Sprint2/Sprint2_Metrics_Endpoint_1.png" alt="Metrics endpoint output" width="100%" height="100%">

<img src="assets/Sprint2/Sprint2_Metrics_Endpoint_2.png" alt="Metrics endpoint output" width="100%" height="100%">

<img src="assets/Sprint2/Sprint2_Metrics_Endpoint_3.png" alt="Metrics endpoint output" width="100%" height="100%">


---

## Step 9 - Installing Prometheus with Alert Rules

Refer to: deploy\prometheus\prometheus-deployment.yaml

One file holds everything Prometheus needs: a ConfigMap (scrape settings and alert rules), the Deployment and the Service. Prometheus scrapes itself and the Health Checker every 15 seconds.

yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

rule_files:
  - /etc/prometheus/alerts.yml

scrape_configs:
  - job_name: prometheus
    static_configs:
      - targets: ["localhost:9090"]
  - job_name: kubernetes-health-monitor
    metrics_path: /metrics
    static_configs:
      - targets: ["health-monitor:8080"]

Six threshold-based alert rules detect critical health issues:

Alert	Condition	Severity
HealthMonitorDown	Monitor target unreachable for 1 minute	critical
UnhealthyPodsDetected	unhealthy_pods > 0 for 1 minute	warning
KubernetesNodeNotReady	A node is not Ready for 2 minutes	critical
KubernetesNodeResourcePressure	Memory, disk or PID pressure for 2 minutes	warning
HighNodeCPU	Node CPU above 85% for 2 minutes	warning
HighNodeMemory	Node memory above 85% for 2 minutes	warning
yaml
- alert: HighNodeCPU
  expr: kubernetes_health_monitor_node_cpu_usage_percent > 85
  for: 2m
  labels:
    severity: warning
  annotations:
    summary: High node CPU usage
    description: "Node {{ $labels.node }} CPU usage is {{ $value }}%."

Prometheus also runs with small resource limits and keeps only 2 days of data to stay lightweight:

yaml
args:
  - --config.file=/etc/prometheus/prometheus.yml
  - --storage.tsdb.path=/prometheus
  - --storage.tsdb.retention.time=2d
resources:
  requests:
    cpu: 100m
    memory: 256Mi
  limits:
    memory: 512Mi

<b>Open the Prometheus UI :</b>

kubectl port-forward svc/prometheus 9090:9090

Open http://localhost:9090 (use another local port such as 19090 if 9090 is blocked).

Under Status > Targets, both the prometheus and kubernetes-health-monitor jobs show UP.

<img src="assets/Sprint2/Sprint2_Prometheus_Targets.png" alt="Prometheus targets UP" width="100%" height="100%">

In the Graph tab, running kubernetes_health_monitor_node_cpu_usage_percent shows the node CPU usage.

<img src="assets/Sprint2/Sprint2_Prometheus_CPU_Graph.png" alt="Node CPU usage graph in Prometheus" width="100%" height="100%">

The Alerts page lists the six rules, all inactive while the cluster is healthy.

<img src="assets/Sprint2/Sprint2_Prometheus_Alerts_Inactive.png" alt="Alert rules in Prometheus" width="100%" height="100%">

---

## Step 10 - Connecting Grafana

Refer to: deploy\prometheus\prometheus-deployment.yaml and deploy\prometheus\grafana-deployment.yaml

Grafana runs as its own Deployment and Service and only talks to Prometheus.

kubectl port-forward svc/grafana 3000:3000

Open http://localhost:3000 (use another local port such as 13000 if 3000 is blocked) and log in with the admin credentials set in grafana-deployment.yaml. Then:

Go to Connections > Data sources > Add data source > Prometheus.
Set the URL to http://prometheus:9090 and click Save & test. This is the in-cluster address of the Prometheus service.
Grafana confirms with "Successfully queried the Prometheus API".

<img src="assets/Sprint2/Sprint2_Grafana_Datasource_1.png" alt="Grafana Prometheus data source" width="100%" height="100%">


<img src="assets/Sprint2/Sprint2_Grafana_Datasource_2.png" alt="Grafana Prometheus data source" width="100%" height="100%">

Then create a dashboard named "Cluster Health" with a panel that uses the query kubernetes_health_monitor_node_cpu_usage_percent.

<img src="assets/Sprint2/Sprint2_Grafana_Panel.png" alt="Grafana dashboard showing node CPU usage" width="100%" height="100%">


---

## Step 11 - Testing the Alerts

The monitor is tested by creating a pod that can never start, because its image does not exist.

kubectl create deployment bad-pod --image=nginx:doesnotexist
kubectl logs deploy/health-monitor --tail=10

The monitor logs the container as waiting with the reason ImagePullBackOff.

<img src="assets/Sprint2/Sprint2_Bad_Pod_Log.png" alt="Monitor detecting the bad pod" width="100%" height="100%">

After about a minute the UnhealthyPodsDetected alert changes from Pending to Firing on the Prometheus Alerts page.

<img src="assets/Sprint2/Sprint2_Bad_Pod_Alert_Firing.png" alt="Monitor detecting the bad pod" width="100%" height="100%">

<b>Clean up :</b>

kubectl delete deployment bad-pod

The alert returns to inactive after about a minute.

<img src="assets/Sprint2/Sprint2_Bad_Pod_Alert_Cleared_1.png" alt="Monitor detecting the bad pod" width="100%" height="100%">


---

## Sprint 2 - Summary of Changes

File	Change
internal/k8sclient/client.go	Added the metrics client, GetNodes() and GetNodeMetricsData()
internal/monitor/monitor.go	Added CheckNodes() and extended CheckPods() (phase, restarts, waiting state)
internal/monitor/metrics.go	Added the Prometheus metrics, CPU and memory percentage calculation and the /metrics handler
cmd/healer/main.go	Runs the metrics endpoint and the monitoring loop with graceful shutdown
go.mod, go.sum	Added the Prometheus client and k8s.io/metrics dependencies
deploy/manifests/rbac.yaml	Read-only access to nodes and the node metrics API
Dockerfile	Multi-stage build running as a non-root user
deploy/manifests/health-monitor.yaml	Deployment, Service, health probes and resource limits
deploy/prometheus/prometheus-deployment.yaml	Prometheus with scrape config, six alert rules, resource limits and 2-day retention
deploy/prometheus/grafana-deployment.yaml	Grafana Deployment and Service

<b>Result :</b> The cluster is monitored for pod and node health, node CPU and memory usage is measured, all data is exposed to Prometheus, Grafana visualizes it, and six threshold-based alerts detect critical problems. The monitor is read-only at this stage.

---
<b></b>
<b></b>
<b></b>
<b></b>
