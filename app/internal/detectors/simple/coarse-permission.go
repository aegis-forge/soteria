package simple

import "tool/app/internal/detector"

var CoarsePermission = detector.Detector{
	Name: "coarse-permission",
	Info: detector.Info{
		Description: "Permissions shouldn't be coarse, they should be finegrained on the specific permissions (no default/read-all/write-all permissions).",
		Message:     "No default (i.e. not defined), read-all, or write-all permissions should be used",
		Severity:    3,
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
			LHS: &detector.Or{
				LHS: &detector.Match{
					LHS: "$.jobs..permissions",
					RHS: "read-all",
				},
				RHS: &detector.Match{
					LHS: "$.jobs..permissions",
					RHS: "write-all",
				},
			},
			RHS: &detector.Exists{
				NOT: true,
				LHS: "$[*]~",
				RHS: "permissions",
			},
		},
	},
}
