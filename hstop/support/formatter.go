package support

import "fmt"

func Size(n uint64) string {
	if n > 1000*1000 {
		return fmt.Sprintf("%.1fM", (float64(n) / (1000 * 1000)))
	} else if n > 1000 {
		return fmt.Sprintf("%.1fk", (float64(n) / 1000))
	} else {
		return fmt.Sprintf("%v", n)
	}
}
