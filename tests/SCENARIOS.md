# Tests

## Basic workload acquisition

Required componenets:
1. Control plane
1. Single agent

Pipeline:
1. Apply workload with mctl
1. Wait for workload acquisition and status ACTIVE

## Multiple workloads

Required componenets:
1. Control plane
1. Single agent

Pipeline:
1. Apply two workloads with mctl
1. Wait for workloads acquisition and status set ACTIVE

## Exclusive leases

Required componenets:
1. Control plane
1. Two agents

Pipeline:
1. Apply workloads with mctl
1. Wait for workload acquisition and status set ACTIVE
1. Ensure that only one agent acquired workload

## Rescheduling

Required componenets:
1. Control plane
1. Two agents

Pipeline:
1. Apply workloads with mctl
1. Wait for workload acquisition and status set ACTIVE
1. Disable agent where workload is acquired
1. Wait for workload to be acquired by other agent and status become ACTIVE

## Concurrent workloads

Required componenets:
1. Control plane
1. Single agent

Pipeline:
1. Simultaneously apply several worklods
1. Wait for workloads acquisition and status set ACTIVE

## Edge

Required componenets:
1. Control plane
1. Single agent
1. Edge

Pipeline:
1. Apply workloads with mctl and set edge config
1. Wait for workloads acquisition and status set ACTIVE
1. Ensure that edge is working and available 


# Required metrics

1. Execution time of each step
1. RTO of lease prolongation
