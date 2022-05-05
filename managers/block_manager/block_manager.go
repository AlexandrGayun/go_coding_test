package block_manager

import (
	"context"
	"github.com/AlexandrGayun/go_test_task/models/block"
	"log"
)

type Repo struct {
	repo block.CacheRepository
}

func CreateNewRepo(repo block.CacheRepository) *Repo {
	return &Repo{repo: repo}
}

func (r *Repo) GetBlockTransactionsInfo(blockNumber uint64) (*block.Block, error) {
	cachedBlock, rErr := r.repo.CacheRead(context.TODO(), blockNumber)
	if rErr == nil {
		return cachedBlock, nil
	}
	newBlock, err := getBlockData(blockNumber)
	if err != nil {
		return nil, err
	}
	wErr := r.repo.CacheWrite(context.TODO(), newBlock)
	if rErr != nil || wErr != nil {
		log.Println("Cache read/write error:", rErr, wErr)
	}
	return newBlock, nil
}
