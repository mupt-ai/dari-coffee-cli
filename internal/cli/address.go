package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/mupt-ai/dari-coffee-cli/internal/api"
	"github.com/mupt-ai/dari-coffee-cli/internal/models"
	"github.com/spf13/cobra"
)

func newCheckAddressCommand() *cobra.Command {
	var outputJSON bool
	cmd := &cobra.Command{
		Use:   "check-address ADDRESS",
		Short: "Check whether an address is eligible for Dari Coffee delivery",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			baseURL := apiBaseURL()
			client, err := api.New(baseURL)
			if err != nil {
				return err
			}
			check, err := client.CheckAddress(cmd.Context(), args[0])
			if err != nil {
				return fmt.Errorf("check address with %s: %w", baseURL, err)
			}
			if outputJSON {
				err = json.NewEncoder(cmd.OutOrStdout()).Encode(check)
			} else {
				_, err = fmt.Fprint(cmd.OutOrStdout(), formatAddressCheck(check))
			}
			if err != nil {
				return err
			}
			return addressCheckExitError(check)
		},
	}
	cmd.Flags().BoolVar(&outputJSON, "json", false, "Print machine-readable JSON")
	return cmd
}

func addressCheckExitError(check models.AddressCheck) error {
	switch check.Status {
	case models.AddressCheckStatusAddressUnresolved, models.AddressCheckStatusGeocoderUnavailable:
		return newSilentError(errors.New(string(check.Status)))
	default:
		return nil
	}
}

func formatAddressCheck(check models.AddressCheck) string {
	var b strings.Builder
	fmt.Fprintln(&b, addressCheckMessage(check.Status))
	if check.FormattedAddress != "" {
		fmt.Fprintf(&b, "Address resolved to: %s\n", check.FormattedAddress)
		fmt.Fprintln(&b, "Is this correct? If not, rerun with the complete delivery address.")
	} else {
		fmt.Fprintf(&b, "Address: %s\n", check.Address)
	}
	return b.String()
}

func addressCheckMessage(status models.AddressCheckStatus) string {
	switch status {
	case models.AddressCheckStatusEligible:
		return "eligible: address is inside the default delivery zone. Checkout can proceed.\nCheckout authorizes payment only; Dari captures payment after accepting the order."
	case models.AddressCheckStatusOutsideDeliveryZone:
		return "outside_delivery_zone: address is outside Dari Coffee's default delivery zone, so approval is less likely.\nCheckout authorizes payment only; Dari captures payment after accepting the order."
	case models.AddressCheckStatusAddressUnresolved:
		return "address_unresolved: could not find a confident address match. Provide a complete and specific delivery address, including street number, street name, city, state, and ZIP if available."
	case models.AddressCheckStatusGeocoderUnavailable:
		return "geocoder_unavailable: address validation is temporarily unavailable. Try again later."
	default:
		return fmt.Sprintf("%s: address check returned an unknown status.", status)
	}
}
