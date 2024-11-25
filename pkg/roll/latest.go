// SPDX-License-Identifier: Apache-2.0

package roll

import (
	"context"
	"fmt"
	"io/fs"
)

var (
	ErrNoMigrationFiles   = fmt.Errorf("no migration files found")
	ErrNoMigrationApplied = fmt.Errorf("no migrations applied")
)

// LatestVersionLocal returns the name of the last migration in `dir`, where the
// migration files are lexicographically ordered by filename.
func (m *Roll) LatestVersionLocal(ctx context.Context, dir fs.FS) (string, error) {
	files, err := fs.Glob(dir, "*.json")
	if err != nil {
		return "", fmt.Errorf("reading directory: %w", err)
	}

	if len(files) == 0 {
		return "", ErrNoMigrationFiles
	}

	latest := files[len(files)-1]

	migration, err := openAndReadMigrationFile(dir, latest)
	if err != nil {
		return "", fmt.Errorf("reading migration file %q: %w", latest, err)
	}

	return migration.Name, nil
}

// LatestVersionRemote returns the name of the last migration to have been
// applied to the target schema.
func (m *Roll) LatestVersionRemote(ctx context.Context) (string, error) {
	latestVersion, err := m.State().LatestVersion(ctx, m.Schema())
	if err != nil {
		return "", fmt.Errorf("failed to get latest version: %w", err)
	}

	if latestVersion == nil {
		return "", ErrNoMigrationApplied
	}

	return *latestVersion, nil
}