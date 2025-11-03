package store

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/dankski/learn-asyncapi/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func init() {
	godotenv.Load(".envrc")
}

func TestUserStore(t *testing.T) {
	os.Setenv("ENV", string(config.EnvTesting))
	conf, err := config.New()
	require.NoError(t, err)

	db, err := NewPostgressDB(conf)
	require.NoError(t, err)
	defer db.Close()

	m, err := migrate.New(
		fmt.Sprintf("file:///%s/migrations", conf.ProjectRoot),
		conf.DatabaseUrl())
	require.NoError(t, err)

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		require.NoError(t, err)
	}

	userStore := NewUserStore(db)
	user, err := userStore.CreateUser(context.Background(), "john doe", "test@test.com", "testingpasssword")
	require.NoError(t, err)

	require.Equal(t, "john doe", user.Username)
	require.Equal(t, "test@test.com", user.Email)
	require.NoError(t, user.ComparePassword("testingpassword"))
}
