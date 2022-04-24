/*
 * Echosphere
 * Copyright (C) 2022 The Echosphere Devs
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
)

// WebAppInfo contains information about a Web App.
type WebAppInfo struct {
	URL string `json:"url"`
}

// SentWebAppMessage contains information about an inline message sent
// by a Web App on behalf of a user.
type SentWebAppMessage struct {
	InlineMessageID string `json:"inline_message_id,omitempty"`
}

// WebAppData contains data sent from a Web App to the bot.
type WebAppData struct {
	Data       string `json:"data"`
	ButtonText string `json:"button_text"`
}

// AnswerWebAppQuery is used to set the result of an interaction with a Web App
// and send a corresponding message on behalf of the user to the chat from which
// the query originated.
func (a API) AnswerWebAppQuery(webAppQueryID string, result InlineQueryResult) (res APIResponseSentWebAppMessage, err error) {
	resultJson, err := json.Marshal(result)
	if err != nil {
		return
	}

	var url = fmt.Sprintf(
		"%sanswerWebAppQuery?web_app_query_id=%s&result=%s",
		a.base,
		webAppQueryID,
		encode(string(resultJson)),
	)

	cnt, err := sendGetRequest(url)
	if err != nil {
		return
	}

	if err = json.Unmarshal(cnt, &res); err != nil {
		return
	}

	err = check(res)
	return
}
