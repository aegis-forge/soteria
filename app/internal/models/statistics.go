package models

type GlobalStatistics struct {
	Jobs       IntStatistics `json:"jobs"`
	Steps      IntStatistics `json:"steps"`
	Containers IntStatistics `json:"containers"`
}

type Statistics struct {
	Workflow   WorkflowStatistics   `json:"workflow"`
	Jobs       JobsStatistics       `json:"jobs"`
	Containers ContainersStatistics `json:"containers"`
	Steps      StepsStatistics      `json:"steps"`
}

// Workflow Sections Statistics

type WorkflowStatistics struct {
	Events      IntStatistics         `json:"events"`
	Permissions PermissionsStatistics `json:"permissions"`
	Environment EnvironmentStatistics `json:"environment"`
	Jobs        IntStatistics         `json:"jobs"`
	Defaults    IntStatistics         `json:"defaults"`
}

type JobsStatistics struct {
	Permissions       PermissionsStatistics `json:"permissions"`
	Blocked           IntStatistics         `json:"blocked"`
	Conditionals      IntStatistics         `json:"conditionals"`
	CustomRunners     IntStatistics         `json:"custom-runners"`
	LocalEnvironments IntStatistics         `json:"local-environments"`
	Environments      EnvironmentStatistics `json:"environment"`
	Defaults          IntStatistics         `json:"defaults"`
	CustomContainers  IntStatistics         `json:"custom-containers"`
	Services          IntStatistics         `json:"services"`
	CustomWorkflows   IntStatistics         `json:"custom-workflows"`
	Secrets           EnvironmentStatistics `json:"secrets"`
	Steps             IntStatistics         `json:"steps"`
	Count             IntStatistics         `json:"count"`
}

type ContainersStatistics struct {
	Credentials  EnvironmentStatistics `json:"credentials"`
	Environments EnvironmentStatistics `json:"environment"`
	Count        IntStatistics         `json:"count"`
}

type StepsStatistics struct {
	Conditionals  IntStatistics         `json:"conditionals"`
	CustomActions IntStatistics         `json:"custom-commands"`
	RunScripts    IntStatistics         `json:"run-scripts"`
	Environments  EnvironmentStatistics `json:"environments"`
	Count         IntStatistics         `json:"count"`
}

// Permissions Statistics

type PermissionsStatistics struct {
	FineGrained   FineGrainedStatistics   `json:"fine-grained"`
	CoarseGrained CoarseGrainedStatistics `json:"coarse-grained"`
	Count         IntStatistics           `json:"count"`
}

type FineGrainedStatistics struct {
	Read  IntStatistics `json:"read"`
	Write IntStatistics `json:"write"`
	None  IntStatistics `json:"none"`
}

type CoarseGrainedStatistics struct {
	ReadAll  IntStatistics `json:"read-all"`
	WriteAll IntStatistics `json:"write-all"`
	NoneAll  IntStatistics `json:"none-all"`
	Default  IntStatistics `json:"default"`
}

// Environment Statistics

type EnvironmentStatistics struct {
	Inherited IntStatistics `json:"inherited"`
	Hardcoded IntStatistics `json:"hardcoded"`
	Variables IntStatistics `json:"variables"`
	Count     IntStatistics `json:"count"`
}

// Integer Statistics

type IntStatistics struct {
	Total  int     `json:"total"`
	Min    int     `json:"min"`
	Max    int     `json:"max"`
	Mean   int     `json:"mean"`
	Median int     `json:"median"`
	Std    float64 `json:"std"`
}
