package models

type Flags struct {
	Check Check
}

type Check struct {
	Stats   bool
	Verbose bool
	Output  string
	MaxRows int
}
