package logs

import (
	"bufio"
	"os"
	"regexp"

	"github.com/nikhil478/knowledgebuilder/internal/models"
)

var errorRegex = regexp.MustCompile(
	`ERROR|Exception|Caused by|FATAL`,
)

func Analyze(path string) ([]models.Fact, error) {

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

		if errorRegex.MatchString(txt) {
			facts = append(facts, models.Fact{
				Type:       "error",
				Value:      txt,
				Confidence: 0.90,
				SourceFile: path,
				SourceLine: line,
			})
		}
	}

	return facts, scanner.Err()
}
