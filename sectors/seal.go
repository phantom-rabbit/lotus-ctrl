package sectors

import (
	"github.com/urfave/cli/v2"
)

var SealCmd = &cli.Command{
	Name:  "seal",
	Usage: "Manage chain",
	Subcommands: []*cli.Command{
		reversalUnsealedCmd,
	},
}
