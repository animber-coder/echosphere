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
	"fmt"
	"testing"
	"time"
)

var (
	stickerFile    *File
	stickerSet     *StickerSet
	stickerSetName = fmt.Sprintf("set%d_by_echosphere_coverage_bot", time.Now().Unix())
)

func TestUploadStickerFile(t *testing.T) {
	resp, err := api.UploadStickerFile(
		chatID,
		StickerFile{
			NewInputFilePath("assets/tests/echosphere_test.png"),
			PNGSticker,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	stickerFile = resp.Result
}

func TestCreateNewStickerSet(t *testing.T) {
	_, err := api.CreateNewStickerSet(
		chatID,
		stickerSetName,
		"Echosphere Coverage Pack",
		"🤖",
		StickerFile{
			NewInputFileID(stickerFile.FileID),
			PNGSticker,
		},
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}
}

func TestAddStickerToSet(t *testing.T) {
	_, err := api.AddStickerToSet(
		chatID,
		stickerSetName,
		"🤖",
		StickerFile{
			NewInputFilePath("assets/tests/echosphere_sticker.png"),
			PNGSticker,
		},
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}
}

func TestGetStickerSet(t *testing.T) {
	resp, err := api.GetStickerSet(
		stickerSetName,
	)

	if err != nil {
		t.Fatal(err)
	}

	stickerSet = resp.Result
}

func TestSetStickerPositionInSet(t *testing.T) {
	_, err := api.SetStickerPositionInSet(
		stickerSet.Stickers[1].FileID,
		0,
	)

	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteStickerFromSet(t *testing.T) {
	_, err := api.DeleteStickerFromSet(
		stickerSet.Stickers[0].FileID,
	)

	if err != nil {
		t.Fatal(err)
	}
}

func TestSendSticker(t *testing.T) {
	_, err := api.SendSticker(
		stickerSet.Stickers[0].FileID,
		chatID,
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}
}

func TestSetStickerSetThumb(t *testing.T) {
	_, err := api.SetStickerSetThumb(
		stickerSetName,
		chatID,
		NewInputFilePath("assets/tests/echosphere_thumb.png"),
	)

	if err != nil {
		t.Fatal(err)
	}
}
