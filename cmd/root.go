/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

var (
	source_path      string
	destination_path string
)

var rootCmd = &cobra.Command{
	Use:   "files-merge",
	Short: "Create new folders structure without duplicated files",
	Long:  `Create new folders structure without duplicated files`,
	Run:   run,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	rootCmd.Flags().StringVarP(&source_path, "source", "s", pwd, "source folder")
	rootCmd.Flags().StringVarP(&destination_path, "destination", "d", "destination folder", "destination folder")
	if err := rootCmd.MarkFlagRequired("destination"); err != nil {
		panic(err)
	}
}

func run(cmd *cobra.Command, args []string) {

	fmt.Println(source_path)
	fmt.Println(destination_path)

	dirs, err := os.ReadDir(source_path)
	if err != nil {
		panic(err)
	}

	for _, dir := range dirs {
		duplicates := getDuplicates(dir, dirs)

		fmt.Println(dir.Name())
		for _, duplicate := range duplicates {
			fmt.Printf("%s - duplicate\n", duplicate.Name())
		}
	}
}

func getDuplicates(dir os.DirEntry, dirs []os.DirEntry) []os.DirEntry {
	duplicatePattern := regexp.MustCompile(`^` + regexp.QuoteMeta(dir.Name()) + `\(\d+\)$`)

	var duplicates []os.DirEntry
	for _, candidate := range dirs {
		if candidate.Name() == dir.Name() {
			continue
		}
		if duplicatePattern.MatchString(candidate.Name()) {
			duplicates = append(duplicates, candidate)
		}
	}

	return duplicates
}
