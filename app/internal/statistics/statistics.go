package statistics

import (
	"slices"
	"tool/app/internal/helpers"
	"tool/app/internal/models"
)

func AggregateStatistics(statistics []models.Statistics) models.GlobalStatistics {
	var jobs []int
	var steps []int
	var containers []int

	for _, stat := range statistics {
		jobs = append(jobs, stat.Jobs.Count.Total)
		steps = append(steps, stat.Steps.Count.Total)
		containers = append(containers, stat.Containers.Count.Total)
	}

	return models.GlobalStatistics{
		Jobs:       BuildIntStatistics(jobs),
		Steps:      BuildIntStatistics(steps),
		Containers: BuildIntStatistics(containers),
	}
}

func ComputeStatistics(workflow models.Workflow) models.Statistics {
	var steps []models.Step
	var containers []models.Container

	for _, job := range workflow.Jobs {
		steps = slices.Concat(steps, job.Steps)

		if helpers.CheckPresence(job.Container) == 1 {
			containers = append(containers, job.Container)
		}
	}

	return models.Statistics{
		Workflow:   computeWorkflowStatistics(workflow),
		Jobs:       computeJobsStatistics(workflow.Jobs),
		Steps:      computeStepsStatistics(steps),
		Containers: computeContainersStatistics(containers),
	}
}

func computeWorkflowStatistics(workflow models.Workflow) models.WorkflowStatistics {
	return models.WorkflowStatistics{
		Events:      eventsCount(workflow.On),
		Permissions: permissionsCount(workflow.Permissions),
		Environment: environmentCount(workflow.Env),
		Jobs:        models.IntStatistics{Total: len(workflow.Jobs)},
		Defaults:    defaultsCount(workflow.Defaults),
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
	var steps []models.IntStatistics

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
		steps = append(steps, models.IntStatistics{Total: len(job.Steps)})
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
		Steps:             intStatisticsArrayCount(steps),
		Count:             models.IntStatistics{Total: count},
	}
}

func computeStepsStatistics(steps []models.Step) models.StepsStatistics {
	var conditionals []models.IntStatistics
	var customActions []models.IntStatistics
	var runScripts []models.IntStatistics
	var environments []models.EnvironmentStatistics

	count := 0

	for _, step := range steps {
		conditionals = append(conditionals, models.IntStatistics{Total: helpers.CheckPresence(step.If)})
		customActions = append(customActions, models.IntStatistics{Total: helpers.CheckPresence(step.Uses)})
		runScripts = append(runScripts, models.IntStatistics{Total: helpers.CheckPresence(step.Run)})
		environments = append(environments, environmentCount(step.Env))
		count++
	}

	return models.StepsStatistics{
		Conditionals:  intStatisticsArrayCount(conditionals),
		CustomActions: intStatisticsArrayCount(customActions),
		RunScripts:    intStatisticsArrayCount(runScripts),
		Environments:  environmentArrayCount(environments),
		Count:         models.IntStatistics{Total: count},
	}
}

func computeContainersStatistics(containers []models.Container) models.ContainersStatistics {
	var credentials []models.EnvironmentStatistics
	var environments []models.EnvironmentStatistics

	count := 0

	for _, container := range containers {
		credentials = append(credentials, environmentCount(container.Credentials))
		environments = append(environments, environmentCount(container.Env))

		count++
	}

	return models.ContainersStatistics{
		Credentials:  environmentArrayCount(credentials),
		Environments: environmentArrayCount(environments),
		Count:        models.IntStatistics{Total: count},
	}
}
