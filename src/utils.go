package main

import (
    "os"
    "path/filepath"
)

//文件是否存在
func isFileExist(path string) bool {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return false
    }
    return true
}

func isDirExists(path string) bool {
    fi, err := os.Stat(path)

    if err != nil {
        return os.IsExist(err)
    } else {
        return fi.IsDir()
    }

    panic("not reached")
}


//获得主程序的路径
func getMainExePath() string {
    dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
    return dir
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}