package main

import (
	"context"
	"fmt"
	lambdaengine "lambda-calc-gui/lambda-calc"
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
	lambdaengine.Start()
}

type CalculationResult struct {
	Equation string `json:equation`
	Error    string `json:error`
	ErrorID  int    `json:error_id`
}

func (a *App) Calculate(equation string) CalculationResult {
	fmt.Println("Entered: " + equation)

	res, eq := lambdaengine.Input(equation)

	fmt.Printf("%v\n", res)

	return CalculationResult{eq, res.Str, res.ErrID}
}
