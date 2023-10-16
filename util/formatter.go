package util

import (
	"fmt"
	"math"
)

func FormatBytes(bytes float64) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}

	if bytes == 0 {
		return "0 B"
	}

	exp := int(math.Log(bytes) / math.Log(1024))
	size := bytes / math.Pow(1024, float64(exp))

	unit := units[exp]
	return fmt.Sprintf("%.2f %s", size, unit)
}
