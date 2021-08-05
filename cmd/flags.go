package cmd

import "github.com/spf13/cobra"

func addDBFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("db.connection-string", "", "Postgres connection string.")
	cmd.PersistentFlags().String("db.host", "localhost", "Database host")
	cmd.PersistentFlags().String("db.port", "5432", "Database port")
	cmd.PersistentFlags().String("db.sslmode", "disable", "Database sslmode")
	cmd.PersistentFlags().String("db.dbname", "name", "Database name")
	cmd.PersistentFlags().String("db.user", "", "Database user (also allowed via PG_USER env)")
	cmd.PersistentFlags().String("db.password", "", "Database password (also allowed via PG_PASSWORD env)")
}

func addAPIFlags(cmd *cobra.Command) {
	cmd.Flags().String("api.port", "3001", "HTTP API port")
	cmd.Flags().Bool("api.dev-cors", false, "Enable development cors for HTTP API")
	cmd.Flags().String("api.dev-cors-host", "", "Allowed host for HTTP API dev cors")
}

func addMetricsFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().Int64("metrics.port", 9909, "Port on which to serve Prometheus metrics")
}

func addAddressesFlags(cmd *cobra.Command) {
	cmd.Flags().String("addresses.bond", "", "Address of the $BOND token")
	cmd.Flags().StringSlice("addresses.exclude-transfers", []string{}, "Exclude transfers from these addresses when computing holders of BOND")
}
