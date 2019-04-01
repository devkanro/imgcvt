package cmd

import (
	"github.com/spf13/cobra"
	"golang.org/x/image/bmp"
	"image"
	"io"
)

func init() {
	rootCmd.AddCommand(bmpCmd)
	bmpCmd.Flags().StringVarP(&output, "output", "o", "",
		`output file/folder.`)
}

var bmpCmd = &cobra.Command{
	Use:   "bmp",
	Short: "Convert images to bmp format",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleArgs(args, func(writer io.Writer, image image.Image) error {
			return bmp.Encode(writer, image)
		})
	},
}
