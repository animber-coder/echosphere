package echosphere

import (
	"testing"
)

var (
	msgTmp     *Message
	api        = NewAPI("1713461126:AAEV5sgVo513Vz4PT33mpp0ZykJqrnSluzM")
	chatID     = int64(41876271)
	photoID    = "AgACAgQAAxkDAAMrYFtODxV2LL6-kR_6qSbG9n8dIOIAAti1MRug29lSkNq_9o8PC5uMd7EnXQADAQADAgADbQADeooGAAEeBA"
	audioID    = "CQACAgQAAxkDAAIBCmBbamz_DqKk2GmrzmoM0SrzRN6wAAK9CAACoNvZUgPyk-87OM_YHgQ"
	documentID = "BQACAgQAAxkDAANmYFtSXcF5kTtwgHeqVUngyuuJMx4AAnQIAAKg29lSb4HP4x-qMT8eBA"
	videoID    = "BAACAgQAAxkDAANxYFtaxF1kfc7nVY_Mtfba3u5dMooAAoYIAAKg29lSpwABJrcveXZlHgQ"
	voiceID    = "AwACAgQAAxkDAAPXYFtmoFriwJFVGDgPPpfUBljgnYAAAq8IAAKg29lStEWfrNMMAxgeBA"
)

func TestGetChat(t *testing.T) {
	resp, err := api.GetChat(chatID)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Result.Type != "private" && resp.Result.Type != "group" &&
		resp.Result.Type != "supergroup" && resp.Result.Type != "channel" {

		t.Fatalf("wrong chat type, got: %s", resp.Result.Type)
	}
}

func TestGetStickerSet(t *testing.T) {
	resp, err := api.GetStickerSet("RickAndMorty")
	if err != nil {
		t.Fatal(err)
	}

	if !resp.Ok {
		t.Fatal(resp.ErrorCode, resp.Description)
	}
}

func TestSendMessage(t *testing.T) {
	resp, err := api.SendMessage("TestSendMessage", chatID)
	if err != nil {
		t.Fatal(err)
	}

	if !resp.Ok {
		t.Fatal(resp.ErrorCode, resp.Description)
	}
	msgTmp = resp.Result
}

func TestSendMessageReply(t *testing.T) {
	resp, err := api.SendMessageReply("TestSendMessageReply", chatID, msgTmp.ID)
	if err != nil {
		t.Fatal(err)
	}

	if !resp.Ok {
		t.Fatal(resp.ErrorCode, resp.Description)
	}
}

func TestSendMessageWithKeyboard(t *testing.T) {
	kbd := api.KeyboardMarkup(false, true, false,
		api.KeyboardRow(
			api.KeyboardButton("test 1", false, false),
			api.KeyboardButton("test 2", false, false),
		),
		api.KeyboardRow(
			api.KeyboardButton("test 3", false, false),
			api.KeyboardButton("test 4", false, false),
		),
	)

	resp, err := api.SendMessageWithKeyboard("TestSendMessageWithKeyboard", chatID, kbd)
	if err != nil {
		t.Fatal(err)
	}

	if !resp.Ok {
		t.Fatal(resp.ErrorCode, resp.Description)
	}
}

func TestDeleteMessage(t *testing.T) {
	resp, err := api.DeleteMessage(chatID, msgTmp.ID)
	if err != nil {
		t.Fatal(err)
	}

	if !resp.Ok {
		t.Fatal(resp.ErrorCode, resp.Description)
	}
}

func TestSendPhoto(t *testing.T) {
	resp, err := api.SendPhoto("tests/echosphere_test.png", "TestSendPhoto", chatID)
	if err != nil {
		t.Fatal(err)
	}

	if !resp.Ok {
		t.Fatal(resp.ErrorCode, resp.Description)
	}
}

func TestSendPhotoByID(t *testing.T) {
	resp, err := api.SendPhotoByID(photoID, "TestSendPhotoByID", chatID)
	if err != nil {
		t.Fatal(err)
	}

	if !resp.Ok {
		t.Fatal(resp.ErrorCode, resp.Description)
	}
}

func TestSendPhotoWithKeyboard(t *testing.T) {
	kbd := api.KeyboardMarkup(false, true, false,
		api.KeyboardRow(
			api.KeyboardButton("test 1", false, false),
			api.KeyboardButton("test 2", false, false),
		),
		api.KeyboardRow(
			api.KeyboardButton("test 3", false, false),
			api.KeyboardButton("test 4", false, false),
		),
	)

	resp, err := api.SendPhotoWithKeyboard("tests/echosphere_test.png", "TestSendPhotoWithKeyboard", chatID, kbd)
	if err != nil {
		t.Fatal(err)
	}

	if !resp.Ok {
		t.Fatal(resp.ErrorCode, resp.Description)
	}
}

func TestSendAudio(t *testing.T) {
	resp, err := api.SendAudio("tests/audio.mp3", "TestSendAudio", chatID)
	if err != nil {
		t.Fatal(err)
	}

	if !resp.Ok {
		t.Fatal(resp.ErrorCode, resp.Description)
	}
}

func TestSendAudioByID(t *testing.T) {
	resp, err := api.SendAudioByID(audioID, "TestSendAudioByID", chatID)
	if err != nil {
		t.Fatal(err)
	}

	if !resp.Ok {
		t.Fatal(resp.ErrorCode, resp.Description)
	}
}

func TestSendDocument(t *testing.T) {
	resp, err := api.SendDocument("tests/document.pdf", "TestSendDocument", chatID)
	if err != nil {
		t.Fatal(err)
	}

	if !resp.Ok {
		t.Fatal(resp.ErrorCode, resp.Description)
	}
}

func TestSendDocumentByID(t *testing.T) {
	resp, err := api.SendDocumentByID(documentID, "TestSendDocumentByID", chatID)
	if err != nil {
		t.Fatal(err)
	}

	if !resp.Ok {
		t.Fatal(resp.ErrorCode, resp.Description)
	}
}

func TestSendVideo(t *testing.T) {
	resp, err := api.SendVideo("tests/video.webm", "TestSendVideo", chatID)
	if err != nil {
		t.Fatal(err)
	}

	if !resp.Ok {
		t.Fatal(resp.ErrorCode, resp.Description)
	}
}

func TestSendVideoByID(t *testing.T) {
	resp, err := api.SendVideoByID(videoID, "TestSendVideoByID", chatID)
	if err != nil {
		t.Fatal(err)
	}

	if !resp.Ok {
		t.Fatal(resp.ErrorCode, resp.Description)
	}
}

// func TestSendVideoNote(t *testing.T) {
// 	resp, err := api.SendVideoNote(videoID, chatID)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if !resp.Ok {
// 		t.Fatal(resp.ErrorCode, resp.Description)
// 	}
// }

func TestSendVoice(t *testing.T) {
	resp, err := api.SendVoice("tests/audio.mp3", "TestSendVoice", chatID)
	if err != nil {
		t.Fatal(err)
	}

	if !resp.Ok {
		t.Fatal(resp.ErrorCode, resp.Description)
	}
}

func TestSendVoiceByID(t *testing.T) {
	resp, err := api.SendVoiceByID(voiceID, "TestSendVoiceByID", chatID)
	if err != nil {
		t.Fatal(err)
	}

	if !resp.Ok {
		t.Fatal(resp.ErrorCode, resp.Description)
	}
}

func TestSendContact(t *testing.T) {
	resp, err := api.SendContact("1234567890", "Name", "Surname", chatID)
	if err != nil {
		t.Fatal(err)
	}

	if !resp.Ok {
		t.Fatal(resp.ErrorCode, resp.Description)
	}
}
