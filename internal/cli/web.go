package cli

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/outbox"
	"github.com/BeMuCa/jaira/core/refsync"
	"github.com/BeMuCa/jaira/core/web"
)

// webFetchEvery is how often 'jaira web' fetches. A minute rather than the
// CLI's ten: the page is open for as long as somebody is sorting, and a
// teammate's drop should show up while they still look at it.
const webFetchEvery = time.Minute

func newWebCmd() *cobra.Command {
	var port int
	var noOpen bool
	cmd := &cobra.Command{
		Use:   "web",
		Short: "Open the review portal in the browser: epics, their progress, and the tickets in none",
		Long: `Starts a page on 127.0.0.1 and opens it in the browser. On the right are the
epics, each with its description and how many of its tickets are done; on the
left are the open tickets that belong to no epic. Drag a ticket onto an epic to
put it there.

There is no server: the page lives as long as this command runs, and it reads
and writes the same refs the CLI does. A drop writes the epic's ref and pushes
it at once; every minute the command fetches, so what a teammate sorted arrives
without a reload. The page shows when it last synchronised. A push that fails
stays queued and goes out with the next change.

The page answers only on loopback and only to requests that name a loopback
host, so neither another machine nor another web site can reach it.

Stop it with Ctrl-C.`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			l, err := web.Listen(port)
			if err != nil {
				return fail(ExitError, "listen", "cannot listen on 127.0.0.1:%d: %v", port, err)
			}
			srv := &web.Server{Source: newWebSource(s)}
			if refs.Usable() == nil {
				root := s.Root
				srv.Fetch = func() error { return fetchInChild(root) }
				srv.Push = webPush
			}
			url := "http://" + l.Addr().String() + "/"
			fmt.Fprintf(cmd.OutOrStdout(), "jaira web: %s  (Ctrl-C to stop)\n", url)
			if !noOpen {
				if err := web.OpenBrowser(url); err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "jaira web: could not open a browser (%v); open the address above\n", err)
				}
			}

			stop := make(chan struct{})
			go srv.Loop(webFetchEvery, stop)
			hs := &http.Server{Handler: srv.Handler(), ReadHeaderTimeout: 10 * time.Second}
			served := make(chan error, 1)
			go func() { served <- hs.Serve(l) }()

			sig := make(chan os.Signal, 1)
			signal.Notify(sig, os.Interrupt)
			defer signal.Stop(sig)
			select {
			case <-sig:
			case err = <-served:
			}
			close(stop)
			_ = hs.Close()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				return err
			}
			return nil
		},
	}
	cmd.Flags().IntVar(&port, "port", 0, "port on 127.0.0.1 to listen on (default: a free one)")
	cmd.Flags().BoolVar(&noOpen, "no-open", false, "print the address but do not open a browser")
	return cmd
}

// webPush sends what a drop queued and says whether it reached the remote.
//
// refs.Flush alone does not: an unreachable remote is the state the outbox
// exists for, so Flush answers it with no error and leaves the entry queued,
// and an entry git refused is a report, not an error. A command prints those
// reports as warnings (flushRefs); the page has nowhere to print, so every way
// of not arriving becomes the error it shows as "push failed".
func webPush() error {
	reports, err := refs.Flush()
	if err != nil {
		return err
	}
	for _, r := range reports {
		if r.Outcome != outbox.Sent {
			_, name := refSubject(r.Kind, r.ID)
			if r.Err != nil {
				return fmt.Errorf("%s is %s: %v", name, r.Outcome, r.Err)
			}
			return fmt.Errorf("%s is %s", name, r.Outcome)
		}
	}
	if queued, err := refs.Box.List(); err == nil && len(queued) > 0 {
		return fmt.Errorf("the remote was not reached; %d change(s) wait to be sent", len(queued))
	}
	return nil
}

// fetchInChild runs 'jaira fetch' as its own process, exactly as the
// background fetch after every command does. The child is what writes what a
// fetch brings and raises the notifications, so the page's fetch and a typed
// one cannot drift apart, and this long-running process holds no lock across
// the network round trip.
func fetchInChild(root string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	c := exec.Command(exe, "-C", root, "fetch")
	// The child would otherwise schedule a background fetch of its own.
	c.Env = append(os.Environ(), refsync.FetchDisableEnv+"=1")
	var stderr bytes.Buffer
	c.Stderr = &stderr
	if err := c.Run(); err != nil {
		return fmt.Errorf("jaira fetch: %v: %s", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}
