package ipfs

import (
	"context"

	shell "github.com/ipfs/go-ipfs-api"
)

func NewSH(ctx context.Context, addr string) *shell.Shell {
	return shell.NewShell(addr)
}