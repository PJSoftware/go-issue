package main

import (
	"bufio"
	"cmp"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	find "github.com/pjsoftware/go-find-files"
	"github.com/pjsoftware/go-issue/internal/util"
)

func FindIssueFiles(status string, idList []int) []find.FileData {
	issueFiles := []find.FileData{}
	if status == "open" || status == "all" || len(options.IssueIDs) > 0 {
		issueFiles = append(issueFiles, issuesIn(FolderIssuesOpen)...)
	}
	if status == "closed" || status == "all" || len(options.IssueIDs) > 0 {
		issueFiles = append(issueFiles, issuesIn(FolderIssuesClosed)...)
	}

	slices.SortFunc(issueFiles, sortFilesByName)
	return issueFiles
}

func ReadIssue(path string) *Issue {
	lines := readLines(path)

	if reHeaderOld.MatchString(lines[0]) {
		issue := NewIssueFromOldFormat(lines, path)
		return issue
	} else if reHeader.MatchString(lines[0]) {
		issue := NewIssueFrom(lines[1:])
		issue.path = path
		return issue
	} else {
		panic(fmt.Sprintf("file %s not recognised issue format", path))
	}
}

func NewIssueFromOldFormat(lines []string, path string) *Issue {
	issue := &Issue{}

	matchHdr := reHeaderOld.FindStringSubmatch(lines[0])
	issue.id, _ = strconv.Atoi(matchHdr[1])
	issue.title = matchHdr[2]
	issue.path = path

	path = strings.ReplaceAll(path, "\\", "/")
	if strings.Contains(path, FolderIssuesClosed) {
		issue.status = StatusClosed
	} else {
		issue.status = StatusOpen
	}

	issue.priority = PriorityMedium
	issue.issueType = TypeTask

	matchPath := reFileDate.FindStringSubmatch(path)
	if matchPath == nil {
		panic(fmt.Sprintf("could not extract date from '%s'", path))
	}
	t, err := time.Parse("20060102", matchPath[1])
	util.PanicOn(err)
	issue.created = t
	issue.updated = t

	issue.creator, err = getUserName()
	util.PanicOn(err)
	issue.assignee = issue.creator
	issue.milestone = "-"

	if len(lines) > 1 {
		nextLine := ""
		for _, nextLine = range lines[1:] {
			if nextLine != "" {
				continue
			}
		}

		if tagStr, trimmed := strings.CutPrefix(nextLine, "Tags: "); trimmed {
			for tag := range strings.SplitSeq(tagStr, ", ") {
				tag = strings.TrimPrefix(tag, "`#")
				tag = strings.TrimSuffix(tag, "`")
				issue.tags = append(issue.tags, tag)
			}
			slices.Sort(issue.tags)
		} else {
			panic(fmt.Sprintf("new format with old header: %s", path))
		}
	}

	issue.description = append(issue.description, issue.title)

	act1 := Activity{issue.created, issue.creator, "Issue raised on Github"}
	issue.activity = append(issue.activity, act1)

	act2 := Activity{time.Now(), issue.creator, "Converted to new format"}
	issue.activity = append(issue.activity, act2)

	return issue
}

func NewIssueFrom(lines []string) *Issue {
	state := stateStarting
	issue := &Issue{}

	for _, line := range lines {
		switch state {
		case stateStarting:
			if line == "----" {
				state = stateHeader
			}

		case stateHeader:
			if line == "----" {
				state = stateIntermediate
			} else {
				issue.parseHeader(line)
			}

		case stateIntermediate:
			if newState, change := stateChange(line); change {
				state = newState
			}

		case stateDescription:
			if newState, change := stateChange(line); change {
				state = newState
				break
			}
			issue.description = append(issue.description, line)

		case stateChildTasks:
			if newState, change := stateChange(line); change {
				state = newState
				break
			}
			task, err := parseTask(line)
			util.PanicOn(err)
			issue.childTasks = append(issue.childTasks, task)

		case stateSteps:
			if newState, change := stateChange(line); change {
				state = newState
				break
			}
			issue.steps = append(issue.steps, line)

		case stateBehaviour:
			if newState, change := stateChange(line); change {
				state = newState
				break
			}
			issue.behaviour = append(issue.behaviour, line)

		case stateNotes:
			if newState, change := stateChange(line); change {
				state = newState
				break
			}
			issue.notes = append(issue.notes, line)

		case stateActivity:
			if newState, change := stateChange(line); change {
				state = newState
				break
			}
			act, err := parseActivity(line)
			util.PanicOn(err)
			issue.activity = append(issue.activity, act)

		case stateAddendum:
			issue.addendum = append(issue.addendum, line)

		}
	}

	return issue
}

func parseTask(line string) (Task, error) {
	task := Task{}
	if match := reChildTask.FindStringSubmatch(line); match != nil {
		doneStr := match[1]
		task.done = doneStr != " "
		task.desc = match[2]
		return task, nil

	} else {
		fmt.Printf("* Warning: could not parse task '%s'\n", line)
		return task, fmt.Errorf("could not parse task: %s", line)
	}
}

func parseActivity(line string) (Activity, error) {
	act := Activity{}
	if match := reActivity.FindStringSubmatch(line); match != nil {
		t, err := time.Parse("2006-01-02", match[1])
		if err != nil {
			fmt.Printf("* Warning: could not parse date '%s' from activity '%s'\n", match[1], line)
			return act, err
		}

		act.date = t
		act.by = match[2]
		act.desc = match[3]
		return act, nil

	} else {
		fmt.Printf("* Warning: could not parse activity '%s'\n", line)
		return act, fmt.Errorf("could not parse activity: %s", line)
	}
}

func stateChange(line string) (int, bool) {
	switch line {
	case HdrDescription:
		return stateDescription, true

	case HdrChildTasks:
		return stateChildTasks, true

	case HdrSteps:
		return stateSteps, true

	case HdrBehaviour:
		return stateBehaviour, true

	case HdrNotes:
		return stateNotes, true

	case HdrActivity:
		return stateActivity, true

	case HdrAddendum:
		return stateAddendum, true

	case "":
		return stateIntermediate, true

	default:
		return 0, false
	}
}

func (iss *Issue) parseHeader(line string) {
	if match := reHeaderEntry.FindStringSubmatch(line); match != nil {
		key := match[1]
		val := match[2]
		val = strings.TrimPrefix(val, `"`)
		val = strings.TrimSuffix(val, `"`)

		switch key {
		case "id":
			id, err := strconv.Atoi(val)
			util.PanicOn(err)
			iss.id = id

		case "title":
			iss.title = val

		case "status":
			iss.status = IssueStatusFrom(val)

		case "priority":
			iss.priority = IssuePriorityFrom(val)

		case "type":
			iss.issueType = IssueTypeFrom(val)

		case "created":
			t, err := time.Parse("2006-01-02", val)
			util.PanicOn(err)
			iss.created = t

		case "updated":
			t, err := time.Parse("2006-01-02", val)
			util.PanicOn(err)
			iss.updated = t

		case "creator":
			iss.creator = val

		case "assignee":
			iss.assignee = val

		case "ref":
			iss.ref = val

		case "milestone":
			if val != "-" {
				iss.milestone = val
			}

		case "tags":
			iss.tags = strings.Split(val, ", ")

		case "refer":
			refs := strings.Split(val, ", ")
			for _, refStr := range refs {
				refStr = strings.TrimPrefix(refStr, "#")
				ref, err := strconv.Atoi(refStr)
				if err == nil {
					iss.issueRefs = append(iss.issueRefs, ref)
				} else {
					fmt.Printf("* Warning: could not parse '%s' as ref in issue %d\n", refStr, iss.id)
				}
			}

		default:
			fmt.Printf("* Warning: Unknown header '%s' in issue %d\n", key, iss.id)

		}
	}
}

func readLines(path string) []string {
	file, err := os.Open(path)
	util.PanicOn(err)
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "// ") {
			lines = append(lines, line)
		}
	}
	util.PanicOn(scanner.Err())

	return lines
}

func sortFilesByName(a, b find.FileData) int {
	return cmp.Compare(a.NameOnly, b.NameOnly)
}
