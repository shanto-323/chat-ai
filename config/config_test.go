package config

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	err := godotenv.Load("../.env")
	assert.Nil(t, err)

	cfg, err := LoadConfig()
	assert.Nil(t, err)

	t.Log(cfg)
}
