package main

import (
	"fmt"
	"os"
)

type EditOptions struct {
}

func (opt *EditOptions) Execute(args []string) error {
	if len(options.IssueIDs) == 0 {
		fmt.Printf("edit: no issue numbers specified\n")
		os.Exit(ecNoIssuesSpecified)
	}

	issueFiles := FindIssueFiles("open", options.IssueIDs)
	if len(issueFiles) == 0 {
		fmt.Printf("edit: only 'open' issues may be edited\n")
		os.Exit(ecNoOpenIssuesFound)
	}

	count := 0

	for _, file := range issueFiles {
		issue := ReadIssue(file.FilePath)
		fmt.Printf("editing issue %d in 'code':\n", issue.id)
		issue.Edit()
		count++
	}
	return nil
}
