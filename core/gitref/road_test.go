package gitref_test

import (
	"strings"
	"testing"
)

const roadFile = "---\ntitle: Auth rollout\ndescription: one line\ncreated-at: 2026-11-04T10:00:00Z\n---\n\n- 01TEST\n"

// A road travels the way an epic does: on its own ref, on the same fetch, with
// no branch merged. Otherwise the grouping the roadmap draws would exist only
// on the branch that happened to carry it.
func TestRoadArrivesWithoutASharedBranch(t *testing.T) {
	ada, grace := twoClones(t)

	if _, err := ada.WriteRoad("auth-rollout", []byte(roadFile), ""); err != nil {
		t.Fatalf("write: %v", err)
	}
	if out := git(t, grace.Dir, "branch", "-a"); strings.Contains(out, "auth-rollout") {
		t.Fatalf("the test is not testing what it says; grace has a branch:\n%s", out)
	}
	if err := grace.Fetch(); err != nil {
		t.Fatalf("fetch: %v", err)
	}
	got, lease, err := grace.ReadRoad("auth-rollout")
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(got) != roadFile {
		t.Errorf("road came back changed:\n%q", got)
	}
	if lease == "" {
		t.Error("read gave no lease to write back against")
	}
	if names, err := grace.ListRoads(); err != nil || len(names) != 1 || names[0] != "auth-rollout" {
		t.Errorf("ListRoads() = %v, %v; want [auth-rollout]", names, err)
	}
	if at, err := grace.ReadRoadAt("auth-rollout", lease); err != nil || string(at) != roadFile {
		t.Errorf("ReadRoadAt() = %q, %v; want the file as it was at the lease", at, err)
	}
}
