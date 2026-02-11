sequenceDiagram
    participant A1 as Agent A (Node 1)
    participant A2 as Agent B (Node 2)
    participant CP as Control Plane
    participant DB as PostgreSQL

    Note over A1,DB: Normal Operation
    A1->>CP: 1. Lease renewal
    CP->>DB: 2. UPDATE lease

    Note over A1: Node 1 FAILS
    Note over A1: No more renewals

    Note over CP: 30s later...
    CP->>DB: 3. Background: Find stale leases
    DB-->>CP: 4. Stale lease found
    CP->>DB: 5. DELETE stale lease

    Note over A2,DB: Recovery
    A2->>CP: 6. Poll for workloads
    CP-->>A2: 7. Return orphaned workload
    A2->>CP: 8. PUT /workloads/{id}/lease
    CP->>DB: 9. INSERT new lease