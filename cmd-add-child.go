package main

import (
	"fmt"
	"os"

	"github.com/pjsoftware/go-issue/internal/util"
)

type AddChildOptions struct {
	Title string `short:"t" long:"title" description:"title for child item" required:"true"`
}

func (opt *AddChildOptions) Execute(args []string) error {
	if len(options.IssueIDs) == 0 {
		fmt.Printf("add-child: no issue numbers specified\n")
		os.Exit(ecNoIssuesSpecified)
	}

	issueFiles := FindIssueFiles("open", options.IssueIDs)
	if len(issueFiles) == 0 {
		fmt.Printf("add-child: none of the specified issues are actually 'open'\n")
		os.Exit(ecNoOpenIssuesFound)
	}

	if len(issueFiles) > 1 {
		fmt.Printf("add-child: only one issue may be specified\n")
		os.Exit(ecTooManyIssuesFound)
	}

	for _, file := range issueFiles {
		issue := ReadIssue(file.FilePath)
		task := Task{desc: opt.Title}
		issue.childTasks = append(issue.childTasks, task)
		issue.PrintSummary()

		name, err := getUserName()
		util.PanicOn(err)
		issue.AppendActivity(name, "Add child task: "+opt.Title)

		issue.WriteIssue()
		fmt.Printf("new child task added to issue %d\n", issue.id)
	}

	return nil
}
