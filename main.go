package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/mdp/qrterminal"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "nostrconnect",
		Usage: "",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "sec",
				Usage:       "secret key to sign the event",
				DefaultText: "the key '1'",
				Value:       "0000000000000000000000000000000000000000000000000000000000000001",
			},
			&cli.StringFlag{
				Name:     "relay",
				Aliases:  []string{"r"},
				Usage:    "relay to listen on for commands",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			pubkey, err := nostr.GetPublicKey(c.String("sec"))
			if err != nil {
				return fmt.Errorf("given private key is invalid: %w", err)
			}
			relay := nostr.NormalizeURL(c.String("relay"))

			full := (&url.URL{
				Scheme: "nostrconnect",
				Host:   pubkey,
				RawQuery: (url.Values{
					"relay":    {relay},
					"metadata": {fmt.Sprintf(`{"name": "nostrconnect CLI tool"}`)},
				}).Encode(),
			}).String()
			npub, _ := nip19.EncodePublicKey(pubkey)

			fmt.Println("")
			fmt.Println("  \033[1mrelay:\033[0m " + relay)
			fmt.Println("  \033[1mnpub:\033[0m " + npub)
			fmt.Println("")
			qrterminal.GenerateHalfBlock(full, qrterminal.M, os.Stdout)
			fmt.Println(full)

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
