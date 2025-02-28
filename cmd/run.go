package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/matt-dunleavy/stackmachine-go/internal/vm"
	"github.com/matt-dunleavy/stackmachine-go/pkg/utils"
)

var showHelp bool

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [file...]",
	Short: "Run compiled stack machine bytecode",
	Long: `Run compiled stack machine bytecode.
If no files are specified, execution reads from standard input.`,
	Run: func(cmd *cobra.Command, args []string) {
		if showHelp {
			printInstructions()
			return
		}

		foundFile := false
		for _, filename := range args {
			if filename == "-" {
				runStdin()
			} else {
				runFile(filename)
				foundFile = true
			}
		}

		if len(args) == 0 || !foundFile {
			runStdin()
		}
	},
	Aliases: []string{"smr"},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVarP(&showHelp, "help", "h", false, "Show instruction set")
}

func printInstructions() {
	fmt.Println("stack-machine-go -- stack-machine run")
	fmt.Println("Public domain, 2010-2011 by Christian Stigen Larsen")
	fmt.Println("Go port by Matt Dunleavy")
	fmt.Printf("\nOpcodes:\n")

	for i := vm.Op(0); i < vm.NOP_END; i++ {
		fmt.Printf("0x%x = %s\n", i, i)
	}

	fmt.Printf("\nTo halt program, jump to current position:\n")
	fmt.Printf("0x0 PUSH 0x%x\n", 4)
	fmt.Printf("0x%x JMP\n\n", 4)
	fmt.Printf("Word size is %d bytes\n", 4)
}

func runFile(filename string) {
	file, err := utils.OpenFileForReading(filename)
	if err != nil {
		utils.StandardError("Error opening file %s: %v", filename, err)
	}
	defer file.Close()

	errorFn := func(msg string) {
		fmt.Fprintf(os.Stderr, "%s\n", msg)
	}

	m := vm.NewMachine(errorFn)
	if err := m.LoadImage(file); err != nil {
		utils.StandardError("Error loading program from %s: %v", filename, err)
	}

	m.Run(0)
}

func runStdin() {
	errorFn := func(msg string) {
		fmt.Fprintf(os.Stderr, "%s\n", msg)
	}

	m := vm.NewMachine(errorFn)
	if err := m.LoadImage(os.Stdin); err != nil {
		utils.StandardError("Error loading program from stdin: %v", err)
	}

	m.Run(0)
}
