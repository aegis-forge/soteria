package insecure_system_configuration

import "tool/app/internal/detector"

var SelfHostedRunner = detector.Detector{
	Name: "self-hosted-runner",
	Info: detector.Info{
		Description: "When using self-hosted runners (especially in public repos), a user could fork the repo and send malicious pull requests to try and escape the sandbox.",
		Message:     "Self-hosted runners should never be used in public repos",
		Severity:    4,
		CICDSEC:     7,
	},
	CountAll: true,
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
