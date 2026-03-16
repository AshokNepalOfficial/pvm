package common

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Version struct {
	Major      int
	Minor      int
	Patch      int
	Url        string
	ThreadSafe bool
}

func (v Version) Semantic() string {
	return fmt.Sprintf("%v.%v.%v", v.Major, v.Minor, v.Patch)
}

func (v Version) StringShort() string {
	semantic := v.Semantic()
	if v.ThreadSafe {
		return semantic
	}
	return semantic + " nts"
}

func (v Version) String() string {
	semantic := v.Semantic()
	if v.ThreadSafe {
		return semantic + " thread safe"
	}
	return semantic + " non-thread safe"
}

func ComputeVersion(text string, safe bool, url string) Version {
	versionRe := regexp.MustCompile(`([0-9]{1,3})(?:.([0-9]{1,3}))?(?:.([0-9]{1,3}))?`)
	matches := versionRe.FindAllStringSubmatch(text, -1)
	if len(matches) == 0 {
		return Version{}
	}

	major, err := strconv.Atoi(matches[0][1])
	if err != nil {
		major = -1
	}

	minor, err := strconv.Atoi(matches[0][2])
	if err != nil {
		minor = -1
	}

	patch, err := strconv.Atoi(matches[0][3])
	if err != nil {
		patch = -1
	}

	return Version{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		ThreadSafe: safe,
		Url:        url,
	}
}

func (v Version) Compare(o Version) int {
	if v.Major == -1 || o.Major == -1 {
		return 0
	}
	if v.Major != o.Major {
		if v.Major < o.Major {
			return -1
		}
		return 1
	}

	if v.Minor == -1 || o.Minor == -1 {
		return 0
	}
	if v.Minor != o.Minor {
		if v.Minor < o.Minor {
			return -1
		}
		return 1
	}

	if v.Patch == -1 || o.Patch == -1 {
		return 0
	}
	if v.Patch != o.Patch {
		if v.Patch < o.Patch {
			return -1
		}
		return 1
	}

	return 0
}

func (v Version) CompareThreadSafe(o Version) int {
	result := v.Compare(o)
	if result != 0 {
		return result
	}

	if v.ThreadSafe == o.ThreadSafe {
		return 0
	}

	if v.ThreadSafe {
		return -1
	}
	return 1
}

func (v Version) LessThan(o Version) bool {
	return v.CompareThreadSafe(o) == -1
}

func (v Version) Same(o Version) bool {
	return v.CompareThreadSafe(o) == 0
}

func SortVersions(input []Version) []Version {
	sort.SliceStable(input, func(i, j int) bool {
		return input[i].LessThan(input[j])
	})
	return input
}

func RetrievePHPVersions() ([]Version, error) {
	resp, err := http.Get("https://windows.php.net/downloads/releases/archives/")
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	sb := string(body)

	re := regexp.MustCompile(`<A HREF="([a-zA-Z0-9./-]+)">([a-zA-Z0-9./-]+)</A>`)
	matches := re.FindAllStringSubmatch(sb, -1)

	versions := make([]Version, 0)

	for _, match := range matches {
		url := match[1]
		name := match[2]

		if name != "" && len(name) > 15 && (name[:15] == "php-devel-pack-" || name[:15] == "php-debug-pack-") {
			continue
		}
		if name != "" && len(name) > 14 && name[:14] == "php-test-pack-" {
			continue
		}
		if name != "" && strings.Contains(name, "src") {
			continue
		}
		if name != "" && !strings.HasSuffix(name, ".zip") {
			continue
		}

		threadSafe := true
		if name != "" && (strings.Contains(strings.ToLower(name), "nts")) {
			threadSafe = false
		}

		if name != "" && !strings.Contains(name, "x64") {
			continue
		}

		versions = append(versions, ComputeVersion(name, threadSafe, url))
	}
	return versions, nil
}

func RetrieveInstalledPHPVersions() ([]Version, error) {
	versions := make([]Version, 0)

	// Fetch the installation path from the environment variable set by Inno Setup
	pvmPath := os.Getenv("PVM_HOME")
	if pvmPath == "" {
		return versions, errors.New("PVM_HOME environment variable is not set. Please restart your terminal or reinstall PVM")
	}

	// Versions are stored in the 'versions' subdirectory of the app root
	versionsPath := filepath.Join(pvmPath, "versions")
	if _, err := os.Stat(versionsPath); os.IsNotExist(err) {
		return versions, nil // Return empty list if no versions installed yet
	}

	folders, err := os.ReadDir(versionsPath)
	if err != nil {
		return versions, err
	}

	for _, folder := range folders {
		// Only look at directories, skip files (like zips) and the 'current' symlink
		if !folder.IsDir() || folder.Name() == "current" {
			continue
		}

		folderName := folder.Name()
		safe := true
		if strings.Contains(strings.ToLower(folderName), "nts") {
			safe = false
		}

		versions = append(versions, ComputeVersion(folderName, safe, ""))
	}

	SortVersions(versions)
	return versions, nil
}