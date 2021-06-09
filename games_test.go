/*
 * Echosphere
 * Copyright (C) 2021  The Echosphere Devs
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

import "testing"

var gameMsgTmp *Message

func TestSendGame(t *testing.T) {
	resp, err := api.SendGame(
		"echosphere_coverage_game",
		chatID,
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	if !resp.Ok {
		t.Fatalf("%d %s", resp.ErrorCode, resp.Description)
	}

	gameMsgTmp = resp.Result
}

func TestSetGameScore(t *testing.T) {
	resp, err := api.SetGameScore(
		chatID,
		545,
		NewMessageID(chatID, gameMsgTmp.ID),
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	if !resp.Ok {
		t.Fatalf("%d %s", resp.ErrorCode, resp.Description)
	}
}

func TestGameHighScores(t *testing.T) {
	resp, err := api.GetGameHighScores(
		chatID,
		NewMessageID(chatID, gameMsgTmp.ID),
	)

	if err != nil {
		t.Fatal(err)
	}

	if !resp.Ok {
		t.Fatalf("%d %s", resp.ErrorCode, resp.Description)
	}
}
