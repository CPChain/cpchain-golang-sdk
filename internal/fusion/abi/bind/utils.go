package bind

import (
	"context"
	"time"

	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/types"
	"github.com/zgljl2012/slog"
)

// WaitMined waits for tx to be mined on the blockchain.
// It stops waiting when the context is canceled.
func WaitMined(ctx context.Context, b DeployBackend, tx *types.Transaction) (*types.Receipt, error) {
	queryTicker := time.NewTicker(time.Second)
	defer queryTicker.Stop()

	for {
		receipt, err := b.TransactionReceipt(ctx, tx.Hash())
		if receipt != nil {
			return receipt, nil
		}
		if err != nil {
			slog.Debug("Receipt retrieval failed", "err", err)
			// fmt.Println("Receipt retrieval failed", "err", err)
		} else {
			slog.Debug("Transaction not yet mined")
			// fmt.Printf("Transaction not yet mined")
		}
		// Wait for the next round.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-queryTicker.C:
		}
	}
}
