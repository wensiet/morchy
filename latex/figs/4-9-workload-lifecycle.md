sequenceDiagram
    participant O as Operator
    participant A1 as Agent 1
    participant A2 as Agent 2
    participant CP as Control Plane
    participant DB as Database

    Note over O,DB: Phase 1: Workload Creation
    O->>A1: Create
    A1->>CP: POST /workloads
    CP->>DB: INSERT
    DB-->>CP: ID: wl-123

    Note over A1,DB: Phase 2: Acquisition
    A1->>CP: GET /workloads
    CP-->>A1: [wl-123]
    A1->>CP: PUT /lease
    CP->>DB: INSERT lease
    Note over A1: Start container

    Note over A1,DB: Phase 3: Normal Operation
    A1->>CP: PUT /lease
    CP->>DB: UPDATE lease
    A1->>CP: POST /events
    CP->>DB: INSERT event

    Note over A1: Phase 4: Failure - Node 1 FAILS

    Note over A2,DB: Phase 5: Recovery
    CP->>DB: DELETE stale lease
    A2->>CP: GET /workloads
    CP-->>A2: [wl-123]
    A2->>CP: PUT /lease
    CP->>DB: INSERT lease
    Note over A2: Reschedule

    Note over O,DB: RTO = 90s (Detection: 30s, Acquisition: 10s, Start: 40s, Margin: 10s)