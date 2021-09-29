package main

import (
	"web-autogen/generator"
)

func main() {
	generator.Generate(generator.GenerateOptions{
		Directory:         "./examples/project-one/",
		ControllersFolder: "controllers",
		ResolveModelFolder: func(packageName string, fileName string) string {
			return "models/" + packageName
		},
	})
}
