package codespace

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/cli/cli/v2/internal/codespaces/rpc"
	"github.com/cli/cli/v2/pkg/liveshare"
	"github.com/spf13/cobra"
)

func newJupyterCmd(app *App) *cobra.Command {
	var selector CodespaceSelector

	jupyterCmd := &cobra.Command{
		Use:   "jupyter",
		Short: "Open a codespace in JupyterLab",
		Args:  noArgsConstraint,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.Jupyter(cmd.Context(), selector)
		},
	}

	selector = AddCodespaceSelector(jupyterCmd, app)

	return jupyterCmd
}

func (a *App) Jupyter(ctx context.Context, selector CodespaceSelector) (err error) {
	// Ensure all child tasks (e.g. port forwarding) terminate before return.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	codespace, err := selector.Select(ctx)
	if err != nil {
		return err
	}

	session, err := startLiveShareSession(ctx, codespace, a, false, "")
	if err != nil {
		return err
	}
	defer safeClose(session, &err)

	a.StartProgressIndicatorWithLabel("Starting JupyterLab on codespace")
	invoker, err := rpc.CreateInvoker(ctx, session)
	if err != nil {
		return err
	}
	defer safeClose(invoker, &err)

	serverPort, serverUrl, err := invoker.StartJupyterServer(ctx)
	if err != nil {
		return fmt.Errorf("failed to start JupyterLab server: %w", err)
	}
	a.StopProgressIndicator()

	// Pass 0 to pick a random port
	listen, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", 0))
	if err != nil {
		return err
	}
	defer listen.Close()
	destPort := listen.Addr().(*net.TCPAddr).Port

	tunnelClosed := make(chan error, 1)
	go func() {
		fwd := liveshare.NewPortForwarder(session, "jupyter", serverPort, true)
		tunnelClosed <- fwd.ForwardToListener(ctx, listen) // always non-nil
	}()

	// Server URL contains an authentication token that must be preserved
	targetUrl := strings.Replace(serverUrl, fmt.Sprintf("%d", serverPort), fmt.Sprintf("%d", destPort), 1)
	err = a.browser.Browse(targetUrl)
	if err != nil {
		return fmt.Errorf("failed to open JupyterLab in browser: %w", err)
	}

	fmt.Fprintln(a.io.Out, targetUrl)

	select {
	case err := <-tunnelClosed:
		return fmt.Errorf("tunnel closed: %w", err)
	case <-ctx.Done():
		return nil // success
	}
}
