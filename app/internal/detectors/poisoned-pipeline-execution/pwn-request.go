package poisoned_pipeline_execution

import "tool/app/internal/detector"

var regexTarget = `\$\{\{\s*github.event.pull_request.(head.ref|head.sha|merge_commit_sha|head.repo.id|id|head.repo.full_name)\s*\}\}`
var regexRefs = `/refs/pull/\$\{\{\s*github.event.pull_request.number\s*\}\}/merge/`

var PwnRequest = detector.Detector{
	Name: "pwn-request",
	Info: detector.Info{
		Description: "Using 'pull request target' as a trigger, together with a checkout of the code from the pull request branch, can make the workflow susceptible to remote code execution",
		Message:     "Do not checkout code coming from pull_request_target",
		Severity:    5,
		CICDSEC:     4,
	},
	Rule: &detector.And{
		LHS: &detector.And{
			LHS: &detector.Or{
				LHS: &detector.Or{
					LHS: &detector.Match{
						LHS: "$.on",
						RHS: "pull_request_target",
					},
					RHS: &detector.Match{
						LHS: "$.on[*]",
						RHS: "pull_request_target",
					},
				},
				RHS: &detector.Match{
					LHS: "$.on..[*]~",
					RHS: "pull_request_target",
				},
			},
			RHS: &detector.Match{
				LHS: "$.jobs..steps[*].uses",
				RHS: "actions/checkout",
			},
		},
		RHS: &detector.Or{
			LHS: &detector.Or{
				LHS: &detector.Match{
					LHS: "$.jobs..steps[*].with.ref",
					RHS: regexTarget,
				},
				RHS: &detector.Match{
					LHS: "$.jobs..steps[*].with.ref",
					RHS: regexRefs,
				},
			},
			RHS: &detector.Or{
				LHS: &detector.Match{
					LHS: "$.jobs..steps[*].with.repository",
					RHS: regexTarget,
				},
				RHS: &detector.Match{
					LHS: "$.jobs..steps[*].with.repository",
					RHS: regexRefs,
				},
			},
		},
	},
}
