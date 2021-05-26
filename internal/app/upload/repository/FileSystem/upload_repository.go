package FileSystem

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/nfnt/resize"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/upload"
)

type UploadRepository struct {
}

func NewUploadRepository() upload.UploadRepository {
	return &UploadRepository{}
}

func (ur *UploadRepository) InsertPhoto(fileHeader *multipart.FileHeader, photoPath string) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	str, err := os.Getwd()
	if err != nil {
		return "", err
	}

	os.Chdir(photoPath)

	photoID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	extension := filepath.Ext(fileHeader.Filename)

	newFile, err := os.OpenFile(photoID.String()+extension, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer newFile.Close()

	os.Chdir(str)

	_, err = io.Copy(newFile, file)
	if err != nil {
		_ = os.Remove(photoID.String() + extension)
		return "", err
	}

	photo := "/" + photoPath + photoID.String() + extension
	return photo, nil
}

func (ur *UploadRepository) InsertPhotos(filesHeaders []*multipart.FileHeader, photoPath string) ([]string, error) {
	imgUrls := make(map[string][]string)

	for i := range filesHeaders {
		url, err := ur.InsertPhoto(filesHeaders[i], photoPath)
		if err != nil {
			return nil, err
		}

		imgUrls["img"] = append(imgUrls["img"], url)
	}

	return imgUrls["img"], nil
}

func (ur *UploadRepository) RemovePhoto(imgUrl string) error {
	if imgUrl == "" {
		return nil
	}

	origWd, _ := os.Getwd()
	err := os.Remove(origWd + imgUrl)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UploadRepository) RemovePhotos(imgUrls []string) error {
	if len(imgUrls) == 0 {
		return nil
	}

	for _, photo := range imgUrls {
		err := ur.RemovePhoto(photo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ur *UploadRepository) ResizePhoto(imgUrl string) error {
	imgIn, err := os.Open(imgUrl)
	if err != nil {
		return err
	}
	defer imgIn.Close()

	imgForDecode, err := os.Open(imgUrl)
	if err != nil {
		return err
	}
	defer imgForDecode.Close()

	var imgJpg image.Image

	switch filepath.Ext(imgUrl) {
	case ".jpeg", ".jpg":
		imgJpg, err = jpeg.Decode(imgForDecode)
		if err != nil {
			return err
		}

	case ".png":
		imgJpg, err = png.Decode(imgForDecode)
		if err != nil {
			return err
		}

	default:
		return nil
	}

	width, height := imgJpg.Bounds().Dx(), imgJpg.Bounds().Dy()
	imgJpg = resize.Resize(uint(width), uint(height), imgJpg, resize.Bicubic)

	imgOut, err := os.Create(imgUrl)
	if err != nil {
		return err
	}
	switch filepath.Ext(imgUrl) {
	case ".jpeg":
	case ".jpg":
		err = jpeg.Encode(imgOut, imgJpg, nil)
		if err != nil {
			return err
		}

	case ".png":
		err = png.Encode(imgOut, imgJpg)
		if err != nil {
			return err
		}

	default:
		return nil
	}
	defer imgOut.Close()

	return nil
}

func (ur *UploadRepository) ResizePhotos(imgUrls []string) error {
	if len(imgUrls) == 0 {
		return nil
	}

	for _, photo := range imgUrls {
		err := ur.ResizePhoto(photo)
		if err != nil {
			return err
		}
	}

	return nil
}
