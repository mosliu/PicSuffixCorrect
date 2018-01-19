package main

import (
    "github.com/urfave/cli"
    "github.com/fatih/color"
    "time"
    "os"
    "encoding/binary"
    "path/filepath"
    "strings"
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
    app.Copyright = "(c) 2017 Dreamaker Studio"
    app.Usage = "Tools box for daily usage"
    app.ArgsUsage = "[args and such]"
    app.Action = func(c *cli.Context) error {
        var target string
        if c.NArg() > 0 {
            target = c.Args().Get(0)
        } else {
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
        suffix := judgePicFileType(target)
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

func judgePicFileType(path string) string {
    fi, err := os.Open(path)
    defer fi.Close()
    checkErr(err)
    const NBUF = 32
    var buf [NBUF]byte
    nr, err := fi.Read(buf[:])
    checkErr(err)
    if nr > 0 {
        //startcode := binary.LittleEndian.Uint16(buf[0:4])
        cutbuf := buf[0:8]
        fileHead := binary.BigEndian.Uint64(cutbuf)
        switch fileHead {
        //https://www.filesignatures.net/index.php?search=gif&mode=EXT
        case 0x89504E470D0A1A0A:
            //	89 50 4E 47 0D 0A 1A 0A
            return "png"
        }
        cutbuf[7] = 0
        cutbuf[6] = 0
        cutbuf[5] = 0
        cutbuf[4] = 0
        //四位的头
        fileHead = binary.BigEndian.Uint64(cutbuf)
        switch fileHead {
        case 0x4749463800000000:
            return "gif"
        case 0x3842505300000000:
            return "psd"
        case 0xFFD8FFE000000000, 0xFFD8FFE100000000, 0xFFD8FFE800000000:
            return "jpg"
        case 0xFFD8FFE200000000, 0xFFD8FFE300000000:
            return "jpg"
        case 0x49492A0000000000, 0x4D4D002A00000000, 0x4D4D002B00000000:
            return "tiff"
        }
        //三位的
        cutbuf[3] = 0
        fileHead = binary.BigEndian.Uint64(cutbuf)
        switch fileHead {
        case 0x4920490000000000:
            return "tiff"
        }
        cutbuf[2] = 0
        fileHead = binary.BigEndian.Uint64(cutbuf)
        switch fileHead {
        case 0x424D000000000000:
            return "BMP"
        }
    }
    return ""
}
