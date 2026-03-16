package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUseCommand(t *testing.T) {
	// 1. Setup temporary PVM_HOME
	tempDir, err := os.MkdirTemp("", "pvm_use_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// 2. Create the folder structure
	versionsDir := filepath.Join(tempDir, "versions")
	binDir := filepath.Join(tempDir, "bin")
	os.MkdirAll(versionsDir, 0755)
	os.MkdirAll(binDir, 0755)

	// 3. Create dummy PHP versions
	php81Path := filepath.Join(versionsDir, "php-8.1.0-Win32-vs16-x64")
	php82Path := filepath.Join(versionsDir, "php-8.2.0-nts-Win32-vs16-x64")
	os.MkdirAll(php81Path, 0755)
	os.MkdirAll(php82Path, 0755)
	
	// Create dummy php.exe for the logic to find
	os.WriteFile(filepath.Join(php81Path, "php.exe"), []byte(""), 0755)
	os.WriteFile(filepath.Join(php82Path, "php.exe"), []byte(""), 0755)

	// 4. Mock the Environment Variable
	oldEnv := os.Getenv("PVM_HOME")
	os.Setenv("PVM_HOME", tempDir)
	defer os.Setenv("PVM_HOME", oldEnv)

	// 5. Run the "Use" command for 8.2 NTS
	// Arguments: [version, type]
	Use([]string{"8.2", "nts"})

	// 6. Assertions
	
	// Check if php.bat was created in the bin folder
	batFile := filepath.Join(binDir, "php.bat")
	assert.FileExists(t, batFile)

	// Verify the content of the bat file points to the NTS version
	content, _ := os.ReadFile(batFile)
	assert.Contains(t, string(content), "php-8.2.0-nts")

	// Verify the "current" junction exists
	// Note: On Windows, os.Stat follows symlinks/junctions
	currentJunction := filepath.Join(versionsDir, "current")
	info, err := os.Lstat(currentJunction)
	assert.Nil(t, err, "The 'current' junction should exist")
	
	// Verify it is a directory (Junctions appear as directories with the ReparsePoint bit)
	assert.True(t, info.IsDir())
}

func TestUseCommand_VersionNotFound(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "pvm_fail_test")
	defer os.RemoveAll(tempDir)
	
	os.MkdirAll(filepath.Join(tempDir, "versions"), 0755)
	
	os.Setenv("PVM_HOME", tempDir)
	
	// Try to use a version that doesn't exist
	// We check that it doesn't crash
	Use([]string{"9.9"})
	
	// Verify no wrapper was created
	assert.NoFileExists(t, filepath.Join(tempDir, "bin", "php.bat"))
}

// Keep your mock if you need it for other internal logic tests
type fakeDirEntry struct {
	name  string
	isDir bool
}

func (f fakeDirEntry) Name() string               { return f.name }
func (f fakeDirEntry) IsDir() bool                { return f.isDir }
func (f fakeDirEntry) Type() os.FileMode          { return os.ModeDir }
func (f fakeDirEntry) Info() (os.FileInfo, error) { return nil, nil }