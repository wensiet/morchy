sequenceDiagram
    participant O as Operator
    participant M as mctl CLI
    participant CP as Control Plane API
    participant DB as PostgreSQL

    O->>M: 1. Apply YAML
    M->>CP: 2. POST /workloads
    CP->>DB: 3. INSERT workload
    DB-->>CP: 4. Return workload ID
    CP-->>M: 5. Return 201 + ID
    M-->>O: 6. Display workload ID