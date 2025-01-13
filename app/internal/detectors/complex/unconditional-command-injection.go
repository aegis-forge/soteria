package complex

import "tool/app/internal/detector"

var UnconditionalCommandInjection = detector.Detector{
	Name: "unconditional-command-injection",
	Info: detector.Info{
		Description: "Using 'issues' as trigger, no conditional statement, and Github context or local env variables in run sections can lead to code injection.",
		Message:     "Do not use GitHub context or local env variables in scripts together with an 'issues' trigger",
		Severity:    4,
		CWE:         82,
	},
	Rule: &detector.And{
		LHS: &detector.Exists{
			NOT: true,
			LHS: "$.jobs..[*]~",
			RHS: "if",
		},
		RHS: &detector.And{
			LHS: &detector.Match{
				LHS: "$.on..[*]~",
				RHS: "issues",
			},
			RHS: &detector.Or{
				LHS: &detector.Or{
					LHS: &detector.Match{
						LHS: "$.jobs..steps[*].with.script",
						RHS: regexBad,
					},
					RHS: &detector.Match{
						LHS: "$.jobs..steps[*].run",
						RHS: regexBad,
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
	},
}
