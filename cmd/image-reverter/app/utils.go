package app

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"regexp"
)

var (
	re = regexp.MustCompile(`https?://[^\s]+?\.jpg`)
)

func createDir(name string) error {
	return os.Mkdir(name, os.FileMode(0777))
}

func removeDir(name string) error {
	return os.RemoveAll(name)
}

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, nil
	}
	return io.ReadAll(resp.Body)
}

func saveDataToFile(data []byte, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	r := bytes.NewReader(data)
	_, err = io.Copy(file, r)
	if err != nil {
		return err
	}
	return nil
}

func convertImage(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode: %s, with err: %w", path, err)
	}

	result := image.NewGray(img.Bounds())
	draw.Draw(result, result.Bounds(), img, img.Bounds().Min, draw.Src)

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, result, nil)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func getPathsFromUrl(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	matches := re.FindAllString(string(body), -1)

	return matches, nil
}
