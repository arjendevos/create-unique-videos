# Create Unique Videos

My attempt to create unique videos especially for TikTok.

## How to use

1. Install `ffmpeg` and `ffprobe` on your system.
2. Clone this repository.
3. Run `go run main.go` with `-input=` flags. Input can both be a folder and file.
4. Enjoy!

## Implemented

- [x] Remove metadata
- [x] Remove audio
- [x] Zoom
- [x] Add fake metadata
- [x] Change image width/height to tiktok's standard

### TikTok In-App Metadata

- [x] Create Date: 2023:08:05 07:35:47
- [x] Modify Date: 2023:08:05 07:35:47
- [x] Track Create Date: 2023:08:05 07:35:47
- [x] Track Modify Date: 2023:08:05 07:35:47
- [x] Media Create Date: 2023:08:05 07:35:47
- [x] Media Modify Date: 2023:08:05 07:35:47
- [x] Handler Type: Metadata Tags
- [x] Handler Description: DataHandler
- [x] Audio Format: mp4a
- [x] Audio Bits Per Sample: 16
- [x] Audio Sample Rate: 44100
- [x] Compressor ID: hvc1
- [x] Vendor ID: FFmpeg
- [ ] Handler Vendor ID: (appears to be not writable)
- [x] Balance: 0
- [x] Purchase File Format: mp4a
- [x] Layout Flags: Stereo
- [x] Audio Channels: 2
- [x] Writer Type: -1
- [x] Copyright: cb2e8b799d13edf6c2e5b21baa16c737
- [x] Mdat Pos: 28
- [x] Minor Version: 512
- [x] Encoder: Lavf57.71.100
- [x] Software: {"publicMode":"1","TEEditor":"2","isFastImport":"0","transType":"2","te_is_reencode":"1","source":"2"}
- [x] Source: 2
- [x] Hw: 1
- [x] Major Brand: qt
- [x] Compatible Brands: qt
- [x] Te Is Reencode: 1
- [x] Creation Time: 2023-08-05T07:35:47Z
- [ ] Video Frame Rate: 30
- [ ] Next Track ID: 3
