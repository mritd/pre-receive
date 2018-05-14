package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type CommitType string

const (
	FEAT     CommitType = "feat"
	FIX      CommitType = "fix"
	DOCS     CommitType = "docs"
	STYLE    CommitType = "style"
	REFACTOR CommitType = "refactor"
	TEST     CommitType = "test"
	CHORE    CommitType = "chore"
	PERF     CommitType = "perf"
	HOTFIX   CommitType = "hotfix"
)
const CommitMessagePattern = `^(?:fixup!\s*)?(\w*)(\(([\w\$\.\*/-].*)\))?\: (.*)|^Merge\ branch(.*)`

const checkFailedMeassge = `##############################################################################
##                                                                          ##
## Commit message style check failed!                                       ##
##                                                                          ##
## Commit message style must satisfy this regular:                          ##
##   ^(?:fixup!\s*)?(\w*)(\(([\w\$\.\*/-].*)\))?\: (. *)|^Merge\ branch(.*) ##
##                                                                          ##
## Example:                                                                 ##
##   feat(test): test commit style check.                                   ##
##                                                                          ##
##############################################################################`

const strictMode = false

var commitMsgReg = regexp.MustCompile(CommitMessagePattern)

func main() {

	input, _ := ioutil.ReadAll(os.Stdin)
	param := strings.Fields(string(input))

	// allow branch/tag delete
	if param[1] == "0000000000000000000000000000000000000000" {
		os.Exit(0)
	}

	commitMsg := getCommitMsg(param[0], param[1])
	for _, tmpStr := range commitMsg {
		commitTypes := commitMsgReg.FindAllStringSubmatch(tmpStr, -1)

		if len(commitTypes) != 1 {
			checkFailed()
		} else {
			switch commitTypes[0][1] {
			case string(FEAT):
			case string(FIX):
			case string(DOCS):
			case string(STYLE):
			case string(REFACTOR):
			case string(TEST):
			case string(CHORE):
			case string(PERF):
			case string(HOTFIX):
			default:
				if !strings.HasPrefix(tmpStr, "Merge branch") {
					checkFailed()
				}
			}
		}
		if !strictMode {
			os.Exit(0)
		}
	}

}

func getCommitMsg(odlCommitID, commitID string) []string {
	getCommitMsgCmd := exec.Command("git", "log", odlCommitID+".."+commitID, "--pretty=format:%s")
	getCommitMsgCmd.Stdin = os.Stdin
	getCommitMsgCmd.Stderr = os.Stderr
	b, err := getCommitMsgCmd.Output()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	commitMsg := strings.Split(string(b), "\n")
	return commitMsg
}

func checkFailed() {
	fmt.Fprintln(os.Stderr, checkFailedMeassge)
	os.Exit(1)
}
