package block

import (
	"context"
	"time"
)

type Block struct {
	Number            uint64 `gorm:"primaryKey"`
	TransactionsCount int
	TotalAmount       float64
	CreatedAt         time.Time `json:"-"`
}

type CacheRepository interface {
	CacheRead(ctx context.Context, blockNumber uint64) (*Block, error)
	CacheWrite(ctx context.Context, block *Block) error
}
