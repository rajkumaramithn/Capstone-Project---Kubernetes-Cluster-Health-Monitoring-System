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

## Sprint 3 - Self-Healing Mechanisms (Pod Recovery)

In Sprint 3 the Health Checker stops only observing and starts acting. When it finds a pod that has failed, been evicted, or is stuck in a crash loop, it removes the pod so that the owning controller (a Deployment or a Job) creates a fresh replacement, which Kubernetes can place on a healthy node. Every action is written to the log with a timestamp and a reason, and is counted as a Prometheus metric, so the recovery is transparent and traceable.

<b>Cluster used :</b> Minikube (local test cluster acting as the staging environment)

<b>Sprint 3 Overview</b>

📁 internal/k8sclient
</br><i>New check that decides why a pod needs healing (Failed, Evicted or CrashLoopBackOff) and a safe deletion method</i>

📁 internal/healer
</br><i>Heals each unhealthy pod and writes an audit log line for every action</i>

📁 internal/monitor
</br><i>New Prometheus counter for healing actions</i>

📁 cmd/healer
</br><i>The monitoring loop now calls the healer</i>

📁 deploy/manifests
</br><i>RBAC permission to delete pods, plus two test workloads</i>

<b>How Self-Healing Works</b>

1. Every 10 seconds the monitor reads the pods in the default namespace.
2. For each pod it asks: is this pod Failed, Evicted, or crash-looping (CrashLoopBackOff with 3 or more restarts)?
3. If yes, the pod is re-read from the API to avoid acting on stale information.
4. The pod is deleted only if it is controller-managed and still the same pod (matching UID).
5. The Deployment or Job controller notices the missing pod and creates a replacement.
6. The action is logged and counted.

Pods that are not managed by a controller are never deleted, because nothing would recreate them.

---

## Step 1 - Deciding Which Pods Need Healing

</br><b>Kubernetes Client (internal/k8sclient)</b>

Refer to: internal\k8sclient\client.go

A new function, UnhealthyReason(), returns the reason a pod should be healed, or an empty string if the pod should be left alone. It covers the three cases in the sprint goal: failed pods, evicted pods (Kubernetes marks them Failed with the reason Evicted) and containers stuck in CrashLoopBackOff after at least 3 restarts.

```go
// UnhealthyReason returns why a pod should be healed, or "" if it should be left alone.
func UnhealthyReason(pod *v1.Pod) string {
	if pod.Status.Phase == v1.PodFailed {
		if pod.Status.Reason == "Evicted" {
			return "Evicted"
		}
		return "Failed"
	}

	for _, cs := range pod.Status.ContainerStatuses {
		if w := cs.State.Waiting; w != nil &&
			w.Reason == "CrashLoopBackOff" &&
			cs.RestartCount >= 3 {
			return "CrashLoopBackOff"
		}
	}

	return ""
}
```

The restart threshold of 3 gives a pod a few chances to recover on its own before the healer steps in.

---


## Step 2 - Deleting Pods Safely

Refer to: internal\k8sclient\client.go

The existing DeleteFailedControllerPod() method already re-read the pod, checked its UID, refused to delete standalone pods, and deleted with a UID precondition. In Sprint 3 only one condition changed: it used to heal only pods in the Failed phase, and now it heals any pod that UnhealthyReason() flags.

<b>Before :</b>

```go
// Confirm that this is still the same pod and is still failed.
if current.UID != pod.UID || current.Status.Phase != v1.PodFailed {
	return false, nil
}
```

<b>After :</b>

```go
// Confirm that this is still the same pod and that it is still unhealthy
// (Failed, Evicted or crash-looping).
if current.UID != pod.UID || UnhealthyReason(current) == "" {
	return false, nil
}
```

The safety checks that remain in the method:

| Safety check | Purpose |
|---|---|
| Re-read the pod before acting | Never act on stale information |
| Compare the UID | Make sure it is still the same pod and not a replacement |
| Require a controller owner | Never delete a standalone pod that nothing would recreate |
| UID precondition on the delete call | Never delete a replacement pod by mistake |

---


## Step 3 - The Healer with an Audit Log

</br><b>Healer (internal/healer)</b>

Refer to: internal\healer\healer.go

HealFailedPods() loops over the pods, asks for the reason, requests the safe deletion, and records what happened. Each action is written with log.Printf, which adds a timestamp, so the log works as the audit trail. Each successful action also increases the healing counter. An error on one pod does not stop the healer from handling the others.

```go
func HealFailedPods(client podHealer, pods []v1.Pod) {
	for _, pod := range pods {
		reason := k8sclient.UnhealthyReason(&pod)
		if reason == "" {
			continue
		}

		deleted, err := client.DeleteFailedControllerPod(
			context.Background(),
			pod,
		)
		if err != nil {
			log.Printf("HEAL ERROR ns=%s pod=%s err=%v", pod.Namespace, pod.Name, err)
			continue
		}

		if deleted {
			log.Printf("HEAL action=delete_pod ns=%s pod=%s reason=%s", pod.Namespace, pod.Name, reason)
			monitor.RecordHealingAction(pod.Namespace, reason)
		}
	}
}
```
---


## Step 4 - Counting Healing Actions

</br><b>Metrics Module</b>

Refer to: internal\monitor\metrics.go

A new Prometheus counter records every healing action, labelled by namespace and reason. This feeds the dashboard and alerts in later sprints.

```go
healingActions = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "kubernetes_health_monitor_healing_actions_total",
		Help: "Self-healing actions taken, by namespace and reason.",
	},
	[]string{"namespace", "reason"},
)
```

It is registered in init() together with the other metrics, and updated through a small helper:

```go
// RecordHealingAction counts a self-healing action.
func RecordHealingAction(namespace, reason string) {
	healingActions.WithLabelValues(namespace, reason).Inc()
}
```

A counter only appears on the metrics page after its first healing action has happened.

---


## Step 5 - Wiring the Healer into the Main Loop

</br><b>Main Entry (cmd/healer/main.go)</b>

Refer to: cmd\healer\main.go

The healer package is imported again and called at the end of each monitoring cycle:

```go
monitor.CheckPods(pods)
monitor.CheckNodes(nodes)
monitor.UpdateMetrics(pods)
monitor.UpdateNodeMetrics(nodes)
healer.HealFailedPods(client, pods)
```

In Sprint 2 the monitor was read-only. From this sprint on it can change the cluster, but only by deleting unhealthy, controller-managed pods.

---


## Step 6 - Allowing the Monitor to Delete Pods (RBAC)

Refer to: deploy\manifests\rbac.yaml

Only the pods rule gains the delete verb. Nodes and the node metrics API stay read-only, following the principle of least privilege.

```yaml
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "list", "watch", "delete"]
- apiGroups: [""]
  resources: ["nodes"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["metrics.k8s.io"]
  resources: ["nodes"]
  verbs: ["get", "list", "watch"]
```

---


## Step 7 - Test Workloads for the Staging Environment

Refer to: deploy\manifests\crashloop-test.yaml and deploy\manifests\failed-job-test.yaml

Two small workloads create the problems the healer must fix. Both use a busybox container that exits with an error.

<b>A Deployment whose container always crashes (produces CrashLoopBackOff) :</b>

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: crashloop-test
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: crashloop-test
  template:
    metadata:
      labels:
        app: crashloop-test
    spec:
      containers:
        - name: crash
          image: busybox:1.36
          command: ["sh", "-c", "echo crashing; exit 1"]
```

<b>A Job whose pod fails and is not restarted (produces a pod in the Failed phase) :</b>

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: failed-job-test
  namespace: default
spec:
  backoffLimit: 2
  template:
    spec:
      restartPolicy: Never
      containers:
        - name: fail
          image: busybox:1.36
          command: ["sh", "-c", "echo failing; exit 1"]
```

---


## Step 8 - Unit Tests and Build

Refer to: internal\k8sclient\client_test.go and internal\healer\healer_test.go

New tests were added for the new behaviour:

| Test | What it proves |
|---|---|
| TestUnhealthyReason | Crash-looping, Evicted and Failed pods are flagged with the right reason, and a healthy pod is not |
| TestHealFailedPodsRequestsRecoveryForCrashLoopPod | The healer asks for recovery of a crash-looping pod |

The existing tests still cover deleting a failed controller-managed pod, skipping a stale pod (UID changed), skipping a healthy pod and continuing after an error.

Evicted pods are difficult to create safely on a local cluster, so the unit test is the evidence for them.

<b>On the Command Line Terminal Run :</b>

```
go mod tidy
go build ./...
go test ./...
```

All packages build and the tests pass.

<img src="assets/Sprint3/Sprint3-Go-Test-Pass.png" alt="Go tests passing" width="100%" height="100%">

---


## Step 9 - Deploying the Updated Monitor

Start a fresh cluster, build the image, load it into Minikube and deploy:

```
minikube start --driver=docker --cpus 2 --memory 2200
minikube addons enable metrics-server
docker build -t health-monitor:latest .
minikube image load health-monitor:latest
kubectl apply -f deploy/manifests/rbac.yaml
kubectl apply -f deploy/manifests/health-monitor.yaml
kubectl apply -f deploy/prometheus/prometheus-deployment.yaml
kubectl apply -f deploy/prometheus/grafana-deployment.yaml
kubectl get pods
```

Wait until the health-monitor, prometheus and grafana pods show 1/1 Running.

<img src="assets/Sprint3/Sprint3-Baseline-Pods.png" alt="Baseline pods running" width="100%" height="100%">

Verify that the monitor is allowed to delete pods:

```
kubectl auth can-i delete pods --as=system:serviceaccount:default:health-monitor-sa
```

The command prints `yes`.


<img src="assets/Sprint3/Sprint3-Rbac-Delete-Permission.png" alt="RBAC delete permission check" width="100%" height="100%">

---


## Step 10 - Test 1: A Crash-Looping Pod

```
kubectl apply -f deploy/manifests/crashloop-test.yaml
kubectl get pods -w
```

The pod fails, restarts and enters CrashLoopBackOff. When its restart count reaches 3, the healer deletes it (the pod shows Terminating) and the Deployment creates a replacement with a new name and a restart count of 0. The replacement crashes in the same way and is healed in turn, so the cycle repeats.

In the test run the first pod reached CrashLoopBackOff with 3 restarts about 49 seconds after creation and was terminated 5 seconds later.

<img src="assets/Sprint3/Sprint3-Crashloop-Pod-Replaced.png" alt="Crash-looping pod terminated and replaced" width="100%" height="100%">


<b>Read the audit trail.</b> The filter is case-sensitive so that only healing lines match:

```
kubectl logs deploy/health-monitor | Select-String "HEAL action" -CaseSensitive
```

Example output from the test run:

```
2026/11/21 11:54:09 HEAL action=delete_pod ns=default pod=crashloop-test-d9b878ddb-x6qfx reason=CrashLoopBackOff
2026/11/21 11:58:29 HEAL action=delete_pod ns=default pod=crashloop-test-d9b878ddb-mvqx7 reason=CrashLoopBackOff
```

Each line records when the action happened, which pod was removed and why. The second line is the replacement pod being healed, which shows that recovery repeats.


<img src="assets/Sprint3/Sprint3-Heal-Audit-Log.png" alt="Healing audit log" width="100%" height="100%">

---


## Step 11 - Test 2: A Failed Job Pod

Remove the first test and create a pod that ends in the Failed phase:

```
kubectl delete -f deploy/manifests/crashloop-test.yaml
kubectl apply -f deploy/manifests/failed-job-test.yaml
kubectl get pods -w
```

The Job's pod shows Error. The healer deletes it and the Job controller creates a new pod, up to the Job's backoff limit of 2. The audit log then contains a healing line with reason=Failed.

```
kubectl logs deploy/health-monitor | Select-String "HEAL action" -CaseSensitive
```


<img src="assets/Sprint3/Sprint3-Failed-Job-Healed-1.png" alt="Failed job pod healed" width="100%" height="100%">

<img src="assets/Sprint3/Sprint3-Failed-Job-Healed-2.png" alt="Failed job pod healed" width="100%" height="100%">

---


## Step 12 - Test 3: A Standalone Pod Is Left Alone

This is the safety check. A pod with no controller would never be recreated, so the healer must not delete it.

```
kubectl run solo --image=busybox:1.36 --restart=Never --command -- sh -c "exit 1"
```

After 30 to 40 seconds, check the pods and the log:

```
kubectl get pods
kubectl logs deploy/health-monitor --tail=10
```

The solo pod is still present with the status Error, and there is no healing line for it in the log.


<img src="assets/Sprint3/Sprint3-Standalone-Pod-Skipped.png" alt="Standalone pod left alone" width="100%" height="100%">

---


## Step 13 - Healing Metrics

Open the monitor's metrics endpoint:

```
kubectl port-forward svc/health-monitor 18080:8080
```

Browse to http://localhost:18080/metrics and search for healing_actions_total. The counter shows the number of actions per reason (CrashLoopBackOff and Failed). Local port 18080 is used because Windows blocked port 8080.


<img src="assets/Sprint3/Sprint3-Healing-Metric-1.png" alt="Healing actions metric on the metrics endpoint" width="100%" height="100%">

<img src="assets/Sprint3/Sprint3-Healing-Metric-2.png" alt="Healing actions metric on the metrics endpoint" width="100%" height="100%">

The same metric is available in Prometheus. With the Prometheus port-forward running, open its UI, go to the Graph tab and run:

```
kubernetes_health_monitor_healing_actions_total
```

<img src="assets/Sprint3/Sprint3-Prometheus-Healing-Metric-1.png" alt="Healing actions metric in Prometheus" width="100%" height="100%">

<img src="assets/Sprint3/Sprint3-Prometheus-Healing-Metric-2.png" alt="Healing actions metric in Prometheus" width="100%" height="100%">

---


## Step 14 - Clean Up

```
kubectl delete -f deploy/manifests/failed-job-test.yaml
kubectl delete pod solo
```

---


## Limitations and Notes

- The healer watches the default namespace only.
- Healing means deleting the pod. Kubernetes cannot restart a pod in place, so the controller creates a replacement, which the scheduler may place on a different node. That is how workloads are rescheduled.
- If the cause of a crash is permanent, such as a bad image or a bug, the replacement crashes too and is healed again. The healer fixes transient failures. Persistent failures show up as repeated healing actions, which later sprints turn into alerts.
- Scaling and resource balancing are delivered in Sprint 4.

---


## Sprint 3 - Summary of Changes

| File | Change |
|---|---|
| internal/k8sclient/client.go | Added UnhealthyReason() and changed DeleteFailedControllerPod() to heal Failed, Evicted and crash-looping pods |
| internal/healer/healer.go | Heals by reason, logs every action with a timestamp and reason, and records the healing metric |
| internal/monitor/metrics.go | Added the kubernetes_health_monitor_healing_actions_total counter and RecordHealingAction() |
| cmd/healer/main.go | Calls the healer in every monitoring cycle |
| deploy/manifests/rbac.yaml | Added the delete verb for pods only |
| deploy/manifests/crashloop-test.yaml | New test Deployment that crash-loops |
| deploy/manifests/failed-job-test.yaml | New test Job whose pod fails |
| internal/k8sclient/client_test.go | Added TestUnhealthyReason |
| internal/healer/healer_test.go | Added a test for crash-looping pods |

<b>Result :</b> The system now recovers automatically from failed, evicted and crash-looping pods. Recovery is safe (controller-managed pods only, with UID checks), transparent (a timestamped audit log) and measurable (a Prometheus counter). It was verified on a local staging cluster with three tests: a crash-looping pod, a failed Job pod and a standalone pod that is correctly left alone.

---

## Sprint 4 - Advanced Self-Healing (Node Scaling & Resource Balancing)

In Sprint 4 the system learns to scale and balance itself. It adds worker nodes when the cluster is busy and removes them when it is idle, scales workloads automatically with a Horizontal Pod Autoscaler (HPA), and spreads pods evenly across the nodes. All of it was tested under simulated load. The decisions for node scaling use the node CPU and memory metrics that the Health Checker has published since Sprint 2, so no Go code changes were needed in this sprint.

<b>Cluster used :</b> Minikube with 2 nodes at the start (control plane and one worker) and the metrics-server addon. The scaler grows it to 3 nodes and shrinks it back.

<b>Sprint 4 Overview</b>

📁 scripts
</br><i>New Python node scaler that adds or removes a worker node based on utilization</i>

📁 deploy/manifests
</br><i>Resource balancing (topology spread), HPA settings and a load generator for the load test</i>

📁 requirements.txt
</br><i>Documents that the scaler uses only the Python standard library</i>

<b>How Scaling Works</b>

| Layer | Mechanism | Trigger |
|---|---|---|
| Pods | Horizontal Pod Autoscaler | Average pod CPU above 50% of its request |
| Nodes | scripts/node_scaler.py | Average node utilization (the higher of CPU and memory) at or above 85% to add a node, or at or below 30% to remove one |
| Placement | Topology spread constraint | Always: pods are spread evenly across nodes |


---


## Step 1 - Starting a Multi-Node Cluster

Node scaling and pod redistribution cannot be shown on a single node, so the cluster starts with two nodes.

<b>On the Command Line Terminal Run :</b>

```
minikube start --driver=docker --nodes 2 --cpus 2 --memory 2000
minikube addons enable metrics-server
kubectl get nodes
kubectl top nodes
```

Wait until `minikube-metrics-server` is running and `kubectl top nodes` shows CPU and memory for both nodes. This can take a few minutes after the addon is enabled.

<img src="assets/Sprint4/Sprint4-Cluster-Two-Nodes.png" alt="Two node cluster and node resource usage" width="100%" height="100%">

---


## Step 2 - Building and Testing the Code

No Go code was changed in this sprint. The build and tests are run to confirm that the code carried over from Sprint 3 is intact.

```
go mod tidy
go build ./...
go test ./...
```

<img src="assets/Sprint4/Sprint4-Go-Test-Pass.png" alt="Go tests passing" width="100%" height="100%">

---


## Step 3 - Deploying the Monitoring Stack

Build the image, load it into Minikube and deploy the monitor, Prometheus and Grafana:

```
docker build -t health-monitor:latest .
minikube image load health-monitor:latest
kubectl apply -f deploy/manifests/rbac.yaml
kubectl apply -f deploy/manifests/health-monitor.yaml
kubectl apply -f deploy/prometheus/prometheus-deployment.yaml
kubectl apply -f deploy/prometheus/grafana-deployment.yaml
kubectl get pods -o wide
```

The `-o wide` option also shows which node each pod runs on.


<img src="assets/Sprint4/Sprint4-Baseline-Pods.png" alt="Baseline pods and their nodes" width="100%" height="100%">

---


## Step 4 - Resource Balancing Across Nodes

Refer to: deploy\manifests\resource-balanced-deployment.yaml

The demo Deployment runs 3 replicas with CPU and memory requests and limits. The new `topologySpreadConstraints` block tells the scheduler to keep the number of pods on each node within 1 of each other, so no node ends up overloaded while another sits idle.

```yaml
spec:
  template:
    spec:
      topologySpreadConstraints:
      - maxSkew: 1
        topologyKey: kubernetes.io/hostname
        whenUnsatisfiable: ScheduleAnyway
        labelSelector:
          matchLabels:
            app: resource-balanced-demo
```

`ScheduleAnyway` means the constraint is a strong preference and not a hard rule, so pods are never left Pending when a node is full.

<b>Deploy and check where the pods landed :</b>

```
kubectl apply -f deploy/manifests/resource-balanced-deployment.yaml
kubectl get pods -o wide -l app=resource-balanced-demo
```

The 3 pods are spread over the two nodes (2 and 1).


<img src="assets/Sprint4/Sprint4-Balanced-Pods-Two-Nodes.png" alt="Pods balanced across two nodes" width="100%" height="100%">

---


## Step 5 - Pod Autoscaling with a Load Test (HPA)

Refer to: deploy\manifests\hpa-test-deployment.yaml, hpa-demo.yaml and load-generator.yaml

<b>The workload (hpa-test-deployment.yaml) :</b> an nginx Deployment with a CPU request of 100m, which the HPA uses as its baseline. A Service named hpa-demo was added so the load generator can reach it.

```yaml
resources:
  requests:
    cpu: 100m
    memory: 64Mi
  limits:
    cpu: 300m
    memory: 128Mi
```

<b>The autoscaler (hpa-demo.yaml) :</b> keeps between 1 and 4 replicas and adds pods when average CPU use goes above 50% of the request. The new `behavior` block shortens the scale-down wait from the default 5 minutes to 60 seconds, so unused pods are released sooner and cost less.

```yaml
spec:
  minReplicas: 1
  maxReplicas: 4
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 50
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 60
```

<b>The simulated load (load-generator.yaml) :</b> three busybox pods that send requests to hpa-demo in a loop.

```yaml
command: ["sh", "-c", "while true; do wget -q -O- http://hpa-demo > /dev/null; done"]
```

<b>Deploy the workload and the autoscaler :</b>

```
kubectl apply -f deploy/manifests/hpa-test-deployment.yaml
kubectl apply -f deploy/manifests/hpa-demo.yaml
kubectl get hpa
```

Wait until the TARGETS column shows a percentage and not `<unknown>`.


<img src="assets/Sprint4/Sprint4-Hpa-Before-Load.png" alt="HPA before load" width="100%" height="100%">

<b>Start the load and watch the HPA scale up :</b>

```
kubectl apply -f deploy/manifests/load-generator.yaml
kubectl get hpa -w
```

CPU use rises above the 50% target and the replica count increases, up to the maximum of 4.


<img src="assets/Sprint4/Sprint4-Hpa-Scale-Up.png" alt="HPA scaling up under load" width="100%" height="100%">

<b>Stop the load and watch it scale down :</b>

```
kubectl delete -f deploy/manifests/load-generator.yaml
kubectl get hpa -w
```

After the 60-second stabilization window, the replicas return to the minimum.

<img src="assets/Sprint4/Sprint4-Hpa-Scale-Down.png" alt="HPA scaling down after the load stops" width="100%" height="100%">

---


## Step 6 - Automatic Node Scaling

Refer to: scripts\node_scaler.py

The scaler reads the node CPU and memory percentages that the Health Checker publishes to Prometheus (kubernetes_health_monitor_node_cpu_usage_percent and kubernetes_health_monitor_node_memory_usage_percent) and decides what to do:

| Condition | Decision | Action |
|---|---|---|
| Average utilization at or above 85% and fewer than 3 nodes | Scale up | `minikube node add`, then load the monitor image onto the new node |
| Average utilization at or below 30% and more than 2 nodes | Scale down | Drain the newest worker node, then `minikube node delete` it |
| Anything else | No change | Nothing |

The control plane node is never removed. The script does nothing unless it is run with `--apply` (a dry run is the default), and with `--loop` it keeps checking on its own. It uses only the Python standard library, so there is nothing to install.

```python
PROM_URL = os.getenv("PROM_URL", "http://localhost:9090")
SCALE_UP = float(os.getenv("SCALE_UP", "85"))      # % utilization
SCALE_DOWN = float(os.getenv("SCALE_DOWN", "30"))  # % utilization
MIN_NODES = int(os.getenv("MIN_NODES", "2"))
MAX_NODES = int(os.getenv("MAX_NODES", "3"))
INTERVAL = int(os.getenv("INTERVAL", "60"))        # seconds between checks
COOLDOWN = int(os.getenv("COOLDOWN", "300"))       # seconds after an action
```

The decision logic:

```python
nodes = sorted(s["metric"]["node"] for s in cpu_series)
cpu = sum(float(s["value"][1]) for s in cpu_series) / len(cpu_series)
mem = sum(float(s["value"][1]) for s in mem_series) / len(mem_series)
load = max(cpu, mem)

if load >= SCALE_UP and len(nodes) < MAX_NODES:
    print("DECISION: SCALE UP")
    run(["minikube", "node", "add"])
    # a new node does not have the monitor image yet
    run(["minikube", "image", "load", "health-monitor:latest"])
elif load <= SCALE_DOWN and len(nodes) > MIN_NODES:
    print("DECISION: SCALE DOWN")
    victim = [n for n in nodes if n != "minikube"][-1]  # never the control plane
    run(["kubectl", "drain", victim, "--ignore-daemonsets",
         "--delete-emptydir-data", "--timeout=120s"])
    run(["minikube", "node", "delete", victim])
else:
    print("DECISION: NO CHANGE")
```

On a real cloud cluster, the two `minikube` commands would be replaced by the cloud provider's node-group API. The decision logic stays the same.

The `requirements.txt` file records that no packages are needed:

```
# node_scaler.py uses only the Python standard library (Python 3.8+), no packages to install.
```

<b>Open Prometheus so the scaler can read it (second terminal) :</b>

```
kubectl port-forward svc/prometheus 19090:9090
```

<b>Dry run (nothing is changed) :</b>

```
$env:PROM_URL="http://localhost:19090"
python scripts/node_scaler.py
```

The script prints the node utilization and `DECISION: NO CHANGE`, because a lightly loaded cluster is between the two thresholds.


<img src="assets/Sprint4/Sprint4-Scaler-No-Change.png" alt="Node scaler dry run" width="100%" height="100%">

<b>Scale up.</b> A lightly loaded test cluster never reaches 85%, so the threshold is lowered for the test only:

```
$env:SCALE_UP="5"
python scripts/node_scaler.py --apply
kubectl get nodes
```

The script prints `DECISION: SCALE UP`, adds the worker node `minikube-m03` and loads the monitor image onto it. The cluster now has 3 nodes.

<img src="assets/Sprint4/Sprint4-Scale-Up-Three-Nodes.png" alt="Scale up to three nodes" width="100%" height="100%">

<b>Redistribute the pods onto the new node :</b>

```
kubectl rollout restart deployment/resource-balanced-demo
kubectl get pods -o wide -l app=resource-balanced-demo
```

The topology spread constraint places the new pods evenly, one on each of the three nodes.

<img src="assets/Sprint4/Sprint4-Balanced-Pods-Three-Nodes.png" alt="Pods balanced across three nodes" width="100%" height="100%">

<b>Node metrics in Prometheus.</b> Open http://localhost:19090 and run these queries in the Graph tab. Each shows one line per node, now including the new node:

```
kubernetes_health_monitor_node_cpu_usage_percent
kubernetes_health_monitor_node_memory_usage_percent
```

<img src="assets/Sprint4/Sprint4-Prometheus-Node-Metrics.png" alt="Node CPU and memory metrics in Prometheus" width="100%" height="100%">

<b>Scale down.</b> This time the scale-down threshold is raised for the test only:

```
Remove-Item Env:SCALE_UP
$env:SCALE_DOWN="90"
python scripts/node_scaler.py --apply
kubectl get nodes
```

The script prints `DECISION: SCALE DOWN`, drains `minikube-m03` (its pods move to the other nodes) and deletes it. The cluster is back to 2 nodes.


<img src="assets/Sprint4/Sprint4-Scale-Down-Two-Nodes.png" alt="Scale down to two nodes" width="100%" height="100%">


## Step 7 - Cost Optimization Evidence

Several design choices in this sprint keep the cluster small and cheap:

| Measure | Effect |
|---|---|
| Node scaling down at 30% utilization or below | Idle nodes are removed instead of running for nothing |
| A minimum of 2 nodes and a maximum of 3 | Keeps cost and capacity within known limits |
| HPA with a 60-second scale-down window | Extra pods are released quickly after a load spike |
| CPU and memory requests and limits on every workload | Pods are right-sized and the scheduler packs them efficiently |
| Even spreading across nodes | No node is overloaded while another sits idle |

Reset the test thresholds and capture the cluster's resource usage:

```
Remove-Item Env:SCALE_DOWN
Remove-Item Env:PROM_URL
kubectl top nodes
kubectl top pods
```

<img src="assets/Sprint4/Sprint4-Cost-Resource-Usage.png" alt="Cluster resource usage" width="100%" height="100%">

<b>Clean up :</b>

```
kubectl delete -f deploy/manifests/hpa-demo.yaml
kubectl delete -f deploy/manifests/hpa-test-deployment.yaml
kubectl delete -f deploy/manifests/resource-balanced-deployment.yaml
```

---


## Notes

- The test thresholds (5% to force a scale up and 90% to force a scale down) were temporary overrides, because a lightly loaded test cluster never reaches the real ones. The defaults are 85% and 30%, which match the 85% recommendation threshold already used by the Health Checker.
- A node added later does not have the locally built monitor image, which is why the scaler loads it right after adding a node.
- Metrics-server needs a few minutes to become ready after it is enabled. The HPA and the node CPU and memory metrics only work once `kubectl top nodes` shows values.
- On a real cloud cluster, the node scaling step maps to the provider's autoscaler or node group API.

---


## Sprint 4 - Summary of Changes

| File | Change |
|---|---|
| scripts/node_scaler.py | New. Reads node utilization from Prometheus and adds or removes a worker node, with dry-run and loop modes |
| deploy/manifests/resource-balanced-deployment.yaml | Added a topology spread constraint to balance pods across nodes |
| deploy/manifests/hpa-demo.yaml | Added a 60-second scale-down stabilization window |
| deploy/manifests/hpa-test-deployment.yaml | Added a Service so the load generator can reach the workload |
| deploy/manifests/load-generator.yaml | New. Simulated load for the HPA test |
| requirements.txt | Documents that the scaler needs no packages |

<b>Result :</b> The system now scales and balances itself. Pods are scaled by the HPA under load, nodes are added and removed automatically from the utilization metrics, and pods are spread evenly across nodes. All of it was verified under simulated load on a multi-node staging cluster, and the sizing choices keep the cluster cost low.

---

<b></b>
<b></b>
<b></b>
<b></b>
