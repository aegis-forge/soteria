package insufficient_pbac

import "tool/app/internal/detector"

var regexSecret = `\$\{\{\s*secrets\.|\$\{\{\s*github.token`

var GlobalSecret = detector.Detector{
	Name: "global-secret",
	Info: detector.Info{
		Description: "When declaring a secret, always declare it locally (step/container scope) and not globally (workflow/job scope).",
		Message:     "Secrets should only be defined in steps or containers",
		Severity:    3,
		CICDSEC:     5,
	},
	CountAll: true,
	Rule: &detector.Or{
		LHS: &detector.Or{
			LHS: &detector.Match{
				LHS: "$.env[*]",
				RHS: regexSecret,
			},
			RHS: &detector.Equals{
				LHS: "$.env",
				RHS: "inherit",
			},
		},
		RHS: &detector.Or{
			LHS: &detector.Or{
				LHS: &detector.Match{
					LHS: "$.jobs..env[*]",
					RHS: regexSecret,
				},
				RHS: &detector.Equals{
					LHS: "$.jobs..env",
					RHS: "inherit",
				},
			},
			RHS: &detector.Or{
				LHS: &detector.Match{
					LHS: "$.jobs..secrets[*]",
					RHS: regexSecret,
				},
				RHS: &detector.Equals{
					LHS: "$.jobs..secrets",
					RHS: "inherit",
				},
			},
		},
	},
}
