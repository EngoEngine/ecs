//go:build actiontesting

package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

func main() {
	var err error

	test := exec.Command("go", "test", "--coverprofile=.github/coverage.out", "--bench=.", ".")
	if err = test.Run(); err != nil {
		log.Fatalf("Unable to run go test. Error was: %v", err.Error())
	}

	cover := exec.Command("go", "tool", "cover", "--func=.github/coverage.out", "-o=.github/coverage.out")
	if err = cover.Run(); err != nil {
		log.Fatalf("Unable to run go coverage tool. Error was: %v", err.Error())
	}

	output, err := os.Open(".github/coverage.out")
	if err != nil {
		log.Fatalf("Unable to open [[ coverage.out ]]. Error was: %v", err.Error())
	}
	defer output.Close()

	outstat, err := output.Stat()
	if err != nil {
		log.Fatalf("Unable to get stats for [[ coverage.out ]]. Error was: %v", err.Error())
	}

	coverageOutput := make([]byte, outstat.Size())
	_, err = output.Read(coverageOutput)
	if err != nil {
		log.Fatalf("Unable to read [[ coverage.out ]]. Error was: %v", err.Error())
	}

	exp, err := regexp.Compile(`\(statements\)\s+\d+.\d+`)
	if err != nil {
		log.Fatalf("Regular Expression 1 failed to compile. Error was: %v", err.Error())
	}

	exp2, err := regexp.Compile(`\d+.\d+`)
	if err != nil {
		log.Fatalf("Regular Expression 2 failed to compile. Error was: %v", err.Error())
	}

	coverage := exp.Find(coverageOutput)
	if coverage == nil {
		log.Fatal("Unable to find coverage match 1 from file output.")
	}
	coverage = exp2.Find(coverage)
	if coverage == nil {
		log.Fatal("Unable to find coverage match 2 from file output.")
	}

	coverageValue, err := strconv.ParseFloat(string(coverage), 64)
	if err != nil {
		log.Fatalf("Unable to convert match to value. Error was: %v", err.Error())
	}

	color := "000000"
	if coverageValue > 90.0 {
		color = "00ff00"
	} else if coverageValue > 80.0 {
		color = "ffff00"
	} else if coverageValue > 70.0 {
		color = "ff0000"
	}
	url := "https://img.shields.io/badge/coverage-" + string(coverage) + "%25-" + color
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Unable to get badge from [[ %v ]]. Error was: %v", url, err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Unable to read response body from [[ %v ]]. Error was: %v", url, err.Error())
	}

	svg, err := os.Create(".github/coverage.svg")
	if err != nil {
		log.Fatalf("Unable to create [[ .github/coverage.svg ]]. Error was: %v", err.Error())
	}
	defer svg.Close()

	_, err = svg.Write(body)
	if err != nil {
		log.Fatalf("Unable to write to [[ .github/coverage.svg ]]. Error was: %v", err.Error())
	}

	checkForCoverageChanged := exec.Command("git", "diff-index", "--cached", "HEAD")
	gitOutput, err := checkForCoverageChanged.StdoutPipe()
	if err != nil {
		log.Fatalf("Unable to get stdout pipe for coverage changed command. Error was: %v", err.Error())
	}
	if err = checkForCoverageChanged.Start(); err != nil {
		log.Fatalf("Unable to start git diff-index. Error was: %v", err.Error())
	}
	gitBytes, _ := io.ReadAll(gitOutput)
	coverMatch, err := regexp.Match(`\.github\/coverage.out`, gitBytes)
	if err != nil {
		log.Fatalf("Unable to match git diff-index output with [[ coverage.out ]]. Error was: %v", err.Error())
	}
	svgMatch, err := regexp.Match(`\.github\/coverage\.svg`, gitBytes)
	if err != nil {
		log.Fatalf("Unable to match git diff-index output with [[ coverage.svg ]]. Error was: %v", err.Error())
	}
	if err = checkForCoverageChanged.Wait(); err != nil {
		log.Fatalf("Error waiting for git diff-index. Error was: %v", err.Error())
	}
	if coverMatch || svgMatch {
		if err = commitToPR(); err != nil {
			log.Fatalf("Unable to commit to PR. Error was: %v", err.Error())
		}
	}
}

func commitToPR() error {
	var err error
	gitUName := exec.Command("git", "config", "--global", "user.name", "'CoverageBot'")
	if err = gitUName.Run(); err != nil {
		return err
	}
	gitEmail := exec.Command("git", "config", "--global", "email", "'coverageBot@users.noreply.github.com'")
	if err = gitEmail.Run(); err != nil {
		return err
	}
	gitAddAll := exec.Command("git", "add", ".github/coverage.out", ".github/coverage.svg")
	if err = gitAddAll.Run(); err != nil {
		return err
	}
	gitCommit := exec.Command("git", "commit", "-m", "'[[BOT]] Coverage changed. Updating badge and coverage output. [skip ci]'")
	if err = gitCommit.Run(); err != nil {
		return err
	}
	gitPush := exec.Command("git", "push")
	if err = gitPush.Run(); err != nil {
		return err
	}
	return nil
}
