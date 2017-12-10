package alfred

import (
	"bytes"
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	event "github.com/kcmerrill/hook"
)

// NewTask will execute a task
func NewTask(task string, context *Context, tasks map[string]Task) {
	t, exists := tasks[task]
	if !exists {
		// TODO, Lookup ... then exit
		return
	}
	c := context
	c.TaskName = task
	event.Trigger("task.new", t, c)
	event.Trigger("task.summary.header", t, c)
	event.Trigger("task.command", t, c)
	event.Trigger("task.summary.footer", t, c)
	event.Trigger("task.group", t.Ok, t, c, tasks)
	//event.Trigger("task.group", t.Fail, t, c, tasks)
}

// Task holds all of our task components
type Task struct {
	Aliases     string
	Summary     string
	Description string
	Args        []string
	Dir         string
	Command     string
	Script      string
	Ok          string
	Fail        string
	ExitCode    int
}

// TaskGroup contains a task name and it's arguments
type TaskGroup struct {
	Name string
	Args []string
}

// ParseTaskGroup takes in a string, and parses it into a TaskGroup
func (t *Task) ParseTaskGroup(group string) []TaskGroup {
	tg := make([]TaskGroup, 0)
	group = strings.TrimSpace(group)

	if group == "" {
		return tg
	}

	if strings.Index(group, "\n") == -1 {
		// This means we have a regular space delimited list
		tasks := strings.Split(group, " ")
		for _, task := range tasks {
			tg = append(tg, TaskGroup{Name: task, Args: []string{}})
		}
	} else {
		// mix and match here
	}

	return tg
}

// Exit determins whether a task should exit or not
func (t *Task) Exit() {
	if t.ExitCode != 0 {
		os.Exit(t.ExitCode)
	}
}

// Template is a helper function to translate a string to a template
func (t *Task) Template(translate string, context *Context) string {
	if translate == "" {
		// Nothing to translate, move along
		return translate
	}
	fmap := sprig.TxtFuncMap()
	te := template.Must(template.New("template").Funcs(fmap).Parse(translate))
	var b bytes.Buffer
	err := te.Execute(&b, context)
	if err != nil {
		event.Trigger("speak", "{{ .Text.Failure }} Bad Template: "+err.Error()+"{{ .Text.Reset }}", t, context)
		return ""
	}
	return b.String()
}
