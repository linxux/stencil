package {{project_name}}

const Version = "{{version}}"

// App represents the application
type App struct {
	Name string
}

// NewApp creates a new application instance
func NewApp(name string) *App {
	return &App{
		Name: name,
	}
}

// Run starts the application
func (a *App) Run() error {
	// Implementation here
	return nil
}
