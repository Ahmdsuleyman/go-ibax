/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/
package types

import (
	"fmt"
	"strings"

	"github.com/IBAX-io/go-ibax/packages/common/crypto"
	"github.com/IBAX-io/go-ibax/packages/conf"
	"github.com/IBAX-io/go-ibax/packages/consts"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack/v5"
)

// Transaction types.
const (
	FirstBlockTxType = iota + 1
	StopNetworkTxType
	SmartContractTxType
)

// FirstBlock is the header of first block transaction
type FirstBlock struct {
	KeyID                 int64
	Time                  int64
	PublicKey             []byte
	NodePublicKey         []byte
	StopNetworkCertBundle []byte
	Test                  int64
	PrivateBlockchain     uint64
}

func (t *FirstBlock) TxType() byte { return FirstBlockTxType }

type StopNetwork struct {
	KeyID           int64
	Time            int64
	StopNetworkCert []byte
}

func (t *StopNetwork) TxType() byte { return StopNetworkTxType }

// Header is contain header data
type Header struct {
	ID          int
	EcosystemID int64
	KeyID       int64
	Time        int64
	NetworkID   int64
	PublicKey   []byte
}

// SmartTransaction is storing smart contract data
type SmartTransaction struct {
	*Header
	TokenEcosystems map[int64]interface{}
	MaxSum          string
	PayOver         string
	Lang            string
	Expedite        string
	SignedBy        int64
	Params          map[string]interface{}
}

func (s *SmartTransaction) TxType() byte { return SmartContractTxType }

func (s *SmartTransaction) WithPrivate(privateKey []byte, internal bool) error {
	var (
		publicKey []byte
		err       error
	)
	if publicKey, err = crypto.PrivateToPublic(privateKey); err != nil {
		log.WithFields(log.Fields{"type": consts.CryptoError, "error": err}).Error("converting node private key to public")
		return err
	}
	s.PublicKey = publicKey
	if internal {
		s.SignedBy = crypto.Address(publicKey)
	}
	if s.NetworkID != conf.Config.LocalConf.NetworkID {
		return fmt.Errorf("error networkid invalid")
	}
	return nil
}

func (s *SmartTransaction) Unmarshal(buffer []byte) error {
	return msgpack.Unmarshal(buffer, s)
}

func (s *SmartTransaction) Marshal() ([]byte, error) {
	return msgpack.Marshal(s)
}

func (t SmartTransaction) Hash() ([]byte, error) {
	b, err := t.Marshal()
	if err != nil {
		return nil, err
	}
	return crypto.DoubleHash(b), nil
}

func (txSmart *SmartTransaction) Validate() error {
	if len(txSmart.Expedite) > 0 {
		expedite, _ := decimal.NewFromString(txSmart.Expedite)
		if expedite.LessThan(decimal.Zero) {
			return fmt.Errorf("expedite fee %s must be greater than 0", expedite)
		}
	}
	if len(strings.TrimSpace(txSmart.Lang)) > 2 {
		return fmt.Errorf(`localization size is greater than 2`)
	}
	if txSmart.NetworkID != conf.Config.LocalConf.NetworkID {
		return fmt.Errorf("error networkid invalid")
	}
	return nil
}
