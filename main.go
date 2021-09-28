package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
	"web-autogen/decoratorparser"
)

var fileHeader = `// Package autogen_web
// This file was automatically generated by go-web-autogen
// It should not be edited manually
`

func main() {
	// types map
	var typeToMethodMap = make(map[string]string)
	typeToMethodMap["int"] = "GetInt"
	typeToMethodMap["int64"] = "GetInt64"
	typeToMethodMap["uint"] = "GetUint"
	typeToMethodMap["uint64"] = "GetUint64"
	typeToMethodMap["string"] = "GetString"

	var currentTime = time.Now()
	fileHeader += "// Generated At: " + currentTime.Format(time.RFC3339) + "\n"
	var functionNameRegex = regexp.MustCompile("[^a-zA-Z0-9]")
	var numberRegex = regexp.MustCompile("[0-9]")
	var packageNameRegex = regexp.MustCompile("package ([a-zA-Z_]+)")
	var controllerMethods = make([]string, 0)
	var controllerMethodDefinitions = make([]string, 0)

	// targets
	target := "examples/project-one/"
	controllersFolder := target + "/controllers"
	// controllers
	controllers, e := os.ReadDir(controllersFolder)
	if e != nil {
		panic(e)
	}
	for _, entry := range controllers {
		targetFile := target + "controllers/" + entry.Name()

		var rawFile, readErr = os.ReadFile(targetFile)
		if readErr != nil {
			log.Fatal(readErr)
		}
		var rawFileStr = string(rawFile)
		var packageName = packageNameRegex.FindAllStringSubmatch(rawFileStr, -1)[0][1]
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, "", rawFile, parser.ParseComments)
		if err != nil {
			log.Fatal(err)
		}
		var baseUrl = ""
		var apiControllerName = ""
		var apiController *decoratorparser.FunctionEntry = nil
		var foundConstructor = false
		for _, f := range node.Decls {

			fn, ok := f.(*ast.FuncDecl)
			if !ok {
				gen, ok := f.(*ast.GenDecl)
				if ok {
					for _, spec := range gen.Specs {
						switch spec.(type) {
						case *ast.TypeSpec:
							typeSpec := spec.(*ast.TypeSpec)
							// todo this might cause issues. see: https://github.com/golang/go/issues/27477
							parsedComments := decoratorparser.New(gen.Doc.List[0].Text, typeSpec.Name.Name, "constructor")
							decoratorFunctions := parsedComments.GetFunctions()
							if apiController != nil {
								log.Fatal("Found 2+ controllers in ", targetFile, " : ", typeSpec.Name.Name, ", ", apiControllerName)
							}
							apiController = decoratorFunctions.GetFirstCallByName("Controller")
							if apiController == nil {
								fmt.Print("[info] Struct", typeSpec.Name.Name, "is not a controller - skipping")
								continue
							}
							apiControllerName = typeSpec.Name.Name
							baseUrl = apiController.Arguments[0]
							fmt.Println("[info] controller baseUrl =", baseUrl)
							if !strings.HasSuffix(baseUrl, "/") {
								baseUrl = baseUrl + "/"
							}

						}
					}
				}
				continue
			}
			if !foundConstructor && fn.Name.Name == "New"+apiControllerName {
				foundConstructor = true
				fmt.Println("[info] Found constructor. Name=" + fn.Name.Name)
				continue
			}
			// todo: add di
			var parametersList = make([]string, 0)
			var comments = fn.Doc.Text()
			var parser = decoratorparser.New(comments, entry.Name(), fn.Name.Name)
			var functionCalls = parser.GetFunctions()
			var _requiredCalls = functionCalls.GetCallsByName("Required")
			var requiredParams = make([]string, 0)
			for _, item := range _requiredCalls {
				for _, p := range item.Arguments {
					requiredParams = append(requiredParams, p)
				}
			}

			var httpMethod = parser.ParseMethod()
			if !httpMethod.IsControllerAction {
				continue
			}
			httpMethod.Url = baseUrl + httpMethod.Url // add controller base url
			var stringHttpMethod = decoratorparser.StringifyMethod(httpMethod)

			var possibleResponses = parser.ParseResponse()

			var hasSuccess = false

			var autogenComments = make([]string, 0)
			autogenComments = append(autogenComments, "@Produce json")
			autogenComments = append(autogenComments, stringHttpMethod.Comment) // add method
			for _, resp := range possibleResponses {
				if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
					hasSuccess = true
				}
				if resp.Model == "" {
					continue
				}
				asString := decoratorparser.StringifyResponse(resp)
				autogenComments = append(autogenComments, asString)
			}
			var returns decoratorparser.ReturnType
			if !hasSuccess {
				var returnType = fn.Type.Results
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

				var respToAdd = decoratorparser.StringifyResponse(decoratorparser.ResponseEntry{
					StatusCode: 200,
					Model:      returns.Prefix + "" + returns.ImportedPackage + "." + returns.Type,
				})
				autogenComments = append(autogenComments, respToAdd)
			}
			for _, param := range fn.Type.Params.List {
				var n = param.Names[0].Name
				var t = param.Type
				// check if exists in url
				var inUrl = strings.Contains(httpMethod.Url, "{"+n+"}")
				typeAsString := fmt.Sprintf("%v", t)
				var method = typeToMethodMap[typeAsString]
				var required = "false"
				for _, item := range requiredParams {
					if item == n {
						required = "true"
					}
				}
				var pType = "query"
				if inUrl {
					pType = "path"
				}
				parametersList = append(parametersList, `NewArgumentParser(c, "`+pType+`", "`+n+`", `+required+`).`+method+`()`)
				autogenComments = append(autogenComments, `@Param `+n+` `+pType+` `+typeAsString+` `+required+` "The `+n+`"`)
			}
			for i := range autogenComments {
				autogenComments[i] = "// " + autogenComments[i]
			}
			var functionName = functionNameRegex.ReplaceAllString(httpMethod.Url, "")
			functionName = httpMethod.HttpMethod + "_" + functionName
			for numberRegex.MatchString(functionName[0:1]) {
				functionName = functionName[1:]
			}
			functionName = strings.ToLower(functionName)
			if len(functionName) == 0 {
				panic(errors.New("bad function name"))
			}
			var methodCall = fn.Name.Name + `(` + strings.Join(parametersList, ",") + `)`
			if apiController != nil {
				methodCall = "New" + apiControllerName + "()." + methodCall
			}
			methodCall = packageName + "." + methodCall
			var template = strings.Join(autogenComments, "\n") + `
func ` + functionName + `(app *fiber.App) {
	app.` + httpMethod.HttpMethod + `("` + httpMethod.Url + `", func(c *fiber.Ctx) error {
		var result = ` + methodCall + `
		return c.JSON(result)
	})
}`
			controllerMethodDefinitions = append(controllerMethodDefinitions, template)
			controllerMethods = append(controllerMethods, functionName)
		}
	}

	var autogenFile = fileHeader
	var autogenBytes, autogenErr = os.ReadFile("./templates/autogen_web/autogen_web.go")
	if autogenErr != nil {
		log.Fatal(autogenErr)
	}
	autogenFile += string(autogenBytes)
	var regStr = make([]string, 0)
	for _, str := range controllerMethods {
		regStr = append(regStr, "    "+str+"(app)")
	}
	var defStr = make([]string, 0)
	for _, str := range controllerMethodDefinitions {
		defStr = append(defStr, str)
	}
	// templates
	autogenFile = strings.ReplaceAll(autogenFile, "	// [StartupRegistrations]", strings.Join(regStr, "\n"))
	autogenFile = strings.ReplaceAll(autogenFile, "// [FunctionRegistrations]", strings.Join(defStr, "\n\n"))
	autogenFile = strings.ReplaceAll(autogenFile, "	// [ImportRegistrations]", `	"project-one/controllers"`)
	writeTestErr := os.WriteFile(target+"/autogen/autogen_web.go", []byte(autogenFile), 0666)
	if writeTestErr != nil {
		log.Fatal(writeTestErr)
	}
	// add runtime code
	runtimeLibBytes, runtimeBytesErr := os.ReadFile("./templates/autogen_web/autogen_web_lib.go")
	if runtimeBytesErr != nil {
		log.Fatal(runtimeBytesErr)
	}
	writeFileErr := os.WriteFile(target+"/autogen/autogen_web_lib.go", runtimeLibBytes, 0666)
	if writeFileErr != nil {
		log.Fatal(writeFileErr)
	}
	// we're done!
}
