package improper_artifact_integrity_validation

import "tool/app/internal/detector"

var regexCache = `Release|release`

var CachingInRelease = detector.Detector{
	Name: "caching-in-release",
	Info: detector.Info{
		Description: "Caching in a release workflow can lead to supply chain attacks such as cache poisoning. This is especially dangerous when using self-hosted runners.",
		Message:     "Caching should never be done in release workflows",
		Severity:    5,
		CICDSEC:     9,
	},
	Rule: &detector.And{
		LHS: &detector.Match{
			LHS: "$.jobs..steps[*].uses",
			RHS: "actions/cache",
		},
		RHS: &detector.Or{
			LHS: &detector.Match{
				LHS: "$.name",
				RHS: regexCache,
			},
			RHS: &detector.Or{
				LHS: &detector.Match{
					LHS: "$.on..branches[*]",
					RHS: regexCache,
				},
				RHS: &detector.Or{
					LHS: &detector.Match{
						LHS: "$.on..types[*]",
						RHS: regexCache,
					},
					RHS: &detector.Match{
						LHS: "$.on..tags[*]",
						RHS: regexCache,
					},
				},
			},
		},
	},
}
