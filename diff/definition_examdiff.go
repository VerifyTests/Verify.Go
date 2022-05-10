package diff

import "github.com/VerifyTests/Verify.Go/utils"

func defineExamDiff() *ToolDefinition {
	leftArgs := func(temp, target string) []string {
		tempTitle := utils.File.GetFileName(temp)
		targetTitle := utils.File.GetFileName(target)
		return []string{target, temp, "/nh", "/diffonly", "/dn1:" + targetTitle, "/dn2:" + tempTitle}
	}

	rightArgs := func(temp, target string) []string {
		tempTitle := utils.File.GetFileName(temp)
		targetTitle := utils.File.GetFileName(target)
		return []string{temp, target, "/nh", "/diffonly", "/dn1:" + tempTitle, "/dn2:" + targetTitle}
	}

	return &ToolDefinition{
		Kind:             ExamDiff,
		Url:              "https://www.prestosoft.com/edp_examdiffpro.asp",
		Cost:             Paid,
		AutoRefresh:      true,
		IsMdi:            false,
		SupportsText:     true,
		RequiresTarget:   true,
		BinaryExtensions: []string{},
		Windows: OsSettings{
			TargetLeftArguments:  leftArgs,
			TargetRightArguments: rightArgs,
			ExePaths: []string{
				"%ProgramFiles%\\ExamDiff Pro\\ExamDiff.exe",
			},
		},
		Notes: " * [Command line reference](https://www.prestosoft.com/ps.asp?page=htmlhelp/edp/command_line_options)\n * `/nh`: do not add files or directories to comparison history\n * `/diffonly`: diff-only merge mode: hide the merge pane",
	}
}