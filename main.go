package main

import (
	"flag"
	"landlord/program/connection"
	"landlord/program/game"
	"landlord/program/game/msg"
	"landlord/program/game/player"
	"landlord/program/model"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/google/logger"
)

var addr = flag.String("addr", "localhost:8888", "http service address")

var upgrader = websocket.Upgrader{} // use default options
// func init() {
// 	log.SetReportCaller(true)

// 	//customize the formatter
// 	log.SetFormatter(&log.TextFormatter{
// 		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
// 			fileName := path.Base(frame.File) + ", line:" + strconv.Itoa(frame.Line)
// 			return frame.Function, fileName
// 		},
// 	})
// }

type TempUser struct {
	userID int
	sync.RWMutex
}

var user TempUser

func echo(w http.ResponseWriter, r *http.Request) {

	con, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("upgrade:", err)
		return
	}
	defer con.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	// 暂时用变量模拟用户登陆，后续从数据库读取用户信息，实例化用户，游戏过程中用redis来暂存游戏信息，用户推出后持久化到数据库
	user.Lock()
	user.userID++
	nowId := user.userID - 1
	user.Unlock()
	logger.Infof("玩家：%v 登陆游戏", strconv.Itoa(nowId))
	currUser := &model.User{
		Id:       nowId,
		NickName: "玩家" + strconv.Itoa(nowId),
		Avatar:   "no_avatar",
	}
	currPlayer := player.NewPlayer(currUser, connection.NewWebSocketConnection(con))

	player.SendMsgToPlayer(currPlayer, msg.MSG_TYPE_OF_LOGIN, "用户登陆")

	shang := currPlayer.User.Id / 3

	if currPlayer.User.Id%3 == 0 {
		currPlayer.CreateGame(game.GAME_TYPE_OF_DOUDOZHU, 10)
		logger.Info("玩家：" + strconv.Itoa(currPlayer.GetPlayerUser().Id) + "创建游戏：" + game.GetGameName(game.GAME_TYPE_OF_DOUDOZHU))
	} else {
		currPlayer.JoinGame(game.GAME_TYPE_OF_DOUDOZHU, shang)
		logger.Info("玩家：" + strconv.Itoa(currPlayer.GetPlayerUser().Id) + "加入游戏：" + game.GetGameName(game.GAME_TYPE_OF_DOUDOZHU) +
			strconv.Itoa(game.GAME_TYPE_OF_DOUDOZHU) + ":" + strconv.Itoa(shang))
	}

	wg.Add(1)
	//启动一个goroutine监听该客户端发来的消息
	go player.HandlerUserMsg(&wg, con, currPlayer)

	wg.Wait()
}

func home(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/pages/index.html", http.StatusMovedPermanently)
}

func main() {
	// const logPath = "example.log"
	// lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	var verbose = flag.Bool("verbose", false, "print info level logs to stdout")
	defer logger.Init("LoggerExample", *verbose, false, os.Stdout).Close()
	flag.Parse()

	// db, err := gorm.Open("mysql", "root:password@/games?charset=utf8&parseTime=True&loc=Local")
	// if err != nil {
	// 	logger.Fatal(err.Error())
	// }
	// db.AutoMigrate(&model.User{})

	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./views/static"))))
	http.Handle("/pages/", http.StripPrefix("/pages/", http.FileServer(http.Dir("./views/pages"))))
	logger.Infof("game start, listen %v", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		logger.Fatal(err.Error())
	}
}
