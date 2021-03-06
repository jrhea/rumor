package metadata

import (
	"context"
	"github.com/protolambda/rumor/control/actor/base"
	"github.com/protolambda/rumor/p2p/rpc/methods"
	"github.com/protolambda/rumor/p2p/types"
	"github.com/sirupsen/logrus"
)

type PeerMetadataSetCmd struct {
	*base.Base
	*PeerMetadataState
	SeqNumber methods.SeqNr    `ask:"--seq-number" help:"Seq Number of metadata"`
	Attnets   types.AttnetBits `ask:"--attnets" help:"Attestation nets bitfield as bytes"`
	Merge     bool             `ask:"--merge" help:"If true, only apply non-zero options to state"`
}

func (c *PeerMetadataSetCmd) Default() {
	c.Merge = true
}

func (c *PeerMetadataSetCmd) Help() string {
	return "Set (a part of) the current status and if following the chain or not."
}

func (c *PeerMetadataSetCmd) Run(ctx context.Context, args ...string) error {
	st := methods.MetaData{}
	if c.Merge {
		st = c.PeerMetadataState.Local
	}
	if !c.Merge || c.Attnets != (types.AttnetBits{}) {
		st.Attnets = c.Attnets
	}
	if !c.Merge || c.SeqNumber != 0 {
		st.SeqNumber = c.SeqNumber
	}
	// mutate full metadata at once (TODO: do we need a lock?)
	c.PeerMetadataState.Local = st

	c.Log.WithFields(logrus.Fields{
		"following": c.PeerMetadataState.Following,
		"metadata":  c.PeerMetadataState.Local.Data(),
	}).Info("Metadata settings")
	return nil
}
