package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	pg_query "github.com/lfittl/pg_query_go"
	"github.com/pganalyze/collector/output/pganalyze_collector"
	"github.com/pganalyze/collector/state"
)

func RunExplain(db *sql.DB, connectedDbName string, inputs []state.PostgresQuerySample) (outputs []state.PostgresQuerySample) {
	for _, sample := range inputs {
		// EXPLAIN was already collected, e.g. from auto_explain
		if sample.HasExplain {
			continue
		}

		// TODO: We should run EXPLAIN for other databases here, which means we actually
		// need to split the query samples per database, and then make one connection for each DB
		if sample.Database != "" && sample.Database != connectedDbName {
			continue
		}

		// To be on the safe side never EXPLAIN a statement that can't be parsed,
		// or multiple statements in one (leading to accidental execution)
		parsetree, err := pg_query.Parse(sample.Query)
		if err == nil && len(parsetree.Statements) == 1 && !strings.HasPrefix(sample.Query, QueryMarkerSQL) {
			sample.HasExplain = true
			sample.ExplainSource = pganalyze_collector.QuerySample_STATEMENT_LOG_EXPLAIN_SOURCE
			sample.ExplainFormat = pganalyze_collector.QuerySample_JSON_EXPLAIN_FORMAT
			err = db.QueryRow(QueryMarkerSQL + "BEGIN READ ONLY; EXPLAIN (VERBOSE, FORMAT JSON) " + sample.Query + "; ROLLBACK").Scan(&sample.ExplainOutput)
			if err != nil {
				sample.ExplainError = fmt.Sprintf("%s", err)
			}
		}

		outputs = append(outputs, sample)
	}

	return
}
