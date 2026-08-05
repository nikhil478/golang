package main

import (
	"fmt"

	"github.com/nikhil478/golang/projects/db/projects/db/internal/engine"
)

func main() {

	engine := engine.NewEngine("greet.json")
	if err := engine.Set("greet", []byte("HELLO NIKHIL !")); err != nil {
		fmt.Printf("error while setting value in engine : %v \n", err)
	}
	if val, err := engine.Get("greet"); err != nil {
		fmt.Printf("error while getting value in engine  { key : %v || value : %v } \n", "greet", err)
	} else {
		fmt.Printf("value while getting value in engine  { key : %v || value : %s } \n", "greet", val)
	}
	if err := engine.Delete("greet"); err != nil {
		fmt.Printf("error while deleting value in engine : %v \n", err)
	}
}
