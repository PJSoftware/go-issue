package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/jessevdk/go-flags"
	find "github.com/pjsoftware/go-find-files"
	"github.com/pjsoftware/go-find-files/opt"
	git "github.com/pjsoftware/go-git"
	"github.com/pjsoftware/go-issue/internal/util"
)

const (
	IssuesFolder       = ".issues"
	FolderIssuesOpen   = IssuesFolder + "/open"
	FolderIssuesClosed = IssuesFolder + "/closed"

	IssuesExtension = "adoc"
)

type GlobalOptions struct {
	IssueIDs []int
	Version  bool            `short:"v" long:"version" description:"show version information"`
	ListCmd  ListOptions     `command:"list" description:"list issues (may include '1 3 7-20' issue selection)"`
	AddCmd   AddOptions      `command:"add" description:"open new issue"`
	CloseCmd CloseOptions    `command:"close" description:"close specified issue(s) (requires issue selection)"`
	AddChild AddChildOptions `command:"add-child" description:"add child task to specified issue (requires issue selection; open issues only)"`
	EditCmd  EditOptions     `command:"edit" description:"edit specified issue(s) (requires issue selection; open issues only)"`
}

var options GlobalOptions

func initFolders() {
	err := os.MkdirAll(git.FolderInRepo(FolderIssuesOpen), 0777)
	util.PanicOn(err)

	err = os.MkdirAll(git.FolderInRepo(FolderIssuesClosed), 0777)
	util.PanicOn(err)
}

func main() {
	initFolders()

	args, issueNumbers := readArgs()
	options.IssueIDs = issueNumbers

	_, err := flags.ParseArgs(&options, args)
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok {
			if flagsErr.Type == flags.ErrHelp {
				os.Exit(ecExitWithoutError) // exit without error on --help
			}
		}

		if options.Version {
			fmt.Printf("issue (Issue Handler) version %s\n", VERSION)
			os.Exit(ecExitWithoutError)
		}

		fmt.Printf("Hint: specify the -h or --help parameter for help\n")
		os.Exit(ecUnknownParameter) // exit with error on non-recognised parameter
	}
}

func issuesIn(folder string) []find.FileData {
	opts := []opt.FFOpt{opt.HasExtension(IssuesExtension)}
	for _, issue := range options.IssueIDs {
		opts = append(opts, opt.HasPrefix(fmt.Sprintf("%04d", issue)))
	}

	return find.InFolder(git.FolderInRepo(folder), opts...)
}

var reArgRange = regexp.MustCompile(`(\d+)-(\d+)`)

func readArgs() ([]string, []int) {
	args := []string{}
	iNum := []int{}

	for _, arg := range os.Args[1:] {
		if reArgRange.MatchString(arg) {
			match := reArgRange.FindStringSubmatch(arg)
			min, _ := strconv.Atoi(match[1])
			max, _ := strconv.Atoi(match[2])
			if min > max {
				min, max = max, min
			}

			for i := min; i <= max; i++ {
				iNum = append(iNum, i)
			}

		} else {
			num, err := strconv.Atoi(arg)
			if err == nil {
				iNum = append(iNum, num)
			} else {
				args = append(args, arg)
			}
		}
	}

	return args, iNum
}
