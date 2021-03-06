package main

import (
	"bytes"
	"errors"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func getUser(c *gin.Context) *user {
	u, exists := c.Get("user")

	if !exists {
		return nil
	}

	return u.(*user)
}

func jsonError(c *gin.Context, status int, err error) {
	c.JSON(status, map[string]string{
		"error": err.Error(),
	})
}

func validateTokenMiddleware(c *gin.Context) {
	bearer := strings.Split(c.GetHeader("Authorization"), " ")
	if len(bearer) < 2 {
		jsonError(c, http.StatusUnauthorized, errors.New("Invalid Token"))
		c.Abort()
		return
	}
	bearerToken := bearer[1]
	if err := verifyToken(&token{AccessToken: bearerToken}); err != nil {
		jsonError(c, http.StatusUnauthorized, errors.New("Invalid Token"))
		c.Abort()
		return
	}

	globMut.Lock()

	user := users[bearerToken]
	userExists := user != nil

	if userExists {
		latestUser, err := uRepo.findUserByUsername(user.Username)
		if err != nil {
			jsonError(c, http.StatusInternalServerError, err)
			c.Abort()
			return
		}
		user = latestUser
	}

	globMut.Unlock()

	if !userExists {
		jsonError(c, http.StatusUnauthorized, errors.New("Invalid Token"))
		c.Abort()
		return
	}

	c.Set("user", user)

	c.Next()
}

type authenticatePayload struct {
	Username string
	Password string
}

func authenticate(c *gin.Context) {
	req := authenticatePayload{}

	if err := c.BindJSON(&req); err != nil {
		jsonError(c, http.StatusBadRequest, err)
		return
	}

	t, u, err := authenticateUser(req.Username, req.Password)
	if err != nil {
		jsonError(c, http.StatusUnauthorized, err)
		return
	}

	globMut.Lock()

	users[t.AccessToken] = u

	globMut.Unlock()

	c.JSON(http.StatusOK, t)
}

type saveSettingsPayload struct {
	Settings *botSettings
}

func saveSettings(c *gin.Context) {
	req := saveSettingsPayload{}

	userValue, _ := c.Get("user")

	user := userValue.(*user)

	if err := c.BindJSON(&req); err != nil {
		jsonError(c, http.StatusBadRequest, err)
		return
	}

	user.Settings = req.Settings

	globMut.Lock()

	uRepo.save(user)

	globMut.Unlock()

	c.JSON(http.StatusOK, user)
}

func getTicketLogs(c *gin.Context) {
	ticketID := c.Param("ticketID")

	if ticketID == "" {
		jsonError(c, http.StatusBadRequest, errors.New("Empty ticket id"))
		return
	}

	globMut.Lock()

	ticket := tickets[ticketID]

	if ticket == nil {
		globMut.Unlock()
		jsonError(c, http.StatusBadRequest, errors.New("Ticket not found"))
		return
	}

	t := *ticket

	globMut.Unlock()

	c.JSON(http.StatusOK, map[string]interface{}{
		"Logs":    t.Logs,
		"ErrLogs": t.ErrLogs,
	})
}

func getTicketStatus(c *gin.Context) {
	ticketID := c.Param("ticketID")

	if ticketID == "" {
		jsonError(c, http.StatusBadRequest, errors.New("Empty ticket id"))
		return
	}

	globMut.Lock()

	ticket := tickets[ticketID]

	if ticket == nil {
		globMut.Unlock()
		jsonError(c, http.StatusBadRequest, errors.New("Ticket not found"))
		return
	}

	done := ticket.Done
	err := ticket.Err

	globMut.Unlock()

	var errMessage string

	if err != nil {
		errMessage = err.Error()
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"running": !done,
		"error":   errMessage,
	})
}

type runJobRequest struct {
	Label      string
	IGPassword string
	Settings   *runSettings
}

type runJobResponse struct {
	id string
}

func runJob(c *gin.Context) {
	req := runJobRequest{}

	user := getUser(c)

	if err := c.BindJSON(&req); err != nil {
		jsonError(c, http.StatusBadRequest, err)
		return
	}

	globMut.Lock()

	ticket, err := globBot.run(req.Label, user.Username, req.IGPassword, req.Settings)

	if err != nil {
		globMut.Unlock()
		logrus.Errorf("Error running the job: %v", err)
		jsonError(c, http.StatusInternalServerError, err)
		return
	}

	tickets[ticket.ID] = ticket

	globMut.Unlock()

	c.JSON(http.StatusOK, ticket)
}

func myTickets(c *gin.Context) {
	user := getUser(c)

	ret := []runTicket{}

	globMut.Lock()

	for _, ticket := range tickets {
		if ticket.Username == user.Username {
			t := *ticket
			t.Logs = []string{}
			t.ErrLogs = []string{}
			ret = append(ret, t)
		}
	}

	globMut.Unlock()

	c.JSON(http.StatusOK, ret)
}

func me(c *gin.Context) {
	user := getUser(c)
	c.JSON(http.StatusOK, user)
}

func script(c *gin.Context) {
	t, err := template.ParseFiles("./assets/script.js")

	if err != nil {
		jsonError(c, http.StatusInternalServerError, err)
		return
	}

	scriptJS := bytes.Buffer{}
	t.Execute(&scriptJS, map[string]string{
		"Host": os.Getenv("INSTABOT_URL"),
	})

	c.Writer.Write(scriptJS.Bytes())
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.WriteHeaderNow()
}

func noRoute(c *gin.Context) {
	jsonError(c, http.StatusNotFound, errors.New("Route not found"))
}
