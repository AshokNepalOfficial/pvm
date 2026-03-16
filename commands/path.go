package commands

import (
	"fmt"
	"ashoknepalofficial/pvm/theme"
	"os"
	"path/filepath"
)

func Path() {
	theme.Title("pvm: PHP Version Manager")

	// 1. Get the installation root from the environment variable set by the installer
	pvmRoot := os.Getenv("PVM_HOME")

	if pvmRoot == "" {
		theme.Error("PVM_HOME is not set. Please restart your terminal or reinstall PVM.")
		return
	}

	// 2. Point to the bin directory where your wrapper scripts (php.bat, etc.) are generated
	binPath := filepath.Join(pvmRoot, "bin")

	fmt.Println("PVM is installed at: " + pvmRoot)
	fmt.Println("\nTo use PHP and Composer from the command line, ensure this directory is in your PATH:")
	fmt.Println("    " + binPath)
	
	// Optional: Check if the current Symlink (C:\php) is also relevant
	pvmSymlink := os.Getenv("PVM_SYMLINK")
	if pvmSymlink != "" {
		fmt.Println("\nYour active PHP symlink is located at:")
		fmt.Println("    " + pvmSymlink)
	}
}