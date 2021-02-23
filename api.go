/*
 * Echosphere
 * Copyright (C) 2019  Nicolò Santamaria, Michele Dimaggio, Alessandro Ianne
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
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
)

type Api string

type Option string

const (
	PARSE_MARKDOWN           Option = "&parse_mode=markdown"
	PARSE_HTML                      = "&parse_mode=html"
	DISABLE_WEB_PAGE_PREVIEW        = "&disable_web_page_preview=true"
	DISABLE_NOTIFICATION            = "&disable_notification=true"
)

type ChatAction string

const (
	TYPING            ChatAction = "typing"
	UPLOAD_PHOTO                 = "upload_photo"
	RECORD_VIDEO                 = "record_video"
	UPLOAD_VIDEO                 = "upload_video"
	RECORD_AUDIO                 = "record_audio"
	UPLOAD_AUDIO                 = "upload_audio"
	UPLOAD_DOCUMENT              = "upload_document"
	FIND_LOCATION                = "find_location"
	RECORD_VIDEO_NOTE            = "record_video_note"
	UPLOAD_VIDEO_NOTE            = "upload_video_note"
)

func encode(s string) string {
	return url.QueryEscape(s)
}

func parseOpts(opts ...Option) string {
	var buf strings.Builder

	for _, o := range opts {
		buf.WriteString(string(o))
	}
	return buf.String()
}

func makeInlineKeyboard(rows ...InlineKbdRow) InlineKeyboard {
	return InlineKeyboard{rows}
}

// NewApi returns a new Api object.
func NewApi(token string) Api {
	return Api(fmt.Sprintf("https://api.telegram.org/bot%s/", token))
}

// DeleteWebhook deletes webhook
func (a Api) DeleteWebhook() (response APIResponseUpdate) {
	content := SendGetRequest(string(a) + "deleteWebhook")
	json.Unmarshal(content, &response)
	return
}

// SetWebhook sets the webhook to bot on Telegram servers
func (a Api) SetWebhook(url string) (response APIResponseUpdate) {
	keyVal := map[string]string{"url": url}
	content, err := SendPostForm(fmt.Sprintf("%ssetWebhook", string(a)), keyVal)
	if err != nil {
		log.Println(err)
		return
	}
	json.Unmarshal(content, &response)
	return
}

// GetResponse returns the incoming updates from telegram.
func (a Api) GetUpdates(offset, timeout int) (response APIResponseUpdate) {
	var url = fmt.Sprintf("%sgetUpdates?timeout=%d", string(a), timeout)

	if offset != 0 {
		url = fmt.Sprintf("%s&offset=%d", url, offset)
	}
	content := SendGetRequest(url)
	json.Unmarshal(content, &response)
	return
}

// Returns the current chat in use.
func (a Api) GetChat(chatId int64) (response Chat) {
	var url = fmt.Sprintf("%sgetChat?chat_id=%d", string(a), chatId)

	content := SendGetRequest(url)
	json.Unmarshal(content, &response)
	return
}

func (a Api) GetStickerSet(name string) (response StickerSet) {
	var url = fmt.Sprintf("%sgetStickerSet?name=%s", string(a), encode(name))

	content := SendGetRequest(url)
	json.Unmarshal(content, &response)
	return
}

func (a Api) SendMessage(text string, chatId int64, opts ...Option) (response APIResponseMessage) {
	var url = fmt.Sprintf(
		"%ssendMessage?text=%s&chat_id=%d%s",
		string(a),
		encode(text),
		chatId,
		parseOpts(opts...),
	)

	content := SendGetRequest(url)
	json.Unmarshal(content, &response)
	return
}

// Sends a message as a reply to a previously received one.
func (a Api) SendMessageReply(text string, chatId int64, messageId int, opts ...Option) (response APIResponseMessage) {
	var url = fmt.Sprintf(
		"%ssendMessage?text=%s&chat_id=%d&reply_to_message_id=%d%s",
		string(a),
		encode(text),
		chatId,
		messageId,
		parseOpts(opts...),
	)

	content := SendGetRequest(url)
	json.Unmarshal(content, &response)
	return
}

func (a Api) SendMessageWithKeyboard(text string, chatId int64, keyboard []byte, opts ...Option) (response APIResponseMessage) {
	var url = fmt.Sprintf(
		"%ssendMessage?text=%s&chat_id=%d&reply_markup=%s%s",
		string(a),
		encode(text),
		chatId,
		keyboard,
		parseOpts(opts...),
	)

	content := SendGetRequest(url)
	json.Unmarshal(content, &response)
	return
}

func (a Api) DeleteMessage(chatId int64, messageId int) (response APIResponseMessage) {
	var url = fmt.Sprintf(
		"%sdeleteMessage?chat_id=%d&message_id=%d",
		string(a),
		chatId,
		messageId,
	)

	content := SendGetRequest(url)
	json.Unmarshal(content, &response)
	return
}

func (a Api) SendPhoto(filename, caption string, chatId int64, opts ...Option) (response APIResponseMessage) {
	var url = fmt.Sprintf(
		"%ssendPhoto?chat_id=%d&caption=%s%s",
		string(a),
		chatId,
		encode(caption),
		parseOpts(opts...),
	)

	content := SendPostRequest(url, filename, "photo")
	json.Unmarshal(content, &response)
	return
}

func (a Api) SendPhotoByID(photoId, caption string, chatId int64, opts ...Option) (response APIResponseMessage) {
	var url = fmt.Sprintf(
		"%ssendPhoto?chat_id=%d&photo=%s&caption=%s%s",
		string(a),
		chatId,
		encode(photoId),
		encode(caption),
		parseOpts(opts...),
	)

	content := SendGetRequest(url)
	json.Unmarshal(content, &response)
	return
}

func (a Api) SendPhotoWithKeyboard(filename, caption string, chatId int64, keyboard []byte, opts ...Option) (response APIResponseMessage) {
	var url = fmt.Sprintf(
		"%ssendPhoto?chat_id=%d&caption=%s&reply_markup=%s%s",
		string(a),
		chatId,
		encode(caption),
		keyboard,
		parseOpts(opts...),
	)

	content := SendPostRequest(url, filename, "photo")
	json.Unmarshal(content, &response)
	return
}

func (a Api) SendAudio(filename, caption string, chatId int64, opts ...Option) (response APIResponseMessage) {
	var url = fmt.Sprintf(
		"%ssendAudio?chat_id=%d&caption=%s%s",
		string(a),
		chatId,
		encode(caption),
		parseOpts(opts...),
	)

	content := SendPostRequest(url, filename, "audio")
	json.Unmarshal(content, &response)
	return
}

func (a Api) SendAudioByID(audioId, caption string, chatId int64, opts ...Option) (response APIResponseMessage) {
	var url = fmt.Sprintf(
		"%ssendAudio?chat_id=%d&audio=%s&caption=%s%s",
		string(a),
		chatId,
		encode(audioId),
		encode(caption),
		parseOpts(opts...),
	)

	content := SendGetRequest(url)
	json.Unmarshal(content, &response)
	return
}

func (a Api) SendDocument(filename, caption string, chatId int64, opts ...Option) (response APIResponseMessage) {
	var url = fmt.Sprintf(
		"%ssendDocument?chat_id=%d&caption=%s%s",
		string(a),
		chatId,
		encode(caption),
		parseOpts(opts...),
	)

	content := SendPostRequest(url, filename, "document")
	json.Unmarshal(content, &response)
	return
}

func (a Api) SendDocumentByID(documentId, caption string, chatId int64, opts ...Option) (response APIResponseMessage) {
	var url = fmt.Sprintf(
		"%ssendDocument?chat_id=%d&document=%s&caption=%s%s",
		string(a),
		chatId,
		encode(documentId),
		encode(caption),
		parseOpts(opts...),
	)

	content := SendGetRequest(url)
	json.Unmarshal(content, &response)
	return
}

func (a Api) SendVideo(filename, caption string, chatId int64, opts ...Option) (response APIResponseMessage) {
	var url = fmt.Sprintf(
		"%ssendVideo?chat_id=%d&caption=%s%s",
		string(a),
		chatId,
		encode(caption),
		parseOpts(opts...),
	)

	content := SendPostRequest(url, filename, "video")
	json.Unmarshal(content, &response)
	return
}

func (a Api) SendVideoByID(videoId, caption string, chatId int64, opts ...Option) (response APIResponseMessage) {
	var url = fmt.Sprintf(
		"%ssendVideo?chat_id=%d&video=%s&caption=%s%s",
		string(a),
		chatId,
		encode(videoId),
		encode(caption),
		parseOpts(opts...),
	)

	content := SendGetRequest(url)
	json.Unmarshal(content, &response)
	return
}

func (a Api) SendVideoNoteByID(videoId string, chatId int64) (response APIResponseMessage) {
	var url = fmt.Sprintf(
		"%ssendVideoNote?chat_id=%d&video_note=%s",
		string(a),
		chatId,
		encode(videoId),
	)

	content := SendGetRequest(url)
	json.Unmarshal(content, &response)
	return
}

func (a Api) SendVoice(filename, caption string, chatId int64, opts ...Option) (response APIResponseMessage) {
	var url = fmt.Sprintf(
		"%ssendVoice?chat_id=%d&caption=%s%s",
		string(a),
		chatId,
		encode(caption),
		parseOpts(opts...),
	)

	content := SendPostRequest(url, filename, "voice")
	json.Unmarshal(content, &response)
	return
}

func (a Api) SendVoiceByID(voiceId, caption string, chatId int64, opts ...Option) (response APIResponseMessage) {
	var url = fmt.Sprintf(
		"%ssendVoice?chat_id=%d&voice=%s%s",
		string(a),
		chatId,
		encode(voiceId),
		parseOpts(opts...),
	)

	content := SendGetRequest(url)
	json.Unmarshal(content, &response)
	return
}

func (a Api) SendContact(phoneNumber, firstName, lastName string, chatId int64) (response APIResponseMessage) {
	var url = fmt.Sprintf(
		"%ssendContact?chat_id=%d&phone_number=%s&first_name=%s&last_name=%s",
		string(a),
		chatId,
		encode(phoneNumber),
		encode(firstName),
		encode(lastName),
	)

	content := SendGetRequest(url)
	json.Unmarshal(content, &response)
	return
}

func (a Api) SendStickerByID(stickerId string, chatId int64) (response APIResponseMessage) {
	var url = fmt.Sprintf(
		"%ssendSticker?chat_id=%d&sticker=%s",
		string(a),
		chatId,
		encode(stickerId),
	)

	content := SendGetRequest(url)
	json.Unmarshal(content, &response)
	return
}

func (a Api) SendChatAction(action ChatAction, chatId int64) (response APIResponseMessage) {
	var url = fmt.Sprintf(
		"%ssendChatAction?chat_id=%d&action=%s",
		string(a),
		chatId,
		action,
	)

	content := SendGetRequest(url)
	json.Unmarshal(content, &response)
	return
}

func (a Api) KeyboardButton(text string, requestContact, requestLocation bool) Button {
	return Button{text, requestContact, requestLocation}
}

func (a Api) KeyboardRow(buttons ...Button) (kbdRow KbdRow) {
	for _, button := range buttons {
		kbdRow = append(kbdRow, button)
	}

	return
}

func (a Api) KeyboardMarkup(resizeKeyboard, oneTimeKeyboard, selective bool, keyboardRows ...KbdRow) (kbd []byte) {
	kbd, _ = json.Marshal(Keyboard{
		keyboardRows,
		resizeKeyboard,
		oneTimeKeyboard,
		selective,
	})
	return
}

func (a Api) KeyboardRemove(selective bool) (kbdrmv []byte) {
	kbdrmv, _ = json.Marshal(KeyboardRemove{true, selective})
	return
}

// Returns a new inline keyboard button with the provided data.
func (a Api) InlineKbdBtn(text, url, callbackData string) InlineButton {
	return InlineButton{
		encode(text),
		url,
		callbackData,
	}
}

// Same as InlineKbdBtn, but only with url.
func (a Api) InlineKbdBtnUrl(text, url string) InlineButton {
	return a.InlineKbdBtn(text, url, "")
}

// Same as InlineKbdBtn, but only with callbackData.
func (a Api) InlineKbdBtnCbd(text, callbackData string) InlineButton {
	return a.InlineKbdBtn(text, "", callbackData)
}

// Returns a new inline keyboard row with the given buttons.
func (a Api) InlineKbdRow(inlineButtons ...InlineButton) InlineKbdRow {
	return inlineButtons
}

// Returns a byte slice containing the inline keyboard json data.
func (a Api) InlineKbdMarkup(inlineKbdRows ...InlineKbdRow) (jsn []byte) {
	jsn, _ = json.Marshal(makeInlineKeyboard(inlineKbdRows...))
	return
}

func (a Api) EditMessageReplyMarkup(chatId int64, messageId int, keyboard []byte) (response APIResponseMessage) {
	var url = fmt.Sprintf(
		"%seditMessageReplyMarkup?chat_id=%d&message_id=%d&reply_markup=%s",
		string(a),
		chatId,
		messageId,
		keyboard,
	)

	content := SendGetRequest(url)
	json.Unmarshal(content, &response)
	return
}

func (a Api) EditMessageText(chatId int64, messageId int, text string, opts ...Option) (response APIResponseMessage) {
	var url = fmt.Sprintf(
		"%seditMessageText?chat_id=%d&message_id=%d&text=%s%s",
		string(a),
		chatId,
		messageId,
		encode(text),
		parseOpts(opts...),
	)

	content := SendGetRequest(url)
	json.Unmarshal(content, &response)
	return
}

func (a Api) EditMessageTextWithKeyboard(chatId int64, messageId int, text string, keyboard []byte, opts ...Option) (response APIResponseMessage) {
	var url = fmt.Sprintf(
		"%seditMessageText?chat_id=%d&message_id=%d&text=%s&reply_markup=%s%s",
		string(a),
		chatId,
		messageId,
		encode(text),
		keyboard,
		parseOpts(opts...),
	)

	content := SendGetRequest(url)
	json.Unmarshal(content, &response)
	return
}

func (a Api) AnswerCallbackQuery(id, text string, showAlert bool) (response APIResponseMessage) {
	var url = fmt.Sprintf(
		"%sanswerCallbackQuery?callback_query_id=%s&text=%s&show_alert=%s",
		string(a),
		id,
		text,
		strconv.FormatBool(showAlert),
	)

	content := SendGetRequest(url)
	json.Unmarshal(content, &response)
	return
}

func (a Api) GetMyCommands() (response APIResponseCommands) {
	var url = fmt.Sprintf(
		"%sgetMyCommands",
		string(a),
	)

	content := SendGetRequest(url)
	json.Unmarshal(content, &response)
	return
}

func (a Api) SetMyCommands(commands ...BotCommand) (response APIResponseCommands) {
	jsn, _ := json.Marshal(commands)

	var url = fmt.Sprintf(
		"%ssetMyCommands?commands=%s",
		string(a),
		jsn,
	)

	content := SendGetRequest(url)
	json.Unmarshal(content, &response)
	return
}

func (a Api) Command(command, description string) BotCommand {
	return BotCommand{command, description}
}

func (a Api) SendAnimation(filename, caption string, chatId int64, opts ...Option) (response APIResponseMessage) {
	var url = fmt.Sprintf(
		"%ssendAnimation?chat_id=%d&caption=%s%s",
		string(a),
		chatId,
		encode(caption),
		parseOpts(opts...),
	)

	content := SendPostRequest(url, filename, "animation")
	json.Unmarshal(content, &response)
	return
}

func (a Api) SendAnimationByID(animationId string, chatId int64) (response APIResponseMessage) {
	var url = fmt.Sprintf(
		"%ssendAnimation?chat_id=%d&animation=%s",
		string(a),
		chatId,
		encode(animationId),
	)

	content := SendGetRequest(url)
	json.Unmarshal(content, &response)
	return
}