package models

type Flags struct {
	Check Check
	Stats Stats
}

type Check struct {
	Config  string
	Repo    string
	String  bool
	Verbose bool
	MaxRows int
	Global  bool
	Output  string
}

type Stats struct {
	Repo    string
	String  bool
	MaxRows int
	Global  bool
	Output  string
}

type Detectors struct {
	Config string
}
