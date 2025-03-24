package poisoned_pipeline_execution

import (
	"tool/app/internal/detector"
)

var regexContext = `(github\.event\.issue\.title|github\.event\.issue\.body|github\.event\.pull_request\.title|github\.event\.pull_request\.body|github\.event\.comment\.body|github\.event\.review\.body|github\.event\.pages\.[A-z]+\.page_name|github\.event\.commits\.[A-z]+\.message|github\.event\.head_commit\.message|github\.event\.head_commit\.author\.email|github\.event\.head_commit\.author\.name|github\.event\.commits\.[A-z]+\.author\.email|github\.event\.commits\.[A-z]+\.author\.name|github\.event\.pull_request\.head\.ref|github\.event\.pull_request\.head\.label|github\.event\.pull_request\.head\.repo\.default_branch|github\.head_ref)`

var ConditionalInjection = detector.Detector{
	Name: "conditional-injection",
	Info: detector.Info{
		Description: "Using 'issues' as trigger, no conditional statement, and Github context or local env variables in run sections can lead to code injection.",
		Message:     "Do not use GitHub context or local env variables in scripts together with an 'issues' trigger",
		Severity:    3,
		CICDSEC:     4,
	},
	Rule: &detector.And{
		LHS: &detector.Exists{
			LHS: "$.jobs..[*]~",
			RHS: "if",
		},
		RHS: &detector.And{
			LHS: &detector.Or{
				LHS: &detector.Or{
					LHS: &detector.Match{
						LHS: "$.on",
						RHS: "issues",
					},
					RHS: &detector.Match{
						LHS: "$.on[*]",
						RHS: "issues",
					},
				},
				RHS: &detector.Match{
					LHS: "$.on..[*]~",
					RHS: "issues",
				},
			},
			RHS: &detector.Or{
				LHS: &detector.Match{
					LHS: "$.jobs..steps[*].with.script",
					RHS: regexContext,
				},
				RHS: &detector.Match{
					LHS: "$.jobs..steps[*].run",
					RHS: regexContext,
				},
			},
		},
	},
}
