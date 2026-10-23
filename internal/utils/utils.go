package utils

import (
	"fmt"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

func CreateProgressBar(length int, which string, outOf string, description string) *progressbar.ProgressBar {
	descriptionString := fmt.Sprintf("[cyan][%s/%s][reset] %s...", which, outOf, description)
	bar := progressbar.NewOptions(length,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()), //you should install "github.com/k0kubun/go-ansi"
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(100),
		progressbar.OptionSetDescription(descriptionString),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	return bar
}

func CreateTabProgressBar(length int, which string, outOf string, description string) *progressbar.ProgressBar {
	descriptionString := fmt.Sprintf("	[cyan][%s/%s][reset] %s...", which, outOf, description)
	bar := progressbar.NewOptions(length,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()), //you should install "github.com/k0kubun/go-ansi"
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(100),
		progressbar.OptionSetDescription(descriptionString),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "	[green]=[reset]",
			SaucerHead:    "	[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	return bar
}
