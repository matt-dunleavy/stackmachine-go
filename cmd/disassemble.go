package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/matt-dunleavy/stackmachine-go/internal/vm"
	"github.com/matt-dunleavy/stackmachine-go/pkg/utils"
)

// disassembleCmd represents the disassemble command
var disassembleCmd = &cobra.Command{
	Use:   "disassemble [file...]",
	Short: "Disassemble compiled bytecode",
	Long: `Disassemble compiled bytecode to human-readable form.
If no files are specified, input is read from standard input.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			disassembleStdin()
		} else {
			for _, filename := range args {
				if filename == "-" {
					disassembleStdin()
				} else {
					disassembleFile(filename)
				}
			}
		}
	},
	Aliases: []string{"smd"},
}

func init() {
	rootCmd.AddCommand(disassembleCmd)
}

func isPrintable(c int) bool {
	return (c >= 32 && c <= 127) || c == '\n' || c == '\r' || c == '\t'
}

func toString(c byte) string {
	switch c {
	case '\t':
		return "\\t"
	case '\n':
		return "\\n"
	case '\r':
		return "\\r"
	default:
		return string(c)
	}
}

func disassemble(m *vm.VM) {
	end := m.Size()

	for m.Pos() <= end {
		op := vm.Op(m.Cur())
		fmt.Printf("0x%x %s", m.Pos(), op)

		if (op == vm.PUSH || op == vm.PUSHIP) && m.Pos() <= end {
			m.Next()
			val := m.Cur()
			fmt.Printf(" 0x%x", val)

			if isPrintable(int(val)) {
				fmt.Printf(" ('%s')", toString(byte(val)))
			}
		}

		fmt.Println()
		m.Next()
	}
}

func disassembleFile(filename string) {
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

	fmt.Printf("; File %s --- %d bytes\n", filename, m.Size())
	disassemble(m)
}

func disassembleStdin() {
	errorFn := func(msg string) {
		fmt.Fprintf(os.Stderr, "%s\n", msg)
	}

	m := vm.NewMachine(errorFn)
	if err := m.LoadImage(os.Stdin); err != nil {
		utils.StandardError("Error loading program from stdin: %v", err)
	}

	fmt.Printf("; From stdin --- %d bytes\n", m.Size())
	disassemble(m)
}
