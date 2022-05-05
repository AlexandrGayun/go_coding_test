package store

import (
	"context"
	models "github.com/AlexandrGayun/go_test_task/models/block"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type pgsql struct {
	q *gorm.DB
}

func CreateNewPGStore(pgconfig string) *pgsql {
	db, err := gorm.Open(postgres.Open(pgconfig), &gorm.Config{})
	if err != nil {
		log.Println("Cannot instantiate DB", err)
	}
	return &pgsql{q: db}
}

func (db *pgsql) CacheWrite(ctx context.Context, block *models.Block) error {
	result := db.q.WithContext(ctx).Create(block)
	err := result.Error
	if err != nil {
		return err
	}
	return nil
}

func (db *pgsql) CacheRead(ctx context.Context, blockNumber uint64) (*models.Block, error) {
	res := &models.Block{}
	result := db.q.WithContext(ctx).First(res, blockNumber)
	err := result.Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
