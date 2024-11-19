package models

type Workflow struct {
	Name        string                 `mapstructure:"name"`
	RunName     string                 `mapstructure:"run_name"`
	On          interface{}            `mapstructure:"on"`
	Permissions interface{}            `mapstructure:"permissions"`
	Env         interface{}            `mapstructure:"env"`
	Defaults    map[string]interface{} `mapstructure:"defaults"`
	Concurrency map[string]interface{} `mapstructure:"concurrency"`
	Jobs        map[string]Job         `mapstructure:"jobs"`
}
