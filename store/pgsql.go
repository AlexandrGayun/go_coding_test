package store

import (
	"context"
	models "github.com/AlexandrGayun/go_test_task/models/block"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

type Pgsql struct {
	q *gorm.DB
}

func CreateNewPGStore(pgconfig string) *Pgsql {
	db, err := gorm.Open(postgres.Open(pgconfig), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Println("Cannot instantiate DB", err)
	}
	return &Pgsql{q: db}
}

func (db *Pgsql) CacheWrite(ctx context.Context, block *models.Block) error {
	result := db.q.WithContext(ctx).Create(block)
	err := result.Error
	if err != nil {
		return err
	}
	return nil
}

func (db *Pgsql) CacheRead(ctx context.Context, blockNumber uint64) (*models.Block, error) {
	res := &models.Block{}
	result := db.q.WithContext(ctx).First(res, blockNumber)
	err := result.Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
