package models

import "io"

type ITable interface {
	String() string
	StreamCSV(writer io.Writer)
}

func ConvertToITables(tablesIn []*Table) []ITable {
	tablesOut := make([]ITable, len(tablesIn))
	for i, t := range tablesIn {
		tablesOut[i] = t
	}
	return tablesOut
}
