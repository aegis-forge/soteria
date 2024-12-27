package detectors

import "tool/app/internal/detector"

var SelfHostedRunner = detector.Detector{
	Name: "self-hosted-runner",
	Info: detector.Info{
		Description: "When using self-hosted runners (especially in public repos), a user could fork the repo and send malicious pull requests to try and escape the sandbox.",
		Severity:    3,
		CWE:         -1,
	},
	Rule: &detector.Or{
		LHS: &detector.Match{
			LHS: "$.jobs..runs-on",
			RHS: "self-hosted",
		},
		RHS: &detector.Match{
			LHS: "$.jobs..runs-on[*]",
			RHS: "self-hosted",
		},
	},
}
