package detectors

import "tool/app/internal/detector"

var NoHashVersionPin = detector.Detector{
	Name: "no-hash-version-pin",
	Info: detector.Info{
		Description: "Always use the full SHA when referring to the version of an external Github Action (especially third-party ones).",
		Severity:    2,
		CWE:         -1,
	},
	Rule: &detector.Match{
		LHS: "$..uses",
		RHS: `^[^@]+@\S{1,39}$`,
	},
}
