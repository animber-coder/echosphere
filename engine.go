/*
 * Echosphere
 * Copyright (C) 2019  Nicolò Santamaria, Michele Dimaggio, Alessandro Ianne
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package echosphere

import (
	"fmt"
	"strings"
	"encoding/json"
)


type Engine struct {
	url string
}


const (
	PARSE_MARKDOWN = 1 << iota
	PARSE_HTML = 1 << iota
	DISABLE_WEB_PAGE_PREVIEW = 1 << iota
	DISABLE_NOTIFICATION = 1 << iota
)


func NewEngine(token string) Engine {
	return Engine{
		url: fmt.Sprintf("https://api.telegram.org/bot%s/", token),
	}
}


// GetResponse returns the incoming updates from telegram.
func (e Engine) GetResponse(offset int, timeout int) APIResponse {
	var url = fmt.Sprintf("%sgetUpdates?timeout=%d", e.url, timeout)

	if offset != 0 {
		url = fmt.Sprintf("%s&offset=%d", url, offset)
	}
	var content []byte = SendGetRequest(url)
	var response APIResponse

	json.Unmarshal(content, &response)
	return response
}


func (e Engine) GetChat(chatId int64) Chat {
	var url = fmt.Sprintf("%sgetChat?chat_id=%d", e.url, chatId)
	var content []byte = SendGetRequest(url)
	var response Chat

	json.Unmarshal(content, &response)
	return response
}


func (e Engine) GetStickerSet(name string) StickerSet {
	var url = fmt.Sprintf("%sgetStickerSet?name=%s", e.url, name)

	var content []byte = SendGetRequest(url)
	var response StickerSet

	json.Unmarshal(content, &response)

	return response
}


func (e Engine) SendMessage(text string, chatId int64) APIResponse {
	var url = fmt.Sprintf("%ssendMessage?text=%s&chat_id=%d", e.url, strings.Replace(text, "\n", "%0A", -1), chatId)
	var content []byte = SendGetRequest(url)
	var response APIResponse

	json.Unmarshal(content, &response)
	return response
}


func (e Engine) SendMessageOptions(text string, chatId int64, options int) APIResponse {
	var url = fmt.Sprintf("%ssendMessage?text=%s&chat_id=%d", e.url, strings.Replace(text, "\n", "%0A", -1), chatId)

	if options & PARSE_MARKDOWN != 0 {
		url += "&parse_mode=markdown"
	}

	if options & PARSE_HTML != 0 {
		url += "&parse_mode=html"
	}

	if options & DISABLE_WEB_PAGE_PREVIEW != 0 {
		url += "&disable_web_page_preview=true"
	}

	if options & DISABLE_NOTIFICATION != 0 {
		url += "&disable_notification=true"
	}

	var content []byte = SendGetRequest(url)
	var response APIResponse

	json.Unmarshal(content, &response)
	return response
}


func (e Engine) SendMessageWithKeyboard(text string, chatId int64, keyboard []byte) APIResponse {
	var url = fmt.Sprintf("%ssendMessage?text=%s&chat_id=%d&parse_mode=markdown&reply_markup=%s", e.url, strings.Replace(text, "\n", "%0A", -1), chatId, keyboard)
	var content []byte = SendGetRequest(url)
	var response APIResponse

	json.Unmarshal(content, &response)
	return response
}


func (e Engine) SendPhoto(filename string, chatId int64, caption string) APIResponse {
	var url = fmt.Sprintf("%ssendPhoto?chat_id=%d&caption=%s", e.url, chatId, caption)
	var content []byte = SendPostRequest(url, filename, "photo")
	var response APIResponse

	json.Unmarshal(content, &response)
	return response
}


func (e Engine) SendPhotoByID(photoId string, chatId int64) APIResponse {
	var url = fmt.Sprintf("%ssendPhoto?chat_id=%d&photo=%s", e.url, chatId, photoId)
	var content []byte = SendGetRequest(url)
	var response APIResponse

	json.Unmarshal(content, &response)
	return response
}


func (e Engine) SendAudio(filename string, chatId int64, caption string) APIResponse {
	var url = fmt.Sprintf("%ssendAudio?chat_id=%d&caption=%s", e.url, chatId, caption)
	var content = SendPostRequest(url, filename, "audio")
	var response APIResponse

	json.Unmarshal(content, &response)
	return response
}


func (e Engine) SendAudioByID(audioId string, chatId int64) APIResponse {
	var url = fmt.Sprintf("%ssendAudio?chat_id=%d&audio=%s", e.url, chatId, audioId)
	var content []byte = SendGetRequest(url)
	var response APIResponse

	json.Unmarshal(content, &response)
	return response
}


func (e Engine) SendDocument(filename string, caption string, chatId int64) APIResponse {
	var url = fmt.Sprintf("%ssendDocument?chat_id=%d&caption=%s", e.url, chatId, caption)
	var content = SendPostRequest(url, filename, "document")
	var response APIResponse

	json.Unmarshal(content, &response)
	return response
}


func (e Engine) SendDocumentByID(documentId string, caption string, chatId int64) APIResponse {
	var url = fmt.Sprintf("%ssendDocument?chat_id=%d&document=%s&caption=%s", e.url, chatId, documentId, caption)
	var content []byte = SendGetRequest(url)
	var response APIResponse

	json.Unmarshal(content, &response)
	return response
}


func (e Engine) SendVideo(filename string, chatId int64, caption string) APIResponse {
	var url = fmt.Sprintf("%ssendVideo?chat_id=%d&caption=%s", e.url, chatId, caption)
	var content = SendPostRequest(url, filename, "video")
	var response APIResponse

	json.Unmarshal(content, &response)
	return response
}


func (e Engine) SendVideoByID(videoId string, chatId int64) APIResponse {
	var url = fmt.Sprintf("%ssendVideo?chat_id=%d&video=%s", e.url, chatId, videoId)
	var content []byte = SendGetRequest(url)
	var response APIResponse

	json.Unmarshal(content, &response)
	return response
}


func (e Engine) SendVideoNoteByID(videoId string, chatId int64) APIResponse {
	var url = fmt.Sprintf("%ssendVideoNote?chat_id=%d&video_note=%s", e.url, chatId, videoId)
	var content []byte = SendGetRequest(url)
	var response APIResponse

	json.Unmarshal(content, &response)
	return response
}


func (e Engine) SendVoice(filename string, chatId int64, caption string) APIResponse {
	var url = fmt.Sprintf("%ssendVoice?chat_id=%d&caption=%s", e.url, chatId, caption)
	var content = SendPostRequest(url, filename, "voice")
	var response APIResponse

	json.Unmarshal(content, &response)
	return response
}


func (e Engine) SendVoiceByID(voiceId string, chatId int64) APIResponse {
	var url = fmt.Sprintf("%ssendVoice?chat_id=%d&voice=%s", e.url, chatId, voiceId)
	var content []byte = SendGetRequest(url)
	var response APIResponse

	json.Unmarshal(content, &response)
	return response
}


func (e Engine) SendContact(phoneNumber string, firstName string, lastName string, chatId int64) APIResponse {
	var url = fmt.Sprintf("%ssendContact?chat_id=%d&phone_number=%s&first_name=%s&last_name=%s", e.url, chatId, phoneNumber, firstName, lastName)
	var content = SendGetRequest(url)
	var response APIResponse

	json.Unmarshal(content, &response)
	return response
}


func (e Engine) SendStickerByID(stickerId string, chatId int64) APIResponse {
	var url = fmt.Sprintf("%ssendSticker?chat_id=%d&sticker=%s", e.url, chatId, stickerId)
	var content []byte = SendGetRequest(url)
	var response APIResponse

	json.Unmarshal(content, &response)
	return response
}


func (e Engine) KeyboardButton(text string, requestContact bool, requestLocation bool) Button {
	return Button{
		text,
		requestContact,
		requestLocation,
	}
}


func (e Engine) KeyboardRow(buttons ...Button) KbdRow {
	var kbdRow KbdRow

	for _, button := range buttons {
		kbdRow = append(kbdRow, button)
	}

	return kbdRow
}


func (e Engine) KeyboardMarkup(resizeKeyboard bool, oneTimeKeyboard bool, selective bool, keyboardRows ...KbdRow) []byte {
	keyboard := Keyboard{
		nil,
		resizeKeyboard,
		oneTimeKeyboard,
		selective,
	}

	for _, row := range keyboardRows {
		keyboard.Keyboard = append(keyboard.Keyboard, row)
	}

	kbd, _ := json.Marshal(keyboard)
	return kbd
}


func (e Engine) KeyboardRemove(selective bool) []byte {
	kbdrmv, _ := json.Marshal(KeyboardRemove{
		true,
		selective,
	})

	return kbdrmv
}


func (e Engine) InlineKbdBtn(text string, url string, callbackData string) InlineButton {
	return InlineButton{
		text,
		url,
		callbackData,
	}
}


func (e Engine) InlineKbdRow(inlineButtons ...InlineButton) InlineKbdRow {
	var inlineKbdRow InlineKbdRow

	for _, inlineButton := range inlineButtons {
		inlineKbdRow = append(inlineKbdRow, inlineButton)
	}

	return inlineKbdRow
}


func (e Engine) InlineKbdMarkup(inlineKbdRows ...InlineKbdRow) []byte {
	var inlineKeyboard InlineKeyboard

	for _, row := range inlineKbdRows {
		inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard, row)
	}

	kbd, _ := json.Marshal(inlineKeyboard)
	return kbd
}
