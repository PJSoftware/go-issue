package main

import (
	"fmt"
	"slices"
)

type ListOptions struct {
	Status string `long:"status" choice:"open" choice:"closed" choice:"all" description:"issue status to list" default:"open"`
	TopTen bool   `long:"top" description:"list top ten open issues in order of priority"`
}

func (opt *ListOptions) Execute(args []string) error {
	if opt.TopTen {
		return listTopTen()
	}

	issueFiles := FindIssueFiles(opt.Status, options.IssueIDs)

	count := 0
	for _, file := range issueFiles {
		issue := ReadIssue(file.FilePath)
		issue.PrintSummary()
		count++
	}

	if count == 0 {
		fmt.Printf("No issues found in '%s'\n", IssuesFolder)
	}

	return nil
}

func listTopTen() error {
	issueFiles := FindIssueFiles(StatusOpen.String(), nil)
	if len(issueFiles) == 0 {
		fmt.Printf("no open issues found\n")
		return nil
	}

	issues := []*Issue{}
	for _, file := range issueFiles {
		issue := ReadIssue(file.FilePath)
		issues = append(issues, issue)
	}
	slices.SortStableFunc(issues, sortTopTen)

	if len(issues) > 10 {
		issues = issues[:10]
	}

	for _, issue := range issues {
		issue.PrintSummary()
	}

	return nil
}

// sortTopTen() compares two issues such that the issues will be sorted from highest to lowest priority.
//
// If two issues have the _same_ priority, earlier issues will be considered to have the higher priority.
func sortTopTen(a, b *Issue) int {
	switch {
	case a.priority.value < b.priority.value:
		return 1
	case a.priority.value > b.priority.value:
		return -1
	default:
		switch {
		case a.created.After(b.created):
			return 1
		case a.created.Before(b.created):
			return -1
		default:
			return 0
		}
	}
}
