package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mupt-ai/dari-coffee-cli/internal/api"
	"github.com/mupt-ai/dari-coffee-cli/internal/models"
	"github.com/spf13/cobra"
)

func newMenuCommand() *cobra.Command {
	var outputJSON bool
	cmd := &cobra.Command{
		Use:   "menu",
		Short: "Show Dari's fixed Coffee CLI menu",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			baseURL := apiBaseURL()
			client, err := api.New(baseURL)
			if err != nil {
				return err
			}
			m, err := client.Menu(cmd.Context())
			if err != nil {
				return fmt.Errorf("fetch menu from %s: %w", baseURL, err)
			}
			if outputJSON {
				err = json.NewEncoder(cmd.OutOrStdout()).Encode(m)
			} else {
				_, err = fmt.Fprint(cmd.OutOrStdout(), formatMenu(m))
			}
			return err
		},
	}
	cmd.Flags().BoolVar(&outputJSON, "json", false, "Print machine-readable JSON")
	return cmd
}

func formatMenu(m models.Menu) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s\n", m.Title)
	if m.Description != "" {
		fmt.Fprintf(&b, "%s\n", m.Description)
	}
	fmt.Fprintln(&b)
	if len(m.Shops) == 0 {
		fmt.Fprintln(&b, "No shops are currently taking Dari orders.")
		return b.String()
	}

	for shopIndex, shop := range m.Shops {
		if shopIndex > 0 {
			fmt.Fprintln(&b)
		}
		fmt.Fprintf(&b, "%s\n", shop.Name)
		if shop.Hours.OpenTime != "" && shop.Hours.CloseTime != "" {
			fmt.Fprintf(&b, "  Hours: %s-%s", shop.Hours.OpenTime, shop.Hours.CloseTime)
			if shop.Hours.LastOrderTime != "" {
				fmt.Fprintf(&b, " (orders until %s)", shop.Hours.LastOrderTime)
			}
			fmt.Fprintln(&b)
		}
		for _, drink := range shop.Drinks {
			fmt.Fprintf(&b, "  %s\n", drink.Name)
			if drink.Description != "" {
				fmt.Fprintf(&b, "    %s\n", drink.Description)
			}
			fmt.Fprintf(&b, "    Sizes: %s\n", formatSizes(drink.Sizes))
		}
	}

	return b.String()
}

func formatSizes(sizes []models.SizePrice) string {
	parts := make([]string, 0, len(sizes))
	for _, size := range sizes {
		parts = append(parts, fmt.Sprintf("%s %s", size.Name, formatUSD(size.PriceCents)))
	}
	return strings.Join(parts, ", ")
}

func formatUSD(cents int) string {
	sign := ""
	if cents < 0 {
		sign = "-"
		cents = -cents
	}
	return fmt.Sprintf("%s$%d.%02d", sign, cents/100, cents%100)
}
