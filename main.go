package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/LAKuEN/detect-stickynotes"
	"gocv.io/x/gocv"
)

func main() {
	filePath := flag.String("filepath", "", "target file name")
	flag.Parse()

	substrFilePath := strings.Split(*filePath, ".")
	dstDirPath := strings.Join(substrFilePath[0:len(substrFilePath)-1], ".")
	if _, err := os.Stat(dstDirPath); os.IsNotExist(err) {
		os.Mkdir(dstDirPath, 0777)
	}
	orig := gocv.IMRead(*filePath, gocv.IMReadColor)
	if orig.Empty() {
		onError(fmt.Sprintf("cannot read file as image: %s", *filePath))
	}

	stickies, err := stickynote.CutNDraw(orig)
	if err != nil {
		onError(err.Error())
	}
	gocv.IMWrite(filepath.Join(dstDirPath, "0_drawed.jpeg"), stickies.DrawedImg)
	for idx, cropped := range stickies.CroppedImgs {
		gocv.IMWrite(filepath.Join(dstDirPath, fmt.Sprintf("%v.jpeg", idx+1)),
			cropped)
	}

	fmt.Printf("saved as %v\n", dstDirPath)
}

// FIXME 共通関数モジュールに切り出した方がいいかも？
//       画像処理系の処理とか、処理の区分が見えているなら、区分毎にモジュール分割してもいいかも
// onError はエラー発生時の共通処理です。
func onError(errMsg string) {
	fmt.Fprintln(os.Stderr, errMsg)
	os.Exit(1)
}
