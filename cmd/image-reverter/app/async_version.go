package app

import (
	"fmt"
	"path/filepath"
	"sync"
)

const (
	queueSize    = 10
	workersCount = 10
)

type AsyncImageConvertor struct {
}

func (a *AsyncImageConvertor) Convert(url string) (string, error) {
	// prepare folders
	folderForColorImages := "./tmp/color"
	err := createDir(folderForColorImages)
	if err != nil {
		return "", fmt.Errorf("failed to create folder for color image, err: %w", err)
	}
	defer removeDir(folderForColorImages)

	folderForGrayImages := "./tmp/gray"
	err = createDir(folderForGrayImages)
	if err != nil {
		return "", fmt.Errorf("failed to create folder for gray image, err: %w", err)
	}
	defer removeDir(folderForGrayImages)

	converter := converter{
		folderForColorImages: folderForColorImages,
		folderForGrayImages:  folderForGrayImages,
		resultCh:             make(chan string),
		errCh:                make(chan error),
	}

	// run convert
	go converter.Run(url)

	select {
	case res := <-converter.resultCh:
		return res, nil
	case err := <-converter.errCh:
		return "", err
	}
}

type converter struct {
	folderForColorImages string
	folderForGrayImages  string
	resultCh             chan string
	errCh                chan error
}

type img struct {
	name      string
	url       string
	colorPath string
	grayPath  string
}

func (c *converter) Run(url string) {

	imagesFromSite, err := getPathsFromUrl(url)
	if err != nil {
		c.errCh <- err
		return
	}
	linksChan := make(chan *img, len(imagesFromSite))
	for _, link := range imagesFromSite {
		linksChan <- &img{
			name: filepath.Base(link),
			url:  link,
		}
	}
	close(linksChan)

	downloadedImgs := c.downloadAndSaveImage(linksChan)

	convertedImgs := c.convertAndSave(downloadedImgs)

	c.collectResult(convertedImgs)
}

func (c *converter) downloadAndSaveImage(imgs chan *img) chan *img {
	savedImagePath := make(chan *img, queueSize)

	worker := func(wg *sync.WaitGroup, in, out chan *img) {
		defer wg.Done()
		for img := range in {
			byteImage, err := downloadFile(img.url)
			if err != nil {
				c.errCh <- fmt.Errorf("failed to download image, err: %w", err)
				return
			}
			pathToColor := fmt.Sprintf("%s/%s", c.folderForColorImages, img.name)
			err = saveDataToFile(byteImage, pathToColor)
			if err != nil {
				c.errCh <- fmt.Errorf("failed to save color image, err: %w", err)
				return
			}
			img.colorPath = pathToColor

			out <- img
		}
	}

	wg := &sync.WaitGroup{}

	go func() {
		defer close(savedImagePath)
		for i := 0; i < workersCount; i++ {
			wg.Add(1)
			go worker(wg, imgs, savedImagePath)
		}
	}()

	return savedImagePath
}

func (c *converter) convertAndSave(in chan *img) chan *img {
	convertedImgs := make(chan *img, queueSize)

	worker := func(wg *sync.WaitGroup, in, out chan *img) {
		defer wg.Done()
		for image := range in {
			conv, err := convertImage(image.colorPath)
			if err != nil {
				c.errCh <- fmt.Errorf("failed to convert image, err: %w", err)
				return
			}
			pathToGray := fmt.Sprintf("%s/%s", c.folderForGrayImages, image.name)
			err = saveDataToFile(conv, pathToGray)
			if err != nil {
				c.errCh <- fmt.Errorf("failed to save gray image, err: %w", err)
				return
			}

			image.grayPath = pathToGray
			out <- image
		}
	}

	go func() {
		defer close(convertedImgs)

		wg := &sync.WaitGroup{}

		for i := 0; i < workersCount; i++ {
			wg.Add(1)
			go worker(wg, in, convertedImgs)
		}
	}()

	return convertedImgs
}

func (c *converter) collectResult(in chan *img) {
	i := 0
	for img := range in {
		if img.grayPath != "" {
			i++
		}
	}
	c.resultCh <- fmt.Sprintf("was converted %d images", i)
}
