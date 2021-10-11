package main

import (
	"github.com/floatzeI/go-web-autogen/generator"
)

func main() {
	generator.Generate(generator.GenerateOptions{
		// Directory:         "./examples/project-one/",
		Directory:         "/home/buizel/sites/RblxTrade/RblxTrade.OwnershipService/",
		ControllersFolder: "controllers",
		ResolveModelFolder: func(packageName string, fileName string) string {
			return "models/" + packageName
		},
	})
}
