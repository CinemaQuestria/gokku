/*
Command gokku takes a gogs commit hook and turns it into a dokku push.
*/
package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/Xe/gokku/lib"
	"github.com/codegangsta/negroni"
	"github.com/thoj/go-ircevent"
)

var (
	bot         *irc.Connection
	dir         = flag.String("dir", "/gokku/repo", "directory to clone to and stuff")
	repo        string
	dokkuremote string
)

func init() {
	bot = irc.IRC("CommitBot-CQ", "gokku")
	bot.UseTLS = true

	bot.AddCallback("001", func(e *irc.Event) {
		bot.Privmsg("NickServ", "IDENTIFY "+os.Getenv("BOT_PASS"))
		time.Sleep(5 * time.Second)
		bot.Join("#" + os.Getenv("BOT_CHANNEL"))
	})

	err := bot.Connect("irc.ponychat.net:6697")
	if err != nil {
		panic(err)
	}

	repo = os.Getenv("GOKKU_REPO")
	dokkuremote = os.Getenv("GOKKU_DOKKU_REMOTE")

	cmd := exec.Command("git", "clone", repo)
	cmd.Dir = *dir
	err = cmd.Run()

	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()

	mux := http.NewServeMux()

	mux.HandleFunc("/", handleHTTP)

	n := negroni.Classic()
	n.UseHandler(mux)

	port := "3000"
	if foo := os.Getenv("PORT"); foo != "" {
		port = os.Getenv("PORT")
	}

	n.Run(":" + port)
}

func handleHTTP(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.Write([]byte("Invalid request: " + err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var hook lib.CommitHook

	err = json.Unmarshal(body, hook)
	if err != nil {
		w.Write([]byte("Invalid request: " + err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	go doCommand(hook)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("all is good friend"))

	bot.Privmsgf(
		"#"+os.Getenv("BOT_CHANNEL"),
		"%s (%s) made %d commits to %s of the site, deploying...",
		hook.Pusher.Name,
		hook.Pusher.Username,
		len(hook.Commits),
		hook.Ref,
	)
	bot.Privmsg("#"+os.Getenv("BOT_CHANNEL"), hook.CompareURL)

	for _, commit := range hook.Commits {
		bot.Privmsgf(
			"#"+os.Getenv("BOT_CHANNEL"),
			"[%d] %s - %s",
			commit.ID[:8],
			commit.Author,
			commit.Message,
		)
	}
}

func doCommand(hook lib.CommitHook) {
	if hook.Ref != "refs/heads/master" {
		return
	}

	cmd := exec.Command("git", "pull")
	cmd.Dir = *dir
	err := cmd.Run()

	if err != nil {
		bot.Privmsgf(
			"#"+os.Getenv("BOT_CHANNEL"),
			"Had some error pulling repo: %s",
			err.Error(),
		)
	}

	pushcmd := exec.Command("git", "push", dokkuremote, "master")
	pushcmd.Dir = *dir

	err = pushcmd.Run()
	if err != nil {
		bot.Privmsgf(
			"#"+os.Getenv("BOT_CHANNEL"),
			"Had some error pushing repo: %s",
			err.Error(),
		)
	}
}
