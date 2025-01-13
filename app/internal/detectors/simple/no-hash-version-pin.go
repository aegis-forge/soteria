package simple

import "tool/app/internal/detector"

var NoHashVersionPin = detector.Detector{
	Name: "no-hash-version-pin",
	Info: detector.Info{
		Description: "Always use the full hash when referring to the version of an external Github Action (especially third-party ones).",
		Message:     "Full commit hash should be used when referring to external GitHub Actions",
		Severity:    2,
		CWE:         -1,
	},
	CountAll: true,
	Rule: &detector.Match{
		LHS: "$..uses",
		RHS: `^[^@]+@\S{1,39}$`,
	},
}
