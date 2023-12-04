package main

import (
	"github.com/aleckramarczyk/AdventOfCode/utils"
)

func main() {
	day, year, part, session := utils.ParseFlags()
	utils.GetInput(day, year, session)
	utils.Run(day, year, part)
}
