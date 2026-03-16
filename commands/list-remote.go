package commands

import (
	"ashoknepalofficial/pvm/common"
	"ashoknepalofficial/pvm/theme"
	"log"
	"os"
	"slices"

	"github.com/fatih/color"
)

func ListRemote() {
	// 1. Verify Environment Variable
	pvmRoot := os.Getenv("PVM_HOME")
	if pvmRoot == "" {
		theme.Error("PVM_HOME environment variable is not set. Please restart your terminal or reinstall PVM.")
		return
	}

	// 2. Fetch remote versions from windows.php.net
	versions, err := common.RetrievePHPVersions()
	if err != nil {
		log.Fatalln(err)
	}

	common.SortVersions(versions)

	// 3. Fetch installed versions (this now uses PVM_HOME internally)
	installedVersions, _ := common.RetrieveInstalledPHPVersions()

	theme.Title("PHP versions available (Remote)")
	
	for _, version := range versions {
		// Check if this remote version exists in our installedVersions slice
		idx := slices.IndexFunc(installedVersions, func(v common.Version) bool { 
			return v.Same(version) 
		})
		
		if idx != -1 {
			// Version is installed
			color.Green("*   " + version.StringShort() + " (installed)")
		} else {
			// Version is only available remotely
			color.White("    " + version.StringShort())
		}
	}
	
	theme.Info("\nUse 'pvm install <version>' to download a version.")
}