package image_test

import (
	"bytes"
	"database/sql"
	"mime/multipart"
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

func newFileHeader(t *testing.T, filename string, content []byte) *multipart.FileHeader {
	t.Helper()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", filename)
	require.NoError(t, err)
	_, err = part.Write(content)
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	reader := multipart.NewReader(&buf, writer.Boundary())
	form, err := reader.ReadForm(int64(buf.Len()))
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, form.RemoveAll())
	})

	return form.File["file"][0]
}