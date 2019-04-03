package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/go-cmd/cmd"
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

type res struct {
	Hashtags []string `json:"hashtags" form:"hashtags"`
	Comments []string `json:"comments" form:"comments"`
	Sample   int      `json:"sample" forn:"sample"`
	Potency  string   `json:"potency" form:"potency"`
}

func (r res) HastagStr() string {
	return strings.Join(r.Hashtags, "\n")
}

func (r res) CommentStr() string {
	return strings.Join(r.Comments, "\n")
}

type saveReq struct {
	Hashtags string `json:"hashtags" form:"hashtags"`
	Comments string `json:"comments" form:"comments"`
	Sample   int    `json:"sample" form:"sample"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Potency  string `json:"potency" form:"potency"`
}

var command *cmd.Cmd

func runBot(r *res) {
	// Start a long-running process, capture stdout and stderr
	if command != nil {
		return
	}
	command = cmd.NewCmd("python3", "main.py")
	statusChan := command.Start() // non-blocking

	ticker := time.NewTicker(2 * time.Second)

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
	case finalStatus := <-statusChan:
		logrus.Info(finalStatus)
		// done
	default:
		// no, still running
		logrus.Info("Still running")
	}

	// Block waiting for command to exit, be stopped, or be killed
	<-statusChan
}

func main() {
	r := gin.Default()

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
	r.GET("/ui", func(c *gin.Context) {
		c.HTML(http.StatusOK, "ui.html", tmplcontent{
			creds: credentials,
			res:   resources,
		})
	})
	r.POST("/run", func(c *gin.Context) {
		go runBot(&resources)
		c.Redirect(http.StatusMovedPermanently, "/ui")
	})
	r.POST("/stop", func(c *gin.Context) {
		if command != nil {
			command.Stop()
			command = nil
		}
		c.Redirect(http.StatusMovedPermanently, "/ui")
	})
	r.POST("/save", func(c *gin.Context) {
		req := saveReq{}
		if err := c.Bind(&req); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		spew.Dump(req)
		newRes := res{
			Hashtags: strings.Split(strings.Trim(req.Hashtags, "\r\n"), "\r\n"),
			Comments: strings.Split(strings.Trim(req.Comments, "\r\n"), "\r\n"),
			Sample:   req.Sample,
			Potency:  req.Potency,
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
	http.ListenAndServe(":8080", r)
}
