package state

import "gopkg.in/guregu/null.v3"

// PostgresStatement - Specific kind of statement that has run one or multiple times
// on the PostgreSQL server.
//
// See also https://www.postgresql.org/docs/9.5/static/pgstatstatements.html
type PostgresStatement struct {
	DatabaseOid       Oid     // OID of database in which the statement was executed
	UserOid           Oid     // OID of user who executed the statement
	NormalizedQuery   string  // Text of a representative statement (normalized)
	Calls             int64   // Number of times executed
	TotalTime         float64 // Total time spent in the statement, in milliseconds
	Rows              int64   // Total number of rows retrieved or affected by the statement
	SharedBlksHit     int64   // Total number of shared block cache hits by the statement
	SharedBlksRead    int64   // Total number of shared blocks read by the statement
	SharedBlksDirtied int64   // Total number of shared blocks dirtied by the statement
	SharedBlksWritten int64   // Total number of shared blocks written by the statement
	LocalBlksHit      int64   // Total number of local block cache hits by the statement
	LocalBlksRead     int64   // Total number of local blocks read by the statement
	LocalBlksDirtied  int64   // Total number of local blocks dirtied by the statement
	LocalBlksWritten  int64   // Total number of local blocks written by the statement
	TempBlksRead      int64   // Total number of temp blocks read by the statement
	TempBlksWritten   int64   // Total number of temp blocks written by the statement
	BlkReadTime       float64 // Total time the statement spent reading blocks, in milliseconds (if track_io_timing is enabled, otherwise zero)
	BlkWriteTime      float64 // Total time the statement spent writing blocks, in milliseconds (if track_io_timing is enabled, otherwise zero)

	// Postgres 9.4+
	QueryId null.Int // Internal hash code, computed from the statement's parse tree

	// Postgres 9.5+
	MinTime    null.Float // Minimum time spent in the statement, in milliseconds
	MaxTime    null.Float // Maximum time spent in the statement, in milliseconds
	MeanTime   null.Float // Mean time spent in the statement, in milliseconds
	StddevTime null.Float // Population standard deviation of time spent in the statement, in milliseconds

	Fingerprint []byte // Fingerprint generated based on the parsetree
}

type PostgresStatementMap map[PostgresStatementKey]PostgresStatement

type DiffedPostgresStatement PostgresStatement

type PostgresStatementKey struct {
	DatabaseOid Oid
	UserOid     Oid
	QueryId     int64
}

func (stmt PostgresStatement) Key() PostgresStatementKey {
	return PostgresStatementKey{DatabaseOid: stmt.DatabaseOid, UserOid: stmt.UserOid, QueryId: stmt.QueryId.Int64}
}

func (curr PostgresStatement) DiffSince(prev PostgresStatement) DiffedPostgresStatement {
	return DiffedPostgresStatement{
		DatabaseOid:       curr.DatabaseOid,
		UserOid:           curr.UserOid,
		NormalizedQuery:   curr.NormalizedQuery,
		QueryId:           curr.QueryId,
		Calls:             curr.Calls - prev.Calls,
		TotalTime:         curr.TotalTime - prev.TotalTime,
		Rows:              curr.Rows - prev.Rows,
		SharedBlksHit:     curr.SharedBlksHit - prev.SharedBlksHit,
		SharedBlksRead:    curr.SharedBlksRead - prev.SharedBlksRead,
		SharedBlksDirtied: curr.SharedBlksDirtied - prev.SharedBlksDirtied,
		SharedBlksWritten: curr.SharedBlksWritten - prev.SharedBlksWritten,
		LocalBlksHit:      curr.LocalBlksHit - prev.LocalBlksHit,
		LocalBlksRead:     curr.LocalBlksRead - prev.LocalBlksRead,
		LocalBlksDirtied:  curr.LocalBlksDirtied - prev.LocalBlksDirtied,
		LocalBlksWritten:  curr.LocalBlksWritten - prev.LocalBlksWritten,
		TempBlksRead:      curr.TempBlksRead - prev.TempBlksRead,
		TempBlksWritten:   curr.TempBlksWritten - prev.TempBlksWritten,
		BlkReadTime:       curr.BlkReadTime - prev.BlkReadTime,
		BlkWriteTime:      curr.BlkWriteTime - prev.BlkWriteTime,
	}
}