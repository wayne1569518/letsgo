/*=============================================================================
#     FileName: proc2011.go
#         Desc: server base 
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-17 18:22:56
#      History:
=============================================================================*/
package protos

import (
    "fmt"
    . "github.com/sunminghong/letsgo/net"
    . "github.com/sunminghong/letsgo/log"
)


//func init() {
//Handlers[2011] = Process2011
//}

func Process2011(msgReader LGIMessageReader,c LGIConnection,fromCid int) {
    LGTrace("process 2011 is called")

    md := msgReader.ReadString()

    fmt.Println()
    fmt.Println(md)
    fmt.Print("you> ")

}

