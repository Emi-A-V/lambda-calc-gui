package main

import (
	"context"
	"fmt"
	lambdaengine "lambda-calc-gui/lambda-calc"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	fmt.Printf("Starting Up")
	lambdaengine.Start(a.variableDefined, a.variableDropped)
}

type CalculationResult struct {
	Equation string `json:equation`
	Error    string `json:error`
	ErrorID  int    `json:error_id`
}

type Variable struct {
	Variable   string   `json:variable`
	Equation   string   `json:equation`
	Dependency []string `json:dependency`
}

func (a *App) Calculate(equation string) CalculationResult {
	fmt.Println("Entered: " + equation)

	res, eq := lambdaengine.Input(equation)

	fmt.Printf("%v\n", res)

	return CalculationResult{eq, res.Str, res.ErrID}
}

// Callback-function on defining a variable.
func (a *App) variableDefined(variable lambdaengine.Var) {
	newVariable := Variable{
		variable.Var,
		variable.Equ,
		variable.Dep,
	}

	runtime.EventsEmit(a.ctx, "app:variable_defined", newVariable)
}

// Callback-function on defining a variable.
func (a *App) variableDropped(variable lambdaengine.Var) {
	newVariable := Variable{
		variable.Var,
		variable.Equ,
		variable.Dep,
	}

	runtime.EventsEmit(a.ctx, "app:variable_dropped", newVariable)
}
