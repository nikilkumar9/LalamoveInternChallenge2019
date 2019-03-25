package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/github"
)

func reverse(releases []*semver.Version) []*semver.Version {
	for i := 0; i < len(releases)/2; i++ {
		j := len(releases) - i - 1
		releases[i], releases[j] = releases[j], releases[i]
	}
	return releases
}

// LatestVersions returns a sorted slice with the highest version as its first element and the highest version of the smaller minor versions in a descending order
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {
	var versionSlice []*semver.Version
	// This is just an example structure of the code, if you implement this interface, the test cases in main_test.go are very easy to run

	// Sorting array of version numbers by ascending order and reversing the order
	semver.Sort(releases)
	releases = reverse(releases)

	//Removing versions that are lower than minimum
	for i := 0; i < len(releases); i++ {
		if !(releases[i].LessThan(*minVersion)) {
			versionSlice = append(versionSlice, releases[i])
		} else {
			break
		}
	}

	// Remove Duplicate copies (code not complete)

	return versionSlice
}

// Here we implement the basics of communicating with github through the library as well as printing the version
// You will need to implement LatestVersions function as well as make this application support the file format outlined in the README
// Please use the format defined by the fmt.Printf line at the bottom, as we will define a passing coding challenge as one that outputs
// the correct information, including this line
func main() {
	// Github
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.ListOptions{PerPage: 100}

	//Opening file and getting names
	file, err := os.Open("file.txt")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string
	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}
	file.Close()

	txtlines = txtlines[1:] //Skip header line
	for _, line := range txtlines {
		software := strings.Split(line, ",")[0]
		minVal := strings.Split(line, ",")[1]
		individual := strings.Split(software, "/")

		releases, _, err := client.Repositories.ListReleases(ctx, individual[0], individual[1], opt)
		if err != nil {
			log.Fatal(err)
		}

		minVersion := semver.New(minVal)
		allReleases := make([]*semver.Version, len(releases))
		for _, release := range releases {
			versionString := *release.TagName
			if versionString[0] == 'v' {
				versionString = versionString[1:]
			}
			versionString = strings.Split(versionString, "-")[0]

		}
		versionSlice := LatestVersions(allReleases, minVersion)

		fmt.Printf("latest versions of kubernetes/kubernetes: %s", versionSlice)
	}
}
