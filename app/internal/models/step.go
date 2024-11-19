package models

type Step struct {
	If               string                 `mapstructure:"if"`
	Name             string                 `mapstructure:"name"`
	Uses             string                 `mapstructure:"uses"`
	Run              string                 `mapstructure:"run"`
	WorkingDirectory string                 `mapstructure:"working-directory"`
	Shell            string                 `mapstructure:"shell"`
	With             map[string]interface{} `mapstructure:"with"`
	Env              interface{}            `mapstructure:"env"`
	ContinueOnError  bool                   `mapstructure:"continue-on-error"`
	TimeoutMinutes   int                    `mapstructure:"timeout-minutes"`
}
