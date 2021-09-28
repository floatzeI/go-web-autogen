package decoratorparser

import "strings"

type ReturnType struct {
	ImportedPackage string
	Type string
	Prefix string
}

func ParseReturnType(retType string) ReturnType {
	var res = ReturnType{}
	if retType[0:1] == "&" || retType[0:1] == "{" {
		if retType[0:1] == "&" {
			retType = retType[1:]
		}
		// cut off the beginning and ending "{"
		retType = retType[1:len(retType)-1]
		// 0 = Package, 1 = Return Struct
		split := strings.Split(retType, " ")
		res.ImportedPackage = split[0]
		res.Type = split[1]
	}else{
		// assume it is a builint (e.g. "int" or "string" or "int64")
		res.Type = retType
	}

	return res
}
