package generator

import (
	"fmt"
	"go/ast"
	"web-autogen/decoratorparser"
)

// ParseReturnType returns the return type of the function
func ParseReturnType(fn *ast.FuncDecl) decoratorparser.ReturnType {
	var returns decoratorparser.ReturnType // http return type
	var returnType = fn.Type.Results       // go type

	if returnType == nil {

	} else {
		// return type for basic type: "int"
		// return type for struct: "&{usersmodels GetUserById}"
		// return type for struct slice: ""
		switch xv := returnType.List[0].Type.(type) {
		case *ast.ArrayType:
			{
				var arrayType = xv.Elt
				var isSlice = xv.Len == nil
				// get the type
				var elementResult = decoratorparser.ParseReturnType(fmt.Sprintf("%v", arrayType))
				if isSlice {
					elementResult.Prefix = "[]"
					returns = elementResult
				} else {
					// todo
					panic("Array support is not yet implemented")
				}
				break
			}
		default:
			var str = fmt.Sprintf("%v", returnType.List[0].Type)
			returns = decoratorparser.ParseReturnType(str)
		}
	}
	return returns
}
