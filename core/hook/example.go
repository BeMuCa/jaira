package hook

import _ "embed"

// exampleScript is the notification hook 'jaira hook example' prints. It is
// embedded rather than kept as a string literal so it stays a file that can be
// run, linted and edited as a shell script instead of an escaped blob.
//
//go:embed example/notify.sh
var exampleScript string

// ExampleScript returns a working hook script for someone who has just found
// the empty "hook" field in their settings and has no idea what belongs in it.
//
// It is an example and not an installation: nothing here writes it anywhere or
// turns it on. The script is the user's once it lands on their disk, which is
// the same division Run already makes — jaira calls, the integration is theirs.
func ExampleScript() string { return exampleScript }
