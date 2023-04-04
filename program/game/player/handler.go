package player

import (
	"fmt"
	"github.com/google/logger"
	"github.com/gorilla/websocket"
	"runtime/debug"
	"strconv"
	"sync"
)

func HandlerUserMsg(wg *sync.WaitGroup, con *websocket.Conn, currPlayer *Player) {
	defer wg.Done()
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("panic recover! p: %v", p)
			debug.PrintStack()
		}
	}()
	for {
		msgType, msg, err := con.ReadMessage()
		if err == nil {
			switch msgType {
			case websocket.TextMessage:
				//同桌用户交流，包含对话流程和出牌流程
				currPlayer.ResolveMsg(msg)
			case websocket.CloseMessage:
				logger.Info("玩家：" + strconv.Itoa(currPlayer.GetPlayerUser().Id) + "断开链接")
				break
				//离开桌子流程，后续包含断线保持，自动出牌
			default:

			}
		} else {
			break
		}
	}
}
