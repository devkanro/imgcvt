package cmd

import (
	"github.com/chai2010/webp"
	"github.com/spf13/cobra"
	"image"
	"io"
)

var webpOption = &webp.Options{Quality: 90}

func init() {
	rootCmd.AddCommand(webpCmd)
	webpCmd.Flags().StringVarP(&output, "output", "o", "",
		`output file/folder.`)
	webpCmd.Flags().Float32VarP(&webpOption.Quality, "quality", "q", 90,
		`quality of webp file.`)
	webpCmd.Flags().BoolVarP(&webpOption.Lossless, "lossless", "l", false,
		`create lossless webp file.`)
	webpCmd.Flags().BoolVarP(&webpOption.Exact, "exact", "x", false,
		`preserve RGB values in transparent area.`)
}

var webpCmd = &cobra.Command{
	Use:   "webp",
	Short: "Convert images to webp format",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleArgs(args, func(writer io.Writer, image image.Image) error {
			return webp.Encode(writer, image, webpOption)
		})
	},
}
