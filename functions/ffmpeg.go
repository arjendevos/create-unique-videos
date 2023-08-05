package functions

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

func execFFMPEG(inputFilePath, outputPrefix string, args ...string) (string, error) {
	inputFileName := filepath.Base(inputFilePath)
	outputFilePath := filepath.Join(filepath.Dir(inputFilePath), outputPrefix+inputFileName)

	commandArgs := []string{"-i", inputFilePath}
	commandArgs = append(commandArgs, args...)
	// commandArgs = append(commandArgs, []string{outputFilePath, "-y"}...)
	commandArgs = append(commandArgs, []string{"-c:v", "libx264", "-crf", "23", outputFilePath, "-y"}...)

	fmt.Print("ffmpeg ")

	for _, arg := range commandArgs {
		fmt.Print(arg + " ")
	}

	fmt.Println()
	fmt.Println()

	cmd := exec.Command("ffmpeg", commandArgs...)
	// cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error processing video file %s: %v", inputFilePath, err)
	}

	return outputFilePath, nil
}
