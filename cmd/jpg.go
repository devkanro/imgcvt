package cmd

import (
	"github.com/spf13/cobra"
	"image"
	"image/jpeg"
	"io"
)

var jpgOption = &jpeg.Options{100}

func init() {
	rootCmd.AddCommand(jpgCmd)
	jpgCmd.Flags().StringVarP(&output, "output", "o", "",
		`output file/folder.`)
	jpgCmd.Flags().IntVarP(&jpgOption.Quality, "quality", "q", 100,
		`quality of jpeg file.`)
}

var jpgCmd = &cobra.Command{
	Use:   "jpg",
	Short: "Convert images to jpg format",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleArgs(args, func(writer io.Writer, image image.Image) error {
			return jpeg.Encode(writer, image, jpgOption)
		})
	},
}
