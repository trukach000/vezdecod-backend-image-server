package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type DbWrapper struct {
	*sql.DB
}

func (db *DbWrapper) BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
	return db.DB.BeginTx(ctx, nil)
}

type Db interface {
	ContextExecutor
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error)
	PingContext(ctx context.Context) error
}
type Tx interface {
	ContextExecutor
	Rollback() error
	Commit() error
}

type ContextExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

func StartTransaction(ctx context.Context) (Tx, error) {
	db, err := GetDatabaseFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return db.BeginTx(ctx, nil)
}

func EndTransaction(ctx context.Context, tx Tx, err error) error {
	if err != nil {
		logrus.Infof("Rollback transaction, tx err: %s", err)
		rlErr := tx.Rollback()
		if rlErr != nil {
			return fmt.Errorf("%w: Rollback error: %s", err, rlErr)
		}
		return err
	}
	cmErr := tx.Commit()
	if cmErr != nil {
		return fmt.Errorf("Commit error: %s", cmErr)
	}
	return nil
}

// GetExecutor returns ContextExecutor from transaction object
// if there is, otherwise it gets Db from context
func GetExecutor(ctx context.Context, tx Tx) (ContextExecutor, error) {
	var ce ContextExecutor
	if tx != nil {
		ce = tx
	} else {
		db, err := GetDatabaseFromContext(ctx)
		if err != nil {
			return nil, err
		}
		ce = db
	}
	return ce, nil
}

func InitDatabase(
	host string,
	port string,
	username string,
	password string,
	dbname string,
) (Db, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&collation=utf8_general_ci&parseTime=true",
		username,
		password,
		host,
		port,
		dbname,
	)

	var err error

	err = waitDbReachable(dsn, "2s", 2*time.Minute)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return &DbWrapper{db}, nil
}

func tryConnectToDB(dsn string) bool {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logrus.Warningf("can't establish connection to db: %s", err)
		return false
	}

	err = db.Ping()
	if err != nil {
		logrus.Warningf("can't ping db: %s", err)
		return false
	}

	_ = db.Close()

	return true
}

func waitDbReachable(dsn string, establishmentRetryPeriod string, maxWait time.Duration) error {
	failed := time.Now().Add(maxWait)
	retryPeriod, err := time.ParseDuration(establishmentRetryPeriod)
	if err != nil {
		return err
	}

	for time.Now().Before(failed) {
		if tryConnectToDB(dsn) {
			logrus.Infof("DB connection established")
			return nil
		}
		time.Sleep(retryPeriod)
	}
	return errors.New(fmt.Sprintf("can't establish connection with given dsn: %s", dsn))
}
