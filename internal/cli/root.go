package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newRootCommand(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "dari-coffee",
		Short:         "Order coffee from Dari's FiDi Coffee CLI",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.AddCommand(newVersionCommand(version))
	cmd.AddCommand(newMenuCommand())
	cmd.AddCommand(newOrderCommand())
	cmd.AddCommand(newCheckAddressCommand())
	cmd.AddCommand(newServiceCommand())

	return cmd
}

func Execute(version string) error {
	cmd := newRootCommand(version)
	if err := cmd.Execute(); err != nil {
		if !isSilentError(err) {
			fmt.Fprintln(cmd.ErrOrStderr(), err)
		}
		return err
	}
	return nil
}

type silentError struct {
	err error
}

func (e silentError) Error() string {
	return e.err.Error()
}

func newSilentError(err error) error {
	return silentError{err: err}
}

func isSilentError(err error) bool {
	_, ok := err.(silentError)
	return ok
}

func apiBaseURL() string {
	return defaultAPIBaseURL()
}
