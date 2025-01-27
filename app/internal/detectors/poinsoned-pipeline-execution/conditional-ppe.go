package poinsoned_pipeline_execution

import (
	"tool/app/internal/detector"
)

var ConditionalPPE = detector.Detector{
	Name: "conditional-ppe",
	Info: detector.Info{
		Description:    "Using 'issues' as trigger, no conditional statement, and Github context or local env variables in run sections can lead to code injection.",
		Message:        "Do not use GitHub context or local env variables in scripts together with an 'issues' trigger",
		Severity:       4,
		Exploitability: -1,
		CICDSEC:        4,
	},
	Rule: &detector.And{
		LHS: &detector.Exists{
			LHS: "$.jobs..[*]~",
			RHS: "if",
		},
		RHS: &detector.And{
			LHS: &detector.Match{
				LHS: "$.on..[*]~",
				RHS: "issues",
			},
			//RHS: &detector.Or{
			//	LHS: &detector.Or{
			//		LHS: &detector.Match{
			//			LHS: "$.jobs..steps[*].with.script",
			//			RHS: regexBad,
			//		},
			//		RHS: &detector.Match{
			//			LHS: "$.jobs..steps[*].run",
			//			RHS: regexBad,
			//		},
			//	},
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
			//},
		},
	},
}
