/*
 * Echosphere
 * Copyright (C) 2018-2022 The Echosphere Devs
 *
 * Echosphere is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Echosphere is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package echosphere

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
)

// Bot is the interface that must be implemented by your definition of
// the struct thus it represent each open session with a user on Telegram.
type Bot interface {
	// Update will be called upon receiving any update from Telegram.
	Update(*Update)
}

// NewBotFn is called every time echosphere receives an update with a chat ID never
// encountered before.
type NewBotFn func(chatId int64) Bot

// The Dispatcher passes the updates from the Telegram Bot API to the Bot instance
// associated with each chatID. When a new chat ID is found, the provided function
// of type NewBotFn will be called.
type Dispatcher struct {
	api        API
	sessionMap map[int64]Bot
	newBot     NewBotFn
	updates    chan *Update
	httpServer *http.Server
	mu         sync.Mutex
}

// NewDispatcher returns a new instance of the Dispatcher object.
// Calls the Update function of the bot associated with each chat ID.
// If a new chat ID is found, newBotFn will be called first.
func NewDispatcher(token string, newBotFn NewBotFn) *Dispatcher {
	d := &Dispatcher{
		api:        NewAPI(token),
		sessionMap: make(map[int64]Bot),
		newBot:     newBotFn,
		updates:    make(chan *Update),
	}
	go d.listen()
	return d
}

// DelSession deletes the Bot instance, seen as a session, from the
// map with all of them.
func (d *Dispatcher) DelSession(chatID int64) {
	d.mu.Lock()
	delete(d.sessionMap, chatID)
	d.mu.Unlock()
}

// AddSession allows to arbitrarily create a new Bot instance.
func (d *Dispatcher) AddSession(chatID int64) {
	d.mu.Lock()
	if _, isIn := d.sessionMap[chatID]; !isIn {
		d.sessionMap[chatID] = d.newBot(chatID)
	}
	d.mu.Unlock()
}

// Poll is a wrapper function for PollOptions.
func (d *Dispatcher) Poll() error {
	return d.PollOptions(true, UpdateOptions{Timeout: 120})
}

// PollOptions starts the polling loop so that the dispatcher calls the function Update
// upon receiving any update from Telegram.
func (d *Dispatcher) PollOptions(dropPendingUpdates bool, opts UpdateOptions) error {
	var (
		timeout      = opts.Timeout
		isFirstRun   = true
		lastUpdateID = -1
	)

	// deletes webhook if present to run in long polling mode
	if _, err := d.api.DeleteWebhook(dropPendingUpdates); err != nil {
		return err
	}

	for {
		if isFirstRun {
			opts.Timeout = 0
		}

		opts.Offset = lastUpdateID + 1
		response, err := d.api.GetUpdates(&opts)
		if err != nil {
			return err
		}

		if !dropPendingUpdates || !isFirstRun {
			for _, u := range response.Result {
				d.updates <- u
			}
		}

		if l := len(response.Result); l > 0 {
			lastUpdateID = response.Result[l-1].ID
		}

		if isFirstRun {
			isFirstRun = false
			opts.Timeout = timeout
		}
	}
}

func (d *Dispatcher) instance(chatID int64) Bot {
	bot, ok := d.sessionMap[chatID]
	if !ok {
		bot = d.newBot(chatID)
		d.mu.Lock()
		d.sessionMap[chatID] = bot
		d.mu.Unlock()
	}
	return bot
}

func (d *Dispatcher) listen() {
	for update := range d.updates {
		var chatID int64

		if update.Message != nil {
			chatID = update.Message.Chat.ID
		} else if update.EditedMessage != nil {
			chatID = update.EditedMessage.Chat.ID
		} else if update.ChannelPost != nil {
			chatID = update.ChannelPost.Chat.ID
		} else if update.EditedChannelPost != nil {
			chatID = update.EditedChannelPost.Chat.ID
		} else if update.CallbackQuery != nil {
			chatID = update.CallbackQuery.Message.Chat.ID
		} else if update.InlineQuery != nil {
			chatID = update.InlineQuery.From.ID
		} else {
			continue
		}

		bot := d.instance(chatID)
		go bot.Update(update)
	}
}

// ListenWebhook is a wrapper function for ListenWebhookOptions.
func (d *Dispatcher) ListenWebhook(webhookURL string) error {
	return d.ListenWebhookOptions(webhookURL, false, nil)
}

// ListenWebhookOptions sets a webhook and listens for incoming updates.
// The webhookUrl should be provided in the following format: '<hostname>:<port>/<path>',
// eg: 'https://example.com:443/bot_token'.
// ListenWebhook will then proceed to communicate the webhook url '<hostname>/<path>' to Telegram
// and run a webserver that listens to ':<port>' and handles the path.
func (d *Dispatcher) ListenWebhookOptions(webhookURL string, dropPendingUpdates bool, opts *WebhookOptions) error {
	u, err := url.Parse(webhookURL)
	if err != nil {
		return err
	}

	whURL := fmt.Sprintf("%s%s", u.Hostname(), u.EscapedPath())
	log.Printf("setting webhook for %s\n", whURL)
	if _, err = d.api.SetWebhook(whURL, dropPendingUpdates, opts); err != nil {
		return err
	}

	if d.httpServer != nil {
		mux := http.NewServeMux()
		mux.Handle("/", d.httpServer.Handler)
		mux.HandleFunc(u.EscapedPath(), d.HandleWebhook)
		d.httpServer.Handler = mux
		return d.httpServer.ListenAndServe()
	}
	http.HandleFunc(u.EscapedPath(), d.HandleWebhook)
	log.Printf("listening on :%s\n", u.Port())
	return http.ListenAndServe(fmt.Sprintf(":%s", u.Port()), nil)
}

// SetHTTPServer allows to set a custom http.Server for ListenWebhook and ListenWebhookOptions.
func (d *Dispatcher) SetHTTPServer(s *http.Server) {
	d.httpServer = s
}

// HandleWebhook is the http.HandlerFunc for the webhook URL.
// Useful if you've already a http server running and want to handle the request yourself.
func (d *Dispatcher) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	var jsn []byte

	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err := gzip.NewReader(r.Body)
		if err != nil {
			log.Println(err)
			return
		}
		defer reader.Close()
		if j, err := io.ReadAll(reader); err == nil {
			jsn = j
		} else {
			log.Println(err)
		}

	default:
		if j, err := io.ReadAll(r.Body); err == nil {
			jsn = j
		} else {
			log.Println(err)
		}
	}

	var update Update
	if err := json.Unmarshal(jsn, &update); err != nil {
		log.Println(err)
		return
	}

	d.updates <- &update
}
