/*=============================================================================
#     FileName: chatclient.go
#         Desc: chat client
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-13 17:48:50
#      History:
=============================================================================*/
package main

import (
    "fmt"
    "strings"
    "time"
    "strconv"
    "flag"

    "./protos"
    . "github.com/sunminghong/letsgo/net"
    . "github.com/sunminghong/letsgo/helper"
    . "github.com/sunminghong/letsgo/log"
	"github.com/sbinet/liner"
)

var endian = LGBigEndian

func tabCompleter(line string) []string {
	opts := make([]string, 0)

	if strings.HasPrefix(line, "/") {
		filters := []string{
			"/conn ",
			"/change ",
			"/quit",
			"/reg ",
			"/rereg ",
		}

		for _, cmd := range filters {
			if strings.HasPrefix(cmd, line) {
				opts = append(opts, cmd)
			}
		}
	}

	return opts
}

// clientsender(): read from stdin and send it via network
func clientsender(cid *int,client *LGConnectionPool) {
	term := liner.NewLiner()
	fmt.Println("chat client")
    defer term.Close()

	term.SetCompleter(tabCompleter)
    for {
        if (*cid)==0 {
            fmt.Print("you no connect anyone server,please input conn cmd,\n")
        }

		input, e := term.Prompt("> ")
		if e != nil {
			break
		}

        //cmd := string(input[:len(input)-1])
        cmd := string(input)
        var text string

        if cmd[0] == '/' {
            cmds := strings.Split(cmd," ")
            switch cmds[0]{
            case "/","/conn":
                ///conn s1 :12001 0
                var name,addr string
                var endian int

                if cmds[0] == "/" {
                    name= "s1"
                    addr=":12000"
                    endian=1
                } else {

                    if len(cmds)>2 {
                        name = cmds[1]
                        addr = cmds[2]
                    }else {
                        name = "c_" + strconv.Itoa(*cid)
                        addr = cmds[1]
                    }

                    p := client.Connections.GetByName(name)
                    if p != nil {
                        fmt.Println(name," is exists !")
                        continue
                    }
                }

                if len(cmds)>3 {
                    endian,_= strconv.Atoi(cmds[3])
                }

                LGDebug("connect to server use endian:",endian)
                datagram := LGNewDatagram(endian)
                go client.Start(name,addr,datagram)


                input, e = term.Prompt("please input your name: ")
                if e != nil {
                    break
                }

                //cmd := string(input[:len(input)-1])
                cmd = string(input)

                for true {
                    b := client.Connections.GetByName(name)
                    if b!=nil{
                        change(cid,client,name)
                        break
                    }
                    time.Sleep(2*1e3)
                }

                text = string(input)

            case "/setlog":
                if lv,err := strconv.Atoi(cmds[1]); err == nil {
                    LGSetLevel(lv)
                }
            case "/change":
                name := cmds[1]
                change(cid,client,name)

            case "/quit":
                text = "/quit"

            case "/11012":
                c := client.Connections.Get(*cid)
                msg := protos.NewMessageWriter(c)
                msg.SetCode(1101,0)
                msg.WriteUint(2,0)
                msg.WriteUint(1,0)
                msg.WriteUint(1,0)
                //msg.WriteUints(1,0,0)
                c.SendMessage(0,msg)
                continue


            case "/11011":
                c := client.Connections.Get(*cid)
                msg := protos.NewMessageWriter(c)
                msg.SetCode(1101,0)
                msg.WriteUint(1,0)
                msg.WriteUint(1,0)
                msg.WriteUint(1,0)
                //msg.WriteUints(1,0,0)
                c.SendMessage(0,msg)
                continue

            case "/1001":
                c := client.Connections.Get(*cid)
                msg := protos.NewMessageWriter(c)
                msg.SetCode(1001,0)
                msg.WriteUint(1,0)
                c.SendMessage(0,msg)
                continue

            default:
                //text = string(input[:len(input)-1])
                text = string(input)
            }
        } else {
            //text = string(input[:len(input)-1])
            text = string(input)
        }

        c := client.Connections.Get(*cid)
        msg := protos.NewMessageWriter(c)
        msg.SetCode(1011,0)
        msg.WriteString(text,0)

        LGTrace("has %v clients,text:%s",client.Connections.Len(),text)
        c.SendMessage(0,msg)
    }
}

func change(cid *int,client *LGConnectionPool,name string,) {
    b:= client.Connections.GetByName(name)
    if b!=nil{
        _cid := b.GetTransport().Cid
        *cid = _cid
        fmt.Println("current connection change:")
    }

    for c,p:=range client.Connections.All() {
        if p.GetName() != name {
            fmt.Println(" ",c,p.GetName())
        } else {
            fmt.Println("*",c,p.GetName())
        }
    }
}

var (
    loglevel = flag.Int("loglevel",0,"log level")
)

func main() {
    flag.Parse()

    LGSetLevel(*loglevel)

    datagram := LGNewDatagram(protos.Endian)

    cid := 0
    client := LGNewConnectionPool(protos.NewConnection, datagram)
    go clientsender(&cid,client)

    //client.Start("", 4444)

    quit := make(chan bool)
    <-quit
    //running :=1
    //for running==1 {
    //    time.Sleep(3*time.Second)
    //}
}

