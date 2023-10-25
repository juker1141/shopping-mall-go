package db

import (
	"context"
	"testing"
)

func TestUpdateCartTx(t *testing.T) {
	arg := UpdateCartTxParams{}
	testStore.UpdateCartTx(context.Background(), arg)
}
