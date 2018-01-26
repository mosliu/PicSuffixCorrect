package src

import (
    "os"
    "path/filepath"
    "strings"
    "time"

    "github.com/fatih/color"
    "github.com/urfave/cli"

    "picsuffix"
)

func main() {
    app := cli.NewApp()
    app.Name = "Moses ToolBox"
    app.Version = "0.0.1.2"
    app.Compiled = time.Now()
    app.Authors = []cli.Author{
        {
            Name:  "Moses",
            Email: "mogo@liuxuan.net",
        },
    }
    app.Copyright = "(c) 2018 Dreamaker Studio"
    app.Usage = "Tools for check if the image file's suffix is right.Can be applied to a directory"
    app.UsageText = "Example: picsuffix ./aaa.jpg or picsuffix ./jpg/"
    app.ArgsUsage = "<filename>"
    app.Action = func(c *cli.Context) error {
        var target string
        if c.NArg() > 0 {
            target = c.Args().Get(0)
        } else {
            color.Red("无输入，程序退出")
            target = "./"
        }
        fi, err := os.Stat(target)
        if err != nil {
            color.Red("输入目的路径 不存在")
        } else {
            if fi.IsDir() {
                color.Red("输入的是目录，遍历目录")
                filepath.Walk(target, dirwalk)

            } else {
                judgeAndRename(target)
            }
        }
        return nil
    }
    defer func() {
        if e := recover(); e != nil {
            color.Red("Panicing,error occured: %s\r\n", e)
        }
    }()

    app.Run(os.Args)
}

func dirwalk(path string, f os.FileInfo, err error) error {
    if f == nil {
        return err
    }
    //不遍历子目录的目录
    if f.IsDir() {
        //color.Green("遍历子目录："+path)
        //filepath.Walk(path, dirwalk)
        return nil
    }
    judgeAndRename(path)
    //println(path)
    return nil
}

func judgeAndRename(target string) {
    _, filename := filepath.Split(target)
    oldsuffix := strings.ToLower(filepath.Ext(target))
    //color.Green(filename + ",oldsuffix:" + oldsuffix)
    switch oldsuffix {
    case ".png", ".jpg", ".bmp", ".tiff", ".tif", ".gif", ".psd":
        suffix := picsuffix.JudgePicFileType(target)
        if suffix != "" {
            suffix = "." + suffix
        } else {
            color.Red("没法判断文件类型")
        }
        if strings.EqualFold(oldsuffix, suffix) {
            color.Green("判定原后缀正确：" + suffix)
        } else {
            color.Magenta(filename + "  :判定后缀应该是：" + suffix + ",文件改名")
            err := os.Rename(target, target+suffix)
            checkErr(err)
        }
    default:
        color.Green("Skip file:" + filename)
    }
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}
