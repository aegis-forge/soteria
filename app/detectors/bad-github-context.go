package detectors

import (
	"tool/app/internal/detector"
)

var regexContext = `(github\.event\.issue\.title|github\.event\.issue\.body|github\.event\.pull_request\.title|github\.event\.pull_request\.body|github\.event\.comment\.body|github\.event\.review\.body|github\.event\.pages\.[A-z]+\.page_name|github\.event\.commits\.[A-z]+\.message|github\.event\.head_commit\.message|github\.event\.head_commit\.author\.email|github\.event\.head_commit\.author\.name|github\.event\.commits\.[A-z]+\.author\.email|github\.event\.commits\.[A-z]+\.author\.name|github\.event\.pull_request\.head\.ref|github\.event\.pull_request\.head\.label|github\.event\.pull_request\.head\.repo\.default_branch|github\.head_ref)`

var BadGithubContext = detector.Detector{
	Name: "bad-github-context",
	Info: detector.Info{
		Description: "Using the Github context env variables in run sections can lead to secret disclosure through code injection.",
		Severity:    4,
		CWE:         82,
	},
	Rule: &detector.Or{
		LHS: &detector.Match{
			LHS: "$.jobs..steps[*].with.script",
			RHS: regexContext,
		},
		RHS: &detector.Match{
			LHS: "$.jobs..steps[*].run",
			RHS: regexContext,
		},
	},
}
