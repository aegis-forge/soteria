package detectors

import "tool/app/internal/detector"

var UnsafeArtifactDownload = detector.Detector{
	Name: "unsafe-artifact-download",
	Info: detector.Info{
		Description: "Downloading artifacts without specifying the path and commit/run_id can lead to privacy escalation in the pipeline.",
		Severity:    5,
		CWE:         73,
	},
	Rule: &detector.And{
		LHS: &detector.Match{
			LHS: "$.jobs..steps[*].uses",
			RHS: "dawidd6/action-download-artifact",
		},
		RHS: &detector.Or{
			LHS: &detector.Exists{
				NOT: true,
				LHS: "$.jobs..steps[*].with[*]~",
				RHS: "path",
			},
			RHS: &detector.And{
				LHS: &detector.Exists{
					NOT: true,
					LHS: "$.jobs..steps[*].with[*]~",
					RHS: "commit",
				},
				RHS: &detector.Exists{
					NOT: true,
					LHS: "$.jobs..steps[*].with[*]~",
					RHS: "run_id",
				},
			},
		},
	},
}
