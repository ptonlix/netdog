package logo

import (
	"fmt"
	"io"

	"github.com/ptonlix/spokenai/pkg/color"
)

// see https://patorjk.com/software/taag/#p=testall&f=Graffiti&t=go-gin-api
const ui = `
 __    _  _______  _______  ______   _______  _______ 
|  |  | ||       ||       ||      | |       ||       |
|   |_| ||    ___||_     _||  _    ||   _   ||    ___|
|       ||   |___   |   |  | | |   ||  | |  ||   | __ 
|  _    ||    ___|  |   |  | |_|   ||  |_|  ||   ||  |
| | |   ||   |___   |   |  |       ||       ||   |_| |
|_|  |__||_______|  |___|  |______| |_______||_______|

`

const introduce = `
----------------------------------------------------------------------
Welcome to use 
这是一个监测网络稳定性的工具

作者/author: Baird
联系作者/Contact the author: wechat:cfd0917
`

func PrintLogo(w io.Writer) {
	fmt.Fprint(w, color.Blue(ui))
}

func PrintIntroduce(w io.Writer) {
	fmt.Fprint(w, color.Yellow(introduce))
}
