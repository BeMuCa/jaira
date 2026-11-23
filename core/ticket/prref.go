package ticket

import (
	"fmt"
	"strconv"
	"strings"
)

// PRPrefix starts a blocked-by entry that waits on a pull request (a merge
// request on GitLab) rather than on a ticket: pr#12 clears once pull request
// 12 of this repository's forge is merged. It sits in blocked-by beside
// ticket ids, so every reader of that list tells the two apart with ParsePRRef.
const PRPrefix = "pr#"

// maxPR bounds a pull request number to what GitHub's GraphQL Int and
// GitLab's iid both hold.
const maxPR = 1<<31 - 1

// ParsePRRef reads a blocked-by entry as typed or hand-edited. isPR says the
// entry starts with pr# (any case) and is meant as a pull request; err then
// says what is wrong with its number. Only a positive decimal number with
// nothing around it is one: "pr# 12", "pr#12a", "pr#0" and "pr#-1" are not.
func ParsePRRef(entry string) (n int, isPR bool, err error) {
	s := strings.TrimSpace(entry)
	if len(s) < len(PRPrefix) || !strings.EqualFold(s[:len(PRPrefix)], PRPrefix) {
		return 0, false, nil
	}
	num := s[len(PRPrefix):]
	if num == "" || strings.IndexFunc(num, func(r rune) bool { return r < '0' || r > '9' }) >= 0 {
		return 0, true, fmt.Errorf("%q is not pr#<number>: the number of a pull request, digits only", entry)
	}
	v, perr := strconv.Atoi(num)
	if perr != nil || v <= 0 || v > maxPR {
		return 0, true, fmt.Errorf("%q is not pr#<number>: %s is no pull request number", entry, num)
	}
	return v, true, nil
}

// FormatPRRef is how a pull request is written into blocked-by.
func FormatPRRef(n int) string { return PRPrefix + strconv.Itoa(n) }
