package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Version_Compare(t *testing.T) {
	// v1: 1.2.3 NTS
	v1 := Version{Major: 1, Minor: 2, Patch: 3, ThreadSafe: false}
	// v2: 1.2.4 TS
	v2 := Version{Major: 1, Minor: 2, Patch: 4, ThreadSafe: true}
	// v3: 1.2.3 TS
	v3 := Version{Major: 1, Minor: 2, Patch: 3, ThreadSafe: true}

	// 1.2.3 is less than 1.2.4
	assert.True(t, v1.LessThan(v2))
	
	// 1.2.3 NTS is not the same as 1.2.3 TS
	assert.False(t, v1.Same(v3))

	// Test "nulled" (-1) logic
	v4 := Version{Major: 1, Minor: 2, Patch: -1}
	v5 := Version{Major: 1, Minor: 2, Patch: 3}
	// Per your implementation, if Patch is -1, Compare returns 0 (not less than)
	assert.False(t, v4.LessThan(v5))

	v6 := Version{Major: 1, Minor: -1}
	v7 := Version{Major: 1, Minor: 2}
	assert.False(t, v6.LessThan(v7))
}

func Test_RetrieveInstalledPHPVersions_EnvVar(t *testing.T) {
	// 1. Setup a temporary directory to act as our PVM_HOME
	tempDir, err := os.MkdirTemp("", "pvm_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir) // Clean up after test

	// 2. Create the 'versions' subdirectory and some dummy PHP folders
	versionsDir := filepath.Join(tempDir, "versions")
	os.MkdirAll(versionsDir, 0755)
	
	os.Mkdir(filepath.Join(versionsDir, "php-8.1.0-Win32-vs16-x64"), 0755)
	os.Mkdir(filepath.Join(versionsDir, "php-8.2.0-nts-Win32-vs16-x64"), 0755)
	
	// Create a dummy file (should be ignored)
	os.Create(filepath.Join(versionsDir, "not-a-folder.txt"))
	
	// Create the 'current' symlink folder (should be ignored)
	os.Mkdir(filepath.Join(versionsDir, "current"), 0755)

	// 3. Set the environment variable
	oldEnv := os.Getenv("PVM_HOME")
	os.Setenv("PVM_HOME", tempDir)
	defer os.Setenv("PVM_HOME", oldEnv) // Restore original env after test

	// 4. Run the function
	installed, err := RetrieveInstalledPHPVersions()

	// 5. Assertions
	assert.Nil(t, err)
	assert.Equal(t, 2, len(installed), "Should find exactly 2 versions (ignoring file and 'current')")
	
	// Check first version (8.1.0 TS)
	assert.Equal(t, 8, installed[0].Major)
	assert.True(t, installed[0].ThreadSafe)
	
	// Check second version (8.2.0 NTS)
	assert.Equal(t, 8, installed[1].Major)
	assert.Equal(t, 2, installed[1].Minor)
	assert.False(t, installed[1].ThreadSafe)
}

func Test_RetrieveInstalledPHPVersions_NoEnv(t *testing.T) {
	// Clear environment variable
	oldEnv := os.Getenv("PVM_HOME")
	os.Unsetenv("PVM_HOME")
	defer os.Setenv("PVM_HOME", oldEnv)

	installed, err := RetrieveInstalledPHPVersions()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "PVM_HOME environment variable is not set")
	assert.Equal(t, 0, len(installed))
}