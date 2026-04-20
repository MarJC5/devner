package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/devner/devner/internal/app"
	"github.com/devner/devner/internal/project"
)

type OpenInEditor struct{ D *app.Deps }

type openInEditorArgs struct {
	Project string `json:"project"`
	Editor  string `json:"editor,omitempty"`
}

func (t *OpenInEditor) Name() string { return "open_in_editor" }
func (t *OpenInEditor) Description() string {
	return "Open a project directory in a local GUI editor (VS Code, Cursor, Zed, Sublime Text, JetBrains). Default: first editor found on PATH."
}
func (t *OpenInEditor) Destructive() bool { return false } // launches a local GUI, no state change
func (t *OpenInEditor) Schema() json.RawMessage {
	return json.RawMessage(`{
  "type":"object",
  "properties":{
    "project":{"type":"string","description":"Devner project name (see list_projects)."},
    "editor":{"type":"string","enum":["code","cursor","zed","subl","webstorm","phpstorm","idea"],"description":"Optional editor. Omit to use the first editor found on PATH."}
  },
  "required":["project"],
  "additionalProperties":false
}`)
}
func (t *OpenInEditor) Execute(ctx context.Context, raw json.RawMessage) (Result, error) {
	var a openInEditorArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return Result{Content: "bad args: " + err.Error()}, err
	}
	p, err := t.D.Store.GetProject(ctx, a.Project)
	if err != nil {
		return Result{Content: fmt.Sprintf("project %q not found", a.Project)}, err
	}
	editor := a.Editor
	if editor == "" {
		resolved, err := project.DefaultEditor()
		if err != nil {
			return Result{Content: err.Error()}, err
		}
		editor = resolved
	}
	if err := project.Open(editor, p.Path); err != nil {
		return Result{Content: err.Error()}, err
	}
	return Result{Content: fmt.Sprintf("✓ opened %s in %s (%s)", p.Name, editor, p.Path)}, nil
}
