package runtimelib

func GetTypeToMethodMap() map[string]string {
	var typeToMethodMap = make(map[string]string)
	typeToMethodMap["int"] = "GetInt"
	typeToMethodMap["int64"] = "GetInt64"
	typeToMethodMap["uint"] = "GetUint"
	typeToMethodMap["uint64"] = "GetUint64"
	typeToMethodMap["string"] = "GetString"
	typeToMethodMap["float64"] = "GetFloat64"
	typeToMethodMap["float32"] = "GetFloat32"
	typeToMethodMap["bool"] = "GetBool"
	return typeToMethodMap
}
