package lane

import (
	"reflect"
	"testing"
)

func TestMoveSwapsWithNeighbour(t *testing.T) {
	got := Move([]string{"a", "b", "c"}, "c", -1)
	want := []string{"a", "c", "b"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Move(-1) = %v, want %v", got, want)
	}

	got = Move([]string{"a", "c", "b"}, "c", 1)
	want = []string{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Move(+1) = %v, want %v", got, want)
	}
}

func TestMoveAtEitherEndIsNoOp(t *testing.T) {
	ids := []string{"a", "b", "c"}
	if got := Move(ids, "a", -1); !reflect.DeepEqual(got, ids) {
		t.Errorf("moving the first lane left = %v, want unchanged %v", got, ids)
	}
	if got := Move(ids, "c", 1); !reflect.DeepEqual(got, ids) {
		t.Errorf("moving the last lane right = %v, want unchanged %v", got, ids)
	}
}

func TestMoveSingleElementIsNoOp(t *testing.T) {
	ids := []string{"only"}
	if got := Move(ids, "only", -1); !reflect.DeepEqual(got, ids) {
		t.Errorf("Move(-1) on single element = %v, want %v", got, ids)
	}
	if got := Move(ids, "only", 1); !reflect.DeepEqual(got, ids) {
		t.Errorf("Move(+1) on single element = %v, want %v", got, ids)
	}
}

func TestMoveUnknownIDIsNoOp(t *testing.T) {
	ids := []string{"a", "b"}
	if got := Move(ids, "ghost", -1); !reflect.DeepEqual(got, ids) {
		t.Errorf("Move on unknown id = %v, want unchanged %v", got, ids)
	}
}

func TestLoadOrderAbsentFileIsNotAnError(t *testing.T) {
	root := t.TempDir()
	ids, err := LoadOrder(root)
	if err != nil {
		t.Fatalf("LoadOrder on absent file: %v", err)
	}
	if ids != nil {
		t.Errorf("ids = %v, want nil", ids)
	}
}

func TestSaveThenLoadOrderRoundTrips(t *testing.T) {
	root := t.TempDir()
	want := []string{"backlog", "in-progress", "done"}
	if err := SaveOrder(root, want); err != nil {
		t.Fatalf("SaveOrder: %v", err)
	}
	got, err := LoadOrder(root)
	if err != nil {
		t.Fatalf("LoadOrder: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("LoadOrder = %v, want %v", got, want)
	}
}
