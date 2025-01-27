package models

type Flags struct {
	Check Check
	Stats Stats
}

type Check struct {
	Config  string
	Repo    string
	Verbose bool
	MaxRows int
	Global  bool
	Output  string
}

type Stats struct {
	Repo    string
	MaxRows int
	Global  bool
	Output  string
}
