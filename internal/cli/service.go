package cli

import (
	"fmt"

	"github.com/mupt-ai/dari-coffee-cli/internal/api"
	"github.com/mupt-ai/dari-coffee-cli/internal/models"
	"github.com/spf13/cobra"
)

func newServiceCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "service",
		Short: "Show Dari Coffee service availability",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			baseURL := apiBaseURL()
			client, err := api.New(baseURL)
			if err != nil {
				return err
			}
			status, err := client.Service(cmd.Context())
			if err != nil {
				return fmt.Errorf("fetch service status from %s: %w", baseURL, err)
			}
			_, err = fmt.Fprint(cmd.OutOrStdout(), formatService(status))
			return err
		},
	}
}

func formatService(status models.ServiceStatus) string {
	if status.CheckoutAvailable {
		return fmt.Sprintf(
			"Dari Coffee is accepting checkout requests.\nCheckout authorizes payment only; Dari captures payment after accepting the order.\nHours: %s-%s\n",
			status.OpenTime,
			status.CloseTime,
		)
	}

	reason := status.Reason
	if reason == "" {
		reason = "unavailable"
	}
	return fmt.Sprintf(
		"Dari Coffee is not accepting checkout requests right now.\nReason: %s\nHours: %s-%s\n",
		reason,
		status.OpenTime,
		status.CloseTime,
	)
}
