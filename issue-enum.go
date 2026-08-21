package main

import (
	"fmt"

	"github.com/fatih/color"
)

// IssueStatus enum

type IssueStatus struct {
	slug  string
	equiv string
}

var (
	StatusOpen       = IssueStatus{"open", "open"}
	StatusInProgress = IssueStatus{"in-progress", "open"}
	StatusBlocked    = IssueStatus{"blocked", "open"}

	StatusClosed    = IssueStatus{"closed", "closed"}
	StatusIgnored   = IssueStatus{"ignored", "closed"}
	StatusCancelled = IssueStatus{"cancelled", "closed"}
)

func IssueStatusFrom(enum string) IssueStatus {
	switch enum {
	case "open":
		return StatusOpen
	case "in-progress":
		return StatusInProgress
	case "blocked":
		return StatusBlocked

	case "closed":
		return StatusClosed
	case "ignored":
		return StatusIgnored
	case "cancelled":
		return StatusCancelled

	default:
		panic(fmt.Sprintf("unknown IssueStatus '%s'", enum))
	}
}

func (enum IssueStatus) String() string {
	return enum.slug
}
func (enum IssueStatus) Equiv() string {
	return enum.equiv
}

// IssuePriority enum

type IssuePriority struct {
	value int
	slug  string
	sym   string
	col   *color.Color
}

var (
	PriorityLow      = IssuePriority{10, "low", "↓", color.New(color.FgGreen)}
	PriorityMedium   = IssuePriority{20, "medium", "~", color.New(color.FgBlue)}
	PriorityHigh     = IssuePriority{30, "high", "↑", color.New(color.FgYellow)}
	PriorityCritical = IssuePriority{40, "critical", "⇑", color.New(color.FgRed)}
)

func (enum IssuePriority) SymbolColour() string {
	return enum.col.Sprint(enum.sym)
}

func (enum IssuePriority) String() string {
	return enum.slug
}

func IssuePriorityFrom(enum string) IssuePriority {
	switch enum {
	case "low":
		return PriorityLow
	case "medium":
		return PriorityMedium
	case "high":
		return PriorityHigh
	case "critical":
		return PriorityCritical

	default:
		panic(fmt.Sprintf("unknown IssuePriority '%s'", enum))
	}
}

// IssueType enum

type IssueType struct {
	slug string
	sym  string
	col  *color.Color
}

var (
	TypeBug   = IssueType{"bug", "X", color.New(color.FgRed)}
	TypeFeat  = IssueType{"feature", "!", color.New(color.FgGreen)}
	TypeTask  = IssueType{"task", "#", color.New(color.FgBlue)}
	TypeDocs  = IssueType{"docs", "@", color.New(color.FgMagenta)}
	TypeQuery = IssueType{"query", "?", color.New(color.FgCyan)}
)

func (enum IssueType) String() string {
	return enum.slug
}

func (enum IssueType) SymbolColour() string {
	return enum.col.Sprint(enum.sym)
}

func IssueTypeFrom(enum string) IssueType {
	switch enum {
	case "bug":
		return TypeBug
	case "feature":
		return TypeFeat
	case "task":
		return TypeTask
	case "docs":
		return TypeDocs
	case "query":
		return TypeQuery

	default:
		panic(fmt.Sprintf("unknown IssueType '%s'", enum))
	}
}
