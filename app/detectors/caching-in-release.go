package detectors

import "tool/app/internal/detector"

var CachingInRelease = detector.Detector{
	Name: "caching-in-release",
	Info: detector.Info{
		Description: "Caching in a release workflow can lead to supply chain attacks such as cache poisoning. This is especially dangerous when using self-hosted runners.",
		Severity:    5,
		CWE:         349,
	},
	Rule: &detector.And{
		LHS: &detector.Match{
			LHS: "$.jobs..steps[*].uses",
			RHS: "actions/cache",
		},
		RHS: &detector.Or{
			LHS: &detector.Match{
				LHS: "$.name",
				RHS: "R|release?",
			},
			RHS: &detector.Or{
				LHS: &detector.Match{
					LHS: "$.on..branches[*]",
					RHS: "R|release?",
				},
				RHS: &detector.Or{
					LHS: &detector.Match{
						LHS: "$.on..types[*]",
						RHS: "R|release?",
					},
					RHS: &detector.Match{
						LHS: "$.on..tags[*]",
						RHS: "R|release?",
					},
				},
			},
		},
	},
}
