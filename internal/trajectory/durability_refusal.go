package trajectory

import (
	"fmt"
	"runtime"
)

// refuseWindowsRootedPublication is the FINAL §B pre-mutation refusal for
// Windows-reachable rooted rows (A3, B5, C1, C2): named, user-visible,
// blocking, evaluated BEFORE any file is created or renamed, so nothing is
// published and no stage is left behind. On every other OS it returns nil and
// the publication proceeds byte-identically to today (AC-DUR-5). The owner
// explicitly rejected the path-based MoveFileEx alternative that would trade
// away os.Root containment; no deviation is open.
func refuseWindowsRootedPublication(feature string) error {
	if runtime.GOOS != "windows" {
		return nil
	}
	return fmt.Errorf("%s requires a durable rooted publication (create or rename plus a directory-entry fsync barrier); Windows cannot provide one while os.Root containment is preserved — refusing before any file is written, nothing was published; run this operation on macOS/Linux", feature)
}
