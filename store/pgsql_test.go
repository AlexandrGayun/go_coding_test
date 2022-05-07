package store_test

import (
	"context"
	"errors"
	"github.com/AlexandrGayun/go_test_task/config"
	pgconfig "github.com/AlexandrGayun/go_test_task/config/db"
	"github.com/AlexandrGayun/go_test_task/models/block"
	"github.com/AlexandrGayun/go_test_task/store"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
	"log"
	"testing"
	"time"
)

type testPackSuite struct {
	migrator *migrate.Migrate
	pgsql    *store.Pgsql
}

var testPack = &testPackSuite{}

func TestMain(m *testing.M) {
	config.LoadEnv("../.env")
	testPack.initializeTestInstances()
	if err := testPack.migrator.Up(); err != nil && err != migrate.ErrNoChange {
		log.Println("cannot run up migration on test database ", err)
	}
	m.Run()
	if err := testPack.migrator.Down(); err != nil {
		log.Println("cannot run down migration on test database ", err)
	}
}

func (tp *testPackSuite) initializeTestInstances() {
	dbSettings := pgconfig.NewDbSettings("test_")
	tp.pgsql = store.CreateNewPGStore(dbSettings.AsString())
	m, err := migrate.New(
		"file://../db/migrations",
		dbSettings.AsUrl())
	if err != nil {
		log.Println("cannot instantiate migrator", err)
	}
	tp.migrator = m
}

func errChecker(t *testing.T, err error, expectedErrCode string) {
	t.Helper()
	var pgErrorCode = ""
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErrorCode = pgErr.Code
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			pgErrorCode = "recordNotFound"
		}
	}
	if pgErrorCode != expectedErrCode {
		t.Errorf("get unexpected error code %q, wanted %q error", err, expectedErrCode)
	}
}

func TestPgsql_CacheWrite(t *testing.T) {
	cacheWriteTests := []struct {
		writeBlock *block.Block
		errCode    string
	}{
		{writeBlock: &block.Block{Number: 1, TransactionsCount: 10, TotalAmount: 100, CreatedAt: time.Now()}, errCode: ""},
		{writeBlock: &block.Block{Number: 2, TransactionsCount: 20, TotalAmount: 200, CreatedAt: time.Now()}, errCode: ""},
		{writeBlock: &block.Block{Number: 2, TransactionsCount: 30, TotalAmount: 300, CreatedAt: time.Now()}, errCode: "23505"},
	}
	for _, tt := range cacheWriteTests {
		t.Run("writing to cache db", func(t *testing.T) {
			err := testPack.pgsql.CacheWrite(context.Background(), tt.writeBlock)
			errChecker(t, err, tt.errCode)
		})
	}
}

func TestPgsql_CacheRead(t *testing.T) {
	cacheReadTests := []struct {
		readBlockId uint64
		errCode     string
	}{
		{1, ""},
		{2, ""},
		{3, "recordNotFound"},
	}
	for _, tt := range cacheReadTests {
		t.Run("reading from cache db", func(t *testing.T) {
			_, err := testPack.pgsql.CacheRead(context.Background(), tt.readBlockId)
			errChecker(t, err, tt.errCode)
		})
	}
}
