package common

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/mitchellh/go-homedir"
)

// ListProfilesFromAWSConfig parses default aws config file (~/.aws/config) and returns
// all profiles found in that file.
func ListProfilesFromAWSConfig() []string {
	var profiles []string

	homeDir, err := homedir.Dir()
	if err != nil {
		panic(err.Error())
	}

	var awsConfigFilePath = filepath.Join(homeDir, ".aws", "config")

	file, err := os.Open(awsConfigFilePath)
	if err != nil {
		panic(err.Error())
	}

	defer file.Close()

	reg := regexp.MustCompile(`^\[profile `)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if t := reg.MatchString(scanner.Text()); t == true {
			s := strings.TrimSuffix(reg.ReplaceAllString(scanner.Text(), "${1}"), "]")
			profiles = append(profiles, s)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err.Error())
	}

	sort.Strings(profiles)
	return profiles
}

// SelectProfile prompts to select one AWS Profile from given profile list.
func SelectProfile(profiles []string) string {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "{{ . | cyan }}",
		Inactive: "{{ . | faint }}",
		Details: `
--------- Current profile ----------
{{ "Name:" | faint }}	{{ . }}`,
	}

	searcher := func(input string, index int) bool {
		return strings.Contains(strings.ToLower(profiles[index]), strings.ToLower(input))
	}

	prompt := promptui.Select{
		Label:     "AWS Profiles",
		Items:     profiles,
		Templates: templates,
		Size:      4,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}

	return profiles[i]

}

// SelectBucket prompts to select one bucket from given bucket list.
func SelectBucket(buckets []string) string {
	templates := &promptui.SelectTemplates{
		Active:   `{{ . | cyan | bold }}`,
		Inactive: `{{ . | faint }}`,
		Details: `
--------- Current Bucket ----------
{{ "Name:" | faint }}	{{ . }}`,
	}

	searcher := func(input string, index int) bool {
		return strings.Contains(strings.ToLower(buckets[index]), strings.ToLower(input))
	}

	prompt := promptui.Select{
		Label:        "Buckets: " + strconv.Itoa(len(buckets)),
		Items:        buckets,
		Templates:    templates,
		Size:         15,
		Searcher:     searcher,
		HideSelected: true,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}

	return buckets[i]

}

// ExitWithError exits with an error status and print current error message
func ExitWithError(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
