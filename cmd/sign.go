package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"bflask/pkg/bflask"
)

var signOpts struct {
	secret  string
	payload string
	salt    string
	digest  string
}

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Sign a Flask session cookie payload",
	RunE: func(cmd *cobra.Command, args []string) error {
		if signOpts.secret == "" {
			return errors.New("--secret is required")
		}
		if signOpts.payload == "" {
			return errors.New("--payload is required")
		}

		cookie, err := bflask.SignCookie([]byte(signOpts.payload), signOpts.secret, signOpts.salt, signOpts.digest)
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), cookie)
		return nil
	},
}

func init() {
	flags := signCmd.Flags()
	flags.StringVarP(&signOpts.secret, "secret", "k", "", "SECRET_KEY used to sign the cookie")
	flags.StringVarP(&signOpts.payload, "payload", "p", "", "JSON payload to sign")
	flags.StringVar(&signOpts.salt, "salt", "cookie-session", "itsdangerous signer salt")
	flags.StringVar(&signOpts.digest, "digest", "sha1", "digest algorithm: sha1, sha224, sha256, sha384, sha512, md5")

	rootCmd.AddCommand(signCmd)
}
