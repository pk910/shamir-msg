package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
)

var (
	keysFlag = &cli.IntFlag{
		Name:    "keys",
		Aliases: []string{"k"},
		Usage:   "The number of total key shards to generate",
		Value:   3,
	}
	thresholdFlag = &cli.IntFlag{
		Name:    "threshold",
		Aliases: []string{"t"},
		Usage:   "The min number of shards needed to re-assemble the secret",
		Value:   2,
	}
	groupSizeFlag = &cli.IntFlag{
		Name:    "group-size",
		Aliases: []string{"g"},
		Usage:   "Split output into groups of this size and separate with a space",
		Value:   6,
	}
	secretFlag = &cli.StringFlag{
		Name:     "secret",
		Aliases:  []string{"s"},
		Usage:    "The secret to split into shards",
		Required: true,
	}
	shardsFlag = &cli.StringSliceFlag{
		Name:     "shard",
		Aliases:  []string{"s"},
		Usage:    "The shards to recombine into the secret",
		Required: true,
	}
	quietFlag = &cli.BoolFlag{
		Name:    "quiet",
		Aliases: []string{"q"},
		Usage:   "Suppress everything except the secret output",
	}

	app = &cli.Command{
		Name:  "shamir-msg",
		Usage: "Split and recombine a secret using Shamir's Secret Sharing",
		Commands: []*cli.Command{
			{
				Name:  "split",
				Usage: "Split a secret into shards",
				Flags: []cli.Flag{keysFlag, thresholdFlag, secretFlag, groupSizeFlag, quietFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return runSplit(ctx, cmd)
				},
			},
			{
				Name:  "combine",
				Usage: "Combine shards into the original secret",
				Flags: []cli.Flag{shardsFlag, quietFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return runCombine(ctx, cmd)
				},
			},
			{
				Name:  "run",
				Usage: "Run the tool in interactive mode (default)",
				Flags: []cli.Flag{groupSizeFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return runInteractive(ctx, cmd)
				},
			},
		},
		DefaultCommand: "run",
	}
)

func main() {
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func printHeader() {
	fmt.Printf("##### Shamir Secret Sharing Tool #####\n")
	fmt.Printf("Split and recombine a secret using Shamir's Secret Sharing\n\n")
}

func runSplit(ctx context.Context, cmd *cli.Command) error {
	quiet := cmd.Bool(quietFlag.Name)
	if !quiet {
		printHeader()
	}

	keys := cmd.Int(keysFlag.Name)
	threshold := cmd.Int(thresholdFlag.Name)
	groupSize := cmd.Int(groupSizeFlag.Name)

	shares, err := ShamirSplit(int(keys), int(threshold), cmd.String(secretFlag.Name), int(groupSize))
	if err != nil {
		return err
	}

	if !quiet {
		fmt.Printf("%v %v\n", color.HiBlackString("Total keys shards: "), color.HiWhiteString(fmt.Sprintf("%d", keys)))
		fmt.Printf("%v %v\n", color.HiBlackString("Required for reconstruction: "), color.HiWhiteString(fmt.Sprintf("%d", threshold)))
		fmt.Printf("\n")
	}

	for i, shard := range shares {
		if !quiet {
			fmt.Print(color.HiBlackString(fmt.Sprintf("Shard %d: ", i+1)))
			color.HiWhite(shard)
		} else {
			fmt.Printf("%s\n", shard)
		}
	}

	if !quiet {
		fmt.Printf("\n")
	}

	return nil
}

func runCombine(ctx context.Context, cmd *cli.Command) error {
	quiet := cmd.Bool(quietFlag.Name)
	if !quiet {
		printHeader()
	}

	secret, err := ShamirCombine(cmd.StringSlice(shardsFlag.Name))
	if err != nil {
		return err
	}

	if !quiet {
		fmt.Print(color.HiBlackString("Secret: "))
		color.HiWhite(secret)
		fmt.Printf("\n")
	} else {
		fmt.Printf("%s\n", secret)
	}

	return nil
}
