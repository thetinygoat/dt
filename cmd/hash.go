package cmd

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"strings"

	"dt/internal/cliio"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/sha3"
)

var hashCmd = &cobra.Command{
	Use:   "hash",
	Short: "Generate digests with common hashing algorithms",
}

func init() {
	hashCmd.AddCommand(newHashCommand("md5", "MD5 digest", func() hash.Hash { return md5.New() }))
	hashCmd.AddCommand(newHashCommand("sha1", "SHA-1 digest", func() hash.Hash { return sha1.New() }))
	hashCmd.AddCommand(newHashCommand("sha256", "SHA-256 digest", func() hash.Hash { return sha256.New() }))
	hashCmd.AddCommand(newHashCommand("sha512", "SHA-512 digest", func() hash.Hash { return sha512.New() }))
	hashCmd.AddCommand(newHashCommand("sha3-256", "SHA3-256 digest", func() hash.Hash { return sha3.New256() }))
	hashCmd.AddCommand(newHashCommand("sha3-512", "SHA3-512 digest", func() hash.Hash { return sha3.New512() }))
	rootCmd.AddCommand(hashCmd)
}

func newHashCommand(name, short string, factory func() hash.Hash) *cobra.Command {
	var encoding string
	var salt string
	cmd := &cobra.Command{
		Use:   name,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := cliio.ReadAll(args)
			if err != nil {
				return err
			}
			if salt != "" {
				data = append(data, []byte(salt)...)
			}
			h := factory()
			if _, err := h.Write(data); err != nil {
				return err
			}
			sum := h.Sum(nil)
			switch strings.ToLower(encoding) {
			case "hex":
				fmt.Println(hex.EncodeToString(sum))
			case "base64":
				fmt.Println(base64.StdEncoding.EncodeToString(sum))
			default:
				return fmt.Errorf("unsupported encoding %q (use hex or base64)", encoding)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&encoding, "encoding", "hex", "output encoding (hex|base64)")
	cmd.Flags().StringVar(&salt, "salt", "", "append salt to the input before hashing")
	return cmd
}
