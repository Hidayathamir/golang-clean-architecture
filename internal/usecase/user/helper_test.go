package user_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newFakeDB(t *testing.T) (gormDB *gorm.DB, sqlMockDB sqlmock.Sqlmock) {
	t.Helper()

	var sqlDB *sql.DB
	var err error

	sqlDB, sqlMockDB, err = sqlmock.New()
	require.NoError(t, err)

	gormDB, err = gorm.Open(postgres.New(postgres.Config{Conn: sqlDB, PreferSimpleProtocol: true}), &gorm.Config{})
	require.NoError(t, err)

	return gormDB, sqlMockDB
}
