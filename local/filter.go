package local

type VariableFilter map[string]bool

func (f VariableFilter) Exclude(variable string) bool {
	return f[variable]
}

var (
	//BadResponseVariables is a list of variables that return bad xml responses in some firmwares of the reagle
	BadResponseVariables = VariableFilter(map[string]bool{"zigbee:Multiplier": true, "zigbee:Divisor": true})

	//By Default don't filter
	DefaultFilter = VariableFilter(map[string]bool{})
)
