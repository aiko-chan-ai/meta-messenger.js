/* eslint-disable */

// Node.js 22.12.0+ is required for native ESM support
const { login, Utils } = require("meta-messenger.js");

login(
    // C3C UFC Utility / EditThisCookie
    Utils.parseCookies([
        {
            domain: ".facebook.com",
            expirationDate: 1000000000,
            hostOnly: false,
            httpOnly: false,
            name: "c_user",
            path: "/",
            sameSite: "no_restriction",
            secure: true,
            session: false,
            storeId: "1",
            value: "<uid>",
            id: 1,
        },
        {
            domain: ".facebook.com",
            expirationDate: 1000000000,
            hostOnly: false,
            httpOnly: true,
            name: "datr",
            path: "/",
            sameSite: "no_restriction",
            secure: true,
            session: false,
            storeId: "1",
            value: "<secret>",
            id: 2,
        },
        {
            domain: ".facebook.com",
            expirationDate: 1000000000,
            hostOnly: false,
            httpOnly: true,
            name: "fr",
            path: "/",
            sameSite: "no_restriction",
            secure: true,
            session: false,
            storeId: "1",
            value: "<secret>",
            id: 3,
        },
        {
            domain: ".facebook.com",
            expirationDate: 1000000000,
            hostOnly: false,
            httpOnly: true,
            name: "xs",
            path: "/",
            sameSite: "no_restriction",
            secure: true,
            session: false,
            storeId: "1",
            value: "<secret>",
            id: 4,
        },
    ]),
).then(api => {
    api.on("fullyReady", () => {
        console.log(`Logged in as ${api.user.name} (${api.currentUserId})`);
    });
    api.on("message", message => {
        if (message.senderId === api.currentUserId) return;
        api.sendMessage(message.threadId, `You said: ${message.text}`);
    });
    api.on("e2eeMessage", message => {
        if (message.senderId === api.currentUserId) return;
        api.sendE2EEMessage(message.chatJid, `You said (e2ee): ${message.text}`);
    });
});
