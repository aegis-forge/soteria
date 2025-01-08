package helpers

var SeverityMap = map[int]string{
	0: "Info",
	1: "Warning",
	2: "Low",
	3: "Medium",
	4: "High",
	5: "Critical",
}

var ColorMap = map[int]string{
	0: "\033[34;1m",
	1: "\033[37;1m",
	2: "\033[34;1m",
	3: "\033[33;1m",
	4: "\033[31;1m",
	5: "\033[35;1m",
}
