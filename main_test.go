package main

import (
	"github.com/AlexandrGayun/go_test_task/config"
	pgconfig "github.com/AlexandrGayun/go_test_task/config/db"
	"github.com/AlexandrGayun/go_test_task/managers/block_manager"
	"github.com/AlexandrGayun/go_test_task/store"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var Router *gin.Engine

func TestMain(m *testing.M) {
	config.LoadEnv()
	dbSettings := pgconfig.NewDbSettings("test_")
	str := store.CreateNewPGStore(dbSettings.AsString())
	mng := block_manager.CreateNewRepo(str)
	mngrs := &managers{blockManager: mng}
	Router = setupRouter(mngrs)

	migrator, err := migrate.New(
		"file://db/migrations",
		dbSettings.AsUrl())
	if err != nil {
		log.Println("cannot instantiate migrator", err)
	}
	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		log.Println("cannot run up migration on test database ", err)
	}

	m.Run()

	if err := migrator.Down(); err != nil {
		log.Println("cannot run down migration on test database ", err)
	}
}

func TestGetBlockRoute(t *testing.T) {
	ja := jsonassert.New(t)
	t.Run("should response with valid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/block/1/total", nil)
		w := httptest.NewRecorder()
		Router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		ja.Assertf(w.Body.String(), `{"Number": "<<PRESENCE>>",
						"TransactionsCount": "<<PRESENCE>>",
						"TotalAmount": "<<PRESENCE>>"}`)
	})
	t.Run("should response with error message if parameter invalid", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/block/-1/total", nil)
		w := httptest.NewRecorder()
		Router.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)
		ja.Assertf(w.Body.String(), `{"error_message": "<<PRESENCE>>"}`)
	})
}
