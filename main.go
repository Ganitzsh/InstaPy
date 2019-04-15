package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-cmd/cmd"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type tmplcontent struct {
	creds
	res
}

type creds struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type PotencyMode string

const (
	PotencyModePositive PotencyMode = "positive"
	PotencyModeNegative PotencyMode = "negative"
)

type res struct {
	Hashtags     []string    `json:"hashtags" form:"hashtags"`
	Comments     []string    `json:"comments" form:"comments"`
	TotalLikes   int         `json:"total_likes" forn:"total_likes"`
	Potency      PotencyMode `json:"potency" form:"potency"`
	PerUser      int         `json:"per_user" form:"per_user"`
	MaxFollowers int         `json:"max_followers"`
	MinFollowers int         `json:"min_followers"`
	MaxFollowing int         `json:"max_following"`
	MinFollowing int         `json:"min_following"`
}

func (r res) HastagStr() string {
	return strings.Join(r.Hashtags, "\n")
}

func (r res) CommentStr() string {
	return strings.Join(r.Comments, "\n")
}

type saveReq struct {
	Hashtags     string `json:"hashtags" form:"hashtags"`
	Comments     string `json:"comments" form:"comments"`
	TotalLikes   int    `json:"total_likes" form:"total_likes"`
	PerUser      int    `json:"per_user" form:"per_user"`
	Username     string `json:"username" form:"username"`
	Password     string `json:"password" form:"password"`
	Potency      string `json:"potency" form:"potency"`
	MaxFollowers int    `json:"max_followers" form:"max_followers"`
	MinFollowers int    `json:"min_followers" form:"min_followers"`
	MaxFollowing int    `json:"max_following" form:"max_following"`
	MinFollowing int    `json:"min_following" form:"min_following"`
}

var command *cmd.Cmd

func runBot(r *res, conn *websocket.Conn, clients map[string]*websocket.Conn) {
	// Start a long-running process, capture stdout and stderr
	if command != nil {
		return
	}
	command = cmd.NewCmd("python3", "main.py")
	statusChan := command.Start() // non-blocking

	ticker := time.NewTicker(1 * time.Second)

	// Print last line of stdout every 2s
	go func() {
		var lastString string
		for range ticker.C {
			if command == nil {
				logrus.Info("I'm out")
				return
			}
			status := command.Status()
			if len(status.Stderr) > 0 {
				str := status.Stderr[len(status.Stderr)-1]
				if lastString != str {
					lastString = str
					logrus.Info(lastString)
					if clients != nil {
						for id, c := range clients {
							if c != nil {
								c.WriteJSON(wsCmdRegister{
									ID:      id,
									Running: true,
									Output:  strings.Join(status.Stderr, "<br/>"),
								})
							}
						}
					}
				}
			}
		}
	}()

	// Stop command after 1 hour
	go func() {
		<-time.After(1 * time.Hour)
		command.Stop()
		command = nil
	}()

	// Check if command is done
	select {
	case <-statusChan:
		logrus.Info("Done!")
		command = nil
		return
	default:
		logrus.Info("Still running")
	}

	if clients != nil {
		for _, client := range clients {
			if client != nil {
				client.WriteJSON(wsCmdRegister{
					Running: command != nil,
				})
			}
		}
	}

	<-statusChan
	command = nil
	if clients != nil {
		for _, client := range clients {
			if client != nil {
				client.WriteJSON(wsCmdRegister{
					Running: command != nil,
				})
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type wsCmdRegister struct {
	ID      string `json:"id"`
	Running bool   `json:"running"`
	Output  string `json:"output"`
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func main() {
	r := gin.Default()

	clients := map[string]*websocket.Conn{}
	b, err := ioutil.ReadFile("./config.json")
	if err != nil {
		logrus.Fatalf("Could not read resources.json: %v", err)
	}
	credentials := creds{}
	if e := json.Unmarshal(b, &credentials); e != nil {
		logrus.Fatalf("Could not unmarshal resources: %v", e)
	}

	b, err = ioutil.ReadFile("./resources.json")
	if err != nil {
		logrus.Fatalf("Could not read resources.json: %v", err)
	}
	resources := res{}
	if err := json.Unmarshal(b, &resources); err != nil {
		logrus.Fatalf("Could not unmarshal resources: %v", err)
	}

	r.LoadHTMLFiles("./ui.html")
	r.Any("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		conId := uuid.New().String()
		clients[conId] = conn
		conn.WriteJSON(wsCmdRegister{
			ID:      conId,
			Running: command != nil,
		})
		go func() {
			defer conn.Close()
			conn.SetReadLimit(maxMessageSize)
			conn.SetReadDeadline(time.Now().Add(pongWait))
			conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						log.Printf("error: %v", err)
					}
					break
				}
				message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
				logrus.Info(message)
			}
		}()
	})
	r.GET("/ui", func(c *gin.Context) {
		c.HTML(http.StatusOK, "ui.html", tmplcontent{
			creds: credentials,
			res:   resources,
		})
	})
	r.POST("/run", func(c *gin.Context) {
		go runBot(&resources, clients[c.Query("conId")], clients)
		c.Redirect(http.StatusMovedPermanently, "/ui")
	})
	r.POST("/stop", func(c *gin.Context) {
		if command != nil {
			command.Stop()
			command = nil
		}
		conId := c.Query("conId")
		for _, client := range clients {
			if client != nil {
				client.WriteJSON(wsCmdRegister{
					ID:      conId,
					Running: command != nil,
				})
			}
		}
		c.Redirect(http.StatusMovedPermanently, "/ui")
	})
	r.POST("/save", func(c *gin.Context) {
		req := saveReq{}
		if err := c.Bind(&req); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		newRes := res{
			Hashtags:     strings.Split(strings.Trim(req.Hashtags, "\r\n"), "\r\n"),
			Comments:     strings.Split(strings.Trim(req.Comments, "\r\n"), "\r\n"),
			TotalLikes:   req.TotalLikes,
			PerUser:      req.PerUser,
			Potency:      PotencyMode(req.Potency),
			MaxFollowers: req.MaxFollowers,
			MinFollowers: req.MinFollowers,
			MaxFollowing: req.MaxFollowing,
			MinFollowing: req.MinFollowing,
		}
		b, _ := json.MarshalIndent(newRes, "", "    ")
		ioutil.WriteFile("./resources.json", b, 0655)
		resources = newRes

		credentials.Password = req.Password
		credentials.Username = req.Username
		b, _ = json.MarshalIndent(credentials, "", "    ")
		ioutil.WriteFile("./config.json", b, 0655)
		c.Redirect(http.StatusMovedPermanently, "/ui")
	})
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV)
	go func() {
		<-sigs
		if command != nil {
			command.Stop()
		}
		os.Exit(0)
	}()
	http.ListenAndServe(":8080", r)
}
