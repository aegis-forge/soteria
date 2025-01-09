package models

type Flags struct {
	Check Check
	Stats Stats
}

type Check struct {
	Config  string
	Verbose bool
	MaxRows int
	Output  string
}

type Stats struct {
	MaxRows int
	Output  string
}
