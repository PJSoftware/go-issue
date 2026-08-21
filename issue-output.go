package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

func (iss *Issue) PrintSummary() {
	status := iss.status
	statusStr := ""

	if status == StatusOpen {
		statusStr = color.GreenString(setLength(status.String(), 6))
	} else {
		statusStr = color.RedString(setLength(status.String(), 6))
	}

	title := iss.title
	if iss.ref != "" {
		title = fmt.Sprintf("[%s] %s", iss.ref, title)
	}

	dateStr := iss.created.Format("2006-01-02")
	elapsed := time.Since(iss.created)
	daysSince := elapsed.Hours() / 24

	switch {
	case daysSince > 90:
		dateStr = color.RedString(dateStr)
	case daysSince > 30:
		dateStr = color.YellowString(dateStr)
	case daysSince > 14:
		dateStr = color.BlueString(dateStr)
	}

	ctStr := ""
	ctLen := 0
	if len(iss.childTasks) > 0 {
		open := iss.ChildTasksOpen()
		done := len(iss.childTasks) - open
		ctStr = fmt.Sprintf("[%d/%d] ", done, len(iss.childTasks))
		ctLen = len(ctStr)
		if open == 0 {
			ctStr = color.GreenString(ctStr)
		} else {
			ctStr = color.RedString(ctStr)
		}
	}

	//          ID ST PR DT TY C TT TG
	fmt.Printf("%s %s %s %s %s %s%s %s\n",
		pad(fmt.Sprintf("%d", iss.id), 4),
		statusStr,
		iss.priority.SymbolColour(),
		dateStr,
		iss.issueType.SymbolColour(),
		ctStr,
		color.WhiteString(setLength(title, 80-ctLen)),
		color.BlueString(strings.Join(iss.tags, ",")),
	)
}

func pad(str string, length int) string {
	for len(str) < length {
		str = " " + str
	}
	return str
}

func setLength(str string, length int) string {
	if len(str) < length {
		for len(str) < length {
			str += " "
		}
		return str
	}

	if len(str) > length {
		str = str[0:length-3] + "..."
	}

	return str
}
