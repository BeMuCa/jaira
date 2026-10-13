package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// newJSONCmd gives the role scripts a way to read JSON that every machine
// with the roles has: this binary.
//
// spawn.sh and run-lane.sh read single fields out of Herdr's and jaira's own
// JSON. They did it with python3 -c, and on Windows python3 is either missing
// or the Store's stub, which prints "Python" and exits 49 — the pane id came
// back empty and a tab was left standing with nothing in it (issue #79). jq is
// no better: it is not installed everywhere either. jaira is, by definition.
func newJSONCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "json",
		Short: "Read JSON for shell scripts",
	}
	cmd.AddCommand(newJSONGetCmd())
	return cmd
}

func newJSONGetCmd() *cobra.Command {
	var def string
	cmd := &cobra.Command{
		Use:   "get <path>",
		Short: "Print one field of the JSON on stdin",
		Long: `Reads one JSON document from stdin and prints the value at <path>: keys
joined by dots, a number for an array element (result.pane.pane_id, items.0).

A string prints as it is, without quotes; null prints nothing; anything else
prints as compact JSON. A path that is not there, or input that is no JSON,
exits 5, unless --default names what to print instead.`,
		Args: exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := jsonAt(cmd.InOrStdin(), args[0])
			if err != nil {
				if cmd.Flags().Changed("default") {
					fmt.Fprintln(cmd.OutOrStdout(), def)
					return nil
				}
				return err
			}
			out, err := jsonText(v)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), out)
			return nil
		},
	}
	cmd.Flags().StringVar(&def, "default", "", "print this instead when the path is not there")
	return cmd
}

// jsonAt decodes r and walks path. Unreadable input is a failure like a
// missing key: either way the script has no value to go on with.
func jsonAt(r io.Reader, path string) (any, error) {
	var v any
	dec := json.NewDecoder(r)
	dec.UseNumber()
	if err := dec.Decode(&v); err != nil {
		return nil, fail(ExitNotFound, "not_found", "no JSON on stdin: %v", err)
	}
	for _, key := range strings.Split(path, ".") {
		switch node := v.(type) {
		case map[string]any:
			next, ok := node[key]
			if !ok {
				return nil, fail(ExitNotFound, "not_found", "%s: no key %q", path, key)
			}
			v = next
		case []any:
			i, err := strconv.Atoi(key)
			if err != nil || i < 0 || i >= len(node) {
				return nil, fail(ExitNotFound, "not_found", "%s: no element %q", path, key)
			}
			v = node[i]
		default:
			return nil, fail(ExitNotFound, "not_found", "%s: %q is inside a value that has no fields", path, key)
		}
	}
	return v, nil
}

// jsonText renders a value for a shell: a string bare, null as nothing,
// everything else compact JSON with its characters left as they are.
func jsonText(v any) (string, error) {
	switch x := v.(type) {
	case nil:
		return "", nil
	case string:
		return x, nil
	}
	return compactJSON(v)
}

// compactJSON is one line of JSON without \u escapes for <, > and &: the
// text lands in a prompt a person may read.
func compactJSON(v any) (string, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return "", err
	}
	return strings.TrimSuffix(b.String(), "\n"), nil
}
