package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCities(t *testing.T) {
	ctx := context.Background()

	cities, err := testQueries.GetCities(ctx)
	require.NoError(t, err)
	// Cities should exist in test database
	assert.NotEmpty(t, cities)
}
