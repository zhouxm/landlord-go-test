package main

import (
	"flag"
	"landlord/program/connection"
	"landlord/program/game"
	"landlord/program/game/msg"
	"landlord/program/game/player"
	"landlord/program/model"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
)

var addr = flag.String("addr", "localhost:8888", "http service address")

var upgrader = websocket.Upgrader{} // use default options

type TempUser struct {
	userID int
	sync.RWMutex
}

var user TempUser

func echo(w http.ResponseWriter, r *http.Request) {

	con, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Error("upgrade:", err)
		return
	}
	defer con.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	//暂时用变量模拟用户登陆，后续从数据库读取用户信息，实例化用户，游戏过程中用redis来暂存游戏信息，用户推出后持久化到数据库
	user.Lock()
	user.userID++
	nowId := user.userID - 1
	user.Unlock()
	logrus.Info("玩家：" + strconv.Itoa(nowId) + "登陆游戏")
	currUser := &model.User{
		Id:       nowId,
		NickName: "玩家" + strconv.Itoa(nowId),
		Avatar:   "no_avatar",
	}
	currPlayer := player.NewPlayer(currUser, connection.NewWebSocketConnection(con))

	player.SendMsgToPlayer(currPlayer, msg.MSG_TYPE_OF_LOGIN, "用户登陆")

	shang := currPlayer.User.Id / 3

	if currPlayer.User.Id%3 == 0 {
		currPlayer.CreateGame(game.TypeOfDoudozhu, 10)
		logrus.Info("玩家：" + strconv.Itoa(currPlayer.GetPlayerUser().Id) + "创建游戏：" + game.GetGameName(game.TypeOfDoudozhu))
	} else {
		currPlayer.JoinGame(game.TypeOfDoudozhu, shang)
		logrus.Info("玩家：" + strconv.Itoa(currPlayer.GetPlayerUser().Id) + "加入游戏：" + game.GetGameName(game.TypeOfDoudozhu) +
			strconv.Itoa(game.TypeOfDoudozhu) + ":" + strconv.Itoa(shang))
	}

	wg.Add(1)
	//启动一个goroutine监听该客户端发来的消息
	go player.HandlerUserMsg(&wg, con, currPlayer)

	wg.Wait()
}

func home(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/pages/index.html", 301)
}

func main() {
	flag.Parse()

	db, err := gorm.Open("mysql", "root:password@tcp(127.0.0.1:3306)/demo?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		logrus.Fatal(err.Error())
	}
	db.AutoMigrate(&model.User{})

	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./views/static"))))
	http.Handle("/pages/", http.StripPrefix("/pages/", http.FileServer(http.Dir("./views/pages"))))
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		logrus.Fatal(err.Error())
	}
}
