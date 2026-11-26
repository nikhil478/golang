package handler

import (
	"fmt"
	"strconv"
	"strings"
)

func parseSize(value string) (int64, error) {
	value = strings.ToLower(strings.TrimSpace(value))

	units := []struct {
		suffix     string
		multiplier int64
	}{
		{"gb", 1 << 30},
		{"mb", 1 << 20},
		{"kb", 1 << 10},
		{"b", 1},
	}

	for _, unit := range units {
		if strings.HasSuffix(value, unit.suffix) {
			number := strings.TrimSuffix(value, unit.suffix)
			n, err := strconv.ParseInt(number, 10, 64)
			if err != nil || n <= 0 || n > (1<<63-1)/unit.multiplier {
				return 0, fmt.Errorf("invalid size")
			}
			return n * unit.multiplier, nil
		}
	}

	return 0, fmt.Errorf("size must use B, KB, MB, or GB")
}