package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/MyCode83/bflask/pkg/bflask"
)

var decodeOpts struct {
	cookie string
	raw    bool
}

var decodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "Decode a Flask session cookie payload",
	RunE: func(cmd *cobra.Command, args []string) error {
		if decodeOpts.cookie == "" {
			return errors.New("--cookie is required")
		}

		payload, err := bflask.DecodeCookie(decodeOpts.cookie)
		if err != nil {
			return err
		}
		if decodeOpts.raw || v.GetBool("quiet") {
			fmt.Fprintln(cmd.OutOrStdout(), string(payload))
			return nil
		}
		fmt.Fprintln(cmd.OutOrStdout(), bflask.PrettyPayload(payload))
		return nil
	},
}

func init() {
	flags := decodeCmd.Flags()
	flags.StringVarP(&decodeOpts.cookie, "cookie", "c", "", "Flask signed session cookie")
	flags.BoolVar(&decodeOpts.raw, "raw", false, "print decoded payload without JSON formatting")

	rootCmd.AddCommand(decodeCmd)
}
