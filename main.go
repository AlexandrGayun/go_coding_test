package main

import (
	"errors"
	"fmt"
	"github.com/AlexandrGayun/go_test_task/config"
	pgconfig "github.com/AlexandrGayun/go_test_task/config/db"
	"github.com/AlexandrGayun/go_test_task/managers/block_manager"
	"github.com/AlexandrGayun/go_test_task/store"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var ErrIncorrectBlockNumber = errors.New("incorrect block number")

type managers struct {
	blockManager *block_manager.Repo
}

func injectManagerMiddleware(mngrs managers) gin.HandlerFunc {
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

func handleRequests(mngrs managers) {
	router := gin.Default()
	// possibly should use versioning for api
	blockGroup := router.Group("/")
	blockGroup.Use(injectManagerMiddleware(mngrs))
	{
		blockGroup.GET("api/block/:block_number/total", getBlockTransactions)
	}

	router.Run("localhost:8080")
}

func main() {
	config.LoadEnv()
	str := store.CreateNewPGStore(pgconfig.DBSettingsAsString())
	mng := block_manager.CreateNewRepo(str)
	mngrs := managers{blockManager: mng}

	handleRequests(mngrs)
}
