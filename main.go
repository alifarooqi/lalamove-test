package main

import (
	"context"
	"fmt"
	"sort"
	"os"
	"encoding/csv"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/github"
)

// For sorting in descending order of SemVer
type sortSemVer []*semver.Version

func (s sortSemVer) Len() int {
    return len(s)
}
func (s sortSemVer) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}
func (s sortSemVer) Less(i, j int) bool {
    return !s[i].LessThan(*s[j])
}


// LatestVersions returns a sorted slice with the highest version as its first element and the highest version of the smaller minor versions in a descending order
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {
	var versionSlice []*semver.Version
	// This is just an example structure of the code, if you implement this interface, the test cases in main_test.go are very easy to run

	// Populating versionSlice with versions greater than minVersion
	sort.Sort(sortSemVer(releases)) // Sort the array in descending order
	var prevMajor int64 = -1
	var prevMinor int64 = -1
	for  _, v := range releases{
		if(!v.LessThan(*minVersion) && (v.Major != prevMajor || (v.Major == prevMajor && v.Minor != prevMinor))){
			versionSlice = append(versionSlice, v)
			prevMajor = v.Major
			prevMinor = v.Minor
		}
	}
	return versionSlice
}

// Get input from the file in arg[0] and output it as a 2d array of string
func readInput(filename string) (inputArray [][]string, err error) {
    f, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer f.Close() // this needs to be after the err check

	lines, err := csv.NewReader(f).ReadAll()
    if err != nil {
        return nil, err
    }

    inputArray = make([][]string, len(lines)-1)
    for i := range inputArray {
	    inputArray[i] = make([]string, 2)
	}
    for i, line := range lines {
    	if i==0{
    		continue
    	}
        inputArray[i-1][0] = line[0]
        inputArray[i-1][1] = line[1]
    }
    return inputArray, nil
}


func printLatestVersions(inputArray [][]string){
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.ListOptions{}
	for _, dependency := range inputArray{
		dependencyInfo := strings.Split(dependency[0], "/")
		releases, _, err := client.Repositories.ListReleases(ctx, dependencyInfo[0], dependencyInfo[1], opt)
		if err != nil {
			fmt.Printf("Error! Letest version of %s not found\n%s\n", dependency[0], err)
			continue
		}
		minVersion := semver.New(dependency[1])
		allReleases := make([]*semver.Version, len(releases))
		for i, release := range releases {
			versionString := release.GetTagName()
			if versionString[0] == 'v' {
				versionString = versionString[1:]
			}
			allReleases[i] = semver.New(versionString)
		}
		versionSlice := LatestVersions(allReleases, minVersion)

		fmt.Printf("latest versions of %s: %s\n", dependency[0], versionSlice)
	}
}
// Here we implement the basics of communicating with github through the library as well as printing the version
// You will need to implement LatestVersions function as well as make this application support the file format outlined in the README
// Please use the format defined by the fmt.Printf line at the bottom, as we will define a passing coding challenge as one that outputs
// the correct information, including this line
func main() {
	// Check if filename is provided in the argument
	if len(os.Args) != 2 {
         fmt.Printf("usage: %s [filename]\n", os.Args[0])
         os.Exit(1)
    }
	filename := os.Args[1]
	inputArray, err := readInput(filename)
	if err != nil {
		panic(err) 
	}
	printLatestVersions(inputArray)	
}
