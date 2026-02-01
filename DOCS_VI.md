# Tài liệu API

> [!TIP]
> Thư viện được viết theo style Schmavery/facebook-chat-api (Không sử dụng mã nguồn) để quen thuộc và dễ dùng hơn (và không dùng callback).

* [Bảo mật cookies](#bảo-mật-cookies)
* [Client](#client)
  * [`new Client(cookies, options)`](#constructor)
  * [`client.connect()`](#connect)
  * [`client.disconnect()`](#disconnect)
  * [Thuộc tính](#thuộc-tính)
* [Tin nhắn thường](#tin-nhắn-thường)
  * [`client.sendMessage()`](#sendMessage)
  * [`client.sendReaction()`](#sendReaction)
  * [`client.editMessage()`](#editMessage)
  * [`client.unsendMessage()`](#unsendMessage)
  * [`client.sendTypingIndicator()`](#sendTypingIndicator)
  * [`client.markAsRead()`](#markAsRead)
* [Media](#media)
  * [`client.sendImage()`](#sendImage)
  * [`client.sendVideo()`](#sendVideo)
  * [`client.sendVoice()`](#sendVoice)
  * [`client.sendFile()`](#sendFile)
  * [`client.sendSticker()`](#sendSticker)
  * [`client.uploadMedia()`](#uploadMedia)
* [Quản lý Thread/Group](#quản-lý-threadgroup)
  * [`client.createThread()`](#createThread)
  * [`client.renameThread()`](#renameThread)
  * [`client.setGroupPhoto()`](#setGroupPhoto)
  * [`client.muteThread()`](#muteThread)
  * [`client.unmuteThread()`](#unmuteThread)
  * [`client.deleteThread()`](#deleteThread)
* [Thông tin User](#thông-tin-user)
  * [`client.getUserInfo()`](#getUserInfo)
  * [`client.searchUsers()`](#searchUsers)
* [E2EE (Mã hóa đầu cuối)](#e2ee-mã-hóa-đầu-cuối)
  * [`client.connectE2EE()`](#connectE2EE)
  * [`client.sendE2EEMessage()`](#sendE2EEMessage)
  * [`client.sendE2EEReaction()`](#sendE2EEReaction)
  * [`client.sendE2EETyping()`](#sendE2EETyping)
  * [`client.editE2EEMessage()`](#editE2EEMessage)
  * [`client.unsendE2EEMessage()`](#unsendE2EEMessage)
* [E2EE Media](#e2ee-media)
  * [`client.sendE2EEImage()`](#sendE2EEImage)
  * [`client.sendE2EEVideo()`](#sendE2EEVideo)
  * [`client.sendE2EEAudio()`](#sendE2EEAudio)
  * [`client.sendE2EEDocument()`](#sendE2EEDocument)
  * [`client.sendE2EESticker()`](#sendE2EESticker)
  * [`client.downloadE2EEMedia()`](#downloadE2EEMedia)
  * [`client.getDeviceData()`](#getDeviceData)
* [Khác](#khác)
  * [`client.unloadLibrary()`](#unloadLibrary)
* [Utilities](#utilities)
  * [`Utils.parseCookies()`](#parseCookies)
  * [`Utils.validate()`](#validate)
  * [`Utils.getMissing()`](#getMissing)
* [Events](#events)
  * [`message`](#event-message) 🔵
  * [`e2eeMessage`](#event-e2eeMessage) 🟢
  * [`messageEdit`](#event-messageEdit) 🔵🟢
  * [`messageUnsend`](#event-messageUnsend) 🔵🟢
  * [`reaction`](#event-reaction) 🔵
  * [`e2eeReaction`](#event-e2eeReaction) 🟢
  * [`typing`](#event-typing) 🔵
  * [`readReceipt`](#event-readReceipt) 🔵
  * [`e2eeReceipt`](#event-e2eeReceipt) 🟢
  * [`e2eeConnected`](#event-e2eeConnected) 🟢
  * [`fullyReady`](#event-fullyReady) 🔵🟢
  * [`disconnected`](#event-disconnected) 🔵🟢
* [Types](#types)

---

## Bảo mật cookies

**Đọc kỹ phần này** trước khi copy+paste các ví dụ bên dưới.

Bạn **không nên** lưu cookies trực tiếp trong code. Lý do:
* Người khác có thể nhìn thấy code của bạn và lấy được cookies
* Backup code có thể bị đọc bởi người khác
* Bạn không thể push code lên Github mà không xóa cookies. Lưu ý: Ngay cả khi bạn undo commit chứa cookies, Git vẫn lưu lại và có thể đọc được
* Nếu bạn thay đổi cookies trong tương lai, bạn phải sửa tất cả các nơi trong code

Cách khuyến nghị là lưu cookies vào file riêng:

```typescript
import { readFileSync } from 'fs'
import { Utils } from 'meta-messenger.js'

const cookies = Utils.parseCookies(readFileSync('cookies.json', 'utf-8'))
```

Hoặc sử dụng biến môi trường:
```typescript
const cookies = {
    c_user: process.env.FB_C_USER,
    xs: process.env.FB_XS,
    fr: process.env.FB_FR,
    datr: process.env.FB_DATR,
    // các cookie khác...
}
```

---

# Client

<a name="constructor"></a>
## new Client(cookies, options?)

Tạo một client mới để kết nối đến Messenger.

__Tham số__

* `cookies`: Object chứa các cookies cần thiết (`c_user`, `xs`, `datr`, `fr`, ...).
* `options` (tùy chọn): Object cấu hình:
  * `platform`: `'facebook'` | `'messenger'` | `'instagram'` - Cookie của nền tảng nào [Thư viện hiện chỉ kiểm tra với `'facebook'`] (mặc định: `'facebook'`)
  * `enableE2EE`: Boolean - Bật mã hóa đầu cuối (Cho Messenger) (mặc định: `true`)
  * `devicePath`: String - Đường dẫn file lưu device data (cho E2EE)
  * `deviceData`: String - Device data đã lưu (JSON string) (Được ưu tiên sử dụng)
  * `logLevel`: `'none'` | `'error'` | `'warn'` | `'info'` | `'debug'` | `'trace'` (mặc định: `'none'`)
  * `autoReconnect`: Boolean - Tự động reconnect khi mất kết nối (mặc định: `true`)

__Ví dụ__

```typescript
import { Client } from 'meta-messenger.js'

const cookies = {
    c_user: '100000000000000',
    xs: '48:abc123...',
    datr: 'xyz789...',
    fr: '1QO0u...'
}

const client = new Client(cookies)
```

---

<a name="connect"></a>
## client.connect()

Kết nối đến Messenger. Trả về Promise với thông tin user và dữ liệu ban đầu.

__Trả về__

Promise<{ user: User, initialData: InitialData }>

* `user`: Thông tin người dùng đã đăng nhập
  * `id`: number - Facebook ID
  * `name`: string - Tên hiển thị
  * `username`: string - Username
* `initialData`: Dữ liệu ban đầu
  * `threads`: Thread[] - Danh sách thread gần đây
  * `messages`: Message[] - Tin nhắn gần đây

__Ví dụ__

```typescript
const { user, initialData } = await client.connect()
console.log(`Đã đăng nhập: ${user.name} (${user.id})`)
console.log(`Số threads: ${initialData.threads.length}`)
```

---

<a name="disconnect"></a>
## client.disconnect()

Ngắt kết nối khỏi Messenger.

__Ví dụ__

```typescript
await client.disconnect()
console.log('Đã ngắt kết nối')
```

---

<a name="thuộc-tính"></a>
## Thuộc tính

<a name="user"></a>
### client.user

Thông tin người dùng đã đăng nhập. `null` nếu chưa kết nối.

__Type:__ `User | null`

---

<a name="currentUserId"></a>
### client.currentUserId

Facebook ID của người dùng hiện tại. `null` nếu chưa kết nối.

__Type:__ `number | null`

---

<a name="initialData"></a>
### client.initialData

Dữ liệu ban đầu (threads và messages). `null` nếu chưa kết nối.

__Type:__ `InitialData | null`

---

<a name="isConnected"></a>
### client.isConnected

Kiểm tra client có đang kết nối không.

__Type:__ `boolean`

---

<a name="isE2EEConnected"></a>
### client.isE2EEConnected

Kiểm tra E2EE đã kết nối chưa.

__Type:__ `boolean`

---

<a name="isFullyReady"></a>
### client.isFullyReady

Kiểm tra client đã hoàn toàn sẵn sàng (socket + E2EE nếu enabled).

__Type:__ `boolean`

---

# Tin nhắn thường

<a name="sendMessage"></a>
## client.sendMessage(threadId, options)

Gửi tin nhắn văn bản đến một thread.

__Tham số__

* `threadId`: number - ID của thread.
* `options`: string | SendMessageOptions
  * Nếu là string: Gửi tin nhắn văn bản đơn giản
  * Nếu là object:
    * `text`: string - Nội dung tin nhắn
    * `replyToId?`: string - ID tin nhắn để reply
    * `mentions?`: Mention[] - Danh sách mention
      * `userId`: number - ID user được mention
      * `offset`: number - Vị trí bắt đầu trong text
      * `length`: number - Độ dài của mention

__Trả về__

Promise<SendMessageResult>
* `messageId`: string - ID tin nhắn đã gửi
* `timestampMs`: number - Timestamp (milliseconds)

__Ví dụ__

```typescript
// Tin nhắn đơn giản
await client.sendMessage(threadId, 'Xin chào!')

// Tin nhắn với reply
await client.sendMessage(threadId, {
    text: 'Đây là reply',
    replyToId: 'mid.$abc123'
})

// Tin nhắn với mention
await client.sendMessage(threadId, {
    text: 'Chào @bạn!',
    mentions: [{
        userId: 100000000000001,
        offset: 5,
        length: 4
    }]
})
```

---

<a name="sendReaction"></a>
## client.sendReaction(threadId, messageId, emoji?)

Gửi hoặc xóa reaction cho một tin nhắn.

__Tham số__

* `threadId`: number - ID của thread
* `messageId`: string - ID tin nhắn cần react
* `emoji?`: string - Emoji reaction (bỏ qua để xóa reaction)

__Ví dụ__

```typescript
// Thêm reaction
await client.sendReaction(threadId, messageId, '👍')

// Xóa reaction
await client.sendReaction(threadId, messageId)
```

---

<a name="editMessage"></a>
## client.editMessage(messageId, newText)

Chỉnh sửa một tin nhắn đã gửi.

__Tham số__

* `messageId`: string - ID tin nhắn cần chỉnh sửa
* `newText`: string - Nội dung mới

__Ví dụ__

```typescript
await client.editMessage('mid.$abc123', 'Nội dung đã sửa')
```

---

<a name="unsendMessage"></a>
## client.unsendMessage(messageId)

Thu hồi (xóa) một tin nhắn đã gửi.

__Tham số__

* `messageId`: string - ID tin nhắn cần thu hồi

__Ví dụ__

```typescript
await client.unsendMessage('mid.$abc123')
```

---

<a name="sendTypingIndicator"></a>
## client.sendTypingIndicator(threadId, isTyping?, isGroup?)

Gửi trạng thái đang nhập.

__Tham số__

* `threadId`: number - ID của thread
* `isTyping?`: boolean - `true` để bắt đầu, `false` để dừng (mặc định: `true`)
* `isGroup?`: boolean - `true` nếu là group chat (mặc định: `false`)

__Ví dụ__

```typescript
// Bắt đầu typing
await client.sendTypingIndicator(threadId, true)

// Dừng typing sau 2 giây
setTimeout(async () => {
    await client.sendTypingIndicator(threadId, false)
}, 2000)
```

---

<a name="markAsRead"></a>
## client.markAsRead(threadId, watermarkTs?)

Đánh dấu đã đọc một thread.

__Tham số__

* `threadId`: number - ID của thread
* `watermarkTs?`: number - Timestamp watermark (mặc định: hiện tại)

__Ví dụ__

```typescript
await client.markAsRead(threadId)
```

---

# Media

<a name="sendImage"></a>
## client.sendImage(threadId, data, filename, caption?)

Gửi ảnh.

__Tham số__

* `threadId`: number - ID của thread
* `data`: Buffer - Dữ liệu ảnh
* `filename`: string - Tên file
* `caption?`: string - Caption (tùy chọn)

__Trả về__

Promise<SendMessageResult>

__Ví dụ__

```typescript
import { readFileSync } from 'fs'

const image = readFileSync('photo.jpg')
await client.sendImage(threadId, image, 'photo.jpg', 'Ảnh đẹp!')
```

---

<a name="sendVideo"></a>
## client.sendVideo(threadId, data, filename, caption?)

Gửi video.

__Tham số__

* `threadId`: number - ID của thread
* `data`: Buffer - Dữ liệu video
* `filename`: string - Tên file
* `caption?`: string - Caption (tùy chọn)

__Trả về__

Promise<SendMessageResult>

__Ví dụ__

```typescript
const video = readFileSync('video.mp4')
await client.sendVideo(threadId, video, 'video.mp4', 'Video hay!')
```

---

<a name="sendVoice"></a>
## client.sendVoice(threadId, data, filename)

Gửi tin nhắn thoại.

__Tham số__

* `threadId`: number - ID của thread
* `data`: Buffer - Dữ liệu audio
* `filename`: string - Tên file

__Trả về__

Promise<SendMessageResult>

__Ví dụ__

```typescript
const voice = readFileSync('voice.mp3')
await client.sendVoice(threadId, voice, 'voice.mp3')
```

---

<a name="sendFile"></a>
## client.sendFile(threadId, data, filename, mimeType, caption?)

Gửi file bất kỳ.

__Tham số__

* `threadId`: number - ID của thread
* `data`: Buffer - Dữ liệu file
* `filename`: string - Tên file
* `mimeType`: string - MIME type (ví dụ: 'application/pdf')
* `caption?`: string - Caption (tùy chọn)

__Trả về__

Promise<SendMessageResult>

__Ví dụ__

```typescript
const pdf = readFileSync('document.pdf')
await client.sendFile(threadId, pdf, 'document.pdf', 'application/pdf', 'Tài liệu')
```

---

<a name="sendSticker"></a>
## client.sendSticker(threadId, stickerId)

Gửi sticker.

__Tham số__

* `threadId`: number - ID của thread
* `stickerId`: number - ID của sticker

__Trả về__

Promise<SendMessageResult>

__Ví dụ__

```typescript
// Gửi sticker thumbs up
await client.sendSticker(threadId, 369239263222822)
```

---

<a name="uploadMedia"></a>
## client.uploadMedia(threadId, data, filename, mimeType)

Upload media và lấy ID để sử dụng sau.

__Tham số__

* `threadId`: number - ID của thread
* `data`: Buffer - Dữ liệu file
* `filename`: string - Tên file
* `mimeType`: string - MIME type

__Trả về__

Promise<UploadMediaResult>
* `fbId`: number - Facebook ID của media
* `filename`: string - Tên file

__Ví dụ__

```typescript
const image = readFileSync('photo.jpg')
const result = await client.uploadMedia(threadId, image, 'photo.jpg', 'image/jpeg')
console.log(`Uploaded: ${result.fbId}`)
```

---

# Quản lý Thread/Group

<a name="createThread"></a>
## client.createThread(userId)

Tạo thread 1:1 với một user.

__Tham số__

* `userId`: number - ID của user

__Trả về__

Promise<CreateThreadResult>
* `threadId`: number - ID của thread mới

__Ví dụ__

```typescript
const { threadId } = await client.createThread(100000000000001)
await client.sendMessage(threadId, 'Xin chào!')
```

---

<a name="renameThread"></a>
## client.renameThread(threadId, newName)

Đổi tên group chat.

__Tham số__

* `threadId`: number - ID của thread
* `newName`: string - Tên mới

__Ví dụ__

```typescript
await client.renameThread(threadId, 'Nhóm bạn thân')
```

---

<a name="setGroupPhoto"></a>
## client.setGroupPhoto(threadId, data, mimeType?)

Đổi ảnh đại diện của group.

__Tham số__

* `threadId`: number - ID của thread
* `data`: Buffer | string - Dữ liệu ảnh (Buffer hoặc base64 string)
* `mimeType?`: string - MIME type (mặc định: 'image/jpeg')

__Lưu ý__

Messenger web không hỗ trợ xóa ảnh group, chỉ có thể thay đổi.

__Ví dụ__

```typescript
const photo = readFileSync('group-photo.jpg')
await client.setGroupPhoto(threadId, photo, 'image/jpeg')
```

---

<a name="muteThread"></a>
## client.muteThread(threadId, seconds?)

Tắt thông báo của thread.

__Tham số__

* `threadId`: number - ID của thread
* `seconds?`: number - Thời gian tắt (giây)
  * `-1`: Tắt vĩnh viễn (mặc định)
  * `0`: Bật lại thông báo
  * `> 0`: Tắt trong khoảng thời gian

__Ví dụ__

```typescript
// Tắt vĩnh viễn
await client.muteThread(threadId)

// Tắt trong 1 giờ
await client.muteThread(threadId, 3600)
```

---

<a name="unmuteThread"></a>
## client.unmuteThread(threadId)

Bật lại thông báo của thread.

__Tham số__

* `threadId`: number - ID của thread

__Ví dụ__

```typescript
await client.unmuteThread(threadId)
```

---

<a name="deleteThread"></a>
## client.deleteThread(threadId)

Xóa thread.

__Tham số__

* `threadId`: number - ID của thread

__Cảnh báo__

Hành động này không thể hoàn tác!

__Ví dụ__

```typescript
await client.deleteThread(threadId)
```

---

# Thông tin User

<a name="getUserInfo"></a>
## client.getUserInfo(userId)

Lấy thông tin chi tiết của một user.

__Tham số__

* `userId`: number - ID của user

__Trả về__

Promise<UserInfo>
* `id`: number - Facebook ID
* `name`: string - Tên đầy đủ
* `firstName?`: string - Tên
* `username?`: string - Username
* `profilePictureUrl?`: string - URL ảnh đại diện
* `isMessengerUser?`: boolean - Có sử dụng Messenger không
* `isVerified?`: boolean - Tài khoản đã xác minh chưa
* `gender?`: number - Giới tính
* `canViewerMessage?`: boolean - Có thể nhắn tin không

__Ví dụ__

```typescript
const user = await client.getUserInfo(100000000000001)
console.log(`${user.name} (@${user.username})`)
```

---

<a name="searchUsers"></a>
## client.searchUsers(query)

Tìm kiếm users theo tên hoặc username.

__Tham số__

* `query`: string - Từ khóa tìm kiếm

__Trả về__

Promise<SearchUserResult[]>
* `id`: number - Facebook ID
* `name`: string - Tên
* `username`: string - Username

__Ví dụ__

```typescript
const users = await client.searchUsers('Nguyễn Văn A')
for (const user of users) {
    console.log(`${user.name} (${user.id})`)
}
```

---

# E2EE (Mã hóa đầu cuối)

<a name="connectE2EE"></a>
## client.connectE2EE()

Kết nối E2EE. Thường được gọi tự động nếu `enableE2EE: true`.

__Lưu ý__

Promise này resolve khi hàm hoàn thành, không phải khi E2EE đã kết nối xong. Hãy lắng nghe event `e2eeConnected` hoặc `fullyReady`.

__Ví dụ__

```typescript
await client.connectE2EE()
// Đợi event e2eeConnected
```

---

<a name="sendE2EEMessage"></a>
## client.sendE2EEMessage(chatJid, text, options?)

Gửi tin nhắn E2EE.

__Tham số__

* `chatJid`: string - Chat JID (format: `user_id@msgr.fb`)
* `text`: string - Nội dung tin nhắn
* `options?`: object
  * `replyToId?`: string - ID tin nhắn để reply
  * `replyToSenderJid?`: string - JID người gửi tin nhắn reply

__Trả về__

Promise<SendMessageResult>

__Ví dụ__

```typescript
await client.sendE2EEMessage('100000000000001@msgr.fb', 'Xin chào!')

// Reply
await client.sendE2EEMessage('100000000000001@msgr.fb', 'Đây là reply', {
    replyToId: 'msgid123',
    replyToSenderJid: '100000000000002@msgr.fb'
})
```

---

<a name="sendE2EEReaction"></a>
## client.sendE2EEReaction(chatJid, messageId, senderJid, emoji?)

Gửi/xóa reaction E2EE.

__Tham số__

* `chatJid`: string - Chat JID
* `messageId`: string - ID tin nhắn
* `senderJid`: string - JID người gửi tin nhắn gốc
* `emoji?`: string - Emoji (bỏ qua để xóa)

__Ví dụ__

```typescript
await client.sendE2EEReaction(chatJid, messageId, senderJid, '❤️')
```

---

<a name="sendE2EETyping"></a>
## client.sendE2EETyping(chatJid, isTyping?)

Gửi trạng thái đang nhập trong cuộc trò chuyện E2EE.

__Tham số__

* `chatJid`: string - Chat JID
* `isTyping?`: boolean - Đang nhập hay không (mặc định: true)

__Ví dụ__

```typescript
// Bắt đầu typing
await client.sendE2EETyping(chatJid, true)

// Dừng typing
await client.sendE2EETyping(chatJid, false)
```

---

<a name="editE2EEMessage"></a>
## client.editE2EEMessage(chatJid, messageId, newText)

Chỉnh sửa tin nhắn E2EE.

__Tham số__

* `chatJid`: string - Chat JID
* `messageId`: string - ID tin nhắn
* `newText`: string - Nội dung mới

__Ví dụ__

```typescript
await client.editE2EEMessage(chatJid, messageId, 'Nội dung đã sửa')
```

---

<a name="unsendE2EEMessage"></a>
## client.unsendE2EEMessage(chatJid, messageId)

Thu hồi tin nhắn E2EE.

__Tham số__

* `chatJid`: string - Chat JID
* `messageId`: string - ID tin nhắn

__Ví dụ__

```typescript
await client.unsendE2EEMessage(chatJid, messageId)
```

---

# E2EE Media

<a name="sendE2EEImage"></a>
## client.sendE2EEImage(chatJid, data, mimeType?, options?)

Gửi ảnh E2EE.

__Tham số__

* `chatJid`: string - Chat JID
* `data`: Buffer - Dữ liệu ảnh
* `mimeType?`: string - MIME type (mặc định: 'image/jpeg')
* `options?`: object
  * `caption?`: string - Caption
  * `width?`: number - Chiều rộng
  * `height?`: number - Chiều cao
  * `replyToId?`: string - ID tin nhắn reply
  * `replyToSenderJid?`: string - JID người gửi

__Ví dụ__

```typescript
const image = readFileSync('photo.jpg')
await client.sendE2EEImage(chatJid, image, 'image/jpeg', {
    caption: 'Ảnh đẹp!'
})
```

---

<a name="sendE2EEVideo"></a>
## client.sendE2EEVideo(chatJid, data, mimeType?, options?)

Gửi video E2EE.

__Tham số__

* `chatJid`: string - Chat JID
* `data`: Buffer - Dữ liệu video
* `mimeType?`: string - MIME type (mặc định: 'video/mp4')
* `options?`: object
  * `caption?`: string - Caption
  * `width?`: number - Chiều rộng
  * `height?`: number - Chiều cao
  * `duration?`: number - Thời lượng (giây)
  * `replyToId?`: string - ID tin nhắn reply
  * `replyToSenderJid?`: string - JID người gửi

__Ví dụ__

```typescript
const video = readFileSync('video.mp4')
await client.sendE2EEVideo(chatJid, video, 'video/mp4', {
    caption: 'Video hay!',
    duration: 30
})
```

---

<a name="sendE2EEAudio"></a>
## client.sendE2EEAudio(chatJid, data, mimeType?, options?)

Gửi audio/voice E2EE.

__Tham số__

* `chatJid`: string - Chat JID
* `data`: Buffer - Dữ liệu audio
* `mimeType?`: string - MIME type (mặc định: 'audio/ogg')
* `options?`: object
  * `ptt?`: boolean - Push-to-talk/voice message (mặc định: false)
  * `duration?`: number - Thời lượng (giây)
  * `replyToId?`: string - ID tin nhắn reply
  * `replyToSenderJid?`: string - JID người gửi

__Ví dụ__

```typescript
const voice = readFileSync('voice.ogg')
await client.sendE2EEAudio(chatJid, voice, 'audio/ogg', {
    ptt: true,
    duration: 10
})
```

---

<a name="sendE2EEDocument"></a>
## client.sendE2EEDocument(chatJid, data, filename, mimeType, options?)

Gửi file/document E2EE.

__Tham số__

* `chatJid`: string - Chat JID
* `data`: Buffer - Dữ liệu file
* `filename`: string - Tên file
* `mimeType`: string - MIME type
* `options?`: object
  * `replyToId?`: string - ID tin nhắn reply
  * `replyToSenderJid?`: string - JID người gửi

__Ví dụ__

```typescript
const pdf = readFileSync('document.pdf')
await client.sendE2EEDocument(chatJid, pdf, 'document.pdf', 'application/pdf')
```

---

<a name="sendE2EESticker"></a>
## client.sendE2EESticker(chatJid, data, mimeType?, options?)

Gửi sticker E2EE.

__Tham số__

* `chatJid`: string - Chat JID
* `data`: Buffer - Dữ liệu sticker (WebP format)
* `mimeType?`: string - MIME type (mặc định: 'image/webp')
* `options?`: object
  * `replyToId?`: string - ID tin nhắn reply
  * `replyToSenderJid?`: string - JID người gửi

__Ví dụ__

```typescript
const sticker = readFileSync('sticker.webp')
await client.sendE2EESticker(chatJid, sticker, 'image/webp')
```

---

<a name="downloadE2EEMedia"></a>
## client.downloadE2EEMedia(options)

Tải xuống và giải mã media E2EE từ attachment.

__Tham số__

* `options`: object
  * `directPath`: string - Đường dẫn trực tiếp từ attachment
  * `mediaKey`: string - Media key mã hóa Base64
  * `mediaSha256`: string - SHA256 của file gốc mã hóa Base64
  * `mediaEncSha256?`: string - SHA256 của file đã mã hóa, mã hóa Base64 (khuyến nghị để xác minh)
  * `mediaType`: string - Loại media: `'image'`, `'video'`, `'audio'`, `'document'`, `'sticker'`
  * `mimeType`: string - MIME type (ví dụ: 'image/jpeg')
  * `fileSize`: number - Kích thước file (bytes)

__Trả về__

Promise<{ data: Buffer; mimeType: string; fileSize: number }>
* `data`: Buffer - Dữ liệu media đã giải mã
* `mimeType`: string - MIME type
* `fileSize`: number - Kích thước file

__Ví dụ__

```typescript
import { writeFileSync } from 'fs'

client.on('e2eeMessage', async (message) => {
    if (message.attachments && message.attachments.length > 0) {
        const attachment = message.attachments[0]
        
        // Kiểm tra attachment có metadata E2EE cần thiết không
        if (attachment.mediaKey && attachment.mediaSha256 && attachment.directPath) {
            try {
                const result = await client.downloadE2EEMedia({
                    directPath: attachment.directPath,
                    mediaKey: attachment.mediaKey,
                    mediaSha256: attachment.mediaSha256,
                    mediaEncSha256: attachment.mediaEncSha256, // Tùy chọn nhưng khuyến nghị
                    mediaType: attachment.type,
                    mimeType: attachment.mimeType || 'application/octet-stream',
                    fileSize: attachment.fileSize || 0,
                })
                
                // Lưu vào file
                const extension = result.mimeType.split('/')[1] || 'bin'
                writeFileSync(`downloaded.${extension}`, result.data)
                console.log(`Đã tải ${result.fileSize} bytes`)
            } catch (error) {
                console.error('Không thể tải media E2EE:', error)
            }
        }
    }
})
```

__Lưu ý__

Method này chỉ hoạt động với tin nhắn E2EE (mã hóa đầu cuối). Với tin nhắn thường, hãy sử dụng trường `url` trong attachment thay thế.

---

<a name="getDeviceData"></a>
## client.getDeviceData()

Lấy E2EE device data để lưu trữ.

__Trả về__

string - Device data dưới dạng JSON string

__Lưu ý__

Lưu device data để tránh phải setup E2EE lại mỗi lần khởi động.

__Ví dụ__

```typescript
import { writeFileSync } from 'fs'

// Lưu device data
const deviceData = client.getDeviceData()
writeFileSync('device.json', deviceData)

// Load khi khởi động
const client = new Client(cookies, {
    deviceData: readFileSync('device.json', 'utf-8')
})
```

---

# Khác

<a name="unloadLibrary"></a>
## client.unloadLibrary()

Giải phóng native library khỏi bộ nhớ.

__Cảnh báo__

Sau khi gọi method này, mọi thao tác với client sẽ gây crash. Chỉ sử dụng khi cần cleanup hoàn toàn trước khi tắt ứng dụng.

__Ví dụ__

```typescript
await client.disconnect()
client.unloadLibrary()
// Không sử dụng client sau đây!
```

---

# Utilities

<a name="parseCookies"></a>
## Utils.parseCookies(input)

Parse cookies từ nhiều định dạng khác nhau.

__Tham số__

* `input`: string - Cookies dưới dạng:
  * JSON array: `[{ "name": "c_user", "value": "..." }, ...]`
  * JSON object: `{ "c_user": "...", "xs": "..." }`
  * Cookie string: `"c_user=...; xs=..."`
  * Netscape format
  * Base64 encoded (bất kỳ format trên)

__Trả về__

Cookies - Object với key-value

__Ví dụ__

```typescript
import { Utils } from 'meta-messenger.js'
import { readFileSync } from 'fs'

const cookies = Utils.parseCookies(readFileSync('cookies.json', 'utf-8'))
```

---

<a name="validate"></a>
## Utils.validate(cookies)

Kiểm tra cookies có đầy đủ các trường bắt buộc không.

__Tham số__

* `cookies`: Cookies - Object cookies

__Trả về__

boolean - `true` nếu hợp lệ

__Ví dụ__

```typescript
if (!Utils.validate(cookies)) {
    console.error('Cookies không hợp lệ!')
}
```

---

<a name="getMissing"></a>
## Utils.getMissing(cookies)

Lấy danh sách các cookies bắt buộc đang thiếu.

__Tham số__

* `cookies`: Cookies - Object cookies

__Trả về__

string[] - Danh sách tên cookies đang thiếu

__Ví dụ__

```typescript
const missing = Utils.getMissing(cookies)
if (missing.length > 0) {
    console.error(`Thiếu cookies: ${missing.join(', ')}`)
}
```

---

# Events

> **Chú thích:**
> - 🔵 **Thường** = Tin nhắn thường (không mã hóa)
> - 🟢 **E2EE** = Tin nhắn mã hóa đầu cuối

| Event | Thường | E2EE | Mô tả |
|-------|:------:|:----:|-------|
| `ready` | 🔵 | ❌ | Kết nối socket thành công |
| `reconnected` | 🔵 | ❌ | Tái kết nối thành công |
| `message` | 🔵 | ❌ | Tin nhắn thường mới |
| `e2eeMessage` | ❌ | 🟢 | Tin nhắn E2EE mới |
| `messageEdit` | 🔵 | 🟢 | Tin nhắn được chỉnh sửa |
| `messageUnsend` | 🔵 | 🟢 | Tin nhắn bị thu hồi |
| `reaction` | 🔵 | ❌ | Reaction tin nhắn thường |
| `e2eeReaction` | ❌ | 🟢 | Reaction tin nhắn E2EE |
| `typing` | 🔵 | ❌ | Đang nhập (thường) |
| `readReceipt` | 🔵 | ❌ | Tin nhắn đã đọc (thường) |
| `e2eeReceipt` | ❌ | 🟢 | Tin nhắn đã đọc (E2EE) |
| `e2eeConnected` | ❌ | 🟢 | Kết nối E2EE thành công |
| `deviceDataChanged` | ❌ | 🟢 | Device data thay đổi |
| `fullyReady` | 🔵 | 🟢 | Client hoàn toàn sẵn sàng |
| `disconnected` | 🔵 | 🟢 | Mất kết nối |
| `error` | 🔵 | 🟢 | Có lỗi xảy ra |

---

<a name="event-ready"></a>
## Event: 'ready'

> 🔵 **Kết nối socket**

Phát ra khi kết nối socket thành công (trước E2EE).

```typescript
client.on('ready', (data) => {
    console.log('Socket connected!')
    if (data.isNewSession) {
        console.log('Đây là session mới')
    }
})
```

__Data object__

* `isNewSession`: boolean - `true` nếu là phiên kết nối mới

---

<a name="event-reconnected"></a>
## Event: 'reconnected'

> 🔵 **Tái kết nối socket**

Phát ra khi tái kết nối socket thành công sau khi mất kết nối.

```typescript
client.on('reconnected', () => {
    console.log('Reconnected to Messenger!')
})
```

---

<a name="event-message"></a>
## Event: 'message'

> 🔵 **Chỉ tin nhắn thường**

Phát ra khi có tin nhắn thường mới.

```typescript
client.on('message', (message: Message) => {
    console.log(`${message.senderId}: ${message.text}`)
})
```

__Message object__

* `id`: string - Message ID
* `threadId`: number - Thread ID
* `senderId`: number - Sender ID
* `text`: string - Nội dung
* `timestampMs`: number - Timestamp
* `attachments?`: Attachment[] - Attachments
* `replyTo?`: ReplyTo - Reply info
* `mentions?`: Mention[] - Mentions
* `isAdminMsg?`: boolean - Tin nhắn hệ thống

---

<a name="event-messageEdit"></a>
## Event: 'messageEdit'

> 🔵🟢 **Hỗ trợ cả tin nhắn thường và E2EE**

Phát ra khi tin nhắn được chỉnh sửa (cả thường và E2EE).

```typescript
client.on('messageEdit', (data) => {
    console.log(`Message ${data.messageId} edited to: ${data.newText}`)
})
```

__Data object__

* `messageId`: string - ID tin nhắn
* `newText`: string - Nội dung mới
* `editCount?`: number - Số lần chỉnh sửa
* `timestampMs`: number - Thời gian chỉnh sửa

---

<a name="event-messageUnsend"></a>
## Event: 'messageUnsend'

> 🔵🟢 **Hỗ trợ cả tin nhắn thường và E2EE**

Phát ra khi tin nhắn bị thu hồi (cả thường và E2EE).

```typescript
client.on('messageUnsend', (data) => {
    console.log(`Message ${data.messageId} unsent in thread ${data.threadId}`)
})
```

__Data object__

* `messageId`: string - ID tin nhắn
* `threadId`: number - Thread ID

---

<a name="event-reaction"></a>
## Event: 'reaction'

> 🔵 **Chỉ tin nhắn thường** - Xem [`e2eeReaction`](#event-e2eeReaction) cho E2EE

Phát ra khi có reaction mới trên tin nhắn thường.

```typescript
client.on('reaction', (data) => {
    console.log(`${data.actorId} reacted ${data.reaction} to ${data.messageId}`)
})
```

__Data object__

* `messageId`: string - ID tin nhắn
* `threadId`: number - Thread ID
* `actorId`: number - ID người reaction
* `reaction`: string - Emoji (rỗng = bỏ reaction)

---

<a name="event-typing"></a>
## Event: 'typing'

> 🔵 **Chỉ tin nhắn thường**

Phát ra khi ai đó đang nhập trong thread thường.

```typescript
client.on('typing', (data) => {
    console.log(`${data.senderId} is ${data.isTyping ? 'typing' : 'stopped typing'}`)
})
```

__Data object__

* `threadId`: number - Thread ID
* `senderId`: number - ID người nhập
* `isTyping`: boolean - Đang nhập hay dừng

---

<a name="event-readReceipt"></a>
## Event: 'readReceipt'

> 🔵 **Chỉ tin nhắn thường** - Xem [`e2eeReceipt`](#event-e2eeReceipt) cho E2EE

Phát ra khi tin nhắn thường được đọc.

```typescript
client.on('readReceipt', (data) => {
    console.log(`${data.readerId} read messages in ${data.threadId}`)
})
```

__Data object__

* `threadId`: number - Thread ID
* `readerId`: number - ID người đọc
* `readWatermarkTimestampMs`: number - Timestamp watermark đã đọc
* `timestampMs?`: number - Thời gian đọc

---

<a name="event-e2eeMessage"></a>
## Event: 'e2eeMessage'

> 🟢 **Chỉ tin nhắn E2EE** - Xem [`message`](#event-message) cho tin nhắn thường

Phát ra khi có tin nhắn E2EE mới.

```typescript
client.on('e2eeMessage', (message: E2EEMessage) => {
    console.log(`[E2EE] ${message.senderJid}: ${message.text}`)
})
```

__E2EEMessage object__

* `id`: string - Message ID
* `threadId`: number - Thread ID
* `chatJid`: string - Chat JID
* `senderJid`: string - Sender JID
* `senderId`: number - Sender ID
* `text`: string - Nội dung
* `timestampMs`: number - Timestamp
* `attachments?`: Attachment[]
* `replyTo?`: ReplyTo
* `mentions?`: Mention[]

---

<a name="event-e2eeReaction"></a>
## Event: 'e2eeReaction'

> 🟢 **Chỉ tin nhắn E2EE** - Xem [`reaction`](#event-reaction) cho tin nhắn thường

Phát ra khi có reaction trên tin nhắn E2EE.

```typescript
client.on('e2eeReaction', (data) => {
    console.log(`${data.senderJid} reacted ${data.reaction}`)
})
```

__Data object__

* `messageId`: string - ID tin nhắn
* `chatJid`: string - Chat JID
* `senderJid`: string - JID người reaction
* `senderId`: number - ID người reaction
* `reaction`: string - Emoji (rỗng = bỏ reaction)

---

<a name="event-e2eeReceipt"></a>
## Event: 'e2eeReceipt'

> 🟢 **Chỉ tin nhắn E2EE** - Xem [`readReceipt`](#event-readReceipt) cho tin nhắn thường

Phát ra khi có receipt cho tin nhắn E2EE (đã đọc, đã nhận, v.v.).

```typescript
client.on('e2eeReceipt', (data) => {
    console.log(`[E2EE] Receipt type ${data.type} for messages:`, data.messageIds)
})
```

__Data object__

* `type`: string - Loại receipt (`'read'`, `'delivered'`, v.v.)
* `chat`: string - Chat JID
* `sender`: string - Sender JID
* `messageIds`: string[] - Danh sách message IDs

---

<a name="event-e2eeConnected"></a>
## Event: 'e2eeConnected'

> 🟢 **Chỉ E2EE**

Phát ra khi kết nối E2EE thành công.

```typescript
client.on('e2eeConnected', () => {
    console.log('E2EE connected!')
})
```

---

<a name="event-fullyReady"></a>
## Event: 'fullyReady'

> 🔵🟢 **Hỗ trợ cả thường và E2EE**

Phát ra khi client hoàn toàn sẵn sàng (socket + E2EE nếu enabled).

```typescript
client.on('fullyReady', () => {
    console.log('Client đã sẵn sàng!')
})
```

__Lưu ý__

Các event message/e2eeMessage sẽ được queue cho đến khi `fullyReady` được phát ra.

---

<a name="event-disconnected"></a>
## Event: 'disconnected'

> 🔵🟢 **Hỗ trợ cả thường và E2EE**

Phát ra khi mất kết nối.

```typescript
client.on('disconnected', (data) => {
    if (data?.isE2EE) {
        console.log('E2EE disconnected')
    } else {
        console.log('Socket disconnected')
    }
})
```

__Data object__

* `isE2EE?`: boolean - `true` nếu mất kết nối E2EE

---

<a name="event-error"></a>
## Event: 'error'

> 🔵🟢 **Hỗ trợ cả thường và E2EE**

Phát ra khi có lỗi xảy ra.

```typescript
client.on('error', (error) => {
    console.error(`Error: ${error.message}`)
})
```

__Tham số__

* `error`: Error - Đối tượng Error tiêu chuẩn của JavaScript

---

<a name="event-deviceDataChanged"></a>
## Event: 'deviceDataChanged'

> 🟢 **Chỉ E2EE** - Chỉ khi dùng option `deviceData`

Phát ra khi device data E2EE thay đổi. Sử dụng để lưu device data vào database.

```typescript
client.on('deviceDataChanged', (data) => {
    // Lưu device data vào database
    await saveToDatabase(data.deviceData)
})
```

__Data object__

* `deviceData`: string - Device data dưới dạng JSON string

__Lưu ý__

Event này chỉ được phát ra khi bạn khởi tạo client với option `deviceData`. Nếu dùng `e2eeDeviceDataPath`, device data sẽ tự động lưu vào file.

---

# Types

## Cookies

```typescript
interface Cookies {
    c_user: string
    xs: string
    datr: string
    [key: string]: string
}
```

## Message

```typescript
interface Message {
    id: string
    threadId: number
    senderId: number
    text: string
    timestampMs: number
    isE2EE?: boolean
    chatJid?: string
    senderJid?: string
    attachments?: Attachment[]
    replyTo?: ReplyTo
    mentions?: Mention[]
    isAdminMsg?: boolean
}
```

## Attachment

```typescript
interface Attachment {
    type: 'image' | 'video' | 'audio' | 'file' | 'sticker' | 'gif' | 'voice' | 'link'
    url?: string
    fileName?: string
    mimeType?: string
    fileSize?: number
    width?: number
    height?: number
    duration?: number
    stickerId?: number
    previewUrl?: string
    // Dành cho tải media E2EE (chỉ có trong tin nhắn E2EE)
    mediaKey?: string      // Khóa mã hóa dạng Base64
    mediaSha256?: string   // SHA256 file gốc dạng Base64
    mediaEncSha256?: string // SHA256 file đã mã hóa dạng Base64
    directPath?: string    // Đường dẫn trực tiếp để tải
}
```

## ReplyTo

```typescript
interface ReplyTo {
    messageId: string
    senderId?: number
    text?: string
}
```

## Mention

```typescript
interface Mention {
    userId: number
    offset: number
    length: number
}
```

## Thread

```typescript
interface Thread {
    id: number
    type: number
    name: string
    lastActivityTimestampMs: number
    isGroup?: boolean
    participants?: number[]
}
```

## User

```typescript
interface User {
    id: number
    name: string
    username: string
}
```

## UserInfo

```typescript
interface UserInfo {
    id: number
    name: string
    firstName?: string
    username?: string
    profilePictureUrl?: string
    isMessengerUser?: boolean
    isVerified?: boolean
    gender?: number
    canViewerMessage?: boolean
}
```
