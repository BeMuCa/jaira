package people

import (
	"errors"
	"fmt"

	"github.com/BeMuCa/jaira/core/gitref"
)

// Repo is where the list travels. *gitref.Repo implements it; a test can
// hand in a fake that loses races on demand.
type Repo interface {
	Usable() error
	Fetch() error
	ReadPeople() ([]byte, string, error)
	WritePeople(content []byte, lease string) (string, error)
}

// attempts bounds how often one write is redone after losing a race. Each
// loss means somebody else wrote the list in between; a few rounds are plenty.
const attempts = 4

// ErrBusy means the write lost every round to somebody else's. Nothing of it
// arrived; running the command again does it on the newest list.
var ErrBusy = errors.New("the people list kept changing on the remote while this write was made")

// Load is the list as the local ref carries it: what the remote last
// confirmed. No ref is an empty list.
func Load(r Repo) (*List, error) {
	content, _, err := r.ReadPeople()
	if errors.Is(err, gitref.ErrNoRef) {
		return &List{}, nil
	}
	if err != nil {
		return nil, err
	}
	return Parse(content), nil
}

// Update applies change to the newest list and pushes it, leased on what it
// read. A lost race is not merged text against text: change runs again on the
// list that won, so two people naming at once both arrive — the second
// write is made on top of the first, never over it.
//
// The write is synchronous and needs the remote. There is no queue to wait
// in: a person runs this by hand and is told at once when it did not arrive,
// which is cheaper than a second store holding unsent names.
func Update(r Repo, change func(*List) (bool, error)) (*List, bool, error) {
	if err := r.Usable(); err != nil {
		return nil, false, err
	}
	for range attempts {
		if err := r.Fetch(); err != nil {
			return nil, false, err
		}
		content, lease, err := r.ReadPeople()
		if err != nil && !errors.Is(err, gitref.ErrNoRef) {
			return nil, false, err
		}
		l := Parse(content)
		changed, err := change(l)
		if err != nil || !changed {
			return l, false, err
		}
		_, err = r.WritePeople(l.Bytes(), lease)
		if err == nil {
			return l, true, nil
		}
		if !errors.Is(err, gitref.ErrRaceLost) {
			return nil, false, err
		}
	}
	return nil, false, fmt.Errorf("%w — run it again", ErrBusy)
}
