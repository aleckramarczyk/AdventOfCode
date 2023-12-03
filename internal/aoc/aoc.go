package aoc

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GetInput() (input []byte) {
	day, year, session := parseFlags()

	input, err := getCachedInput(day, year)
	if err != nil {
		input := requestInput(day, year, session)
		err = cacheInput(input, year, day)
		if err != nil {
			log.Printf("failed to cache input: %s", err)
		}
	}

	return input
}

// TODO
func cacheInput(input []byte, year int, day int) (err error) {
	dir := fmt.Sprintf("../cache/%d/%d/", year, day)
	cache, err := os.Create(filepath.Join(dir, filepath.Base("input")))
	if err != nil {
		log.Printf("failed to create input cache file: %s", err)
		return err
	}
	defer cache.Close()

	_, err = io.Copy(cache, bytes.NewReader(input))
	if err != nil {
		log.Printf("failed to copy input to file: %s", err)
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

func getCachedInput(day int, year int) ([]byte, error) {
	input, err := os.ReadFile(fmt.Sprintf("../cache/%d/%d/input", year, day))
	if err != nil {
		return nil, err
	}
	return input, nil
}

func parseFlags() (day int, year int, session string) {
	flag.IntVar(&day, "d", 0, "day")
	flag.IntVar(&year, "y", 0, "year")
	flag.StringVar(&session, "s", os.Getenv("AOC_SESSION_COOKIE"), "session cookie")

	// default to today
	today := time.Now()
	if day == 0 {
		day = today.Day()
	}
	if year == 0 {
		year = today.Year()
	}

	if day > 25 || day < 1 {
		log.Fatalf("day out of range: %d", day)
	}

	if session == "" {
		log.Fatalf("no session cookie set on flag or AOC_SESSION_COOKIE enviroment variable")
	}

	return day, year, session
}
