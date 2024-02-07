package db_test

import (
	"context"
	"testing"

	"github.com/ShvetsovYura/metrics-collector/internal/storage/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddCounter(t *testing.T) {
	db, err := db.NewDBPool(context.Background(), "postgres://mc_user:Dthcbz@localhost:5432/mc_db?sslmode=disable")

	require.NoError(t, err)
	err1 := db.SetGauge("nono", 2.3234234)
	require.NoError(t, err1)
	assert.Equal(t, 2, 2)

}
