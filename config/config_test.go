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
	require.Equal(t, "8080", cfg.Port)
	require.Equal(t, "5", cfg.CacheCapacity)
	require.Equal(t, "scale", cfg.Resize)
}
