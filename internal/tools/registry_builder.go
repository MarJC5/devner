package tools

import "github.com/devner/devner/internal/app"

// BuildDefault returns the tool registry the agent loop exposes to the LLM
// for a given app.Deps bundle.
func BuildDefault(d *app.Deps) *Registry {
	return NewRegistry(
		&ListProjects{D: d},
		&ProjectStatus{D: d},
		&CreateProject{D: d},
		&DeleteProject{D: d},
		&CreateDatabase{D: d},
		&DropDatabase{D: d},
		&StartStack{D: d},
		&StopStack{D: d},
		&RebuildStack{D: d},
		&TailLogs{D: d},
		&StartDevServer{D: d},
		&StopDevServer{D: d},
		&DevServerStatus{D: d},
		&Composer{D: d},
		&NPM{D: d},
		&WPCli{D: d},
		&Artisan{D: d},
		&OpenInEditor{D: d},
		&ExecInProject{D: d},
	)
}
