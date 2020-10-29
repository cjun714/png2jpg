package main

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	if len(os.Args) != 3 {
		panic("usage: png2jpg 90 ./")
	}

	quality, e := strconv.Atoi(os.Args[1])
	if e != nil {
		panic(e)
	}

	dirPath := os.Args[2]
	files, e := ioutil.ReadDir(dirPath)
	if e != nil {
		panic(e)
	}

	totalStart := time.Now()

	var wg sync.WaitGroup
	for _, file := range files {
		fileName := file.Name()
		if !strings.HasSuffix(strings.ToLower(fileName), ".png") {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			srcPath := dirPath + "/" + fileName
			dstPath := srcPath[:len(srcPath)-4] + ".jpg"
			fmt.Println("saving:", dstPath)
			start := time.Now()
			e := toJPG(srcPath, dstPath, quality)
			if e != nil {
				panic(e)
			}
			fmt.Println("done: ", dstPath, time.Since(start))

		}()

	}
	wg.Wait()

	fmt.Println("total cost: ", time.Since(totalStart))
}

func toJPG(imgPath, jpgPath string, quality int) error {
	file, e := os.Open(imgPath)
	defer file.Close()
	if e != nil {
		return e
	}

	img, _, e := image.Decode(file)
	if e != nil {
		fmt.Println("decode file failed:", imgPath)
		return e
	}

	toFile, e := os.Create(jpgPath)
	defer toFile.Close()
	if e != nil {
		return e
	}

	return jpeg.Encode(toFile, img, &jpeg.Options{Quality: quality})
}
