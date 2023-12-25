package sectors

import (
	"fmt"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/storage/sealer/partialfile"
	"github.com/filecoin-project/lotus/storage/sealer/storiface"
	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
	"io"
	"os"
)

var log = logging.Logger("unsealed")

var reversalUnsealedCmd = &cli.Command{
	Name:  "reversal-sealed",
	Usage: "从unsealed文件读取",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "path",
			Usage:    "unsealed文件路径",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "car",
			Usage:    "输出的文件路径",
			Required: true,
		},
		&cli.Int64Flag{
			Name:  "offset",
			Value: 0,
		},
		&cli.Int64Flag{
			Name:  "size",
			Value: 32 << 30,
		},
	},
	Action: func(cctx *cli.Context) error {
		unsealedPath := cctx.String("path")
		carPath := cctx.String("car")
		offset := abi.PaddedPieceSize(cctx.Int64("offset"))
		size := abi.PaddedPieceSize(cctx.Int64("size"))

		log.Infof("Check local %s (+%d,%d)", unsealedPath, offset, size)

		pf, err := partialfile.OpenPartialFile(size, unsealedPath)
		if err != nil {
			return err
		}

		defer pf.Close()

		allocated, err := pf.HasAllocated(storiface.UnpaddedByteIndex(offset.Unpadded()), size.Unpadded())
		if err != nil {
			return err
		}

		if !allocated {
			return fmt.Errorf("miner has unsealed file but not unseal piece, %s (+%d,%d)", unsealedPath, offset, size)
		}

		reader, err := pf.Reader(0, size)
		if err != nil {
			return err
		}

		file, err := os.OpenFile(carPath, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return err
		}

		defer file.Close()

		_, err = io.Copy(file, reader)
		if err != nil {
			return err
		}

		log.Infof("success!!!")
		return nil
	},
}
