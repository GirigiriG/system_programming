package utils

import (
	"fmt"
	"os"
)

func GetTemp() {
	data, _ := os.ReadFile("/proc")
	fmt.Println("CPU temp (m°C):", string(data))
}
