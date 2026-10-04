package app

import (
	"parley-deck-cli/internal/protocol"
	"reflect"
)

// Every independently callable action refreshes its captured identities. This
// also covers a reduction during the same driver tick (e.g. signoff requests).
func (o driverImplOps) quotaCurrent() (driverImplOps, error) {
	ids, _, err := protocol.QuotaMembers(o.ideaDir, o.base.Idea.Participants)
	if err != nil {
		return o, err
	}
	if !reflect.DeepEqual(ids, o.base.Idea.Participants) {
		o = o.WithParticipants(ids).(driverImplOps)
	}
	return o, nil
}
