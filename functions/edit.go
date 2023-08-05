package functions

import (
	"fmt"
	"os"
)

func EditVideo(input string, zoomFactor, trimByInSeconds float64) (string, error) {
	if !isVideoFile(input) {
		return "", fmt.Errorf("input is not a video file")
	}
	// Remove metadata and sound
	outputPath, err := execFFMPEG(input, "meta_", "-map_metadata", "-1", "-c:v", "copy", "-an")
	if err != nil {
		return "", err
	}

	// Convert to tiktok frame
	tempOutputPath, err := execFFMPEG(outputPath, "temp_", "-vf", "scale=1080:1920")
	if err != nil {
		return "", err
	}

	if err := os.Rename(tempOutputPath, outputPath); err != nil {
		return "", err
	}

	cropRatio := fmt.Sprintf("%.2f", 1/zoomFactor)

	// Get video length
	videoLength, err := getVideoLength(outputPath)
	if err != nil {
		return "", err
	}

	// Trim and zoom
	tempOutputPath2, err := execFFMPEG(outputPath, "zoom_", "-vf", fmt.Sprintf("scale=iw*%v:-1, crop=in_w*%v:in_h*%v", zoomFactor, cropRatio, cropRatio), "-ss", "00:00:00", "-to", formatTime(videoLength-trimByInSeconds))
	if err != nil {
		return "", err
	}

	if err := os.Rename(tempOutputPath2, outputPath); err != nil {
		return "", err
	}

	err = updateVideoMetadata(outputPath)
	if err != nil {
		return "", err
	}

	if err := os.Remove(input); err != nil {
		return "", err
	}

	return outputPath, nil
}
