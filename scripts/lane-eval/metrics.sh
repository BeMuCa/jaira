#!/usr/bin/env bash
# Turns one lane-eval run into numbers.
#
#   metrics.sh <out> <transcript-dir>...
#
# Reads <out>/results.json (written by the workflow's last step) and every
# agent-*.jsonl below the transcript directories. An agent counts when the
# first line of its prompt is the marker the workflow puts there:
#   [lane-eval id=<basename of out> arm=<arm> ticket=<t> run=<n> step=<k> lane=<lane>]
# Agents without it (setup, judges, the recorder) are not part of an arm's cost.
#
# Per agent, from its transcript:
#   tokens      context of its last turn: input + cache read + cache write + output
#               (what the Agent tool reports as total_tokens, to within ~1%)
#   duration_ms last timestamp minus first
#   tool_uses   distinct tool calls, the final hand-back included
#
# Writes <out>/metrics.csv (arm,ticket,run,lane,tokens,duration_ms,tool_uses,loops,verdict)
# and <out>/summary.csv (one line per arm), and prints the summary.
set -euo pipefail
if [ $# -lt 2 ]; then sed -n '2,20p' "$0" >&2; exit 2; fi
out=$(cd "$1" && pwd); shift
id=$(basename "$out")
results="$out/results.json"
[ -f "$results" ] || { echo "metrics.sh: $results is missing" >&2; exit 1; }

per_agent() {
  jq -s --arg id "$id" '
    def ms: (sub("\\.[0-9]+Z$"; "Z") | fromdateiso8601) * 1000
            + ((capture("\\.(?<f>[0-9]+)Z$").f // "0") | .[0:3] | tonumber);
    def text: if type == "string" then . else (map(select(.type == "text") | .text) | join("\n")) end;
    (map(select(.type == "user")) | first | .message.content | text) as $prompt
    | ($prompt | capture("^\\s*\\[lane-eval id=(?<id>[^ ]+) arm=(?<arm>[^ ]+) ticket=(?<ticket>[^ ]+) run=(?<run>[0-9]+) step=(?<step>[0-9]+) lane=(?<lane>[^\\]]+)\\]")?) as $m
    | select($m != null and $m.id == $id)
    | (map(select(.type == "assistant" and .message.usage != null)) | last | .message.usage) as $u
    | [ .[] | .timestamp // empty | ms ] as $ts
    | {arm: $m.arm, ticket: $m.ticket, run: ($m.run | tonumber), step: ($m.step | tonumber), lane: $m.lane,
       tokens: (($u.input_tokens // 0) + ($u.cache_read_input_tokens // 0) + ($u.cache_creation_input_tokens // 0) + ($u.output_tokens // 0)),
       duration_ms: (if ($ts | length) > 0 then ($ts | max) - ($ts | min) else 0 end),
       tool_uses: ([ .[] | select(.type == "assistant") | .message.content | arrays | .[] | select(.type == "tool_use") | .id ] | unique | length)}
  ' "$1" 2>/dev/null || true
}

agents=$(find "$@" -name 'agent-*.jsonl' -type f | while read -r f; do per_agent "$f"; done | jq -s 'sort_by(.arm, .ticket, .run, .step)')
[ "$(jq length <<<"$agents")" -gt 0 ] || { echo "metrics.sh: no agent of eval '$id' found under $*" >&2; exit 1; }

jq -r --slurpfile res "$results" '
  ($res[0].runs // []) as $runs
  | "arm,ticket,run,lane,tokens,duration_ms,tool_uses,loops,verdict",
    (.[] | . as $a
     | ($runs | map(select(.arm == $a.arm and .ticket == $a.ticket and .run == $a.run)) | first // {}) as $r
     | [$a.arm, $a.ticket, $a.run, $a.lane, $a.tokens, $a.duration_ms, $a.tool_uses, ($r.loops // ""), ($r.verdict // "")]
     | @csv)
' <<<"$agents" > "$out/metrics.csv"

jq -r --slurpfile res "$results" '
  ($res[0].runs // []) as $runs
  | group_by(.arm)
  | "arm,runs,agents,tokens,tokens_per_run,duration_ms,tool_uses,loops_per_run,win,lose,tie,split,unjudged",
    (.[] | .[0].arm as $arm
     | ($runs | map(select(.arm == $arm))) as $r
     | ([.[] | "\(.ticket)/\(.run)"] | unique | length) as $n
     | def count(v): ($r | map(select(.verdict == v)) | length);
     [ $arm, $n, length, (map(.tokens) | add), ((map(.tokens) | add) / $n | floor),
       (map(.duration_ms) | add), (map(.tool_uses) | add),
       (if ($r | length) > 0 then ((($r | map(.loops // 0) | add) / ($r | length)) * 100 | round / 100) else "" end),
       count("win"), count("lose"), count("tie"), count("split"),
       ($r | map(select((.verdict // "") == "" or .verdict == "unjudged")) | length) ]
     | @csv)
' <<<"$agents" > "$out/summary.csv"

column -s, -t < "$out/summary.csv" | tr -d '"'
echo "wrote $out/metrics.csv and $out/summary.csv"
