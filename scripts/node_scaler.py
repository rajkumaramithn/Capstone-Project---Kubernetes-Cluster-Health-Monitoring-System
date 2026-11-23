#!/usr/bin/env python3
"""
Node scaler for the Kubernetes Health Monitor (Sprint 4).

Reads node CPU/memory utilization from Prometheus (published by the Health
Checker) and adds or removes a Minikube worker node.

  python scripts/node_scaler.py             dry run, prints the decision only
  python scripts/node_scaler.py --apply     performs the scaling action
  python scripts/node_scaler.py --apply --loop   re-evaluates continuously

Settings (environment variables): PROM_URL, SCALE_UP, SCALE_DOWN,
MIN_NODES, MAX_NODES, INTERVAL, COOLDOWN.
Standard library only (Python 3.8+).
"""
import json
import os
import subprocess
import sys
import time
import urllib.parse
import urllib.request

PROM_URL = os.getenv("PROM_URL", "http://localhost:9090")
SCALE_UP = float(os.getenv("SCALE_UP", "85"))      # % utilization
SCALE_DOWN = float(os.getenv("SCALE_DOWN", "30"))  # % utilization
MIN_NODES = int(os.getenv("MIN_NODES", "2"))
MAX_NODES = int(os.getenv("MAX_NODES", "3"))
INTERVAL = int(os.getenv("INTERVAL", "60"))        # seconds between checks
COOLDOWN = int(os.getenv("COOLDOWN", "300"))       # seconds after an action
APPLY = "--apply" in sys.argv
LOOP = "--loop" in sys.argv


def query(expr):
    url = f"{PROM_URL}/api/v1/query?query={urllib.parse.quote(expr)}"
    with urllib.request.urlopen(url, timeout=10) as resp:
        return json.load(resp)["data"]["result"]


def run(cmd):
    print("ACTION:", " ".join(cmd))
    if APPLY:
        subprocess.run(cmd, check=True)


def evaluate():
    cpu_series = query("kubernetes_health_monitor_node_cpu_usage_percent")
    mem_series = query("kubernetes_health_monitor_node_memory_usage_percent")
    if not cpu_series or not mem_series:
        print("No node metrics found in Prometheus. Is the monitor running?")
        return False

    nodes = sorted(s["metric"]["node"] for s in cpu_series)
    cpu = sum(float(s["value"][1]) for s in cpu_series) / len(cpu_series)
    mem = sum(float(s["value"][1]) for s in mem_series) / len(mem_series)
    load = max(cpu, mem)
    print(f"nodes={len(nodes)} {nodes} avg_cpu={cpu:.1f}% avg_mem={mem:.1f}% "
          f"(scale up >= {SCALE_UP}%, scale down <= {SCALE_DOWN}%)")

    if load >= SCALE_UP and len(nodes) < MAX_NODES:
        print("DECISION: SCALE UP")
        run(["minikube", "node", "add"])
        # a new node does not have the monitor image yet
        run(["minikube", "image", "load", "health-monitor:latest"])
        return APPLY
    if load <= SCALE_DOWN and len(nodes) > MIN_NODES:
        print("DECISION: SCALE DOWN")
        victim = [n for n in nodes if n != "minikube"][-1]  # never the control plane
        run(["kubectl", "drain", victim, "--ignore-daemonsets",
             "--delete-emptydir-data", "--timeout=120s"])
        run(["minikube", "node", "delete", victim])
        return APPLY

    print("DECISION: NO CHANGE")
    return False


def main():
    if not LOOP:
        evaluate()
        return
    while True:
        try:
            acted = evaluate()
        except Exception as exc:  # keep the loop alive on temporary errors
            print("ERROR:", exc)
            acted = False
        time.sleep(COOLDOWN if acted else INTERVAL)


if __name__ == "__main__":
    main()