package local

import "encoding/xml"

type VariableNames struct {
	XMLName  xml.Name `xml:"Variables" json:"-"`
	Variable []string `json:"variables"`
}

type Variables struct {
	XMLName  xml.Name   `xml:"Variables" json:"-"`
	Variable []Variable `json:"variables"`
}

type Variable struct {
	XMLName     xml.Name `xml:"Variable" json:"-"`
	Name        string   `json:"name"`
	Value       string   `json:"value"`
	Units       string   `json:"units"`
	Description string   `json:"description"`
}

func NewVariables(variables ...Variable) Variables {
	return Variables{Variable: variables}
}

func NewVariable(variable string) Variable {
	return Variable{Name: variable}
}
