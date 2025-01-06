package database

import (
	"crowdfunding/config"
	"testing"

	"github.com/stretchr/testify/require"
)

func init() {
	filename := "../cmd/api/config.yaml"
	err := config.LoadConfig(filename)
	if err != nil {
		panic(err)
	}
}

func TestConnectionPostgresSqlx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, err := ConnectPostgresSqlx(config.Cfg.DB)
		require.Nil(t, err)
		require.NotNil(t, db)
	})
}

//	func TestConnectionPostgresGorm(t *testing.T) {
//		t.Run("success", func(t *testing.T) {
//			db, err := ConnectPostgresGorm(config.Cfg.DB)
//			require.Nil(t, err)
//			require.NotNil(t, db)
//		})
//	}
