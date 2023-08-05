package api

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/arjendevos/create-unique-video/functions"
	"github.com/gin-gonic/gin"
)

func HandleUpload(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse form"})
		return
	}

	files := form.File["videos"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No videos found in the request"})
		return
	}

	// Get the values of zoomFactor and trimBy from the form data
	zoomFactor, err := strconv.ParseFloat(c.PostForm("zoomFactor"), 64)
	if err != nil {
		zoomFactor = 1.0 // Default value if not provided or invalid
	}

	trimBy, err := strconv.ParseFloat(c.PostForm("trimBy"), 64)
	if err != nil {
		trimBy = 0.0 // Default value if not provided or invalid
	}

	// Create the "assets" directory if it doesn't exist
	if _, err := os.Stat("static/uploads"); os.IsNotExist(err) {
		err = os.Mkdir("static/uploads", 0755)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create 'assets' directory"})
			return
		}
	}

	editedVideos := []string{} // Store URLs of edited videos

	// Process and save each uploaded video
	for _, file := range files {
		// Save the uploaded video to the "assets" folder
		filePath := fmt.Sprintf("static/uploads/%s", file.Filename)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the video"})
			return
		}

		// Call EditVideo function with zoomFactor and trimBy
		outputPath, err := functions.EditVideo(filePath, zoomFactor, trimBy)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit the video"})
			return
		}

		editedVideos = append(editedVideos, strings.ReplaceAll(outputPath, "static/", ""))

	}

	c.JSON(http.StatusOK, gin.H{"message": "Videos uploaded and saved successfully", "editedVideos": editedVideos})
}
