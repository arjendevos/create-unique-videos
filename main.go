package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	var input string

	// Define command-line flags
	flag.StringVar(&input, "input", "", "Input containing video or videos")
	flag.Parse()

	// Check if input and output folders are provided
	if input == "" {
		fmt.Println("Usage: go run main.go -input-folder=input_folder")
		return
	}

	// Check if input folder exists
	fileInfo, err := os.Stat(input)
	if err != nil || os.IsNotExist(err) {
		fmt.Println("Input folder does not exist")
		return
	}

	isInputDir := fileInfo.IsDir()
	input, err = filepath.Abs(input)
	if err != nil {
		fmt.Println("Error getting absolute path of input folder:", err)
		return
	}

	if isInputDir {
		// Traverse the input folder for video files
		if err := filepath.Walk(input, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println("Error accessing file:", err)
				return nil
			}

			if info.IsDir() {
				return nil
			}

			if err := editVideo(path); err != nil {
				return err
			}

			return nil
		}); err != nil {
			fmt.Println("Error traversing input folder:", err)
		}

		return
	}

	if err := editVideo(input); err != nil {
		fmt.Println("Error editing video:", err)
	}

}

func editVideo(input string) error {
	if !isVideoFile(input) {
		return fmt.Errorf("input is not a video file")
	}
	// Remove metadata and sound
	outputPath, err := execFFMPEG(input, "meta_", "-map_metadata", "-1", "-c:v", "copy", "-an")
	if err != nil {
		return err
	}

	// Convert to tiktok frame
	tempOutputPath, err := execFFMPEG(outputPath, "temp_", "-vf", "scale=1080:1920")
	if err != nil {
		return err
	}

	if err := os.Rename(tempOutputPath, outputPath); err != nil {
		return err
	}

	zoomFactor := 1.1
	trimByInSeconds := 0.5
	cropRatio := fmt.Sprintf("%.2f", 1/zoomFactor)

	// Get video length
	videoLength, err := getVideoLength(outputPath)
	if err != nil {
		return err
	}

	// Trim and zoom
	tempOutputPath2, err := execFFMPEG(outputPath, "zoom_", "-vf", fmt.Sprintf("scale=iw*%v:-1, crop=in_w*%v:in_h*%v", zoomFactor, cropRatio, cropRatio), "-ss", "00:00:00", "-to", formatTime(videoLength-trimByInSeconds))
	if err != nil {
		return err
	}

	if err := os.Rename(tempOutputPath2, outputPath); err != nil {
		return err
	}

	return updateVideoMetadata(outputPath)

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

func updateVideoMetadata(videoPath string) error {
	// Generate the current time in the desired format.
	currentTime := time.Now()

	// Prepare the exiftool command with the desired metadata.
	cmdArgs := []string{
		"-overwrite_original",
		"-CreateDate=" + currentTime.Format("2006:01:02 15:04:05"),
		"-ModifyDate=" + currentTime.Format("2006:01:02 15:04:05"),
		"-TrackCreateDate=" + currentTime.Format("2006:01:02 15:04:05"),
		"-TrackModifyDate=" + currentTime.Format("2006:01:02 15:04:05"),
		"-MediaCreateDate=" + currentTime.Format("2006:01:02 15:04:05"),
		"-MediaModifyDate=" + currentTime.Format("2006:01:02 15:04:05"),
		"-HandlerType=Metadata Tags",
		"-HandlerDescription=DataHandler",
		"-AudioFormat=mp4a",
		"-AudioBitsPerSample=16",
		"-AudioSampleRate=44100",
		"-CompressorID=hvc1",
		"-VendorID=FFmpeg",
		// "-HandlerVendorID=",
		"-Balance=0",
		"-PurchaseFileFormat=mp4a",
		"-LayoutFlags=Stereo",
		"-AudioChannels=2",
		"-MoovPos=4700321",
		"-WriterType=-1",
		"-Copyright=cb2e8b799d13edf6c2e5b21baa16c737",
		"-MdatPos=28",
		"-MinorVersion=512",
		"-Encoder=Lavf57.71.100",
		"-Software={\"publicMode\":\"1\",\"TEEditor\":\"2\",\"isFastImport\":\"0\",\"transType\":\"2\",\"te_is_reencode\":\"1\",\"source\":\"2\"}",
		"-Source=2",
		"-Hw=1",
		"-MajorBrand=qt",
		"-CompatibleBrands=qt",
		"-TeIsReencode=1",
		"-CreationTime=" + currentTime.Format("2006-01-02T15:04:05Z"),
		"-XMPToolkit=",
		// "-VideoFrameRate=30",
		videoPath,
	}

	// Execute the exiftool command with the provided video path and arguments.
	cmd := exec.Command("exiftool", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to update video metadata: %s, %v", output, err)
	}

	return nil
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
