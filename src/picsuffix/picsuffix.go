package picsuffix

import (
    "os"
    "encoding/binary"
)

// Used for get a image file's suffix
// return a suffix string in lower case,
// example: "png" or "jpg"
func JudgePicFileType(path string) string {
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
            return "bmp"
        }
    }
    return ""
}