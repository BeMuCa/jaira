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
# Per agent, from its transcript, assistant messages counted once each (a
# message is logged once per content block, all copies carrying the same id and
# usage, so they are de-duplicated by message id):
#   input_tokens        uncached input, summed over all turns
#   cache_write_tokens  cache_creation_input_tokens, summed over all turns
#   cache_read_tokens   cache_read_input_tokens, summed over all turns
#   output_tokens       output, summed over all turns
#   cost_eq             input + cache_write*1.25 + cache_read*0.1, rounded down:
#                       what the agent consumed, in input-token prices
#   peak_context        context of its last turn: input + cache read + cache write
#                       + output (what the Agent tool reports as total_tokens).
#                       A size, not a cost: one long session has one large peak
#                       but pays for its context again on every turn.
#   duration_ms         last timestamp minus first
#   tool_uses           distinct tool calls, the final hand-back included
#
# Writes <out>/metrics.csv (arm,ticket,run,lane, the token columns above,
# duration_ms,tool_uses,loops,verdict) and <out>/summary.csv (one line per arm,
# token columns summed, peak_context the largest), and prints the summary.
set -euo pipefail
if [ $# -lt 2 ]; then sed -n '2,31p' "$0" >&2; exit 2; fi
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
    | [ to_entries[] | select(.value.type == "assistant" and .value.message.usage != null)
        | {k: (.value.message.id // "line \(.key)"), u: .value.message.usage} ] as $turns
    | ($turns | last | .u) as $last
    | ($turns | reduce .[] as $t ({}; .[$t.k] = $t.u) | [.[]]) as $us
    | def sum(f): $us | map(f // 0) | add // 0;
      [ .[] | .timestamp // empty | ms ] as $ts
    | sum(.input_tokens) as $in | sum(.cache_creation_input_tokens) as $cw | sum(.cache_read_input_tokens) as $cr
    | {arm: $m.arm, ticket: $m.ticket, run: ($m.run | tonumber), step: ($m.step | tonumber), lane: $m.lane,
       input_tokens: $in, cache_write_tokens: $cw, cache_read_tokens: $cr, output_tokens: sum(.output_tokens),
       cost_eq: ($in + $cw * 1.25 + $cr * 0.1 | floor),
       peak_context: (($last.input_tokens // 0) + ($last.cache_read_input_tokens // 0) + ($last.cache_creation_input_tokens // 0) + ($last.output_tokens // 0)),
       duration_ms: (if ($ts | length) > 0 then ($ts | max) - ($ts | min) else 0 end),
       tool_uses: ([ .[] | select(.type == "assistant") | .message.content | arrays | .[] | select(.type == "tool_use") | .id ] | unique | length)}
  ' "$1" 2>/dev/null || true
}

agents=$(find "$@" -name 'agent-*.jsonl' -type f | while read -r f; do per_agent "$f"; done | jq -s 'sort_by(.arm, .ticket, .run, .step)')
[ "$(jq length <<<"$agents")" -gt 0 ] || { echo "metrics.sh: no agent of eval '$id' found under $*" >&2; exit 1; }

jq -r --slurpfile res "$results" '
  ($res[0].runs // []) as $runs
  | "arm,ticket,run,lane,input_tokens,cache_write_tokens,cache_read_tokens,output_tokens,cost_eq,peak_context,duration_ms,tool_uses,loops,verdict",
    (.[] | . as $a
     | ($runs | map(select(.arm == $a.arm and .ticket == $a.ticket and .run == $a.run)) | first // {}) as $r
     | [$a.arm, $a.ticket, $a.run, $a.lane, $a.input_tokens, $a.cache_write_tokens, $a.cache_read_tokens,
        $a.output_tokens, $a.cost_eq, $a.peak_context, $a.duration_ms, $a.tool_uses, ($r.loops // ""), ($r.verdict // "")]
     | @csv)
' <<<"$agents" > "$out/metrics.csv"

jq -r --slurpfile res "$results" '
  ($res[0].runs // []) as $runs
  | group_by(.arm)
  | "arm,runs,agents,input_tokens,cache_write_tokens,cache_read_tokens,output_tokens,cost_eq,cost_eq_per_run,peak_context,duration_ms,tool_uses,loops_per_run,win,lose,tie,split,unjudged",
    (.[] | .[0].arm as $arm
     | ($runs | map(select(.arm == $arm))) as $r
     | ([.[] | "\(.ticket)/\(.run)"] | unique | length) as $n
     | def count(v): ($r | map(select(.verdict == v)) | length);
     [ $arm, $n, length, (map(.input_tokens) | add), (map(.cache_write_tokens) | add),
       (map(.cache_read_tokens) | add), (map(.output_tokens) | add),
       (map(.cost_eq) | add), ((map(.cost_eq) | add) / $n | floor), (map(.peak_context) | max),
       (map(.duration_ms) | add), (map(.tool_uses) | add),
       (if ($r | length) > 0 then ((($r | map(.loops // 0) | add) / ($r | length)) * 100 | round / 100) else "" end),
       count("win"), count("lose"), count("tie"), count("split"),
       ($r | map(select((.verdict // "") == "" or .verdict == "unjudged")) | length) ]
     | @csv)
' <<<"$agents" > "$out/summary.csv"

column -s, -t < "$out/summary.csv" | tr -d '"'
echo "wrote $out/metrics.csv and $out/summary.csv"
