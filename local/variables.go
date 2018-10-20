package local

import "encoding/xml"

type Variables struct {
	XMLName  xml.Name `xml:"Variables"`
	Variable []Variable
}

type Variable struct {
	XMLName     xml.Name `xml:"Variable"`
	Name        string
	Value       string
	Units       string
	Description string
}

func NewVariables(variables ...Variable) Variables {
	return Variables{Variable: variables}
}

func NewVariable(variable string) Variable {
	return Variable{Name: variable}
}
