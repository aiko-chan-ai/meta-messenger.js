/*
 * meta-messenger.js, Unofficial Meta Messenger Chat API for NodeJS
 * Copyright (c) 2026 Elysia and contributors
 */

/**
 * meta-messenger.js - TypeScript wrapper for Facebook Messenger with E2EE support
 *
 * @example
 * ```typescript
 * import { Client, login } from 'meta-messenger.js'
 *
 * // Method 1: Using Client class directly
 * const client = new Client({
 *     c_user: 'your_user_id',
 *     xs: 'your_xs_cookie',
 *     datr: 'your_datr_cookie',
 *     fr: 'your_fr_cookie'
 * })
 *
 * client.on('message', (message) => {
 *     console.log('New message:', message)
 *     if (message.text === 'ping') {
 *         client.sendMessage(message.threadId, 'pong')
 *     }
 * })
 *
 * await client.connect()
 *
 * // Method 2: Using login function (facebook-chat-api style)
 * const api = await login({
 *     c_user: 'your_user_id',
 *     xs: 'your_xs_cookie',
 *     datr: 'your_datr_cookie',
 *     fr: 'your_fr_cookie'
 * })
 *
 * api.on('message', (message) => {
 *     console.log('Got message:', message.text)
 * })
 * ```
 *
 * @packageDocumentation
 */

export type { ClientEventMap } from "./client.js";
export { Client } from "./client.js";
export {
    type BaseEvent,
    type ClientEvent,
    type ClientOptions,
    type Cookies,
    type CreateThreadResult,
    type DisconnectedEvent,
    type E2EEConnectedEvent,
    type E2EEMessage,
    type E2EEMessageEvent,
    type E2EEReactionEvent,
    type E2EEReceiptEvent,
    type ErrorEvent,
    // Types
    type EventType,
    type InitialData,
    type LogLevel,
    type Message,
    type MessageEditEvent,
    type MessageEvent,
    type MessageUnsendEvent,
    type Platform,
    type ReactionEvent,
    type ReadyEvent,
    type ReconnectedEvent,
    type SearchUserResult,
    type SendMessageOptions,
    type SendMessageResult,
    type Thread,
    // Enums
    ThreadType,
    type TypingEvent,
    type UploadMediaResult,
    type User,
    type UserInfo,
} from "./types.js";
export { type CookieObject, Utils } from "./utils.js";

import { Client } from "./client.js";
import type { ClientOptions, Cookies } from "./types.js";

/**
 * Login to Facebook Messenger (E2EE disabled for simplicity)
 *
 * @param cookies - Authentication cookies
 * @param options - Client options
 * @returns Connected client instance
 *
 * @example
 * ```typescript
 * const api = await login({
 *     c_user: 'your_user_id',
 *     xs: 'your_xs_cookie',
 *     datr: 'your_datr_cookie',
 *     fr: 'your_fr_cookie'
 * })
 *
 * console.log('Logged in as:', api.user?.name)
 *
 * api.on('message', async (message) => {
 *     if (message.senderId !== api.currentUserId) {
 *         await api.sendMessage(message.threadId, 'Hello!')
 *     }
 * })
 * ```
 */
export async function login(cookies: Cookies, options?: ClientOptions): Promise<Client> {
    const client = new Client(cookies, options);
    await client.connect();
    return client;
}

/**
 * Create a client without connecting
 *
 * @param cookies - Authentication cookies
 * @param options - Client options
 * @returns Client instance (not connected)
 */
export function createClient(cookies: Cookies, options?: ClientOptions): Client {
    return new Client(cookies, options);
}

// Default export for convenience
export default { Client, login, createClient };
