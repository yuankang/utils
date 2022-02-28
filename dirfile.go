package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// 判断文件或者文件夹是否存在
func FileExist(s string) bool {
	_, err := os.Stat(s)
	return err == nil || os.IsExist(err)
}

// 判断文件或者文件夹是否存在, 可以指定是否创建
func DirExist(name string, isMake bool) error {
	var err error
	if FileExist(name) == false {
		if isMake {
			err = os.MkdirAll(name, 0755)
			if err != nil {
				log.Printf("create %s fail: %s", name, err.Error())
			} else {
				log.Printf("create %s succ", name)
			}
			return err
		} else {
			return fmt.Errorf("%s not exist", name)
		}
	}
	return err
}

func GetFileModTime(f string) (int64, error) {
	fi, err := os.Stat(f)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return fi.ModTime().Unix(), nil
}

func GetFileModTime0(f string) (int64, error) {
	/*
		fi, err := os.Stat(f)
		if err != nil {
			log.Println(err)
			return 0, err
		}
	*/

	// Sys()返回的是interface{}，所以需要类型断言，不同平台需要的类型不一样，linux上为*syscall.Stat_t
	//stat_t := fi.Sys().(*syscall.Stat_t)
	//log.Println(stat_t)

	// atime，ctime，mtime分别是访问时间，创建时间和修改时间，具体参见man 2 stat
	//log.Println(timespecToTime(stat_t.Atim))
	//log.Println(timespecToTime(stat_t.Ctim))
	//log.Println(timespecToTime(stat_t.Mtim).Unix())
	//return TimespecToTime(stat_t.Mtim).Unix(), nil
	return 0, nil
}

func GetFileSize(f string) (int64, error) {
	fi, err := os.Stat(f)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	/*
		fmt.Println("name:",fi.Name())
		fmt.Println("size:",fi.Size())
		fmt.Println("is dir:",fi.IsDir())
		fmt.Println("mode::",fi.Mode())
		fmt.Println("modTime:",fi.ModTime())
		name: water
		size: 403
		is dir: false
		mode:: -rw-r--r--
		modTime: 2018-05-06 18:52:07 +0800 CST
	*/
	return fi.Size(), nil
}

func ReadAllFile(file string) ([]byte, error) {
	fp, err := os.Open(file)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer fp.Close()

	d, err := ioutil.ReadAll(fp)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return d, nil
}

func SaveToDisk(name string, b []byte) (int, error) {
	f, err := os.Create(name)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer f.Close()

	n, err := f.Write(b)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return n, nil
}

func SaveToDiskAppend(name string, b []byte) error {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func CopyFile(dstName, srcName string) error {
	src, err := os.Open(srcName)
	if err != nil {
		log.Println(err)
		return err
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		log.Println(err)
		return err
	}
	return err
}

func IsDir(s string) bool {
	info, err := os.Stat(s)
	if err != nil {
		log.Println(err)
		return false
	}
	return info.IsDir()
}

type fileInfo struct {
	Name  string
	Mtime int64
}

// 递归获取当前目录下的所有文件
func GetAllFile(pathname string) ([]fileInfo, error) {
	result := []fileInfo{}

	fis, err := ioutil.ReadDir(pathname)
	if err != nil {
		log.Printf("read file directory fail，pathname=%v, err=%v \n", pathname, err)
		return result, err
	}

	//all directories' file
	for _, fi := range fis {
		fullname := pathname + "/" + fi.Name()
		// if directory, then call recursion; if file, then append to slice
		if fi.IsDir() {
			temp, err := GetAllFile(fullname)
			if err != nil {
				log.Printf("read file directory fail,fullname=%v, err=%v", fullname, err)
				return result, err
			}
			result = append(result, temp...)
		} else {
			fileStruct := fileInfo{Name: fullname, Mtime: fi.ModTime().Unix()}
			result = append(result, fileStruct)
		}
	}

	return result, nil
}
