package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "smg",
	Short: "A stack-based virtual machine",
	Long: `Stack Machine is a simple stack-based virtual machine that can run 
Forth/PostScript-like programs. It includes an assembler, interpreter, and 
execution environment.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Add common flags here if needed
}
