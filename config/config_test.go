package config

import (
	"context"
	"testing"

	"github.com/c2fo/testify/require"
)

func TestLoadConfig(t *testing.T) {
	ctx := context.Background()

	cfg := New(&ctx)
	require.Equal(t, "INFO", cfg.Logger.Level)
	require.Equal(t, "localhost", cfg.Host)
	require.Equal(t, "postgres", cfg.Port)
	require.Equal(t, "123456", cfg.CacheCapacity)
	require.Equal(t, "test_db", cfg.Resize)
}
