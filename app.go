package main

import (
	"context"
	"lambda-calc-gui/lambda-calc"
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
	lambdaengine.Start()
}

type Config struct {
	Symbols map[string]string
}

var config Config

func (a *App) Calculate(equation string) string {
	res := lambdaengine.Input(equation)
	return res.Str
}
