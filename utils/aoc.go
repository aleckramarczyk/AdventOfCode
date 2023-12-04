package utils

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func GetInput(day int, year int, session string) {
	inputFilePath := fmt.Sprintf("%d/day%d/input", year, day)
	if exists, err := fileExists(inputFilePath); err == nil {
		if exists {
			log.Printf("using cached input")
			return
		} else {
			log.Printf("requesting input")
			content := requestInput(day, year, session)
			err = cacheInput(content, year, day)
			if err != nil {
				log.Printf("failed to cache input: %s", err)
			}
		}
	} else {
		log.Fatalf("could not validate cached input: %s", err)
	}
}

func getDayString(day int) string {
	dayString := "day" + strconv.Itoa(day)
	return dayString
}

func cacheInput(input []byte, year int, day int) (err error) {
	dir := fmt.Sprintf("%d/%s/", year, getDayString(day))
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Printf("failed to create cache directory %s: %s", dir, err)
		return err
	}

	inputFile, err := os.Create(filepath.Join(dir, filepath.Base("input")))
	if err != nil {
		log.Printf("failed to create input cache file: %s", err)
		return err
	}
	defer inputFile.Close()

	_, err = io.Copy(inputFile, bytes.NewReader(input))
	if err != nil {
		log.Printf("failed to copy input to file: %s", err)
		return err
	}
	return nil
}

func requestInput(day int, year int, session string) []byte {
	// get request with session cookie
	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("creating request %s", err)
	}

	req.AddCookie(&http.Cookie{Name: "session", Value: session})

	// make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("making request %s", err)
	}

	var body []byte
	body, err = io.ReadAll(resp.Body)

	if strings.HasPrefix(string(body), "Please don't repeatedly") {
		log.Fatalf("Repeated request error")
	}

	return body
}

func fileExists(filepath string) (bool, error) {
	if _, err := os.Stat(filepath); err == nil {
		// file exists
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		// file does not exist
		return false, nil
	} else {
		// Some other error
		return false, err
	}
}

func ParseFlags() (day int, year int, part int, session string) {
	flag.IntVar(&day, "d", 0, "day")
	flag.IntVar(&year, "y", 0, "year")
	flag.IntVar(&part, "p", 0, "part 1 or 2")
	flag.StringVar(&session, "s", os.Getenv("AOC_SESSION_COOKIE"), "session cookie")
	flag.Parse()

	// default to today
	today := time.Now()
	if day == 0 {
		day = today.Day()
	}
	if year == 0 {
		year = today.Year()
	}

	// validate flag inputs
	if day > 25 || day < 1 {
		log.Fatalf("day out of range: %d", day)
	}
	if part > 2 || part < 0 {
		log.Printf("part specified is out of range, defaulting to both")
		part = 0
	}
	if session == "" {
		log.Fatalf("no session cookie set on flag or AOC_SESSION_COOKIE environment variable")
	}

	return day, year, part, session
}

func Run(day int, year int, part int) {
	dayDir := fmt.Sprintf("%d/%s/", year, getDayString(day))
	// check if solution exists for this day
	solution := dayDir + "main.go"
	if exists, err := fileExists(solution); err == nil {
		if !exists {
			log.Printf("solution does not exist for %d day %d, please check back later", year, day)
			return
		}
	} else {
		log.Fatalf("could not validate if solution exists: %s", err)
	}

	// build solution
	log.Printf("building %d, %d solution", year, day)
	solutionExecutable := dayDir + "main"
	cmd := exec.Command("go", "build", "-o", solutionExecutable, solution)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("failed to build %s: %s", solution, err)
	}

	// run solution
	log.Printf("running %d, %d, part %d", year, day, part)
	start := time.Now()
	cmd = exec.Command(solutionExecutable, "-p", strconv.Itoa(part))
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("failed to run solution: %s", err)
	}
	duration := time.Since(start)
	fmt.Println(string(out))
	log.Println("solution completed in ", duration)

	// delete built executable
	err = os.Remove(solutionExecutable)
	if err != nil {
		log.Printf("error deleting built executable: %s", err)
	}
}
