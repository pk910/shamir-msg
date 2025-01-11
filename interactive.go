package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v3"
)

func runInteractive(ctx context.Context, cmd *cli.Command) error {
	printHeader()

	prompt := promptui.Select{
		Label: "What do you want to do?",
		Items: []string{
			"Combine shards into the original secret",
			"Split a secret into shards",
		},
		HideSelected: true,
	}

	mode, modeStr, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("mode selection failed %v", err)
	}

	//time.Sleep(50 * time.Millisecond)
	fmt.Printf("%v %v\n", color.HiBlackString("What do you want to do?"), color.HiWhiteString(modeStr))

	switch mode {
	case 0:
		// Combine shards into the original secret

		shards := make([]string, 0)
		for {
			prompt := promptui.Prompt{
				Label:       "Enter a shard",
				HideEntered: true,
			}

			shard, err := prompt.Run()
			if err != nil {
				return fmt.Errorf("shard prompt failed %v", err)
			}

			if shard == "" {
				break
			}

			time.Sleep(50 * time.Millisecond)
			fmt.Printf("%v %v\n", color.HiBlackString(fmt.Sprintf("Shard %d:", len(shards)+1)), color.HiWhiteString(shard))

			shards = append(shards, shard)
		}

		fmt.Printf("\n")

		secret, err := ShamirCombine(shards)
		if err != nil {
			return fmt.Errorf("ShamirCombine failed %v", err)
		}

		fmt.Printf("%v %v\n", color.HiBlackString("Secret:"), color.HiWhiteString(secret))
		fmt.Printf("Press enter to exit\n")
		fmt.Scanln()
		fmt.Printf("\033[1A") // one line up
		fmt.Printf("\033[1A") // one line up
		fmt.Printf("\033[1A") // one line up
		fmt.Printf("\033[K")  // clear line
		fmt.Printf("\033[1B") // one line down
		fmt.Printf("\033[1B") // one line down

	case 1:
		// Split a secret into shards

		prompt := promptui.Prompt{
			Label:       "Enter the secret to split",
			Mask:        '*',
			HideEntered: true,
		}

		secret, err := prompt.Run()
		if err != nil {
			return fmt.Errorf("secret prompt failed %v", err)
		}

		fmt.Printf("%v %v\n", color.HiBlackString("Enter the secret to split:"), color.HiWhiteString(strings.Repeat("*", len(secret))))

		validate := func(input string) error {
			_, err := strconv.ParseInt(input, 10, 64)
			if err != nil {
				return errors.New("invalid number")
			}
			return nil
		}

		prompt = promptui.Prompt{
			Label:       "Total number of shards",
			Validate:    validate,
			HideEntered: true,
		}

		keysStr, err := prompt.Run()
		if err != nil {
			return fmt.Errorf("keys prompt failed %v", err)
		}

		keys, _ := strconv.ParseInt(keysStr, 10, 64)

		fmt.Printf("%v %v\n", color.HiBlackString("Total number of shards:"), color.HiWhiteString(fmt.Sprintf("%d", keys)))

		prompt = promptui.Prompt{
			Label:       "Minimum number of shards for reconstruction",
			Validate:    validate,
			HideEntered: true,
		}

		thresholdStr, err := prompt.Run()
		if err != nil {
			return fmt.Errorf("threshold prompt failed %v", err)
		}

		threshold, _ := strconv.ParseInt(thresholdStr, 10, 64)

		fmt.Printf("%v %v\n", color.HiBlackString("Minimum number of shards for reconstruction:"), color.HiWhiteString(fmt.Sprintf("%d", keys)))

		fmt.Printf("\n")

		groupSize := cmd.Int(groupSizeFlag.Name)
		shares, err := ShamirSplit(int(keys), int(threshold), secret, int(groupSize))
		if err != nil {
			return fmt.Errorf("ShamirSplit failed %v", err)
		}

		for i, shard := range shares {
			fmt.Printf("%v %v\n", color.HiBlackString(fmt.Sprintf("Shard %d:", i+1)), color.HiWhiteString(shard))
		}

		fmt.Printf("\n")
	}

	return nil
}
