package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	var inputFolder string

	// Define command-line flags
	flag.StringVar(&inputFolder, "input", "", "Input containing video or videos")
	flag.Parse()

	fmt.Println()

	// Check if input and output folders are provided
	if inputFolder == "" {
		fmt.Println("Usage: go run main.go -input-folder=input_folder")
		return
	}

	// Check if input folder exists
	fileInfo, err := os.Stat(inputFolder)
	if err != nil || os.IsNotExist(err) {
		fmt.Println("Input folder does not exist:", inputFolder)
		return
	}

	isInputDir := fileInfo.IsDir()
	inputFolder, err = filepath.Abs(inputFolder)
	if err != nil {
		fmt.Println("Error getting absolute path of input folder:", err)
		return
	}

	if isInputDir {

		// Traverse the input folder for video files
		err = filepath.Walk(inputFolder, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println("Error accessing file:", err)
				return nil
			}

			if !info.IsDir() && isVideoFile(path) {
				outputPath, err := trimAndZoom(path, "edited_", 20, 2)
				if err != nil {
					panic(err)
				}

				tempOutputPath, err := removeMetadataAndSound(outputPath, "temporary_")
				if err != nil {
					panic(err)
				}

				if err := os.Rename(tempOutputPath, outputPath); err != nil {
					panic(err)
				}

			}

			return nil
		})

		if err != nil {
			fmt.Println("Error traversing input folder:", err)
		}

	} else {
		outputPath, err := trimAndZoom(inputFolder, "edited_", 20, 2)
		if err != nil {
			panic(err)
		}

		tempOutputPath, err := removeMetadataAndSound(outputPath, "temporary_")
		if err != nil {
			panic(err)
		}

		if err := os.Rename(tempOutputPath, outputPath); err != nil {
			panic(err)
		}

	}
}

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

func trimAndZoom(inputFilePath, outputPrefix string, zoomFactor, trimByInSeconds float64) (string, error) {
	inputFileName := filepath.Base(inputFilePath)

	videoLength, err := getVideoLength(inputFilePath)
	if err != nil {
		return "", err
	}

	//scale=iw/zoomFactor:-1
	outputPath, err := execFfmpeg(inputFilePath, outputPrefix, "-vf", fmt.Sprintf("crop=iw/%v:ih/%v, scale=iw/%v:-1", zoomFactor, zoomFactor, zoomFactor), "-ss", "00:00:00", "-to", strconv.FormatFloat(videoLength-trimByInSeconds, 'f', 0, 64))
	if err != nil {
		return "", err
	}

	//-ss 00:00:00 -to $(echo "$(ffprobe -i input_video.mp4 -show_entries format=duration -v quiet -of csv=p=0) - 1" | bc)

	fmt.Printf("Video file %s trimmed and zoomed by %v\n", inputFileName, zoomFactor)
	return outputPath, nil
}

func removeMetadataAndSound(inputFilePath, outputPrefix string) (string, error) {
	inputFileName := filepath.Base(inputFilePath)

	fmt.Println("Processing video file:", inputFileName)
	// ffmpeg -i input_video.mp4 -metadata:s:v:0 "Lv Meta Info={"data":{"os":"web","videoId":"8cBtNeTkL5kFM6cYlhdhGXLqqcUehvP16AtaJnr","product":"vicutweb","editType":"edit","exportType":"export","original_edittype":"edit","templateId":"","appVersion":"10.5.0"},"source_type":"vicutweb"}" -c copy output_video.mp4

	outputPath, err := execFfmpeg(inputFilePath, outputPrefix, "-map_metadata", "-1", "-c:v", "copy", "-an")
	if err != nil {
		return "", err
	}

	fmt.Printf("Metadata removed from %s\n", inputFileName)
	return outputPath, nil
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

func execFfmpeg(inputFilePath, outputPrefix string, args ...string) (string, error) {
	inputFileName := filepath.Base(inputFilePath)
	outputFilePath := filepath.Join(filepath.Dir(inputFilePath), outputPrefix+inputFileName)

	commandArgs := []string{"-i", inputFilePath}
	commandArgs = append(commandArgs, args...)
	commandArgs = append(commandArgs, []string{outputFilePath, "-y"}...)

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
