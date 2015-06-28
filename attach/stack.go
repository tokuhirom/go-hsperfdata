package attach

import (
	"regexp"
	"strings"
)

type JavaThread struct {
	name string
}

func ParseStack(stack string) ([]*JavaThread, error) {
	re, err := regexp.Compile(`^"([^"]+)"`)
	if err != nil {
		return nil, err
	}
	atre, err := regexp.Compile(`\s*at ([^(]+)\(`)

	threads := make([]*JavaThread, 0)
	lines := strings.Split(stack, "\n")
	name := ""
	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) == 2 {
			name = matches[1]
		} else {
			matches = atre.FindStringSubmatch(line)
			if len(matches) > 0 {
			} else if len(name) > 0 {
				threads = append(threads, &JavaThread{name})
				name = ""
			}
		}
	}
	return threads, nil
}
