sequenceDiagram
    participant A as Agent
    participant CP as Control Plane API
    participant DB as PostgreSQL
    participant D as Docker Runtime

    Note over A,D: Polling Phase
    A->>CP: 1. GET /workloads?cpu=4000&ram=8192
    CP->>DB: 2. Query available workloads
    DB-->>CP: 3. Return workload list
    CP-->>A: 4. Return workloads

    Note over A,D: Lease Acquisition
    A->>CP: 5. PUT /workloads/{id}/lease?node_id=123
    CP->>DB: 6. INSERT lease
    DB-->>CP: 7. Lease confirmed
    CP-->>A: 8. Return 200 OK

    Note over A,D: Container Creation
    A->>D: 9. Create container
    D-->>A: 10. Container started