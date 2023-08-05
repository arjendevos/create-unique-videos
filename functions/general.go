package functions

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func isVideoFile(path string) bool {
	// You can add more video extensions if needed
	videoExtensions := []string{".mp4", ".mov", ".avi", ".mkv", ".wmv"}
	ext := strings.ToLower(filepath.Ext(path))
	for _, vExt := range videoExtensions {
		if ext == vExt {
			return true
		}
	}
	return false
}

func formatTime(seconds float64) string {
	duration := time.Duration(seconds * float64(time.Second))
	return fmt.Sprintf("%02d:%02d:%02d", int(duration.Hours()), int(duration.Minutes())%60, int(duration.Seconds())%60)
}

func getVideoLength(inputFilePath string) (float64, error) {
	cmd := exec.Command("ffprobe", "-i", inputFilePath, "-show_entries", "format=duration", "-v", "quiet", "-of", "csv=p=0")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("error getting video file length: %v", err)
	}

	length := strings.TrimSpace(string(output))
	return strconv.ParseFloat(length, 64)
}
