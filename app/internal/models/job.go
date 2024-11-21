package models

type Job struct {
	Name            string                 `mapstructure:"name"`
	Permissions     interface{}            `mapstructure:"permissions"`
	Needs           []string               `mapstructure:"needs"`
	If              string                 `mapstructure:"if"`
	RunsOn          interface{}            `mapstructure:"runs-on"`
	Environment     map[string]interface{} `mapstructure:"environment"`
	Concurrency     map[string]interface{} `mapstructure:"concurrency"`
	Output          map[string]interface{} `mapstructure:"output"`
	Env             interface{}            `mapstructure:"env"`
	Defaults        map[string]interface{} `mapstructure:"defaults"`
	Steps           []Step                 `mapstructure:"steps"`
	TimeoutMinutes  int                    `mapstructure:"timeout-minutes"`
	Strategy        map[string]interface{} `mapstructure:"strategy"`
	ContinueOnError bool                   `mapstructure:"continue-on-error"`
	Container       Container              `mapstructure:"container"`
	Services        map[string]Container   `mapstructure:"services"`
	Uses            string                 `mapstructure:"uses"`
	With            map[string]interface{} `mapstructure:"with"`
	Secrets         interface{}            `mapstructure:"secrets"`
}

type Container struct {
	Image       string                 `mapstructure:"image"`
	Credentials map[string]interface{} `mapstructure:"credentials"`
	Env         map[string]interface{} `mapstructure:"env"`
	Ports       []interface{}          `mapstructure:"ports"`
	Volumes     []string               `mapstructure:"volumes"`
	Options     []string               `mapstructure:"options"`
}
