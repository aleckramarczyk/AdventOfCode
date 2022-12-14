package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readInput() (inputLines []string) {
	file, err := os.Open("input")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		inputLines = append(inputLines, scanner.Text())
	}
	return
}

func convertInput(inputLines []string) {
	forest = make([][]*tree, len(inputLines))
	for i, line := range inputLines {
		rawTreeHeights := strings.Split(line, "")
		var trees []*tree
		for _, theight := range rawTreeHeights {
			height, _ := strconv.Atoi(theight)
			newTree := &tree{
				height: height,
			}
			trees = append(trees, newTree)
		}
		forest[i] = trees
	}
	mapRelationships()
}

func mapRelationships() {
	for rowIndex, row := range forest {
		for treeIndex, t := range row {
			if rowIndex == 0 {
				t.up = nil
				t.down = forest[1][treeIndex]
			} else if rowIndex == len(forest)-1 {
				t.down = nil
				t.up = forest[rowIndex-1][treeIndex]
			} else {
				t.up = forest[rowIndex-1][treeIndex]
				t.down = forest[rowIndex+1][treeIndex]
			}
			if treeIndex == 0 {
				t.left = nil
				t.right = forest[rowIndex][treeIndex+1]
			} else if treeIndex == len(row)-1 {
				t.right = nil
				t.left = forest[rowIndex][treeIndex-1]
			} else {
				t.right = forest[rowIndex][treeIndex+1]
				t.left = forest[rowIndex][treeIndex-1]
			}
		}
	}
}

type tree struct {
	height           int
	left             *tree
	right            *tree
	up               *tree
	down             *tree
	visited          bool
	visible          bool
	isHighestScoring bool
}

var forest [][]*tree

func main() {
	inputLines := readInput()
	convertInput(inputLines)
	numberOfTreesVisibleFromEdges := getNumberOfTreesVisibleFromEdges()
	highestScenicScore := getHighestScenicScore()
	printForest()
	fmt.Println(numberOfTreesVisibleFromEdges)
	fmt.Println(highestScenicScore)
}

func getHighestScenicScore() (currentHighestScore int) {
	var highestScoringTree = new(tree)
	for _, row := range forest {
		for _, t := range row {
			leftScore := getLeftScore(t, 1, t.height)
			rightScore := getRightScore(t, 1, t.height)
			upScore := getUpScore(t, 1, t.height)
			downScore := getDownScore(t, 1, t.height)

			treeScore := leftScore * rightScore * upScore * downScore
			if treeScore > currentHighestScore {
				if highestScoringTree == nil {
					t.isHighestScoring = true
					highestScoringTree = t
				}
				currentHighestScore = treeScore
				highestScoringTree.isHighestScoring = false
				highestScoringTree = t
				t.isHighestScoring = true
			}
		}
	}
	return currentHighestScore
}

func getLeftScore(t *tree, score int, originalHeight int) int {
	if t.left != nil {
		if originalHeight > t.left.height {
			score += getLeftScore(t.left, score, originalHeight)
		}
	} else {
		return 0
	}
	return score
}

func getRightScore(t *tree, score int, originalHeight int) int {
	if t.right != nil {
		if originalHeight > t.right.height {
			score += getRightScore(t.right, score, originalHeight)
		}
	} else {
		return 0
	}
	return score
}

func getUpScore(t *tree, score int, originalHeight int) int {
	if t.up != nil {
		if originalHeight > t.up.height {
			score += getUpScore(t.up, score, originalHeight)
		}
	} else {
		return 0
	}
	return score
}

func getDownScore(t *tree, score int, originalHeight int) int {
	if t.down != nil {
		if originalHeight > t.down.height {
			score += getDownScore(t.down, score, originalHeight)
		}
	} else {
		return 0
	}
	return score
}

func getNumberOfTreesVisibleFromEdges() (visibleTrees int) {
	for _, row := range forest {
		//Check from left
		visibleHeight := -1
		for _, t := range row {
			if t.height > visibleHeight {
				if !t.visited {
					visibleTrees++
					t.visible = true
					t.visited = true
				}
				visibleHeight = t.height
			}
		}
		//Check from right
		visibleHeight = -1
		for tIndex := len(row) - 1; tIndex >= 0; tIndex-- {
			if row[tIndex].height > visibleHeight {
				if !row[tIndex].visited {
					visibleTrees++
					row[tIndex].visible = true
					row[tIndex].visited = true
				}
				visibleHeight = row[tIndex].height
			}
		}
	}
	for column := 0; column < len(forest[0]); column++ {
		//Check from top
		visibleHeight := -1
		for t := 0; t < len(forest); t++ {
			if forest[t][column].height > visibleHeight {
				if !forest[t][column].visited {
					visibleTrees++
					forest[t][column].visible = true
					forest[t][column].visited = true
				}
				visibleHeight = forest[t][column].height
			}
		}
		//Check from bottom
		visibleHeight = -1
		for t := len(forest) - 1; t >= 0; t-- {
			if forest[t][column].height > visibleHeight {
				if !forest[t][column].visited {
					visibleTrees++
					forest[t][column].visible = true
					forest[t][column].visited = true
				}
				visibleHeight = forest[t][column].height
			}
		}
	}
	return
}

func printForest() {
	for _, row := range forest {
		for _, t := range row {
			if t.isHighestScoring {
				colored := fmt.Sprintf("\x1b[%dm%s\x1b[0m", 34, strconv.Itoa(t.height))
				fmt.Printf(colored + " ")
			} else {
				if t.visible {
					colored := fmt.Sprintf("\x1b[%dm%s\x1b[0m", 32, strconv.Itoa(t.height))
					fmt.Printf(colored + " ")
				} else {
					fmt.Printf("%d ", t.height)
				}
			}
		}
		fmt.Printf("\n")
	}
}
