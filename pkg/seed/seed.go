package seed

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"gorm.io/gorm"
)

func Seed(cfg *config.Config, db *gorm.DB) {
	dir := cfg.GetDatabaseSeeds()

	entries, err := os.ReadDir(dir)
	errkit.PanicIfErr(err)

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)

	for _, name := range files {
		content, err := os.ReadFile(filepath.Join(dir, name))
		errkit.PanicIfErr(err)

		up := upStatement(string(content))
		if strings.TrimSpace(up) == "" {
			continue
		}

		err = db.Exec(up).Error
		errkit.PanicIfErr(err)

		logkit.Logger.Infof("Applied seed %s", name)
	}
}

func upStatement(content string) string {
	const marker = "-- +migrate Up"
	idx := strings.Index(content, marker)
	if idx < 0 {
		return content
	}
	rest := content[idx+len(marker):]
	if down := strings.Index(rest, "-- +migrate Down"); down >= 0 {
		rest = rest[:down]
	}
	return rest
}
