package storage

import (
	"bufio"
	"io/fs"
	"io/ioutil"
	"os"

	"github.com/valyala/fasthttp"
	"idaman.id/storage/pkg/app"
	"idaman.id/storage/pkg/config"
)

type StorageLocal struct {
	StorageDir string
}

func CreateStorageLocal() *StorageLocal {
	storage := StorageLocal{
		StorageDir: "storage/file",
	}
	return &storage
}

func (storage *StorageLocal) SaveFile(fileHeader FileDto) (result *FileEntity, err error) {
	var file FileEntity
	file.New(fileHeader)

	path := storage.StorageDir + "/" + file.UniqueId + "." + file.Extension

	appUrl := config.GetString("APP_URL")
	file.Path = path
	file.Url = appUrl + "/" + path

	err = fasthttp.SaveMultipartFile(fileHeader, path)

	isSaveFailed := err != nil
	if isSaveFailed {
		return nil, err
	}

	return &file, nil
}

func (storage *StorageLocal) RetrieveFile(file *FileEntity) (result BinaryFile, err error) {
	path := storage.StorageDir + "/" + file.UniqueId + "." + file.Extension
	osFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer osFile.Close()

	reader := bufio.NewReader(osFile)
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (storage *StorageLocal) DeleteFile(file *FileEntity) error {
	path := storage.StorageDir + "/" + file.UniqueId + "." + file.Extension
	err := os.Remove(path)

	switch err.(type) {
	case *fs.PathError:
		err = &app.NotFoundError{
			Message: app.STATUS_NOT_FOUND,
			Context: "File",
		}
	}

	return err
}
