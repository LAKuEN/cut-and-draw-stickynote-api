package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/LAKuEN/detect-stickynotes"
	"github.com/LAKuEN/zip-files"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gocv.io/x/gocv"
)

// ImgFile はファイル名とgocv.Matの画像データを包含する構造体です。
type ImgFile struct {
	name string
	mat  gocv.Mat
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", getHandler)
	e.POST("/", postHandler)

	port := ":" + os.Getenv("PORT")
	e.Start(port)
}

func getHandler(c echo.Context) error {
	message := `# How to call the API
$ curl -X POST {URL} -F "file=@{image file path}"`
	return c.String(http.StatusOK, message+"\n")
}

func postHandler(c echo.Context) error {
	imgFile, err := context2mat(c)
	defer imgFile.mat.Close()
	if imgFile.mat.Empty() {
		return c.String(http.StatusBadRequest,
			"It is not possible to process a file other than the image file\n")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	dstDirPath := filepath.Join(cwd, imgFile.name)
	if _, err := os.Stat(dstDirPath); os.IsNotExist(err) {
		os.Mkdir(dstDirPath, 0777)
	}

	stickies, err := stickynote.CutNDraw(imgFile.mat)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	gocv.IMWrite(filepath.Join(dstDirPath, "0_drawed.jpeg"), stickies.DrawedImg)
	for idx, cropped := range stickies.CroppedImgs {
		gocv.IMWrite(filepath.Join(dstDirPath, fmt.Sprintf("%v.jpeg", idx+1)),
			cropped)
	}

	zipFilePath, err := zipfiles.InDir(dstDirPath)
	if err != nil {
		return c.String(http.StatusInternalServerError,
			"Failed to create zip file\n")
	}
	zipFileBytes, err := ioutil.ReadFile(zipFilePath)
	if err != nil {
		return c.String(http.StatusInternalServerError,
			"Failed to convert zip file to bytearray\n")
	}

	return c.Blob(http.StatusOK, "application/zip", zipFileBytes)
}

// context2mat はecho.Contextから画像ファイル名と画像データを
// ImgFile構造体として取り出します。
func context2mat(c echo.Context) (ImgFile, error) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return ImgFile{}, err
	}
	file, err := fileHeader.Open()
	if err != nil {
		return ImgFile{}, err
	}
	defer file.Close()
	imgBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return ImgFile{}, err
	}
	img, err := gocv.IMDecode(imgBytes, gocv.IMReadColor)
	if err != nil {
		return ImgFile{}, err
	}

	return ImgFile{name: fileHeader.Filename, mat: img}, nil
}
