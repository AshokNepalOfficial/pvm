package commands

import (
	"ashoknepalofficial/pvm/common"
	"ashoknepalofficial/pvm/theme"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func List() {
	// 1. Get installed versions (Helper already uses PVM_HOME)
	versions, err := common.RetrieveInstalledPHPVersions()
	if err != nil {
		theme.Error(err.Error())
		return
	}

	if len(versions) == 0 {
		theme.Info("No PHP versions installed. Use 'pvm install <version>' to get started.")
		return
	}

	// 2. Identify the active version by reading the 'current' junction
	pvmRoot := os.Getenv("PVM_HOME")
	activeFolderName := ""
	if pvmRoot != "" {
		currentLink := filepath.Join(pvmRoot, "versions", "current")
		// EvalSymlinks returns the actual path a symlink/junction points to
		if target, err := filepath.EvalSymlinks(currentLink); err == nil {
			activeFolderName = filepath.Base(target)
		}
	}

	theme.Title("Installed PHP versions")
	
	for _, version := range versions {
		// Logic to check if this version folder matches the folder 'current' points to
		// We use StringShort or folder name logic
		isCurrent := false
		
		// This assumes your ComputeVersion creates a semantic string that matches 
		// the folder naming convention used in 'use.go'
		if activeFolderName != "" && strings.Contains(activeFolderName, version.Semantic()) {
			// Double check Thread Safety match if folder name contains 'nts'
			isNtsFolder := strings.Contains(strings.ToLower(activeFolderName), "nts")
			if isNtsFolder == !version.ThreadSafe {
				isCurrent = true
			}
		}

		if isCurrent {
			color.Green("  * " + version.StringShort() + " (active)")
		} else {
			color.White("    " + version.StringShort())
		}
	}
}