package runtimelib

import "errors"

type DependencyResolver struct {
	entries         []dependencyResolverEntry
	interfaceToType []interfaceToRealTypeEntry
}

type dependencyResolverEntry struct {
	// The function to call
	FunctionName string
	// The interface type
	InterfaceType string
	// Actual type of the struct
	RealType             string
	InterfacePackageName string
	RealTypePackageName  string
	Mode                 string
}

type interfaceToRealTypeEntry struct {
	// "singleton", "scoped" - Singleton is created once then shared with all requests, scoped is created once per request
	Mode string
	// The name of the interface
	InterfaceType string
	// The name of the struct
	RealType             string
	InterfacePackageName string
	RealTypePackageName  string
}

func NewDependencyResolver() *DependencyResolver {
	return &DependencyResolver{}
}

func (d *DependencyResolver) AddDependency(interfaceType string, realType string, mode string, interfacePackage string, realPackage string) {
	if mode != "scoped" && mode != "singleton" {
		panic("invalid mode: " + mode)
	}
	d.interfaceToType = append(d.interfaceToType, interfaceToRealTypeEntry{
		Mode:                 mode,
		InterfaceType:        interfaceType,
		RealType:             realType,
		InterfacePackageName: interfacePackage,
		RealTypePackageName:  realPackage,
	})
}

func (d *DependencyResolver) resolveDependency(interfaceType string) interfaceToRealTypeEntry {
	for _, item := range d.interfaceToType {
		if item.InterfaceType == interfaceType {
			return item
		}
	}
	panic(errors.New("cannot resolve dependency \"" + interfaceType + "\" - does it exist?"))
}

// Create a resolver function (or get an existing one), and return it
func (d *DependencyResolver) CreateDependencyResolverFunction(interfaceType string) dependencyResolverEntry {
	resolved := d.resolveDependency(interfaceType)
	for _, item := range d.entries {
		if item.InterfaceType == interfaceType {
			return item
		}
	}
	// create it
	functionName := "resolve" + interfaceType
	created := dependencyResolverEntry{
		FunctionName:         functionName,
		InterfaceType:        interfaceType,
		InterfacePackageName: resolved.InterfacePackageName,
		RealType:             resolved.RealType,
		RealTypePackageName:  resolved.RealTypePackageName,
		Mode:                 resolved.Mode,
	}
	d.entries = append(d.entries, created)
	return created
}

func (d *DependencyResolver) GetFunctions() []dependencyResolverEntry {
	return d.entries
}
