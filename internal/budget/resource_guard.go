package budget

import (
	"context"
	"path/filepath"

	"parley-deck-cli/internal/fsutil"
)

// AcquireResourceGuard reuses the ledger's pinned, verified host-local kernel
// lock for a synchronization-only resource. dir must be a permanent, exclusive
// metadata directory for that resource. No ledger or charge is created. The
// same origin/cache migration restrictions as budget locks apply; removing the
// origin or established local lock while writers exist is unsupported.
func AcquireResourceGuard(ctx context.Context, dir string) (func(), error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := fsutil.MkdirAllResilient(dir, 0700); err != nil {
		return nil, err
	}
	return lock(ctx, filepath.Join(dir, "resource"))
}
