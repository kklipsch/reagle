package local

import (
	"encoding/xml"
)

func NewCommand(name string) Command {
	return Command{Name: name}
}

//Command is the request body structure
type Command struct {
	XMLName xml.Name `xml:"Command"`
	Name    string
}
