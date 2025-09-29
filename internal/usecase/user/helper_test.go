package user_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func newFakeDB(t *testing.T) (gormDB *gorm.DB, sqlMockDB sqlmock.Sqlmock) {
	t.Helper()

	var sqlDB *sql.DB
	var err error

	sqlDB, sqlMockDB, err = sqlmock.New()
	require.NoError(t, err)

	gormDB, err = gorm.Open(mysql.New(mysql.Config{Conn: sqlDB, SkipInitializeWithVersion: true}), &gorm.Config{})
	require.NoError(t, err)

	return gormDB, sqlMockDB
}
