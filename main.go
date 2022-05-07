package main

import (
	"fmt"
	"github.com/AlexandrGayun/go_test_task/config"
	pgconfig "github.com/AlexandrGayun/go_test_task/config/db"
	"github.com/AlexandrGayun/go_test_task/managers/block_manager"
	"github.com/AlexandrGayun/go_test_task/store"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"net/http"
	"strconv"
)

const (
	ErrIncorrectBlockNumber = ValidationError("incorrect block number")
)

type ValidationError string

func (e ValidationError) Error() string {
	return string(e)
}

type managers struct {
	blockManager *block_manager.Repo
}

func injectManagerMiddleware(mngrs *managers) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("blockManager", mngrs.blockManager)
		ctx.Next()
	}
}

func errorResponse(c *gin.Context, err error) {
	c.IndentedJSON(http.StatusBadRequest, gin.H{
		"error_message": fmt.Sprintf("%s", err),
	})
}

func getBlockTransactions(c *gin.Context) {
	blockId, err := strconv.ParseUint(c.Param("block_number"), 10, 64)
	if err != nil {
		errorResponse(c, ErrIncorrectBlockNumber)
		return
	}
	blockManager := c.MustGet("blockManager").(*block_manager.Repo)
	result, err := blockManager.GetBlockTransactionsInfo(blockId)
	if err != nil {
		errorResponse(c, err)
	} else {
		c.IndentedJSON(http.StatusOK, result)
	}
}

func setupRouter(mngrs *managers) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// possibly should use versioning for api
	blockGroup := router.Group("/")
	blockGroup.Use(injectManagerMiddleware(mngrs))
	{
		blockGroup.GET("api/block/:block_number/total", getBlockTransactions)
	}
	return router
}

func runMigrations(settingsUrl string) {
	migrator, err := migrate.New(
		"file://db/migrations",
		settingsUrl)
	if err != nil {
		log.Println("cannot instantiate migrator", err)
	}
	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		log.Println("cannot run up migration on database ", err)
	}
}

func main() {
	config.LoadEnv()
	dbSettings := pgconfig.NewDbSettings("")
	// it's not a good idea to autorun migrate. But it's ok now. Just for easier setup
	runMigrations(dbSettings.AsUrl())
	str := store.CreateNewPGStore(dbSettings.AsString())
	mng := block_manager.CreateNewRepo(str)
	mngrs := &managers{blockManager: mng}
	router := setupRouter(mngrs)
	err := router.Run(":8080")
	if err != nil {
		log.Fatalln("cannot run the server", err)
	}
}
