package ansi

import "fmt"

// TODO add support for no ANSI
const Reset = csi + "0" + sgr

const csi = "\x1b["
const sgr = "m"
const foreground = "38"
const trueColor = "2"

type RGB struct {
    R uint8
    G uint8
    B uint8
}

func (c RGB) ToString() string {
    format := csi + foreground + ";" + trueColor + ";%d;%d;%d" + sgr

    return fmt.Sprintf(format, c.R, c.G, c.B)
}
