package model

import (
	"fmt"
)

type SensorDefinition struct {
	Code		string
	Name		string
	Unit		string
}

func (s SensorDefinition) NameAndUnit() string {
	unit := ""
	if len(s.Unit) > 0 {
		unit = fmt.Sprintf(" (%s)", s.Unit)
	}
	return fmt.Sprintf("%s%s", s.Name, unit)
}

func (m *Model) GetSensorDefinitions() []SensorDefinition {
	return []SensorDefinition{
		{Code: "so2", 					Name: "SO2", 					Unit: ""}, 
		{Code: "nox", 					Name: "NOX", 					Unit: ""}, 
		{Code: "pm", 						Name: "PM", 					Unit: ""}, 
		{Code: "h2s", 					Name: "H2S", 					Unit: ""}, 
		{Code: "opacity", 			Name: "Opacity", 			Unit: ""}, 
		{Code: "flow", 					Name: "Flow", 				Unit: ""}, 
		{Code: "o2", 						Name: "O2", 					Unit: ""}, 
		{Code: "temperature",		Name: "Temperature", 	Unit: ""}, 
		{Code: "pressure", 			Name: "Pressure", 		Unit: ""}, 
	}
}
