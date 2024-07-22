package app

import (
	"fmt"
	"path/filepath"
)

type SyncImageConvertor struct{}

func (s *SyncImageConvertor) Convert(url string) (string, error) {
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

	imagesFromSite, err := getPathsFromUrl(url)
	if err != nil {
		return "", fmt.Errorf("failed to get images from site %w", err)
	}

	downloadedImage := make([]string, 0)
	for _, linkToImage := range imagesFromSite {
		byteImage, err := downloadFile(linkToImage)
		if err != nil {
			return "", fmt.Errorf("failed to download image, err: %w", err)
		}
		pathToColor := fmt.Sprintf("%s/%s", folderForColorImages, filepath.Base(linkToImage))
		err = saveDataToFile(byteImage, pathToColor)
		if err != nil {
			return "", fmt.Errorf("failed to save color image, err: %w", err)
		}
		downloadedImage = append(downloadedImage, pathToColor)
	}

	convertedImage := make([]string, 0)
	for _, image := range downloadedImage {
		conv, err := convertImage(image)
		if err != nil {
			return "", fmt.Errorf("failed to convert image, err: %w", err)
		}
		pathToGray := fmt.Sprintf("%s/%s", folderForGrayImages, filepath.Base(image))
		err = saveDataToFile(conv, pathToGray)
		if err != nil {
			return "", fmt.Errorf("failed to save gray image, err: %w", err)
		}
		convertedImage = append(convertedImage, pathToGray)
	}
	res := fmt.Sprintf("was converted %d images", len(convertedImage))
	return res, nil
}
