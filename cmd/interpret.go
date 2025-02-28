package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/matt-dunleavy/go-vm/internal/compiler"
	"github.com/matt-dunleavy/go-vm/pkg/utils"
)

// interpretCmd represents the interpret command
var interpretCmd = &cobra.Command{
	Use:   "interpret [file...]",
	Short: "Compile and run stack machine source code on-the-fly",
	Long: `Compile and run stack machine source code on-the-fly.
If no files are specified, input is read from standard input.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			interpretStdin()
		} else {
			for _, filename := range args {
				if filename == "-" {
					interpretStdin()
				} else {
					interpretFile(filename)
				}
			}
		}
	},
	Aliases: []string{"sm"},
}

func init() {
	rootCmd.AddCommand(interpretCmd)
}

func interpretFile(filename string) {
	file, err := utils.OpenFileForReading(filename)
	if err != nil {
		utils.StandardError("Error opening file %s: %v", filename, err)
	}
	defer file.Close()

	errorFn := func(msg string) {
		utils.StandardError("%s:%s", filename, msg)
	}

	c := compiler.NewCompiler(errorFn)
	if err := c.CompileSource(file); err != nil {
		utils.StandardError("Error compiling %s: %v", filename, err)
	}

	c.GetProgram().Run(0)
}

func interpretStdin() {
	errorFn := func(msg string) {
		utils.StandardError("<stdin>:%s", msg)
	}

	c := compiler.NewCompiler(errorFn)
	if err := c.CompileSource(os.Stdin); err != nil {
		utils.StandardError("Error compiling from stdin: %v", err)
	}

	c.GetProgram().Run(0)
}
