package statistics

import (
	"tool/app/internal/models"
)

func ComputeStatistics(workflow models.Workflow) models.Statistics {
	return models.Statistics{
		Workflow: computeWorkflowStatistics(workflow),
	}
}

func computeWorkflowStatistics(workflow models.Workflow) models.WorkflowStatistics {
	return models.WorkflowStatistics{
		Events:      eventsCount(workflow.On),
		Defaults:    defaultsCount(workflow.Defaults),
		Permissions: permissionsCount(workflow.Permissions),
		Environment: environmentCount(workflow.Env),
	}
}
