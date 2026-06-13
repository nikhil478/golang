package logs

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/nikhil478/knowledgebuilder/internal/models"
)

var (
	serviceRegex = regexp.MustCompile(`\[(.*?)\]`)
	traceRegex   = regexp.MustCompile(`traceId=([a-zA-Z0-9\-]+)`)
)

func ParseLog(path string) ([]models.Fact, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var facts []models.Fact

	scanner := bufio.NewScanner(file)

	line := 0

	for scanner.Scan() {

		line++

		txt := scanner.Text()

		if strings.Contains(txt, "ERROR") {

			facts = append(facts, models.Fact{
				Type:       "error",
				Value:      txt,
				Confidence: 0.9,
				SourceFile: path,
				SourceLine: line,
			})
		}

		matches := serviceRegex.FindStringSubmatch(txt)

		if len(matches) > 1 {

			facts = append(facts, models.Fact{
				Type:       "service",
				Value:      matches[1],
				Confidence: 0.8,
				SourceFile: path,
				SourceLine: line,
			})
		}

		trace := traceRegex.FindStringSubmatch(txt)

		if len(trace) > 1 {

			facts = append(facts, models.Fact{
				Type:       "trace",
				Value:      trace[1],
				Confidence: 1.0,
				SourceFile: path,
				SourceLine: line,
			})
		}
	}

	return facts, scanner.Err()
}
