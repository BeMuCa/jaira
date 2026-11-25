# JAIRA_USER narrows "me" to that name and its people-list aliases
Status: Accepted
Date: 2026-11-23
## Context
Under JAIRA_USER the checkout's git user.email and ~/.jaira/identity still counted as "me".
A second agent or teammate borrowing a checkout became its owner, and the owner's tickets passed the ownership gate as theirs.
## Decision
With JAIRA_USER set, identity looks up only that name, and Aliases returns it plus its people-list names.
git user.email and the machine alias file are skipped; whoami marks the alias file "(not read under JAIRA_USER)" (`alias_file_read`).
## Rejected
- Keep adding git user.email under JAIRA_USER: it belongs to whoever configured the checkout.
- Keep reading the alias file under JAIRA_USER: it is the machine owner's, like the email.
## Consequences
A person acting by JAIRA_USER who needs another name records it with `jaira people alias`, not in the alias file.
## Sources
Q030WZ (decision note 2026-11-23 21:02), DVPAPE
