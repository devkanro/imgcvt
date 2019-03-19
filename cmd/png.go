package cmd

import (
	"github.com/spf13/cobra"
	"image"
	"image/png"
	"io"
)

func init() {
	rootCmd.AddCommand(pngCmd)
	pngCmd.Flags().StringVarP(&output, "output", "o", "",
		`output file/folder.`)
}

var pngCmd = &cobra.Command{
	Use:   "png",
	Short: "Convert images to png format",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleArgs(args, func(writer io.Writer, image image.Image) error {
			return png.Encode(writer, image)
		})
	},
}
