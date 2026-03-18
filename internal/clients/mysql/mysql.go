package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"team-task-manager/internal/clients/mysql/sqlc"
	"team-task-manager/internal/metrics"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

const (
	defaultRequestTimeout = time.Second * 5
	defaultOpenConns      = 50
	defaultConnLifeTime   = time.Minute * 5
	defaultConnIdleTime   = time.Minute * 3
	defaultMonitorDelay   = time.Minute * 1
)

var interruptTxErr = errors.New("tx interrupted")
var defaultTxOpt = &sql.TxOptions{Isolation: sql.LevelRepeatableRead}

type IQuerier interface {
	sqlc.Querier
}

type IDB interface {
	ExecTx(*sql.TxOptions, func(context.Context, IQuerier) error) error
	Querier() IQuerier
	CtxWithCancel() (context.Context, context.CancelFunc)
}

type ILogger interface {
	InfoKV(message string, argsKV ...any)
}

type DB struct {
	ctx  context.Context
	conn *sql.DB
	sqlc *sqlc.Queries
}

type Client struct {
	db IDB
}

func NewMySQLConn(ctx context.Context,
	Host string,
	Port string,
	User string,
	Password string,
	DBName string) (*sql.DB, error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		User, Password, Host, Port, DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(defaultOpenConns)
	db.SetMaxIdleConns(defaultOpenConns)
	db.SetConnMaxLifetime(defaultConnLifeTime)
	db.SetConnMaxIdleTime(defaultConnIdleTime)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func NewClient(ctx context.Context, conn *sql.DB, log ILogger) *Client {
	go func() {
		ticker := time.NewTicker(defaultMonitorDelay)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				stats := conn.Stats()
				log.InfoKV("STATS", "in_use", stats.InUse, "idle", stats.Idle,
					"wait_count", stats.WaitCount, "wait_duration", stats.WaitDuration,
					"max_open_conns", stats.MaxOpenConnections, "timestamp", time.Now())
				metrics.RepordDBStats(stats)
			case <-ctx.Done():
				return
			}
		}
	}()

	return buildClient(&DB{
		ctx:  ctx,
		sqlc: sqlc.New(conn),
		conn: conn,
	})
}

func buildClient(db IDB) *Client {
	return &Client{
		db: db,
	}
}

func (db *DB) CtxWithCancel() (context.Context, context.CancelFunc) {
	return context.WithTimeout(db.ctx, defaultRequestTimeout)
}

func (db *DB) ExecTx(txOpt *sql.TxOptions, withTx func(context.Context, IQuerier) error) (err error) {
	ctx, cancel := db.CtxWithCancel()
	defer cancel()

	var tx *sql.Tx
	tx, err = db.conn.BeginTx(ctx, txOpt)
	if err != nil {
		return
	}

	defer func() {
		errRB := tx.Rollback()
		if errRB != nil && !errors.Is(errRB, sql.ErrTxDone) {
			if err != nil {
				err = fmt.Errorf("ExecTx error: %w; Rollback error: %w", err, errRB)
			} else {
				err = fmt.Errorf("rollback error: %w", errRB)
			}
		}
	}()

	if err = withTx(ctx, db.sqlc.WithTx(tx)); err != nil {
		return
	}

	if err = tx.Commit(); err != nil {
		return
	}

	return
}

func (db *DB) Querier() IQuerier {
	return db.sqlc
}

func isDuplicate(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == mysqlerr.ER_DUP_ENTRY {
		return true
	}
	return false
}

func isLongData(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == mysqlerr.ER_DATA_TOO_LONG {
		return true
	}
	return false
}

func isNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func nullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{
			Valid: false,
		}
	}
	return sql.NullString{
		String: *s,
		Valid:  true,
	}
}

func fromNullString(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func nullInt64(n *int64) sql.NullInt64 {
	if n == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Valid: true, Int64: *n}
}
