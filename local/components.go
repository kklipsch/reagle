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
	XMLName   xml.Name        `xml:"Components" json:"-"`
	Component []ComponentName `json:"components"`
}

type ComponentName struct {
	XMLName    xml.Name      `xml:"Component" json:"-"`
	Name       string        `json:"name"`
	HardwareID string        `xml:"HarwareId" json:"hardware_id"`
	FixedID    int           `xml:"FixedId" json:"fixed_id"`
	Variables  VariableNames `json:"variables"`
}

type Components struct {
	XMLName   xml.Name    `xml:"Components" json:"-"`
	Component []Component `json:"components"`
}

type Component struct {
	XMLName    xml.Name  `xml:"Component" json:"-"`
	Name       string    `json:"name"`
	HardwareID string    `xml:"HarwareId" json:"hardware_id"`
	FixedID    int       `xml:"FixedId" json:"fixed_id"`
	Variables  Variables `json:"variables"`
}
