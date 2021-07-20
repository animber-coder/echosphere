/*
 * Echosphere
 * Copyright (C) 2018-2021  The Echosphere Devs
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
	"net/url"
	"os"
	"path/filepath"
)

// API is the object that contains all the functions that wrap those of the Telegram Bot API.
type API struct {
	token string
	base  string
}

func encode(s string) string {
	return url.QueryEscape(s)
}

func sendFile(file InputFile, url, fileType string) (cnt []byte, err error) {
	switch {
	case file.id != "":
		cnt, err = sendGetRequest(fmt.Sprintf("%s&%s=%s", url, fileType, file.id))

	case file.path != "" && len(file.content) == 0:
		file.content, err = os.ReadFile(file.path)
		if err != nil {
			return cnt, err
		}
		file.path = filepath.Base(file.path)
		fallthrough

	case file.path != "" && len(file.content) > 0:
		cnt, err = sendPostRequest(url, file.path, fileType, file.content)
	}

	if err != nil {
		return cnt, err
	}

	return cnt, nil
}

// NewAPI returns a new API object.
func NewAPI(token string) API {
	return API{
		token: token,
		base:  fmt.Sprintf("https://api.telegram.org/bot%s/", token),
	}
}

// GetUpdates is used to receive incoming updates using long polling.
func (a API) GetUpdates(opts *UpdateOptions) (APIResponseUpdate, error) {
	var res APIResponseUpdate
	var url = fmt.Sprintf(
		"%sgetUpdates?%s",
		a.base,
		querify(opts),
	)

	cnt, err := sendGetRequest(url)
	if err != nil {
		return res, err
	}

	return res, json.Unmarshal(cnt, &res)
}

// SetWebhook is used to specify a url and receive incoming updates via an outgoing webhook.
func (a API) SetWebhook(webhookURL string, dropPendingUpdates bool, opts *WebhookOptions) (APIResponseBase, error) {
	var res APIResponseBase
	var url = fmt.Sprintf(
		"%ssetWebhook?drop_pending_updates=%t&%s",
		a.base,
		dropPendingUpdates,
		querify(opts),
	)

	keyVal := map[string]string{"url": webhookURL}
	cnt, err := sendPostForm(url, keyVal)
	if err != nil {
		return res, err
	}

	return res, json.Unmarshal(cnt, &res)
}

// DeleteWebhook is used to remove webhook integration if you decide to switch back to GetUpdates.
func (a API) DeleteWebhook(dropPendingUpdates bool) (APIResponseBase, error) {
	var res APIResponseBase
	var url = fmt.Sprintf(
		"%sdeleteWebhook?drop_pending_updates=%t",
		a.base,
		dropPendingUpdates,
	)

	cnt, err := sendGetRequest(url)
	if err != nil {
		return res, err
	}

	return res, json.Unmarshal(cnt, &res)
}

// GetWebhookInfo is used to get current webhook status.
func (a API) GetWebhookInfo() (APIResponseWebhook, error) {
	var res APIResponseWebhook
	var url = fmt.Sprintf(
		"%sgetWebhookInfo",
		a.base,
	)

	cnt, err := sendGetRequest(url)
	if err != nil {
		return res, err
	}

	return res, json.Unmarshal(cnt, &res)
}

// SendMessage is used to send text messages.
func (a API) SendMessage(text string, chatID int64, opts *MessageOptions) (APIResponseMessage, error) {
	var res APIResponseMessage
	var url = fmt.Sprintf(
		"%ssendMessage?text=%s&chat_id=%d&%s",
		a.base,
		encode(text),
		chatID,
		querify(opts),
	)

	cnt, err := sendGetRequest(url)
	if err != nil {
		return res, err
	}

	return res, json.Unmarshal(cnt, &res)
}

// SendPhoto is used to send photos.
func (a API) SendPhoto(file InputFile, chatID int64, opts *PhotoOptions) (APIResponseMessage, error) {
	var res APIResponseMessage
	var url = fmt.Sprintf(
		"%ssendPhoto?chat_id=%d&%s",
		a.base,
		chatID,
		querify(opts),
	)

	cnt, err := sendFile(file, url, "photo")
	if err != nil {
		return res, err
	}

	return res, json.Unmarshal(cnt, &res)
}

// SendAudio is used to send audio files,
// if you want Telegram clients to display them in the music player.
// Your audio must be in the .MP3 or .M4A format.
func (a API) SendAudio(file InputFile, chatID int64, opts *AudioOptions) (APIResponseMessage, error) {
	var res APIResponseMessage
	var url = fmt.Sprintf(
		"%ssendAudio?chat_id=%d&%s",
		a.base,
		chatID,
		querify(opts),
	)

	cnt, err := sendFile(file, url, "audio")
	if err != nil {
		return res, err
	}

	return res, json.Unmarshal(cnt, &res)
}

// SendDocument is used to send general files.
func (a API) SendDocument(file InputFile, chatID int64, opts *DocumentOptions) (APIResponseMessage, error) {
	var res APIResponseMessage
	var url = fmt.Sprintf(
		"%ssendDocument?chat_id=%d&%s",
		a.base,
		chatID,
		querify(opts),
	)

	cnt, err := sendFile(file, url, "document")
	if err != nil {
		return res, err
	}

	return res, json.Unmarshal(cnt, &res)
}

// SendVideo is used to send video files.
// Telegram clients support mp4 videos (other formats may be sent with SendDocument).
func (a API) SendVideo(file InputFile, chatID int64, opts *VideoOptions) (APIResponseMessage, error) {
	var res APIResponseMessage
	var url = fmt.Sprintf(
		"%ssendVideo?chat_id=%d&%s",
		a.base,
		chatID,
		querify(opts),
	)

	cnt, err := sendFile(file, url, "video")
	if err != nil {
		return res, err
	}

	return res, json.Unmarshal(cnt, &res)
}

// SendAnimation is used to send animation files (GIF or H.264/MPEG-4 AVC video without sound).
func (a API) SendAnimation(file InputFile, chatID int64, opts *AnimationOptions) (APIResponseMessage, error) {
	var res APIResponseMessage
	var url = fmt.Sprintf(
		"%ssendAnimation?chat_id=%d&%s",
		a.base,
		chatID,
		querify(opts),
	)

	cnt, err := sendFile(file, url, "animation")
	if err != nil {
		return res, err
	}

	return res, json.Unmarshal(cnt, &res)
}

// SendVoice is used to send audio files, if you want Telegram clients to display the file as a playable voice message.
// For this to work, your audio must be in an .OGG file encoded with OPUS (other formats may be sent as Audio or Document).
func (a API) SendVoice(file InputFile, chatID int64, opts *VoiceOptions) (APIResponseMessage, error) {
	var res APIResponseMessage
	var url = fmt.Sprintf(
		"%ssendVoice?chat_id=%d&%s",
		a.base,
		chatID,
		querify(opts),
	)

	cnt, err := sendFile(file, url, "voice")
	if err != nil {
		return res, err
	}

	return res, json.Unmarshal(cnt, &res)
}

// SendVideoNote is used to send video messages.
func (a API) SendVideoNote(file InputFile, chatID int64, opts *VideoNoteOptions) (APIResponseMessage, error) {
	var res APIResponseMessage
	var url = fmt.Sprintf(
		"%ssendVideoNote?chat_id=%d&%s",
		a.base,
		chatID,
		querify(opts),
	)

	cnt, err := sendFile(file, url, "video_note")
	if err != nil {
		return res, err
	}

	return res, json.Unmarshal(cnt, &res)
}

// SendContact is used to send phone contacts.
func (a API) SendContact(phoneNumber, firstName string, chatID int64, opts *ContactOptions) (APIResponseMessage, error) {
	var res APIResponseMessage
	var url = fmt.Sprintf(
		"%ssendContact?chat_id=%d&phone_number=%s&first_name=%s&%s",
		a.base,
		chatID,
		encode(phoneNumber),
		encode(firstName),
		querify(opts),
	)

	cnt, err := sendGetRequest(url)
	if err != nil {
		return res, err
	}

	return res, json.Unmarshal(cnt, &res)
}

// SendChatAction is used to tell the user that something is happening on the bot's side.
// The status is set for 5 seconds or less (when a message arrives from your bot, Telegram clients clear its typing status).
func (a API) SendChatAction(action ChatAction, chatID int64) (APIResponseBase, error) {
	var res APIResponseBase
	var url = fmt.Sprintf(
		"%ssendChatAction?chat_id=%d&action=%s",
		a.base,
		chatID,
		action,
	)

	cnt, err := sendGetRequest(url)
	if err != nil {
		return res, err
	}

	return res, json.Unmarshal(cnt, &res)
}

// GetChat is used to get up to date information about the chat.
// (current name of the user for one-on-one conversations, current username of a user, group or channel, etc.)
func (a API) GetChat(chatID int64) (APIResponseChat, error) {
	var res APIResponseChat
	var url = fmt.Sprintf("%sgetChat?chat_id=%d", a.base, chatID)

	cnt, err := sendGetRequest(url)
	if err != nil {
		return res, err
	}

	return res, json.Unmarshal(cnt, &res)
}

// GetChatAdministrators is used to get a list of administrators in a chat.
func (a API) GetChatAdministrators(chatID int64) (APIResponseAdministrators, error) {
	var res APIResponseAdministrators
	var url = fmt.Sprintf(
		"%sgetChatAdministrators?chat_id=%d",
		a.base,
		chatID,
	)

	cnt, err := sendGetRequest(url)
	if err != nil {
		return res, err
	}

	return res, json.Unmarshal(cnt, &res)
}

// GetChatMemberCount is used to get the number of members in a chat.
func (a API) GetChatMemberCount(chatID int64) (APIResponseMemberCount, error) {
	var res APIResponseMemberCount
	var url = fmt.Sprintf(
		"%sgetChatMemberCount?chat_id=%d",
		a.base,
		chatID,
	)

	cnt, err := sendGetRequest(url)
	if err != nil {
		return res, err
	}

	return res, json.Unmarshal(cnt, &res)
}

// GetChatMember is used to get information about a member of a chat.
func (a API) GetChatMember(chatID, userID int64) (APIResponseChatMember, error) {
	var res APIResponseChatMember
	var url = fmt.Sprintf(
		"%sgetChatMember?chat_id=%d&user_id=%d",
		a.base,
		chatID,
		userID,
	)

	cnt, err := sendGetRequest(url)
	if err != nil {
		return res, err
	}

	return res, json.Unmarshal(cnt, &res)
}

// AnswerCallbackQuery is used to send answers to callback queries sent from inline keyboards.
// The answer will be displayed to the user as a notification at the top of the chat screen or as an alert.
func (a API) AnswerCallbackQuery(callbackID string, opts *CallbackQueryOptions) (APIResponseMessage, error) {
	var res APIResponseMessage
	var url = fmt.Sprintf(
		"%sanswerCallbackQuery?callback_query_id=%s&%s",
		a.base,
		callbackID,
		querify(opts),
	)

	cnt, err := sendGetRequest(url)
	if err != nil {
		return res, err
	}

	return res, json.Unmarshal(cnt, &res)
}

// SetMyCommands is used to change the list of the bot's commands for the given scope and user language.
func (a API) SetMyCommands(opts *CommandOptions, commands ...BotCommand) (APIResponseBool, error) {
	var res APIResponseBool
	jsn, _ := json.Marshal(commands)

	var url = fmt.Sprintf(
		"%ssetMyCommands?commands=%s&%s",
		a.base,
		jsn,
		querify(opts),
	)

	cnt, err := sendGetRequest(url)
	if err != nil {
		return res, err
	}

	return res, json.Unmarshal(cnt, &res)
}

// DeleteMyCommands is used to delete the list of the bot's commands for the given scope and user language.
func (a API) DeleteMyCommands(opts *CommandOptions) (APIResponseBool, error) {
	var res APIResponseBool
	var url = fmt.Sprintf(
		"%sdeleteMyCommands?%s",
		a.base,
		querify(opts),
	)

	cnt, err := sendGetRequest(url)
	if err != nil {
		return res, err
	}

	return res, json.Unmarshal(cnt, &res)
}

// GetMyCommands is used to get the current list of the bot's commands for the given scope and user language.
func (a API) GetMyCommands(opts *CommandOptions) (APIResponseCommands, error) {
	var res APIResponseCommands
	var url = fmt.Sprintf(
		"%sgetMyCommands?%s",
		a.base,
		querify(opts),
	)

	cnt, err := sendGetRequest(url)
	if err != nil {
		return res, err
	}

	return res, json.Unmarshal(cnt, &res)
}

// EditMessageText is used to edit text and game messages.
func (a API) EditMessageText(text string, msg MessageIDOptions, opts *MessageTextOptions) (APIResponseMessage, error) {
	var res APIResponseMessage
	var url = fmt.Sprintf(
		"%seditMessageText?text=%s&%s&%s",
		a.base,
		encode(text),
		querify(msg),
		querify(opts),
	)

	cnt, err := sendGetRequest(url)
	if err != nil {
		return res, err
	}

	return res, json.Unmarshal(cnt, &res)
}

// EditMessageCaption is used to edit captions of messages.
func (a API) EditMessageCaption(msg MessageIDOptions, opts *MessageCaptionOptions) (APIResponseMessage, error) {
	var res APIResponseMessage
	var url = fmt.Sprintf(
		"%seditMessageCaption?%s&%s",
		a.base,
		querify(msg),
		querify(opts),
	)

	cnt, err := sendGetRequest(url)
	if err != nil {
		return res, err
	}

	return res, json.Unmarshal(cnt, &res)
}

// EditMessageMedia is used to edit animation, audio, document, photo or video messages.
// If a message is part of a message album, then it can be edited only to an audio for audio albums,
// only to a document for document albums and to a photo or a video otherwise.
// When an inline message is edited, a new file can't be uploaded.
// Use a previously uploaded file via its file_id or specify a URL.
func (a API) EditMessageMedia(msg MessageIDOptions, opts *MessageMediaOptions) (APIResponseMessage, error) {
	var res APIResponseMessage
	var url = fmt.Sprintf(
		"%seditMessageMedia?%s&%s",
		a.base,
		querify(msg),
		querify(opts),
	)

	cnt, err := sendGetRequest(url)
	if err != nil {
		return res, err
	}

	return res, json.Unmarshal(cnt, &res)
}

// EditMessageReplyMarkup is used to edit only the reply markup of messages.
func (a API) EditMessageReplyMarkup(msg MessageIDOptions, opts *MessageReplyMarkup) (APIResponseMessage, error) {
	var res APIResponseMessage
	var url = fmt.Sprintf(
		"%seditMessageReplyMarkup?%s&%s",
		a.base,
		querify(msg),
		querify(opts),
	)

	cnt, err := sendGetRequest(url)
	if err != nil {
		return res, err
	}

	return res, json.Unmarshal(cnt, &res)
}

// DeleteMessage is used to delete a message, including service messages, with the following limitations:
// - A message can only be deleted if it was sent less than 48 hours ago.
// - A dice message in a private chat can only be deleted if it was sent more than 24 hours ago.
// - Bots can delete outgoing messages in private chats, groups, and supergroups.
// - Bots can delete incoming messages in private chats.
// - Bots granted can_post_messages permissions can delete outgoing messages in channels.
// - If the bot is an administrator of a group, it can delete any message there.
// - If the bot has can_delete_messages permission in a supergroup or a channel, it can delete any message there.
func (a API) DeleteMessage(chatID int64, messageID int) (APIResponseBase, error) {
	var res APIResponseBase
	var url = fmt.Sprintf(
		"%sdeleteMessage?chat_id=%d&message_id=%d",
		a.base,
		chatID,
		messageID,
	)

	cnt, err := sendGetRequest(url)
	if err != nil {
		return res, err
	}

	return res, json.Unmarshal(cnt, &res)
}

// GetFile returns the basic info about a file and prepares it for downloading.
// For the moment, bots can download files of up to 20MB in size.
// The file can then be downloaded with DownloadFile where filePath is taken from the response.
// It is guaranteed that the file will be downloadable for at least 1 hour.
// When the download file expires, a new one can be requested by calling GetFile again.
func (a API) GetFile(fileID string) (APIResponseFile, error) {
	var res APIResponseFile

	cnt, err := sendGetRequest(fmt.Sprintf(
		"%sgetFile?file_id=%s",
		a.base,
		fileID,
	))
	if err != nil {
		return res, err
	}

	return res, json.Unmarshal(cnt, &res)
}

// DownloadFile returns the bytes of the file corresponding to the given filePath.
// This function is callable for at least 1 hour since the call to GetFile.
// When the download expires a new one can be requested by calling GetFile again.
func (a API) DownloadFile(filePath string) ([]byte, error) {
	return sendGetRequest(fmt.Sprintf(
		"https://api.telegram.org/file/bot%s/%s",
		a.token,
		filePath,
	))
}
