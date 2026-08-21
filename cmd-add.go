package main

import (
	"fmt"

	"github.com/pjsoftware/go-issue/internal/util"
)

type AddOptions struct {
	Title    string `short:"t" long:"title" description:"title of new issue" required:"true"`
	Priority string `short:"p" long:"priority" choice:"low" choice:"medium" choice:"high" choice:"critical" description:"priority of new issue" default:"medium"`
	Ref      string `short:"r" long:"ref" description:"optional reference for new issue"`
	Type     string `long:"type" choice:"bug" choice:"feature" choice:"task" choice:"docs" choice:"query" description:"type of new issue" default:"task"`
	Edit     bool   `long:"edit" description:"edit new issue after creation"`
}

func (opt *AddOptions) Execute(args []string) error {
	issueFiles := FindIssueFiles("all", nil)

	max := 0
	for _, file := range issueFiles {
		issue := ReadIssue(file.FilePath)
		if issue.id > max {
			max = issue.id
		}
	}

	user, err := getUserName()
	util.PanicOn(err)

	issue := NewIssue(max+1).
		SetTitle(opt.Title).
		SetPriority(opt.Priority).
		SetType(opt.Type).
		AppendDescription(opt.Title).
		AppendActivity(user, "Issue opened")

	if opt.Ref != "" {
		issue.ref = opt.Ref
	}

	issue.WriteIssue()
	fmt.Printf("Added issue %d\n", issue.id)
	issue.PrintSummary()

	if opt.Edit {
		fmt.Printf("editing issue %d in 'code':\n", issue.id)
		issue.Edit()
	}

	return nil
}
