package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/MyCode83/bflask/internal/config"
)

var (
	cfgFile string
	v       = viper.New()
)

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return err
	}
	return nil
}

var rootCmd = &cobra.Command{
	Use:   "bflask",
	Short: "Bruteforce Flask signed session cookies",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfgFile == "" {
			return cmd.Help()
		}
		return crackCmd.RunE(cmd, args)
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		config.Defaults(v)
		if err := config.Setup(v, cfgFile); err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		return nil
	},
}

func init() {
	flags := rootCmd.PersistentFlags()
	flags.StringVar(&cfgFile, "config", "", "config file")
	flags.BoolP("quiet", "q", false, "print only command results")
	if err := v.BindPFlag("quiet", flags.Lookup("quiet")); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(crackCmd)
}
