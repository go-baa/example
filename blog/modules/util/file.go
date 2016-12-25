package util

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// IsExist check path is exists, exist return true, not exist return false
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	// Check if error is "no such file or directory"
	if _, ok := err.(*os.PathError); ok {
		return false
	}
	return false
}

// ReadDir read path return files by os.FileInfo
func ReadDir(dirname string) ([]os.FileInfo, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	return list, nil
}

// MkdirAll check the path isexist or mkdir, and check writable
func MkdirAll(path string) error {
	var err error
	// check path exist or create
	if IsExist(path) == false {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}
	// check path writable
	if IsWritable(path) == false {
		return errors.New("path [" + path + "] is not writable!")
	}
	return nil
}

// GetExt get file extenstion, not contains dot
func GetExt(file string) string {
	f := filepath.Ext(file)
	if f == "" {
		return f
	}
	return f[1:]
}

// CopyFile copy source file to destination file
func CopyFile(s, d string) error {
	// first check the source be link
	linfo, err := os.Readlink(s)
	if err == nil || len(linfo) > 0 {
		// if source is link, create link to destination
		return os.Symlink(linfo, d)
	}
	// normal, create file
	sf, err := os.Open(s)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := os.Create(d)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	return err
}

// ReadFile read file content
func ReadFile(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

// Glob 查找指定目录下的文件列表
func Glob(base, pattern string, abs bool) ([]string, error) {
	base, err := filepath.Abs(base)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(base)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		base = filepath.Dir(base)
	}

	return globRecursion(base, base, pattern, abs)
}

// globRecursion 递归查找
func globRecursion(root, base, pattern string, abs bool) ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(base, pattern))
	if err != nil {
		return nil, err
	}

	var items []string
	for i := range matches {
		item := matches[i]
		if !abs {
			rel, err := filepath.Rel(root, item)
			if err != nil {
				return nil, err
			}
			item = rel
		}
		items = append(items, item)
	}

	file, err := os.Open(base)
	if err != nil {
		return nil, err
	}
	files, err := file.Readdir(-1)
	if err != nil {
		return nil, err
	}
	for i := range files {
		if files[i].IsDir() {
			child, err := globRecursion(root, filepath.Join(base, files[i].Name()), pattern, abs)
			if err != nil {
				return nil, err
			}
			items = append(items, child...)
		}
	}

	return items, nil
}
