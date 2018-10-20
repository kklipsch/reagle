package local

import "encoding/xml"

func NewComponents(components ...Component) Components {
	return Components{Component: components}
}

func NewComponent(name string, variable ...string) Component {
	variables := []Variable{}
	for _, v := range variable {
		variables = append(variables, NewVariable(v))
	}

	return Component{
		Name:      name,
		Variables: NewVariables(variables...),
	}
}

type ComponentNames struct {
	XMLName   xml.Name `xml:"Components"`
	Component []ComponentName
}

type ComponentName struct {
	XMLName    xml.Name `xml:"Component"`
	Name       string
	HardwareID string `xml:"HarwareId"`
	FixedID    int    `xml:"FixedId"`
	Variables  VariableNames
}

type Components struct {
	XMLName   xml.Name `xml:"Components"`
	Component []Component
}

type Component struct {
	XMLName    xml.Name `xml:"Component"`
	Name       string
	HardwareID string `xml:"HarwareId"`
	FixedID    int    `xml:"FixedId"`
	Variables  Variables
}
