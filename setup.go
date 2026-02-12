package goqite

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"text/template"

	internalsql "maragu.dev/goqite/internal/sql"
)

type schemaData struct {
	Function string
	Index    string
	Schema   string
	Table    string
	Trigger  string
}

// Setup creates the queue table and related database objects if they don't already exist.
// It is safe to call multiple times (idempotent).
func (q *Queue) Setup(ctx context.Context) error {
	var tmplStr string
	var data schemaData

	switch q.flavor {
	case SQLFlavorSQLite:
		tmplStr = sqliteSchema
		data = schemaData{
			Table:   q.table,
			Trigger: q.table + "_updated_timestamp",
			Index:   q.table + "_queue_priority_created_idx",
		}

	case SQLFlavorPostgreSQL:
		tmplStr = postgresSchema
		function := q.table + "_update_timestamp"
		data = schemaData{
			Table:    q.table,
			Schema:   q.schema,
			Function: function,
			Trigger:  "goqite_updated_timestamp",
			Index:    "goqite_queue_priority_created_idx",
		}

	default:
		return fmt.Errorf("unsupported SQL flavor %d", q.flavor)
	}

	tmpl, err := template.New("schema").Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("parse schema template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("execute schema template: %w", err)
	}

	return internalsql.InTx(ctx, q.db, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, buf.String()); err != nil {
			return fmt.Errorf("execute schema: %w", err)
		}
		return nil
	})
}
