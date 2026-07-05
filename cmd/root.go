package cmd

import (
	"fmt"
	"os"

	"github.com/Marcin99b/files-merge/internal/merge"

	"github.com/spf13/cobra"
)

var (
	sourcePath      string
	destinationPath string
)

var rootCmd = &cobra.Command{
	Use:   "files-merge",
	Short: "Create new folders structure without duplicated files",
	Long:  `Create new folders structure without duplicated files`,
	RunE:  run,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	rootCmd.Flags().StringVarP(&sourcePath, "source", "s", pwd, "source folder")
	rootCmd.Flags().StringVarP(&destinationPath, "destination", "d", "", "destination folder")
	if err := rootCmd.MarkFlagRequired("destination"); err != nil {
		panic(err)
	}
}

func run(cmd *cobra.Command, args []string) error {
	results, err := merge.Directories(sourcePath, destinationPath)
	if err != nil {
		return err
	}

	failureCount := 0

	for _, result := range results {
		fmt.Println(result.FolderName)
		for _, duplicateFolderName := range result.DuplicateFolderNames {
			fmt.Printf("%s - duplicate\n", duplicateFolderName)
		}
		for _, copiedFilePath := range result.CopiedFilePaths {
			fmt.Println(copiedFilePath)
		}
		for _, failure := range result.Failures {
			fmt.Printf("SKIPPED: %s (%v)\n", failure.Path, failure.Err)
		}
		failureCount += len(result.Failures)
	}

	if failureCount > 0 {
		fmt.Printf("\n%d file(s) could not be copied and were skipped, see SKIPPED entries above.\n", failureCount)
	}

	return nil
}
