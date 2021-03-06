package generator

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

	"github.com/floatzeI/go-web-autogen/decoratorparser"
	"github.com/floatzeI/go-web-autogen/runtimelib"

	"golang.org/x/mod/modfile"
)

var fileHeader = `// Package autogen_web
// This file was automatically generated by go-web-autogen
// It should not be edited manually
`

type GenerateOptions struct {
	Directory         string
	ControllersFolder string
	// Give the packageName and importing file's path, this function should return the import path of the package.
	ResolveModelFolder func(string, string) string
}

func Generate(options GenerateOptions) {
	// types map
	var typeToMethodMap = runtimelib.GetTypeToMethodMap()

	var currentTime = time.Now()
	fileHeader += "// Generated At: " + currentTime.Format(time.RFC3339) + "\n"
	var functionNameRegex = regexp.MustCompile("[^a-zA-Z0-9]")
	var numberRegex = regexp.MustCompile("[0-9]")
	var packageNameRegex = regexp.MustCompile("package ([a-zA-Z_]+)")
	var controllerMethods = make([]string, 0)
	var controllerMethodDefinitions = make([]string, 0)

	// targets
	target := options.Directory
	modFileLocation := target + "go.mod"
	modFileData, err := os.ReadFile(modFileLocation)
	if err != nil {
		panic(err)
	}
	mod, modErr := modfile.Parse(string(modFileData), modFileData, nil)
	if modErr != nil {
		panic(modErr)
	}
	moduleName := mod.Module.Mod.Path
	controllersFolder := target + "controllers"
	// controllers
	controllers, e := os.ReadDir(controllersFolder)
	if e != nil {
		panic(e)
	}
	var requiredPackages = make([]decoratorparser.RequiredPackageEntry, 0)
	var dependencies = runtimelib.NewDependencyResolver()
	requiredPackages = append(requiredPackages, decoratorparser.RequiredPackageEntry{
		Name: "services",
		Path: moduleName + "/services",
	})
	dependencies.AddDependency("Users", "Users", "singleton", "services", "services")
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
		var apiControllerConstructor *ast.FuncDecl = nil
		var apiController *decoratorparser.FunctionEntry = nil
		var foundConstructor = false
		for _, f := range node.Decls {
			var hasAccessToContextObject = false

			fn, ok := f.(*ast.FuncDecl)
			if !ok {
				gen, ok := f.(*ast.GenDecl)
				if ok {
					for _, spec := range gen.Specs {
						switch typeSpec := spec.(type) {
						case *ast.TypeSpec:
							// todo this might cause issues. see: https://github.com/golang/go/issues/27477
							if gen.Doc == nil {
								continue
							}
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
							if len(apiController.Arguments) > 0 {
								baseUrl = apiController.Arguments[0]
								fmt.Println("[info] controller baseUrl =", baseUrl)
								if !strings.HasSuffix(baseUrl, "/") {
									baseUrl = baseUrl + "/"
								}
							}

						}
					}
				}
				continue
			}
			if !foundConstructor && fn.Name.Name == "New"+apiControllerName {
				foundConstructor = true
				fmt.Println("[info] Found a constructor for", apiControllerName, "- Name =", fn.Name.Name)
				apiControllerConstructor = fn
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
				requiredParams = append(requiredParams, item.Arguments...)
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
			autogenComments = append(autogenComments, "@Produce json") // todo: only do this if return type is a struct/slice of structs
			autogenComments = append(autogenComments, stringHttpMethod.Comment)
			for _, resp := range possibleResponses {
				// todo: should ">= 300 && <= 399" be included (aka redirects)?
				if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
					hasSuccess = true
				}
				// todo: handle this - some functions may return 204 (or 200 with no body)
				if resp.Model == "" {
					panic("Empty return type is not implemented yet")
				}
				asString := decoratorparser.StringifyResponse(resp)
				autogenComments = append(autogenComments, asString)
			}
			// If no success code is explicitly set, default to the function return type
			if !hasSuccess {
				returns := ParseReturnType(fn)

				var respToAdd = decoratorparser.StringifyResponse(decoratorparser.ResponseEntry{
					StatusCode: 200,
					Model:      returns.Prefix + "" + returns.ImportedPackage + "." + returns.Type,
				})
				autogenComments = append(autogenComments, respToAdd)
			}
			var toJsonDecode = make([]decoratorparser.JsonDecodeEntry, 0)
			for _, param := range fn.Type.Params.List {
				var paramName = param.Names[0].Name
				var paramType = param.Type
				// check if exists in url
				var inUrl = strings.Contains(httpMethod.Url, "{"+paramName+"}")
				typeAsString := fmt.Sprintf("%v", paramType)
				var method = typeToMethodMap[typeAsString]

				startExp, isStartExpr := paramType.(*ast.StarExpr)
				var paramTypeAsString = "unknown"
				if isStartExpr {
					paramTypeAsString = fmt.Sprintf("%v", startExp.X)
					if isStartExpr {
						parsed := decoratorparser.ParseReturnType(paramTypeAsString)
						// todo: DI here (maybe?)
						if parsed.ImportedPackage == "fiber" && parsed.Type == "Ctx" {
							parametersList = append(parametersList, "c")
							hasAccessToContextObject = true
							continue
						}
					}
				}
				if method == "" {
					//try alternative method
					parsed := decoratorparser.ParseReturnType(typeAsString)
					if parsed.ImportedPackage != "" && parsed.Type != "" {
						fullName := parsed.ImportedPackage + "." + parsed.Prefix + parsed.Type
						// todo: convert ImportedPackage to a file path
						modelPackagePath := options.ResolveModelFolder(parsed.ImportedPackage, targetFile)
						requiredPackages = append(requiredPackages, decoratorparser.RequiredPackageEntry{
							Name: parsed.ImportedPackage,
							Path: moduleName + "/" + modelPackagePath,
						})
						var decParamName = paramName + `_decoded`
						toJsonDecode = append(toJsonDecode, decoratorparser.JsonDecodeEntry{
							Type:             fullName,
							DecodedParamName: decParamName,
						})
						autogenComments = append(autogenComments, `@Param `+paramName+` body `+fullName+` true "The `+paramName+`"`)
						parametersList = append(parametersList, decParamName)
						continue
					}
					panic(errors.New("the type \"" + paramTypeAsString + "\" (addr=" + typeAsString + ") for param " + paramName + " in " + targetFile + " is either not supported or not injected"))
				}
				var required = "false"
				for _, item := range requiredParams {
					if item == paramName {
						required = "true"
					}
				}
				var pType = "query"
				if inUrl {
					pType = "path"
				}
				parametersList = append(parametersList, `NewArgumentParser(c, "`+pType+`", "`+paramName+`", `+required+`).`+method+`()`)
				autogenComments = append(autogenComments, `@Param `+paramName+` `+pType+` `+typeAsString+` `+required+` "The `+paramName+`"`)
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
			var methodCall = fn.Name.Name + `(` + strings.Join(parametersList, ", ") + `)`
			if apiController != nil {
				var constructorArgs = make([]string, 0)
				for _, param := range apiControllerConstructor.Type.Params.List {
					var paramName = param.Names[0].Name
					var paramType = param.Type
					var parsedType = decoratorparser.ParseReturnType(fmt.Sprintf("%v", paramType))
					if strings.HasPrefix(parsedType.Type, "0x") {
						// convert "services.Users" to "users"
						fullStr := strings.Split(rawFileStr[param.Type.Pos():param.Type.End()-1], ".")
						if len(fullStr) == 1 {
							// is not in a different package
							parsedType = decoratorparser.ParseReturnType(fmt.Sprintf("%v", fullStr[0]))
						} else {
							// is packaged
							parsedType = decoratorparser.ParseReturnType(fmt.Sprintf("%v", fullStr[1]))
						}
						// todo: this might cause bugs in the unlikely scenerio that someone creates an interface starting with "0x"...
						// parsedType = decoratorparser.ParseReturnType(fmt.Sprintf("%v", paramName))
					}
					var dep = dependencies.CreateDependencyResolverFunction(parsedType.Type)
					constructorArgs = append(constructorArgs, dep.FunctionName+"()")
					fmt.Println("n", paramName, paramType, parsedType)
				}
				methodCall = "New" + apiControllerName + "(" + strings.Join(constructorArgs, ", ") + ")." + methodCall
			}
			methodCall = packageName + "." + methodCall
			var beforeMethodCall = ""
			var afterMethodCall = ""
			if hasAccessToContextObject {
				afterMethodCall = `
		if len(c.Response().Body()) != 0 {
			return nil
		}`
			}
			if len(toJsonDecode) != 0 {
				for _, item := range toJsonDecode {
					beforeMethodCall += `
		var ` + item.DecodedParamName + " " + item.Type + `
		if err := c.BodyParser(&` + item.DecodedParamName + `); err != nil {
            return err
        }
`
				}
			}
			var template = strings.Join(autogenComments, "\n") + `
func ` + functionName + `(app *fiber.App) {
	app.` + httpMethod.HttpMethod + `("` + httpMethod.Url + `", func(c *fiber.Ctx) error {
		` + beforeMethodCall + `
		var result = ` + methodCall + `
		` + afterMethodCall + `
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
	defStr = append(defStr, controllerMethodDefinitions...)
	// templates
	autogenFile = strings.ReplaceAll(autogenFile, "	// [StartupRegistrations]", strings.Join(regStr, "\n"))
	// add deps
	for _, item := range dependencies.GetFunctions() {
		var createServiceCall = item.RealTypePackageName + `.New` + item.RealType + `()`
		var fullIType = item.InterfacePackageName + "." + item.InterfaceType
		var comment = `// ` + item.FunctionName + ` resolves the ` + item.Mode + ` ` + item.InterfaceType + ` dependency`
		if item.Mode == "scoped" {
			defStr = append(defStr, `
`+comment+`
func `+item.FunctionName+`() `+fullIType+` {
	return `+createServiceCall+`
}`)
		} else {
			var singletonVarName = `singleton` + item.RealType
			defStr = append(defStr, `
var `+singletonVarName+` *`+fullIType+` = `+createServiceCall+`
`+comment+`
func `+item.FunctionName+`() *`+fullIType+` {
	return `+singletonVarName+`
}`)
		}
	}
	autogenFile = strings.ReplaceAll(autogenFile, "// [FunctionRegistrations]", strings.Join(defStr, "\n\n"))
	requiredPackages = append(requiredPackages, decoratorparser.RequiredPackageEntry{
		Name: "controllers",
		Path: moduleName + "/controllers",
	})
	// todo: sort packages alphabetically
	var packagesToList = make([]string, 0)
	var added = make(map[string]bool)
	for _, item := range requiredPackages {
		if _, exists := added[item.Path]; exists {
			continue
		}
		packagesToList = append(packagesToList, "    "+item.Name+` "`+item.Path+`"`)
	}
	autogenFile = strings.ReplaceAll(autogenFile, "	// [ImportRegistrations]", strings.Join(packagesToList, "\n"))
	writeTestErr := os.WriteFile(target+"autogen/autogen_web.go", []byte(autogenFile), 0666)
	if writeTestErr != nil {
		log.Fatal(writeTestErr)
	}
	// add runtime code
	runtimeLibBytes, runtimeBytesErr := os.ReadFile("./templates/autogen_web/autogen_web_lib.go")
	if runtimeBytesErr != nil {
		log.Fatal(runtimeBytesErr)
	}
	writeFileErr := os.WriteFile(target+"autogen/autogen_web_lib.go", runtimeLibBytes, 0666)
	if writeFileErr != nil {
		log.Fatal(writeFileErr)
	}
	// we're done!
	fmt.Println("[info] finished")
}
