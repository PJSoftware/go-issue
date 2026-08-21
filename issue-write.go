package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	git "github.com/pjsoftware/go-git"
	"github.com/pjsoftware/go-issue/internal/util"
)

func (iss *Issue) WriteIssue() {
	fn := fmt.Sprintf("%04d.%s", iss.id, IssuesExtension)
	path := git.FileInRepo(filepath.Join(IssuesFolder, iss.status.Equiv(), fn))
	if path != iss.path {
		os.Remove(iss.path)
		iss.path = path
	}

	file, err := os.Create(iss.path)
	util.PanicOn(err)
	defer file.Close()

	iss.writeHeader(file)
	iss.writeBody(file)
	iss.writeAddendum(file)
}

func (iss *Issue) writeHeader(file *os.File) {
	test(fmt.Fprintf(file, "== Issue %d\n\n", iss.id))
	test(fmt.Fprintf(file, "----\n"))

	test(fmt.Fprintf(file, "id: %d\n", iss.id))

	if iss.ref != "" {
		test(fmt.Fprintf(file, "ref: %s\n", iss.ref))
	}

	test(fmt.Fprintf(file, "title: %s\n", iss.title))
	test(fmt.Fprintf(file, "status: %s\n", iss.status.String()))
	test(fmt.Fprintf(file, "priority: %s\n", iss.priority.String()))
	test(fmt.Fprintf(file, "type: %s\n", iss.issueType.String()))
	test(fmt.Fprintf(file, "created: %s\n", iss.created.Format("2006-01-02")))
	test(fmt.Fprintf(file, "updated: %s\n", iss.updated.Format("2006-01-02")))
	test(fmt.Fprintf(file, "creator: %s\n", iss.creator))
	test(fmt.Fprintf(file, "assignee: %s\n", iss.assignee))

	if iss.milestone != "" {
		test(fmt.Fprintf(file, "milestone: %s\n", iss.milestone))
	}

	if len(iss.tags) > 0 {
		test(fmt.Fprintf(file, "tags: %s\n", strings.Join(iss.tags, ", ")))
	}

	if len(iss.issueRefs) > 0 {
		refStr := ""
		for _, ref := range iss.issueRefs {
			if refStr != "" {
				refStr += ", "
			}
			refStr += fmt.Sprintf("%d", ref)
		}
		test(fmt.Fprintf(file, "refer: %s\n", refStr))
	}

	test(fmt.Fprintf(file, "----\n\n"))
}

func writeSectComment(file *os.File) {
	test(fmt.Fprintf(file, "// section header; please do not edit the following line\n"))
	test(fmt.Fprintf(file, "// or insert any blank lines in the section\n"))
}

func (iss *Issue) writeBody(file *os.File) {
	writeSectComment(file)
	test(fmt.Fprintf(file, "%s\n", HdrDescription))
	for _, line := range iss.description {
		test(fmt.Fprintf(file, "%s\n", line))
	}
	test(fmt.Fprintf(file, "\n"))

	writeSectComment(file)
	test(fmt.Fprintf(file, "%s\n", HdrChildTasks))
	for _, task := range iss.childTasks {
		test(fmt.Fprintf(file, "%s\n", task.output()))
	}
	test(fmt.Fprintf(file, "\n"))

	writeSectComment(file)
	test(fmt.Fprintf(file, "%s\n", HdrSteps))
	for _, line := range iss.steps {
		test(fmt.Fprintf(file, "%s\n", line))
	}
	test(fmt.Fprintf(file, "\n"))

	writeSectComment(file)
	test(fmt.Fprintf(file, "%s\n", HdrBehaviour))
	for _, line := range iss.behaviour {
		test(fmt.Fprintf(file, "%s\n", line))
	}
	test(fmt.Fprintf(file, "\n"))

	writeSectComment(file)
	test(fmt.Fprintf(file, "%s\n", HdrNotes))
	for _, line := range iss.notes {
		test(fmt.Fprintf(file, "%s\n", line))
	}
	test(fmt.Fprintf(file, "\n"))

	writeSectComment(file)
	test(fmt.Fprintf(file, "%s\n", HdrActivity))
	for _, act := range iss.activity {
		test(fmt.Fprintf(file, "%s\n", act.output()))
	}
	test(fmt.Fprintf(file, "\n"))
}

func (iss *Issue) writeAddendum(file *os.File) {
	test(fmt.Fprintf(file, "// The extended notes section below is free-form; if you need additional notes,\n"))
	test(fmt.Fprintf(file, "// uncomment the header line. You may add any content you like, but do not edit\n"))
	test(fmt.Fprintf(file, "// the header\n"))
	if len(iss.addendum) == 0 {
		test(fmt.Fprintf(file, "// %s\n", HdrAddendum))
		return
	}

	test(fmt.Fprintf(file, "%s\n", HdrAddendum))
	for _, line := range iss.addendum {
		test(fmt.Fprintf(file, "%s\n", line))
	}
}

func test(n int, err error) {
	if n == 0 {
		fmt.Printf("* Warning: 0 bytes printed\n")
	}
	util.PanicOn(err)
}

func (act *Activity) output() string {
	return fmt.Sprintf("- *%s (%s):* %s", act.date.Format("2006-01-02"), act.by, act.desc)
}

func (task *Task) output() string {
	doneStr := " "
	if task.done {
		doneStr = "x"
	}
	return fmt.Sprintf("- [%s] %s", doneStr, task.desc)
}
