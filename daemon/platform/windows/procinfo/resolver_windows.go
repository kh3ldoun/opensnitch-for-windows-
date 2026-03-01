//go:build windows

package procinfo

import (
	"fmt"

	"github.com/evilsocket/opensnitch/daemon/platform/windows/wfp"
)

type Resolver struct{}

func NewResolver() *Resolver { return &Resolver{} }

func (r *Resolver) Resolve(pid uint32) (wfp.ProcessMetadata, error) {
	// TODO: replace with full Toolhelp32 + token SID + PEB command-line extraction.
	if pid == 0 {
		return wfp.ProcessMetadata{}, fmt.Errorf("invalid pid")
	}
	return wfp.ProcessMetadata{PID: pid}, nil
}
