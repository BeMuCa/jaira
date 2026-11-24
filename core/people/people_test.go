package people

import (
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/BeMuCa/jaira/core/gitref"
)

// Every name of a person finds them, however it is written: the jaira name,
// a git user.name, an email, in any case and with stray space. A name the list
// does not know is itself, so an unnamed board behaves exactly as before.
func TestEveryNameOfAPersonMeansThem(t *testing.T) {
	l := Parse([]byte("# comment\n\nberkc: BeMuCa, berk@example.test\nAlexander Sacharov: alex, AlSa\nnot a line\n: no name\n"))
	cases := []struct{ in, want string }{
		{"berkc", "berkc"},
		{"BeMuCa", "berkc"},
		{"bemuca", "berkc"},
		{"  BeMuCa ", "berkc"},
		{"BERK@example.test", "berkc"},
		{"alex", "Alexander Sacharov"},
		{"alexander sacharov", "Alexander Sacharov"},
		{"grace", "grace"},
		{" grace ", "grace"},
		{"", ""},
	}
	for _, c := range cases {
		if got := l.Canonical(c.in); got != c.want {
			t.Errorf("Canonical(%q) = %q, want %q", c.in, got, c.want)
		}
	}
	if !l.Same("BeMuCa", "berk@example.test") || l.Same("BeMuCa", "alex") || l.Same("", "") {
		t.Error("Same disagrees with Canonical")
	}
	var none *List
	if none.Canonical(" x ") != "x" || !none.Same("X", "x") || none.Find("x") != nil {
		t.Error("a nil list is not an empty one")
	}
}

// The file jaira writes reads back as the same list, sorted, so two writes of
// one list are one file.
func TestTheFileReadsBackAsWritten(t *testing.T) {
	l := Parse([]byte("zed: b, A\nalex: AlSa\n"))
	again := Parse(l.Bytes())
	if !slices.EqualFunc(l.Sorted(), again.Sorted(), func(a, b Person) bool {
		return a.Name == b.Name && slices.Equal(a.Aliases, b.Aliases)
	}) {
		t.Fatalf("round trip changed the list:\n%s", l.Bytes())
	}
	if string(again.Bytes()) != string(l.Bytes()) {
		t.Fatal("writing the same list twice gave two files")
	}
	if got := again.Sorted(); got[0].Name != "alex" || !slices.Equal(got[1].Aliases, []string{"A", "b"}) {
		t.Fatalf("not sorted: %+v", got)
	}
}

func TestAlias(t *testing.T) {
	l := &List{}
	if changed, err := l.Alias("berkc", "BeMuCa"); err != nil || !changed {
		t.Fatalf("first alias: %v %v", changed, err)
	}
	// Named by an alias, the person is still found.
	if changed, err := l.Alias("bemuca", "berk@example.test"); err != nil || !changed {
		t.Fatalf("alias by alias: %v %v", changed, err)
	}
	if got := l.Names("berk@example.test"); !slices.Equal(got, []string{"berkc", "BeMuCa", "berk@example.test"}) {
		t.Fatalf("names = %v", got)
	}
	for _, again := range []string{"BeMuCa", "BEMUCA", "berkc"} {
		if changed, err := l.Alias("berkc", again); err != nil || changed {
			t.Errorf("alias %q again: changed=%v err=%v", again, changed, err)
		}
	}
	if changed, err := l.Alias("alex", "alex"); err != nil || changed || l.Find("alex") != nil {
		t.Errorf("a person named only by their own name was added: %v %v", changed, err)
	}
	if _, err := l.Alias("alex", "AlSa"); err != nil {
		t.Fatal(err)
	}
	// A name of somebody else is not taken: neither their jaira name nor an alias.
	for _, theirs := range []string{"berkc", "BeMuCa"} {
		if _, err := l.Alias("alex", theirs); !errors.Is(err, ErrTaken) {
			t.Errorf("alias %q onto alex: err = %v, want ErrTaken", theirs, err)
		}
	}
	for _, bad := range []string{"", "  ", "a:b", "a,b", "a\nb", "#x"} {
		if _, err := l.Alias("alex", bad); !errors.Is(err, ErrBadName) {
			t.Errorf("alias %q: err = %v, want ErrBadName", bad, err)
		}
		if _, err := l.Alias(bad, "zz"); !errors.Is(err, ErrBadName) {
			t.Errorf("person %q: err = %v, want ErrBadName", bad, err)
		}
	}
}

func TestRenameKeepsTheOldNameAsAnAlias(t *testing.T) {
	l := Parse([]byte("BeMuCa: berk@example.test\nalex: AlSa\n"))
	if err := l.Rename("berk@example.test", "berkc"); err != nil {
		t.Fatal(err)
	}
	if got := l.Names("BeMuCa"); !slices.Equal(got, []string{"berkc", "BeMuCa", "berk@example.test"}) {
		t.Fatalf("names = %v", got)
	}
	// To one of their own aliases: it moves up, and is not listed twice.
	if err := l.Rename("berkc", "bemuca"); err != nil {
		t.Fatal(err)
	}
	if got := l.Names("berkc"); !slices.Equal(got, []string{"bemuca", "berkc", "berk@example.test"}) {
		t.Fatalf("names = %v", got)
	}
	if err := l.Rename("bemuca", "AlSa"); !errors.Is(err, ErrTaken) {
		t.Errorf("rename onto another's alias: %v", err)
	}
	if err := l.Rename("nobody", "x"); !errors.Is(err, ErrUnknown) {
		t.Errorf("rename of nobody: %v", err)
	}
	if err := l.Rename("alex", "a:b"); !errors.Is(err, ErrBadName) {
		t.Errorf("rename to a bad name: %v", err)
	}
}

// fakeRepo is a remote that somebody else writes to between this write's
// read and its push, as often as told.
type fakeRepo struct {
	content []byte
	sha     int
	races   int
	other   func(*List)
	offline bool
}

func (f *fakeRepo) Usable() error { return nil }
func (f *fakeRepo) Fetch() error {
	if f.offline {
		return gitref.ErrOffline
	}
	return nil
}
func (f *fakeRepo) ReadPeople() ([]byte, string, error) {
	if f.sha == 0 {
		return nil, "", gitref.ErrNoRef
	}
	return f.content, fmt.Sprint(f.sha), nil
}
func (f *fakeRepo) WritePeople(content []byte, lease string) (string, error) {
	if f.races > 0 {
		f.races--
		l := Parse(f.content)
		f.other(l)
		f.content, f.sha = l.Bytes(), f.sha+1
	}
	if lease != fmt.Sprint(f.sha) && !(lease == "" && f.sha == 0) {
		return "", gitref.ErrRaceLost
	}
	f.content, f.sha = content, f.sha+1
	return fmt.Sprint(f.sha), nil
}

// Two people naming at once both arrive: the write that lost is made again on
// top of the one that won, never pushed over it.
func TestAWriteThatLosesARaceIsMadeAgainOnTopOfTheWinner(t *testing.T) {
	r := &fakeRepo{races: 2, other: func(l *List) { _, _ = l.Alias("grace", fmt.Sprintf("g%d", len(l.Names("grace")))) }}
	l, changed, err := Update(r, func(l *List) (bool, error) { return l.Alias("ada", "ada@example.test") })
	if err != nil || !changed {
		t.Fatalf("update: %v %v", changed, err)
	}
	final := Parse(r.content)
	if final.Canonical("ada@example.test") != "ada" || final.Canonical("g1") != "grace" || final.Canonical("g2") != "grace" {
		t.Fatalf("a write was lost:\n%s", r.content)
	}
	if l.Canonical("g2") != "grace" {
		t.Fatal("the list handed back is not the one that was pushed")
	}

	r = &fakeRepo{races: attempts, other: func(l *List) { _, _ = l.Alias("grace", fmt.Sprintf("g%d", len(l.Names("grace")))) }}
	if _, _, err := Update(r, func(l *List) (bool, error) { return l.Alias("ada", "x") }); !errors.Is(err, ErrBusy) {
		t.Fatalf("losing every round: %v", err)
	}
	if Parse(r.content).Find("x") != nil {
		t.Fatal("a write that never won arrived anyway")
	}
}

// Offline, nothing is written and the caller is told.
func TestAnOfflineWriteIsRefused(t *testing.T) {
	r := &fakeRepo{offline: true}
	if _, _, err := Update(r, func(l *List) (bool, error) { return l.Alias("ada", "x") }); !errors.Is(err, gitref.ErrOffline) {
		t.Fatalf("err = %v", err)
	}
	if r.sha != 0 {
		t.Fatal("written although offline")
	}
}
