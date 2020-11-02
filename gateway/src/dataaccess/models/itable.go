package models

type ITable interface {
	String() string
}

func ConvertToITables(tablesIn []*Table) []ITable {
	tablesOut := make([]ITable, len(tablesIn))
	for i, t := range tablesIn {
		tablesOut[i] = t
	}
	return tablesOut
}
