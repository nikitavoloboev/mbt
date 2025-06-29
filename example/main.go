// Package main is a fang example.
package main

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

func main() {
	var foo string
	var bar int
	var baz float64
	var d time.Duration
	var eerr bool

	cmd := &cobra.Command{
		Use:   "example [args]",
		Short: "An example program!",
		Long: `A little example program!

It doesn’t really do anything, but that’s the point.™`,
		Example: `
# Run it:
example

# Run it with some arguments:
FOO=bar ZAZ="quoted value" example --name=Carlos -a -s Becker -a

# Run a subcommand with an argument:
example sub --async --foo=xyz --async arguments

# Run with a quoted string:
example sub "quoted string"

# Mix and match:
example sub "multi-word quoted string" --name "another quoted string" -a

# Multi-line:
ENV_A=0 ENV_B=0 ENV_C=0 \
  CERT_FILE=/path/to/chain.pem KEY_FILE=/path/to/key.pem \
  example sub "quoted argument"

# Run a subcommand's subcommand with an argument:
example sub another args --flag
		`,

		RunE: func(c *cobra.Command, _ []string) error {
			if eerr {
				return errors.New("we have an error")
			}
			c.Println("You ran the root command. Now try --help.")
			return nil
		},
	}
	cmd.PersistentFlags().StringVarP(&foo, "surname", "s", "doe", "Your surname")
	cmd.Flags().StringVar(&foo, "name", "jane", "Your name")
	cmd.Flags().DurationVar(&d, "duration", 0, "Time since your last commit")
	cmd.Flags().IntVar(&bar, "age", 0, "Your age")
	cmd.Flags().Float64Var(&baz, "idk", 0.0, "I don't know")
	cmd.Flags().BoolP("async", "a", false, "Run async")
	cmd.Flags().BoolVarP(&eerr, "error", "e", false, "Makes the program exit with error")

	_ = cmd.Flags().MarkHidden("age")
	_ = cmd.Flags().MarkHidden("duration")
	_ = cmd.Flags().MarkHidden("idk")
	_ = cmd.Flags().MarkHidden("error")

	cmd.AddGroup(&cobra.Group{
		ID:    "group1",
		Title: "My Group",
	})
	sub := &cobra.Command{
		Use:     "sub [command] [flags] [args]",
		Short:   "An example subcommand",
		GroupID: "group1",
		Example: `example sub some arguments --and-flags
example sub another --thing`,
		Run: func(c *cobra.Command, _ []string) {
			c.Println("Ran the sub command!")
		},
	}
	cmd.AddCommand(sub)
	sub.AddCommand(&cobra.Command{
		Use:     "another",
		Short:   "another sub command",
		Example: `example sub another --foo=bar`,
		RunE: func(c *cobra.Command, _ []string) error {
			cmd.Println("Working...")
			select {
			case <-time.After(time.Second * 5):
				cmd.Println("Done!")
			case <-c.Context().Done():
				return c.Context().Err()
			}
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:     "throw",
		Short:   "Throws an error",
		GroupID: "group1",
		RunE: func(*cobra.Command, []string) error {
			return errors.New("a super long error string that is meant to test the error handling in fang. It should be long enough to wrap around and test the error styling and formatting capabilities of fang. This is a test to see how well fang handles long error messages and whether it can display them properly without breaking the layout or causing any issues")
		},
	})

	// This is where the magic happens.
	if err := fang.Execute(
		context.Background(),
		cmd,
		fang.WithNotifySignal(os.Interrupt, os.Kill),
	); err != nil {
		os.Exit(1)
	}
}
