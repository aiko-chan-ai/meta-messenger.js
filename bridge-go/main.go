package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"

	"messagix-bridge/bridge"
)

type handle uint64

var nextHandle atomic.Uint64
var clients = make(map[handle]*bridge.Client)
var clientsMu sync.RWMutex

func newHandle() handle { return handle(nextHandle.Add(1)) }

type jsonResp struct {
	OK    bool        `json:"ok"`
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func success(data interface{}) *C.char {
	resp := jsonResp{OK: true, Data: data}
	b, _ := json.Marshal(resp)
	return C.CString(string(b))
}

func fail(err error) *C.char {
	resp := jsonResp{OK: false, Error: err.Error()}
	b, _ := json.Marshal(resp)
	return C.CString(string(b))
}

//export MxFreeCString
func MxFreeCString(s *C.char) {
	C.free(unsafe.Pointer(s))
}

//export MxNewClient
func MxNewClient(input *C.char) *C.char {
	var cfg bridge.ClientConfig
	if err := json.Unmarshal([]byte(C.GoString(input)), &cfg); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	client, err := bridge.NewClient(&cfg)
	if err != nil {
		return fail(err)
	}

	h := newHandle()
	client.ID = uint64(h)

	clientsMu.Lock()
	clients[h] = client
	clientsMu.Unlock()

	return success(map[string]interface{}{
		"handle": h,
	})
}

//export MxConnect
func MxConnect(input *C.char) *C.char {
	var payload struct {
		Handle uint64 `json:"handle"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	userInfo, initialData, err := client.Connect()
	if err != nil {
		return fail(err)
	}

	return success(map[string]interface{}{
		"user":        userInfo,
		"initialData": initialData,
	})
}

//export MxConnectE2EE
func MxConnectE2EE(input *C.char) *C.char {
	var payload struct {
		Handle uint64 `json:"handle"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	if err := client.ConnectE2EE(); err != nil {
		return fail(err)
	}

	return success(map[string]interface{}{})
}

//export MxDisconnect
func MxDisconnect(input *C.char) *C.char {
	var payload struct {
		Handle uint64 `json:"handle"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.Lock()
	client := clients[handle(payload.Handle)]
	if client != nil {
		delete(clients, handle(payload.Handle))
	}
	clientsMu.Unlock()

	if client != nil {
		client.Disconnect()
	}

	return success(map[string]interface{}{})
}

//export MxIsConnected
func MxIsConnected(input *C.char) *C.char {
	var payload struct {
		Handle uint64 `json:"handle"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	return success(map[string]interface{}{
		"connected":     client.IsConnected(),
		"e2eeConnected": client.IsE2EEConnected(),
	})
}

//export MxSendMessage
func MxSendMessage(input *C.char) *C.char {
	var payload struct {
		Handle  uint64                    `json:"handle"`
		Options bridge.SendMessageOptions `json:"options"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	result, err := client.SendMessage(&payload.Options)
	if err != nil {
		return fail(err)
	}

	return success(result)
}

//export MxSendReaction
func MxSendReaction(input *C.char) *C.char {
	var payload struct {
		Handle    uint64 `json:"handle"`
		ThreadID  int64  `json:"threadId"`
		MessageID string `json:"messageId"`
		Emoji     string `json:"emoji"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	if err := client.SendReaction(payload.ThreadID, payload.MessageID, payload.Emoji); err != nil {
		return fail(err)
	}

	return success(map[string]interface{}{})
}

//export MxEditMessage
func MxEditMessage(input *C.char) *C.char {
	var payload struct {
		Handle    uint64 `json:"handle"`
		MessageID string `json:"messageId"`
		NewText   string `json:"newText"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	if err := client.EditMessage(payload.MessageID, payload.NewText); err != nil {
		return fail(err)
	}

	return success(map[string]interface{}{})
}

//export MxUnsendMessage
func MxUnsendMessage(input *C.char) *C.char {
	var payload struct {
		Handle    uint64 `json:"handle"`
		MessageID string `json:"messageId"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	if err := client.UnsendMessage(payload.MessageID); err != nil {
		return fail(err)
	}

	return success(map[string]interface{}{})
}

//export MxSendTyping
func MxSendTyping(input *C.char) *C.char {
	var payload struct {
		Handle     uint64 `json:"handle"`
		ThreadID   int64  `json:"threadId"`
		IsTyping   bool   `json:"isTyping"`
		IsGroup    bool   `json:"isGroup"`
		ThreadType int64  `json:"threadType"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	if err := client.SendTypingIndicator(payload.ThreadID, payload.IsTyping, payload.IsGroup, payload.ThreadType); err != nil {
		return fail(err)
	}

	return success(map[string]interface{}{})
}

//export MxMarkRead
func MxMarkRead(input *C.char) *C.char {
	var payload struct {
		Handle      uint64 `json:"handle"`
		ThreadID    int64  `json:"threadId"`
		WatermarkTs int64  `json:"watermarkTs"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	if err := client.MarkRead(payload.ThreadID, payload.WatermarkTs); err != nil {
		return fail(err)
	}

	return success(map[string]interface{}{})
}

//export MxUploadMedia
func MxUploadMedia(input *C.char) *C.char {
	var payload struct {
		Handle  uint64                    `json:"handle"`
		Options bridge.UploadMediaOptions `json:"options"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	result, err := client.UploadMedia(&payload.Options)
	if err != nil {
		return fail(err)
	}

	return success(result)
}

//export MxSendImage
func MxSendImage(input *C.char) *C.char {
	var payload struct {
		Handle  uint64                  `json:"handle"`
		Options bridge.SendImageOptions `json:"options"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	result, err := client.SendImage(&payload.Options)
	if err != nil {
		return fail(err)
	}

	return success(result)
}

//export MxSendVideo
func MxSendVideo(input *C.char) *C.char {
	var payload struct {
		Handle  uint64                  `json:"handle"`
		Options bridge.SendVideoOptions `json:"options"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	result, err := client.SendVideo(&payload.Options)
	if err != nil {
		return fail(err)
	}

	return success(result)
}

//export MxSendVoice
func MxSendVoice(input *C.char) *C.char {
	var payload struct {
		Handle  uint64                  `json:"handle"`
		Options bridge.SendVoiceOptions `json:"options"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	result, err := client.SendVoice(&payload.Options)
	if err != nil {
		return fail(err)
	}

	return success(result)
}

//export MxSendFile
func MxSendFile(input *C.char) *C.char {
	var payload struct {
		Handle  uint64                 `json:"handle"`
		Options bridge.SendFileOptions `json:"options"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	result, err := client.SendFile(&payload.Options)
	if err != nil {
		return fail(err)
	}

	return success(result)
}

//export MxSendSticker
func MxSendSticker(input *C.char) *C.char {
	var payload struct {
		Handle  uint64                    `json:"handle"`
		Options bridge.SendStickerOptions `json:"options"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	result, err := client.SendSticker(&payload.Options)
	if err != nil {
		return fail(err)
	}

	return success(result)
}

//export MxCreateThread
func MxCreateThread(input *C.char) *C.char {
	var payload struct {
		Handle  uint64                     `json:"handle"`
		Options bridge.CreateThreadOptions `json:"options"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	result, err := client.CreateThread(&payload.Options)
	if err != nil {
		return fail(err)
	}

	return success(result)
}

//export MxGetUserInfo
func MxGetUserInfo(input *C.char) *C.char {
	var payload struct {
		Handle  uint64                    `json:"handle"`
		Options bridge.GetUserInfoOptions `json:"options"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	result, err := client.GetUserInfo(&payload.Options)
	if err != nil {
		return fail(err)
	}

	return success(result)
}

//export MxSetGroupPhoto
func MxSetGroupPhoto(input *C.char) *C.char {
	var payload struct {
		Handle   uint64 `json:"handle"`
		ThreadID int64  `json:"threadId"`
		Data     string `json:"data"` // base64 encoded
		MimeType string `json:"mimeType"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	// Decode base64 data
	data, err := base64.StdEncoding.DecodeString(payload.Data)
	if err != nil {
		return fail(fmt.Errorf("invalid base64 data: %w", err))
	}

	if err := client.SetGroupPhoto(&bridge.SetGroupPhotoOptions{
		ThreadID: payload.ThreadID,
		Data:     data,
		MimeType: payload.MimeType,
	}); err != nil {
		return fail(err)
	}

	return success(map[string]interface{}{})
}

//export MxRenameThread
func MxRenameThread(input *C.char) *C.char {
	var payload struct {
		Handle  uint64                     `json:"handle"`
		Options bridge.RenameThreadOptions `json:"options"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	if err := client.RenameThread(&payload.Options); err != nil {
		return fail(err)
	}

	return success(map[string]interface{}{})
}

//export MxMuteThread
func MxMuteThread(input *C.char) *C.char {
	var payload struct {
		Handle  uint64                   `json:"handle"`
		Options bridge.MuteThreadOptions `json:"options"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	if err := client.MuteThread(&payload.Options); err != nil {
		return fail(err)
	}

	return success(map[string]interface{}{})
}

//export MxDeleteThread
func MxDeleteThread(input *C.char) *C.char {
	var payload struct {
		Handle  uint64                     `json:"handle"`
		Options bridge.DeleteThreadOptions `json:"options"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	if err := client.DeleteThread(&payload.Options); err != nil {
		return fail(err)
	}

	return success(map[string]interface{}{})
}

//export MxSearchUsers
func MxSearchUsers(input *C.char) *C.char {
	var payload struct {
		Handle  uint64                    `json:"handle"`
		Options bridge.SearchUsersOptions `json:"options"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	users, err := client.SearchUsers(&payload.Options)
	if err != nil {
		return fail(err)
	}

	return success(map[string]interface{}{
		"users": users,
	})
}

//export MxPollEvents
func MxPollEvents(input *C.char) *C.char {
	var payload struct {
		Handle    uint64 `json:"handle"`
		TimeoutMs int    `json:"timeoutMs"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	select {
	case evt, ok := <-client.Events():
		if !ok {
			return success(map[string]interface{}{
				"type": "closed",
			})
		}
		return success(evt)
	default:
		return success(map[string]interface{}{
			"type": "timeout",
		})
	}
}

// E2EE functions

//export MxSendE2EEMessage
func MxSendE2EEMessage(input *C.char) *C.char {
	var payload struct {
		Handle           uint64 `json:"handle"`
		ChatJID          string `json:"chatJid"`
		Text             string `json:"text"`
		ReplyToID        string `json:"replyToId,omitempty"`
		ReplyToSenderJID string `json:"replyToSenderJid,omitempty"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	result, err := client.SendMessage(&bridge.SendMessageOptions{
		Text:                 payload.Text,
		IsE2EE:               true,
		E2EEChatJID:          payload.ChatJID,
		E2EEReplyToID:        payload.ReplyToID,
		E2EEReplyToSenderJID: payload.ReplyToSenderJID,
	})
	if err != nil {
		return fail(err)
	}

	return success(result)
}

//export MxSendE2EEReaction
func MxSendE2EEReaction(input *C.char) *C.char {
	var payload struct {
		Handle    uint64 `json:"handle"`
		ChatJID   string `json:"chatJid"`
		MessageID string `json:"messageId"`
		SenderJID string `json:"senderJid"`
		Emoji     string `json:"emoji"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	if err := client.SendE2EEReaction(payload.ChatJID, payload.MessageID, payload.SenderJID, payload.Emoji); err != nil {
		return fail(err)
	}

	return success(map[string]interface{}{})
}

//export MxSendE2EETyping
func MxSendE2EETyping(input *C.char) *C.char {
	var payload struct {
		Handle   uint64 `json:"handle"`
		ChatJID  string `json:"chatJid"`
		IsTyping bool   `json:"isTyping"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	if err := client.SendE2EETyping(payload.ChatJID, payload.IsTyping); err != nil {
		return fail(err)
	}

	return success(map[string]interface{}{})
}

//export MxEditE2EEMessage
func MxEditE2EEMessage(input *C.char) *C.char {
	var payload struct {
		Handle    uint64 `json:"handle"`
		ChatJID   string `json:"chatJid"`
		MessageID string `json:"messageId"`
		NewText   string `json:"newText"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	if err := client.EditE2EEMessage(payload.ChatJID, payload.MessageID, payload.NewText); err != nil {
		return fail(err)
	}

	return success(map[string]interface{}{})
}

//export MxUnsendE2EEMessage
func MxUnsendE2EEMessage(input *C.char) *C.char {
	var payload struct {
		Handle    uint64 `json:"handle"`
		ChatJID   string `json:"chatJid"`
		MessageID string `json:"messageId"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	if err := client.UnsendE2EEMessage(payload.ChatJID, payload.MessageID); err != nil {
		return fail(err)
	}

	return success(map[string]interface{}{})
}

//export MxGetDeviceData
func MxGetDeviceData(input *C.char) *C.char {
	var payload struct {
		Handle uint64 `json:"handle"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	if client.DeviceStore == nil {
		return fail(fmt.Errorf("device store not initialized"))
	}

	data, err := client.DeviceStore.GetDeviceData()
	if err != nil {
		return fail(err)
	}

	return success(map[string]interface{}{
		"deviceData": data,
	})
}

// ==================== E2EE Media Functions ====================

//export MxSendE2EEImage
func MxSendE2EEImage(input *C.char) *C.char {
	var payload struct {
		Handle  uint64                      `json:"handle"`
		Options bridge.SendE2EEImageOptions `json:"options"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	result, err := client.SendE2EEImage(&payload.Options)
	if err != nil {
		return fail(err)
	}

	return success(result)
}

//export MxSendE2EEVideo
func MxSendE2EEVideo(input *C.char) *C.char {
	var payload struct {
		Handle  uint64                      `json:"handle"`
		Options bridge.SendE2EEVideoOptions `json:"options"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	result, err := client.SendE2EEVideo(&payload.Options)
	if err != nil {
		return fail(err)
	}

	return success(result)
}

//export MxSendE2EEAudio
func MxSendE2EEAudio(input *C.char) *C.char {
	var payload struct {
		Handle  uint64                      `json:"handle"`
		Options bridge.SendE2EEAudioOptions `json:"options"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	result, err := client.SendE2EEAudio(&payload.Options)
	if err != nil {
		return fail(err)
	}

	return success(result)
}

//export MxSendE2EEDocument
func MxSendE2EEDocument(input *C.char) *C.char {
	var payload struct {
		Handle  uint64                         `json:"handle"`
		Options bridge.SendE2EEDocumentOptions `json:"options"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	result, err := client.SendE2EEDocument(&payload.Options)
	if err != nil {
		return fail(err)
	}

	return success(result)
}

//export MxSendE2EESticker
func MxSendE2EESticker(input *C.char) *C.char {
	var payload struct {
		Handle  uint64                        `json:"handle"`
		Options bridge.SendE2EEStickerOptions `json:"options"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	result, err := client.SendE2EESticker(&payload.Options)
	if err != nil {
		return fail(err)
	}

	return success(result)
}

//export MxDownloadE2EEMedia
func MxDownloadE2EEMedia(input *C.char) *C.char {
	var payload struct {
		Handle  uint64                          `json:"handle"`
		Options bridge.DownloadE2EEMediaOptions `json:"options"`
	}
	if err := json.Unmarshal([]byte(C.GoString(input)), &payload); err != nil {
		return fail(fmt.Errorf("invalid json: %w", err))
	}

	clientsMu.RLock()
	client := clients[handle(payload.Handle)]
	clientsMu.RUnlock()
	if client == nil {
		return fail(fmt.Errorf("client not found"))
	}

	result, err := client.DownloadE2EEMedia(&payload.Options)
	if err != nil {
		return fail(err)
	}

	// Encode data as base64 for JSON transport
	return success(map[string]interface{}{
		"data":     base64.StdEncoding.EncodeToString(result.Data),
		"mimeType": result.MimeType,
		"fileSize": result.FileSize,
	})
}

func main() {}
