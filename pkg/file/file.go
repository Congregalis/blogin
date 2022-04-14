package file

import (
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
)

func GetSize(f multipart.File) (int, error) {
	b, err := ioutil.ReadAll(f)

	return len(b), err
}

func GetExt(fileName string) string {
	return path.Ext(fileName)
}

func ChectNotExist(src string) bool {
	_, err := os.Stat(src)

	return os.IsNotExist(err)
}

func CheckPermission(src string) bool {
	_, err := os.Stat(src)

	return os.IsPermission(err)
}

func IsNotExistMkDir(src string) error {
	if notExist := ChectNotExist(src); notExist {
		if err := mkDir(src); err != nil {
			return err
		}
	}

	return nil
}

func mkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)

	return err
}

func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return f, nil
}
