package models

type Flags struct {
	Check   Check
	MaxRows int
}

type Check struct {
	Stats   bool
	Verbose bool
	Output  string
}
