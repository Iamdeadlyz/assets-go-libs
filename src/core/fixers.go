package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/trustwallet/assets-go-libs/pkg"
	"github.com/trustwallet/assets-go-libs/pkg/file"
	"github.com/trustwallet/assets-go-libs/pkg/validation"
	"github.com/trustwallet/go-primitives/address"
	"github.com/trustwallet/go-primitives/coin"

	log "github.com/sirupsen/logrus"
)

func (s *Service) FixJSON(file *file.AssetFile) error {
	return pkg.FormatJSONFile(file.Info.Path())
}

func (s *Service) FixETHAddressChecksum(file *file.AssetFile) error {
	if !coin.IsEVM(file.Info.Chain().ID) {
		return nil
	}

	assetDir := filepath.Base(file.Info.Path())

	err := validation.ValidateETHForkAddress(file.Info.Chain(), assetDir)
	if err != nil {
		checksum, e := address.EIP55Checksum(assetDir)
		if e != nil {
			return fmt.Errorf("failed to get checksum: %s", e)
		}

		newName := fmt.Sprintf("blockchains/%s/assets/%s", file.Info.Chain().Handle, checksum)

		if e = os.Rename(file.Info.Path(), newName); e != nil {
			return fmt.Errorf("failed to rename dir: %s", e)
		}

		s.fileService.UpdateFile(file, checksum)

		log.WithField("from", assetDir).
			WithField("to", checksum).
			Debug("Renamed asset")
	}

	return nil
}
