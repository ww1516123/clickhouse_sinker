package util

import (
	"github.com/housepower/clickhouse_sinker/model"
)

//这里对metric的value类型，只有三种情况， （float64，string，map[string]interface{})
func GetValueByType(metric model.Metric, cwt *model.ColumnWithType) interface{} {
	swType := switchType(cwt.Type)
	switch swType {
	case "int":
		return metric.GetInt(cwt.Name)
	case "float":
		return metric.GetFloat(cwt.Name)
	case "string":
		return metric.GetString(cwt.Name)
	case "stringArray":
		return metric.GetArray(cwt.Name, "string")
	case "intArray":
		return metric.GetArray(cwt.Name, "int")
	case "floatArray":
		return metric.GetArray(cwt.Name, "float")
	//never happen
	default:
		return ""
	}
}

func switchType(typ string) string {
	switch typ {
	case "Date", "DateTime", "UInt8", "UInt16", "UInt32", "UInt64", "Int8", "Int16", "Int32", "Int64":
		return "int"
	case "Array(Date)", "Array(DateTime)", "Array(UInt8)", "Array(UInt16)", "Array(UInt32)", "Array(UInt64)", "Array(Int8)", "Array(Int16)", "Array(Int32)", "Array(Int64)":
		return "intArray"
	case "String", "FixString":
		return "string"
	case "Array(String)", "Array(FixString)":
		return "stringArray"
	case "Float32", "Float64":
		return "float"
	case "Array(Float32)", "Array(Float64)":
		return "floatArray"
	default:
		panic("unsupport type " + typ)
	}
}
