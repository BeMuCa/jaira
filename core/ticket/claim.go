package ticket

import "time"

// ClaimTTL is how long a claim is honoured without being renewed.
//
// Claims expire rather than being held until released, because there is no server
// to notice that a session died. A lease that outlives its holder would wedge the
// board and require manual unlocking — which, with no coordinator, means editing a
// file by hand. Expiry makes abandonment self-healing.
const ClaimTTL = 30 * time.Minute
