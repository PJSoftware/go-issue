package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"time"

	"github.com/pjsoftware/go-issue/internal/util"
)

type Activity struct {
	date time.Time
	by   string
	desc string
}

type Task struct {
	done bool
	desc string
}

const (
	stateStarting = iota
	stateHeader
	stateIntermediate
	stateDescription
	stateChildTasks
	stateSteps
	stateBehaviour
	stateNotes
	stateActivity
	stateAddendum
)

const (
	HdrDescription = "Description::"
	HdrChildTasks  = "Child Tasks::"
	HdrSteps       = "Steps to Reproduce::"
	HdrBehaviour   = "Expected Behaviour::"
	HdrNotes       = "Proposed Solution / Notes::"
	HdrActivity    = "Activity Log::"
	HdrAddendum    = "=== Extended Notes"
)

type Issue struct {
	// metadata

	path string

	// required header fields

	id        int
	title     string
	status    IssueStatus
	priority  IssuePriority
	issueType IssueType
	created   time.Time
	updated   time.Time
	creator   string
	assignee  string

	// optional header fields

	ref       string
	milestone string
	tags      []string
	issueRefs []int

	// body fields

	description []string
	childTasks  []Task
	steps       []string
	behaviour   []string
	notes       []string
	activity    []Activity

	// additional notes

	addendum []string
}

var (
	reHeader      = regexp.MustCompile(`^== Issue (\d+)$`)
	reHeaderOld   = regexp.MustCompile(`^== Issue (\d+): (.+)$`)
	reFileDate    = regexp.MustCompile(`/\d{4}-(\d{8})-`)
	reHeaderEntry = regexp.MustCompile(`([^:]+): (.+)`)
	reActivity    = regexp.MustCompile(`- [*](\d{4}-\d{2}-\d{2}) [(](.+)[)]:[*] (.+)`)
	reChildTask   = regexp.MustCompile(`- [[]([ xX])[]] (.+)`)
)

func NewIssue(id int) *Issue {
	iss := &Issue{}
	iss.id = id
	iss.status = StatusOpen

	iss.created = time.Now()
	iss.updated = iss.created

	creator, err := getUserName()
	util.PanicOn(err)
	iss.creator = creator
	iss.assignee = creator

	return iss
}

func (iss *Issue) SetTitle(title string) *Issue {
	iss.title = title
	return iss
}

func (iss *Issue) SetPriority(pr string) *Issue {
	iss.priority = IssuePriorityFrom(pr)
	return iss
}

func (iss *Issue) SetType(t string) *Issue {
	iss.issueType = IssueTypeFrom(t)
	return iss
}

func (iss *Issue) AppendDescription(desc string) *Issue {
	iss.description = append(iss.description, desc)
	return iss
}

func (iss *Issue) AppendActivity(name string, desc string) *Issue {
	iss.updated = time.Now()

	act := Activity{date: iss.updated, by: name, desc: desc}
	iss.activity = append(iss.activity, act)

	return iss
}

func (iss *Issue) Edit() {
	cmd := exec.Command("code", "-r", iss.path)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("error opening issue in code: %v\n", err)
	}
}

func (iss *Issue) ChildTasksOpen() int {
	if len(iss.childTasks) == 0 {
		return 0
	}

	open := 0
	for _, ct := range iss.childTasks {
		if !ct.done {
			open++
		}
	}
	return open
}

// replace with git.GetUserName()
func getUserName() (string, error) {
	return "Peter Jones", nil
}
