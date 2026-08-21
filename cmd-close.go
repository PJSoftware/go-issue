package main

import (
	"fmt"
	"os"

	"github.com/pjsoftware/go-issue/internal/util"
)

type CloseOptions struct {
	Ignore bool `long:"ignore" description:"flag closed issue as 'ignored'"`
	Cancel bool `long:"cancel" description:"flag closed issue as 'cancelled'"`
	Edit   bool `long:"edit" description:"edit new issue after creation"`
}

func (opt *CloseOptions) Execute(args []string) error {
	if len(options.IssueIDs) == 0 {
		fmt.Printf("close: no issue numbers specified\n")
		os.Exit(ecNoIssuesSpecified)
	}

	issueFiles := FindIssueFiles("open", options.IssueIDs)
	if len(issueFiles) == 0 {
		fmt.Printf("close: none of the specified issues are actually 'open'\n")
		os.Exit(ecNoOpenIssuesFound)
	}

	count := 0
	fmt.Printf("closing issues:\n")

	for _, file := range issueFiles {
		issue := ReadIssue(file.FilePath)
		if open := issue.ChildTasksOpen(); open > 0 {
			fmt.Printf("cannot close issue %d with %d open child tasks\n", issue.id, open)
			continue
		}

		if opt.Cancel {
			issue.status = StatusCancelled
		} else if opt.Ignore {
			issue.status = StatusIgnored
		} else {
			issue.status = StatusClosed
		}

		name, err := getUserName()
		util.PanicOn(err)
		issue.AppendActivity(name, "Issue "+issue.status.slug)
		issue.PrintSummary()
		issue.WriteIssue()

		if opt.Edit {
			fmt.Printf("editing issue %d in 'code':\n", issue.id)
			issue.Edit()
		}
		count++
	}

	fmt.Printf("%d closed of %d specified\n", count, len(options.IssueIDs))
	return nil
}
