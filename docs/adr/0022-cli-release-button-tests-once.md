# A release is a button that tests once and tags on merge
Status: Accepted
Date: 2026-10-12

## Context
v0.3.7 ran the test suite four times: cut, CI on the release PR, CI on master, goreleaser before-hook.
The release PR only moves files under core/release/notes/**, so the tree is byte-identical to the tested one.

## Decision
The release-notes workflow cuts notes into notes/<version>/, runs go test ./... and the shipped-manifest check once,
then opens a PR from release/<version>. Merging it runs release-tag, which tags v<version> and starts release.yaml.
ci.yaml skips changes limited to core/release/notes/**; the tag build runs no tests.

## Rejected
- Tagging by hand from a maintainer's machine (no rule, nothing tested at cut time).
- Re-running the full matrix on the release PR, on master and in goreleaser.

## Consequences
A release PR that touches anything outside notes/** gets the full CI run again.
Nothing lands on master directly; the maintainer merges the release PR.

## Sources
AR97K9, 776466, 5HDAJZ, 0K8Y6S
