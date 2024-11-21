package statistics

import (
	"tool/app/internal/helpers"
	"tool/app/internal/models"
)

func ComputeStatistics(workflow models.Workflow) models.Statistics {
	return models.Statistics{
		Workflow: computeWorkflowStatistics(workflow),
		Jobs:     computeJobsStatistics(workflow.Jobs),
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

func computeJobsStatistics(jobs map[string]models.Job) models.JobsStatistics {
	var permissions []models.PermissionsStatistics
	var blocked []models.IntStatistics
	var conditionals []models.IntStatistics
	var customRunners []models.IntStatistics
	var localEnvs []models.IntStatistics
	var environments []models.EnvironmentStatistics
	var defaults []models.IntStatistics
	var customContainers []models.IntStatistics
	var services []models.IntStatistics
	var customWorkflows []models.IntStatistics
	var secrets []models.EnvironmentStatistics

	count := 0

	for _, job := range jobs {
		permissions = append(permissions, permissionsCount(job.Permissions))
		blocked = append(blocked, models.IntStatistics{Total: len(job.Needs)})
		conditionals = append(conditionals, models.IntStatistics{Total: helpers.CheckPresence(job.If)})
		customRunners = append(customRunners, models.IntStatistics{Total: helpers.CheckPresence(job.RunsOn)})
		localEnvs = append(localEnvs, models.IntStatistics{Total: helpers.CheckPresence(job.Environment)})
		environments = append(environments, environmentCount(job.Env))
		defaults = append(defaults, models.IntStatistics{Total: len(job.Defaults)})
		customContainers = append(customContainers, models.IntStatistics{Total: helpers.CheckPresence(job.Container)})
		services = append(services, models.IntStatistics{Total: len(job.Services)})
		customWorkflows = append(customWorkflows, models.IntStatistics{Total: helpers.CheckPresence(job.Uses)})
		secrets = append(secrets, environmentCount(job.Secrets))
		count++
	}

	return models.JobsStatistics{
		Permissions:       permissionsArrayCount(permissions),
		Blocked:           intStatisticsArrayCount(blocked),
		Conditionals:      intStatisticsArrayCount(conditionals),
		CustomRunners:     intStatisticsArrayCount(customRunners),
		LocalEnvironments: intStatisticsArrayCount(localEnvs),
		Environments:      environmentArrayCount(environments),
		Defaults:          intStatisticsArrayCount(defaults),
		CustomContainers:  intStatisticsArrayCount(customContainers),
		Services:          intStatisticsArrayCount(services),
		CustomWorkflows:   intStatisticsArrayCount(customWorkflows),
		Secrets:           environmentArrayCount(secrets),
		Count:             count,
	}
}
