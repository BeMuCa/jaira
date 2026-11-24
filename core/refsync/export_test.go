package refsync

import "github.com/BeMuCa/jaira/core/outbox"

// SetListBox swaps the outbox listing Flush uses, and returns what puts the
// real one back.
func SetListBox(f func(*outbox.Box) ([]outbox.Entry, error)) (restore func()) {
	old := listBox
	listBox = f
	return func() { listBox = old }
}
