package functions

import (
	"fmt"
	"os/exec"
	"time"
)

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
