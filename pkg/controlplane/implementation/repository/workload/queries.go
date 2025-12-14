package workload

import (
	"fmt"
	"strings"

	"github.com/wernsiet/morchy/pkg/runtime"
)

type queries struct {
}

func (q queries) SelectManyWorkloads(statusEq *string, limits *runtime.ResourceLimits) (string, []any) {
	query := `
	SELECT 
		w.id, w.status, w.created_at,
		s.id, s.image, s.cpu, s.ram, s.command, s.env,
		l.workload_id
	FROM workload w
	JOIN spec s ON w.id = s.id
	LEFT JOIN lease l ON w.id = l.workload_id
	`

	baseOrdering := " ORDER BY w.created_at"

	var args []any
	var conditions []string

	if statusEq != nil {
		conditions = append(conditions, fmt.Sprintf("w.status = $%d", len(args)+1))
		args = append(args, *statusEq)
	}

	if limits != nil {
		conditions = append(conditions, fmt.Sprintf("s.cpu <= $%d", len(args)+1))
		args = append(args, limits.CPU)
		conditions = append(conditions, fmt.Sprintf("s.ram <= $%d", len(args)+1))
		args = append(args, limits.RAM)
	}

	conditions = append(conditions, "l.workload_id IS NULL")

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	return query + baseOrdering, args
}

func (q queries) CreateWorkload() string {
	return `
		INSERT INTO workload (id, status)
		 VALUES ($1, $2)
		 RETURNING id, status, created_at
	`
}

func (q queries) CreateWorkloadSpec() string {
	return `
		INSERT INTO spec (id, image, cpu, ram, command, env)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, image, cpu, ram, command, env
	`
}

func (q queries) GetLease() string {
	return "SELECT node_id, workload_id, created_at, updated_at FROM lease WHERE node_id = $1 AND workload_id = $2"
}

func (q queries) CreateLease() string {
	return ("INSERT INTO lease (node_id, workload_id) VALUES ($1, $2) " +
		"RETURNING node_id, workload_id, created_at, updated_at")
}

func (q queries) DeleteExpiredLeases() string {
	return "DELETE FROM lease WHERE updated_at < NOW() - make_interval(secs => $1)"
}

func (q queries) UpdateLeaseUpdatedAt() string {
	return "UPDATE lease SET updated_at = NOW() WHERE node_id = $1 AND workload_id = $2"
}

func (q queries) UpsertLease() string {
	return `
		INSERT INTO lease (node_id, workload_id)
		VALUES ($1, $2)
		ON CONFLICT (workload_id)
		DO UPDATE
		SET updated_at = NOW()
		WHERE lease.node_id = EXCLUDED.node_id
		RETURNING node_id, workload_id, created_at, updated_at;
	`
}

func (q queries) DeleteLease() string {
	return "DELETE FROM lease WHERE node_id = $1 AND workload_id = $2"
}

func (q queries) SelectWorkloadByID() string {
	return `
		SELECT 
			w.id, w.status, w.created_at,
			s.id, s.image, s.cpu, s.ram, s.command, s.env
		FROM workload w
		JOIN spec s ON w.id = s.id
		WHERE w.id = $1
	`
}

func (q queries) SaveEvent() string {
	return "INSERT INTO event (id, source_id, node_id, payload, produced_at) VALUES ($1, $2, $3, $4, $5)"
}

func (q queries) DeleteLeaseByWorkload() string {
	return "DELETE FROM lease WHERE workload_id = $1"
}

func (q queries) DeleteSpecByID() string {
	return "DELETE FROM spec WHERE id = $1"
}

func (q queries) DeleteWorkload() string {
	return "DELETE FROM workload WHERE id = $1"
}
