package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/devner/devner/internal/app"
	"github.com/devner/devner/internal/database"
)

type CreateDatabase struct{ D *app.Deps }

type dbArgs struct {
	Name   string `json:"name"`
	Engine string `json:"engine"`
}

func (t *CreateDatabase) Name() string        { return "create_database" }
func (t *CreateDatabase) Description() string { return "Create a database and a dedicated user on mysql or postgres. Returns credentials." }
func (t *CreateDatabase) Destructive() bool   { return false }
func (t *CreateDatabase) Schema() json.RawMessage {
	return json.RawMessage(`{
  "type":"object",
  "properties":{
    "name":{"type":"string"},
    "engine":{"type":"string","enum":["mysql","postgres"]}
  },
  "required":["name","engine"],
  "additionalProperties":false
}`)
}
func (t *CreateDatabase) Execute(ctx context.Context, raw json.RawMessage) (Result, error) {
	var a dbArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return Result{Content: "bad args: " + err.Error()}, err
	}
	creds, err := t.D.DB.Create(ctx, database.Engine(a.Engine), a.Name)
	if err != nil {
		return Result{Content: "create failed: " + err.Error()}, err
	}
	return Result{
		Content: fmt.Sprintf("✓ %s/%s user=%s host=%s port=%d", a.Engine, creds.Database, creds.User, creds.Host, creds.Port),
		Data:    creds,
	}, nil
}

type DropDatabase struct{ D *app.Deps }

func (t *DropDatabase) Name() string        { return "drop_database" }
func (t *DropDatabase) Description() string { return "Drop a database and its user. DESTRUCTIVE — all data is lost." }
func (t *DropDatabase) Destructive() bool   { return true }
func (t *DropDatabase) Schema() json.RawMessage {
	return json.RawMessage(`{
  "type":"object",
  "properties":{
    "name":{"type":"string"},
    "engine":{"type":"string","enum":["mysql","postgres"]}
  },
  "required":["name","engine"],
  "additionalProperties":false
}`)
}
func (t *DropDatabase) Execute(ctx context.Context, raw json.RawMessage) (Result, error) {
	var a dbArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return Result{Content: "bad args: " + err.Error()}, err
	}
	if err := t.D.DB.Drop(ctx, database.Engine(a.Engine), a.Name); err != nil {
		return Result{Content: "drop failed: " + err.Error()}, err
	}
	return Result{Content: "✓ dropped " + a.Engine + "/" + a.Name}, nil
}
