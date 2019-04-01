package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"image"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

var rootCmd = &cobra.Command{
	Use:     "imgcvt",
	Short:   "imgcvt is universal image converter",
	Long:    `A converter written in go for image formats png, jpg etc.`,
	Version: "v1.2",
}

var output string

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func handleArgs(args []string, converter func(io.Writer, image.Image) error) error {
	for _, path := range args {
		stat, err := os.Stat(path)
		if err != nil {
			return err
		}

		if stat.IsDir() {
			err := handleFolder(path, converter)
			if err != nil {
				fmt.Printf("Skip convert folder '%s' due to error: %s\n", path, err)
				continue
			}
		} else {
			var outputFilePath string

			inputFilePath := filepath.Join(path)
			if output == "" {
				outputFilePath = inputFilePath[:len(inputFilePath)-len(filepath.Ext(inputFilePath))] + ".png"
			} else {
				outputFilePath = filepath.Join(output)
				if filepath.Ext(outputFilePath) != ".png" {
					name := filepath.Base(inputFilePath)
					outputFilePath = filepath.Join(outputFilePath, name[:len(name)-len(filepath.Ext(name))]+".png")
				}
			}

			err := handleFile(inputFilePath, outputFilePath, converter)
			if err != nil {
				fmt.Printf("Skip convert file '%s' due to error: %s\n", path, err)
				continue
			}
		}
	}

	return nil
}

func handleFolder(path string, converter func(io.Writer, image.Image) error) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(files))
	sem := make(chan bool, runtime.NumCPU())

	for _, fileInfo := range files {
		sem <- true
		go func(fileInfo os.FileInfo) {
			defer wg.Done()
			defer func() { <-sem }()

			inputFilePath := filepath.Join(path, fileInfo.Name())
			name := fileInfo.Name()

			var outputFilePath string
			if output == "" {
				outputFilePath = inputFilePath[:len(inputFilePath)-len(filepath.Ext(inputFilePath))] + ".png"
			} else {
				outputFilePath = filepath.Join(output, name[:len(name)-len(filepath.Ext(path))]+".png")
			}

			err := handleFile(inputFilePath, outputFilePath, converter)
			if err != nil {
				fmt.Printf("Skip convert file '%s' due to error: %s\n", fileInfo.Name(), err)
				return
			}
		}(fileInfo)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	wg.Wait()

	return nil
}

func handleFile(path string, outputPath string, converter func(io.Writer, image.Image) error) error {
	name := filepath.Base(path)

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	imageData, format, err := image.Decode(file)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = outputFile.Close()
	}()

	err = converter(outputFile, imageData)
	if err != nil {
		return err
	}

	fmt.Printf("Convert %s from %s to png\n", name, format)
	return nil
}
