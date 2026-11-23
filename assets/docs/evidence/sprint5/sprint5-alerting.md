@'
# Sprint 5 – Alerting and Notification Integration

## Objective
Integrate Prometheus, Alertmanager, and Slack for real-time Kubernetes health notifications.

## Architecture
Kubernetes Health Monitor -> Prometheus -> Alertmanager -> Slack

## Components
- Prometheus: Collects health metrics and evaluates alert rules.
- Alertmanager: Groups and routes firing and resolved alerts.
- Slack: Receives Kubernetes health notifications.
- Kubernetes Secret: Securely stores the Slack webhook.

## Configuration
- deploy/prometheus/prometheus.yml
- deploy/prometheus/prometheus-deployment.yaml
- deploy/alertmanager/alertmanager.yml
- deploy/alertmanager/alertmanager-deployment.yaml

## Alert Rules
- UnhealthyPodsDetected: Warning for unhealthy pods.
- KubernetesNodeNotReady: Critical notification for nodes not Ready.
- KubernetesNodeResourcePressure: Warning for node resource pressure.

## Deployment
1. Create the Alertmanager ConfigMap.
2. Store the Slack webhook in the alertmanager-slack Kubernetes Secret.
3. Deploy Alertmanager and its Service.
4. Configure Prometheus to connect to alertmanager:9093.
5. Apply the Prometheus configuration and restart its deployment.

## Validation
- Slack received firing and resolved notifications.
- Alertmanager successfully routed notifications.
- Unhealthy pod metric returned 0.
- Excessive historical restart metric separately identified Grafana.
- No active Prometheus alerts remained.
- All Go tests passed with go test ./....

## Security
The Slack webhook is stored in a Kubernetes Secret. Never commit webhook credentials or Secret contents to Git.

## Evidence
See docs/evidence/sprint5/ for validation screenshots.
'@ | Set-Content -Encoding utf8 .\docs\sprint5-alerting.md