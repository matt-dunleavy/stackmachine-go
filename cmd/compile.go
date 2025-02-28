package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/matt-dunleavy/go-vm/internal/compiler"
	"github.com/matt-dunleavy/go-vm/pkg/utils"
)

// compileCmd represents the compile command
var compileCmd = &cobra.Command{
	Use:   "compile [file...]",
	Short: "Compile stack machine source code to bytecode",
	Long: `Compile stack machine source code to bytecode.
If no files are specified, compilation reads from standard input.
The default output filename is the input filename with '.bin' extension.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			compileStdin()
		} else {
			for _, filename := range args {
				if filename == "-" {
					compileStdin()
				} else {
					compileFile(filename)
				}
			}
		}
	},
	Aliases: []string{"smc"},
}

func init() {
	rootCmd.AddCommand(compileCmd)
}

func compileFile(filename string) {
	file, err := utils.OpenFileForReading(filename)
	if err != nil {
		utils.StandardError("Error opening file %s: %v", filename, err)
	}
	defer file.Close()

	outFilename := utils.GetOutputFilename(filename, ".bin")
	outFile, err := utils.OpenFileForWriting(outFilename)
	if err != nil {
		utils.StandardError("Error creating output file %s: %v", outFilename, err)
	}
	defer outFile.Close()

	compileFn := func(msg string) {
		utils.StandardError("%s:%s", filename, msg)
	}

	c := compiler.NewCompiler(compileFn)
	if err := c.CompileSource(file); err != nil {
		utils.StandardError("Error compiling %s: %v", filename, err)
	}

	if err := c.GetProgram().SaveImage(outFile); err != nil {
		utils.StandardError("Error saving compiled program: %v", err)
	}

	fmt.Printf("Compiled %s to %s\n", filename, outFilename)
}

func compileStdin() {
	outFilename := "out.bin"
	outFile, err := utils.OpenFileForWriting(outFilename)
	if err != nil {
		utils.StandardError("Error creating output file %s: %v", outFilename, err)
	}
	defer outFile.Close()

	compileFn := func(msg string) {
		utils.StandardError("<stdin>:%s", msg)
	}

	c := compiler.NewCompiler(compileFn)
	if err := c.CompileSource(os.Stdin); err != nil {
		utils.StandardError("Error compiling from stdin: %v", err)
	}

	if err := c.GetProgram().SaveImage(outFile); err != nil {
		utils.StandardError("Error saving compiled program: %v", err)
	}

	fmt.Printf("Compiled from stdin to %s\n", outFilename)
}
