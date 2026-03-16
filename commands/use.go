package commands

import (
	"fmt"
	"ashoknepalofficial/pvm/common"
	"ashoknepalofficial/pvm/theme"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Use(args []string) {
	threadSafe := true

	if len(args) < 1 {
		theme.Error("You must specify a version to use (e.g., pvm use 8.2).")
		return
	}

	if len(args) > 1 {
		if args[1] == "nts" {
			threadSafe = false
		}
	}

	// 1. GET THE INSTALLATION ROOT FROM ENVIRONMENT VARIABLE
	pvmRoot := os.Getenv("PVM_HOME")
	if pvmRoot == "" {
		theme.Error("PVM_HOME environment variable is not set. Please restart your terminal or reinstall PVM.")
		return
	}

	versionsPath := filepath.Join(pvmRoot, "versions")
	binPath := filepath.Join(pvmRoot, "bin")

	// 2. CHECK IF DIRECTORIES EXIST
	if _, err := os.Stat(versionsPath); os.IsNotExist(err) {
		theme.Error("No PHP versions installed in " + versionsPath)
		return
	}

	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		os.MkdirAll(binPath, 0755)
	}

	// 3. FIND THE INSTALLED VERSION
	installedFolders, err := os.ReadDir(versionsPath)
	if err != nil {
		log.Fatalln(err)
	}

	var selectedVersion *versionMeta
	for _, folder := range installedFolders {
		if !folder.IsDir() || folder.Name() == "current" {
			continue
		}

		safe := true
		if strings.Contains(strings.ToLower(folder.Name()), "nts") {
			safe = false
		}

		foundVersion := common.ComputeVersion(folder.Name(), safe, "")
		
		// Match the version string (e.g., "8.2" matches "8.2.10")
		if threadSafe == foundVersion.ThreadSafe && strings.HasPrefix(foundVersion.Semantic(), args[0]) {
			selectedVersion = &versionMeta{
				number: foundVersion,
				folder: folder,
			}
		}
	}

	if selectedVersion == nil {
		theme.Error("The specified version is not installed locally.")
		return
	}

	// 4. CLEAN OLD WRAPPER SCRIPTS IN BIN
	filesToClear := []string{"php.bat", "php", "php-cgi.bat", "php-cgi", "composer.bat", "composer"}
	for _, f := range filesToClear {
		path := filepath.Join(binPath, f)
		if _, err := os.Stat(path); err == nil {
			os.Remove(path)
		}
	}

	versionFolderPath := filepath.Join(versionsPath, selectedVersion.folder.Name())
	phpExePath := filepath.Join(versionFolderPath, "php.exe")
	phpCgiExePath := filepath.Join(versionFolderPath, "php-cgi.exe")
	composerPharPath := filepath.Join(versionFolderPath, "composer", "composer.phar")

	// 5. CREATE WRAPPER SCRIPTS (BAT for CMD, SH for Git Bash/WSL)
	
	// PHP Wrapper
	os.WriteFile(filepath.Join(binPath, "php.bat"), []byte(fmt.Sprintf("@echo off\n\"%s\" %%*\n", phpExePath)), 0755)
	os.WriteFile(filepath.Join(binPath, "php"), []byte(fmt.Sprintf("#!/bin/bash\n\"%s\" \"$@\"\n", phpExePath)), 0755)

	// PHP-CGI Wrapper
	os.WriteFile(filepath.Join(binPath, "php-cgi.bat"), []byte(fmt.Sprintf("@echo off\n\"%s\" %%*\n", phpCgiExePath)), 0755)
	os.WriteFile(filepath.Join(binPath, "php-cgi"), []byte(fmt.Sprintf("#!/bin/bash\n\"%s\" \"$@\"\n", phpCgiExePath)), 0755)

	// Composer Wrapper
	os.WriteFile(filepath.Join(binPath, "composer.bat"), []byte(fmt.Sprintf("@echo off\n\"%s\" \"%s\" %%*\n", phpExePath, composerPharPath)), 0755)
	os.WriteFile(filepath.Join(binPath, "composer"), []byte(fmt.Sprintf("#!/bin/bash\n\"%s\" \"%s\" \"$@\"\n", phpExePath, composerPharPath)), 0755)

	// 6. UPDATE THE "CURRENT" JUNCTION (The most important part)
	// This updates PVM_HOME\versions\current to point to the new version folder
	// Because the Installer linked C:\php to PVM_HOME\versions\current, updating this
	// junction changes the PHP version globally.
	
	currentLinkPath := filepath.Join(versionsPath, "current")

	// Remove old junction/symlink
	// Using 'rmdir' for junctions or 'os.Remove' for symlinks
	os.Remove(currentLinkPath) 

	// Create new Junction (works without special developer mode on Windows)
	cmd := exec.Command("cmd", "/C", "mklink", "/J", currentLinkPath, versionFolderPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		theme.Warning("Could not update the global symlink. You may need to run this terminal as Administrator.")
		fmt.Println(string(output))
	}

	// 7. EXT DIRECTORY LINK
	extensionLinkPath := filepath.Join(binPath, "ext")
	os.Remove(extensionLinkPath) // Clear old
	exec.Command("cmd", "/C", "mklink", "/J", extensionLinkPath, filepath.Join(versionFolderPath, "ext")).Run()

	theme.Success(fmt.Sprintf("Successfully switched to PHP %s", selectedVersion.number.String()))
	theme.Info("Run 'php -v' to verify.")
}

type versionMeta struct {
	number common.Version
	folder os.DirEntry
}