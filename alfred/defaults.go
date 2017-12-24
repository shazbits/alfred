package alfred

import (
	"os"
)

func defaults(task Task, context *Context, tasks map[string]Task) {
	if len(task.Defaults) == 0 {
		return
	}

	if len(task.Defaults) < len(context.Args) {
		// no need to set defaults, as the args are already set
		return
	}

	// ok, so we have some defaults, lets update the context.
	for idx := len(context.Args); idx < len(task.Defaults); idx++ {
		// empty? we should bail ...
		if task.Defaults[idx] == "" {
			outFail("template", "Invalid Argument(s)", context)
			os.Exit(42)
		}
		// set the defaults
		context.Args = append(context.Args, task.Defaults[idx])
	}
}
