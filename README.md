# Unofficial Facebook Messenger API

<a href="https://www.npmjs.com/package/meta-messenger.js"><img alt="npm version" src="https://img.shields.io/npm/v/meta-messenger.js.svg?style=flat-square"></a>
<a href="https://www.npmjs.com/package/meta-messenger.js"><img src="https://img.shields.io/npm/dm/meta-messenger.js.svg?style=flat-square" alt="npm downloads"></a>

Facebook now has an official API for chat bots [here](https://developers.facebook.com/docs/messenger-platform).

This API is one of the ways to automate chat functions on user accounts.
Using the [mautrix-meta](https://github.com/mautrix/meta) library (Go) through FFI binding, this library supports (almost) all features that the original library has implemented through a high-level JavaScript API.

> [!CAUTION]
> **Using this feature on user accounts is prohibited under Meta's Terms of Service and may result in account suspension.**
> **I am not responsible for any Meta accounts that are suspended due to using this library.**

## Installation
```bash
npm install meta-messenger.js
```

### Requirements
- Node.js >= 22
- Cookies from a logged-in Facebook/Messenger account

## Usage Example
```typescript
import { Client, Utils } from 'meta-messenger.js'
import { readFileSync } from 'fs'

// Read cookies from file (supports multiple formats)
const cookies = Utils.parseCookies(readFileSync('cookies.json', 'utf-8'))

// Create a new client
const client = new Client(cookies)

// Listen to events
client.on('message', async (message) => {
    console.log(`Message from ${message.senderId}: ${message.text}`)
    
    // Simple echo bot
    if (message.text === 'ping') {
        await client.sendMessage(message.threadId, 'pong')
    }
})

// Connect
await client.connect()
console.log(`Logged in as: ${client.user?.name}`)
```

![img](screenshots/demo1.png)

## [Documentation](DOCS.md)

## Main Features

### Send Messages
```typescript
// Simple message
await client.sendMessage(threadId, 'Hello!')

// Message with reply
await client.sendMessage(threadId, {
    text: 'This is a reply',
    replyToId: 'mid.xxx'
})

// Message with mention
await client.sendMessage(threadId, {
    text: 'Hello @friend!',
    mentions: [{ userId: 123456, offset: 6, length: 7 }]
})
```

### Send Media
```typescript
import { readFileSync } from 'fs'

// Send image
const image = readFileSync('photo.jpg')
await client.sendImage(threadId, image, 'photo.jpg', 'Caption here')

// Send video
const video = readFileSync('video.mp4')
await client.sendVideo(threadId, video, 'video.mp4')

// Send sticker
await client.sendSticker(threadId, 369239263222822) // thumbs up
```

### Save Session

To avoid having to get new cookies each time, you should save the device data (for E2EE) to a file:

```typescript
import { writeFileSync, readFileSync } from 'fs'

// Save device data when it changes
client.on('deviceDataChanged', ({ deviceData }) => {
    writeFileSync('device.json', deviceData)
})

// Load device data on startup
const deviceData = readFileSync('device.json', 'utf-8')
const client = new Client(cookies, { deviceData })
```

### Listen to Messages

```typescript
client.on('message', (message) => {
    console.log(`[${message.threadId}] ${message.senderId}: ${message.text}`)
    
    // Handle attachments
    for (const att of message.attachments || []) {
        if (att.type === 'image') {
            console.log(`  Image: ${att.url}`)
        } else if (att.type === 'sticker') {
            console.log(`  Sticker: ${att.stickerId}`)
        }
    }
})

// E2EE messages
client.on('e2eeMessage', (message) => {
    console.log(`[E2EE] ${message.senderJid}: ${message.text}`)
})
```

## Supported Cookie Formats

The API supports multiple cookie formats through `Utils.parseCookies()`:

1. **C3C UFC Utility / EditThisCookie**
2. **Simple object**
3. **Cookie header string**
4. **Netscape cookie file format**
5. **Base64 encoded** (any of the above formats)

Required cookies: `c_user`, `xs`, `datr`, `fr`

## FAQs

1. **How do I get cookies?**
> Log in to Facebook in your browser, open DevTools (F12), go to Application > Cookies tab, and copy the values for `c_user`, `xs`, `datr`, `fr`.

2. **Why am I not receiving E2EE messages?**
> Make sure you have enabled E2EE with `enableE2EE: true` and wait for the `fullyReady` event before processing messages.

3. **What's the difference between message and e2eeMessage events?**
> - `message`: Regular messages, metadata visible to Facebook
> - `e2eeMessage`: End-to-end encrypted messages, only sender and receiver can read

4. **Why do I need to wait for fullyReady?**
> The `fullyReady` event is emitted when both socket and E2EE (if enabled) are ready. If you process messages before this, you may miss some E2EE messages.

## License

GNU Affero General Public License v3.0

According to mautrix-meta License

## Credits

- Claude Opus 4.5 - Supports understanding and developing this library
- [mautrix-meta](https://github.com/mautrix/meta) - A Matrix-Facebook Messenger and Instagram DM puppeting bridge.
- [whatsmeow](https://github.com/tulir/whatsmeow) - Go library for the WhatsApp web multidevice API
- [facebook-chat-api](https://github.com/Schmavery/facebook-chat-api) - Inspired the API Style
- [whatsmeow-node](https://github.com/vinikjkkj/whatsmeow-node) - Reference for writing NodeJS library with FFI.
- [koffi](https://codeberg.org/Koromix/rygel/src/branch/master/src/koffi) - Fast and easy-to-use C FFI module for Node.js

## Disclaimer

> All product and company names are trademarks™ or registered® trademarks of their respective holders. Use of them does not imply any affiliation with or endorsement by them.  
> "Facebook", "Messenger", "Meta", "Instagram" is a registered trademark of Meta, Inc., used under license agreement.