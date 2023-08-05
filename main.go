package main

import (
	"fmt"
	"net/http"

	"github.com/arjendevos/create-unique-video/api"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Serve static files from the "static" directory (including the HTML file)
	r.StaticFS("/", http.Dir("static"))

	r.POST("/upload", api.HandleUpload)

	port := 8080
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Server listening on http://localhost%s\n", addr)
	err := r.Run(addr)
	if err != nil {
		panic(err)
	}
}

// func main() {
// 	var input string
// 	var zoomFactor float64
// 	var trimBy float64

// 	// Define command-line flags
// 	flag.StringVar(&input, "input", "", "Input containing video or videos")
// 	flag.Float64Var(&zoomFactor, "zoom", 1.01, "Zoom ratio")
// 	flag.Float64Var(&trimBy, "trim", 0.1, "Trim by seconds")
// 	flag.Parse()

// 	// Check if input and output folders are provided
// 	if input == "" || zoomFactor < 0 || trimBy < 0 {
// 		fmt.Println("Usage: go run main.go -input <input> -zoom <zoom> -trim <trim>")
// 		return
// 	}

// 	// Check if input folder exists
// 	fileInfo, err := os.Stat(input)
// 	if err != nil || os.IsNotExist(err) {
// 		fmt.Println("Input folder does not exist")
// 		return
// 	}

// 	isInputDir := fileInfo.IsDir()
// 	input, err = filepath.Abs(input)
// 	if err != nil {
// 		fmt.Println("Error getting absolute path of input folder:", err)
// 		return
// 	}

// 	if isInputDir {
// 		// Traverse the input folder for video files
// 		if err := filepath.Walk(input, func(path string, info os.FileInfo, err error) error {
// 			if err != nil {
// 				fmt.Println("Error accessing file:", err)
// 				return nil
// 			}

// 			if info.IsDir() {
// 				return nil
// 			}

// 			if err := functions.EditVideo(path, zoomFactor, trimBy); err != nil {
// 				return err
// 			}

// 			return nil
// 		}); err != nil {
// 			fmt.Println("Error traversing input folder:", err)
// 		}

// 		return
// 	}

// 	if err := functions.EditVideo(input, zoomFactor, trimBy); err != nil {
// 		fmt.Println("Error editing video:", err)
// 	}

// }
