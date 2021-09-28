package decoratorparser

type DecoratorParser struct {
	comments         string
	functions        FunctionsResponse
	controller       string
	controllerAction string
}

func New(commentsText string, controller string, controllerAction string) *DecoratorParser {
	p := &DecoratorParser{
		comments:         commentsText,
		controller:       controller,
		controllerAction: controllerAction,
	}
	p.functions = p.GetFunctions()
	return p
}
