package detectors

import "tool/app/internal/detector"

var CoarsePermissions = detector.Detector{
	Name: "coarse-permissions",
	Info: detector.Info{
		Description: "Permissions shouldn't be coarse, they should be finegrained on the specific permissions (no default/read-all/write-all permissions).",
		Severity:    2,
		CWE:         -1,
	},
	Rule: &detector.Or{
		LHS: &detector.Or{
			LHS: &detector.Match{
				LHS: "$.permissions",
				RHS: "read-all",
			},
			RHS: &detector.Match{
				LHS: "$.permissions",
				RHS: "write-all",
			},
		},
		RHS: &detector.Or{
			LHS: &detector.Match{
				LHS: "$.jobs..permissions",
				RHS: "read-all",
			},
			RHS: &detector.Match{
				LHS: "$.jobs..permissions",
				RHS: "write-all",
			},
		},
	},
}
