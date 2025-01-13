package simple

import (
	"tool/app/internal/detector"
)

var regexBad = `\$\{{2}\s*(env\.[A-z]+)\s*}{2}`

var BadLocalEnvironment = detector.Detector{
	Name: "bad-local-environment",
	Info: detector.Info{
		Description: "Using environment variables coming from external files can lead to secret disclosure through code injection.",
		Message:     "Do not use local environment variables (i.e. env.*) in scripts",
		Severity:    4,
		CWE:         94,
	},
	CountAll: true,
	Rule: &detector.Or{
		LHS: &detector.Match{
			LHS: "$.jobs..steps[*].with.script",
			RHS: regexBad,
		},
		RHS: &detector.Match{
			LHS: "$.jobs..steps[*].run",
			RHS: regexBad,
		},
	},
}
