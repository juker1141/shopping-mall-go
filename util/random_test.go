package util

import (
	"fmt"
	"testing"

	"github.com/juker1141/shopping-mall-go/val"
	"github.com/stretchr/testify/require"
)

func TestRandomCellPhone(t *testing.T) {
	n := 10

	for i := 0; i < n; i++ {
		cellphone := RandomCellPhone()
		err := val.ValidateCellPhone(cellphone)
		if err != nil {
			fmt.Println(cellphone)
		}
		require.NoError(t, err)
	}
}
