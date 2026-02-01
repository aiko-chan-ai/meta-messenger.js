/*
 * meta-messenger.js, Unofficial Meta Messenger Chat API for NodeJS
 * Copyright (c) 2026 Elysia and contributors
 */

/**
 * Utility class for cookie parsing and conversion
 *
 * Supports multiple cookie formats:
 * - C3C UFC Utility / EditThisCookie (array of cookie objects)
 * - Simple object format { name: value }
 * - Netscape/HTTP cookie file format
 * - Cookie header string format
 * - Base64 encoded cookies
 */

import type { Cookies } from "./types.js";

/**
 * Cookie object format (C3C UFC Utility / EditThisCookie style)
 */
export interface CookieObject {
    name: string;
    value: string;
    domain?: string;
    path?: string;
    expires?: number | string;
    expirationDate?: number;
    httpOnly?: boolean;
    secure?: boolean;
    sameSite?: string;
    hostOnly?: boolean;
    session?: boolean;
}

/**
 * Netscape cookie file line format
 */
export interface NetscapeCookie {
    domain: string;
    flag: boolean;
    path: string;
    secure: boolean;
    expiration: number;
    name: string;
    value: string;
}

/**
 * Static utility class for cookie operations
 * @example
 * ```typescript
 * import { Utils } from 'meta-messenger.js'
 *
 * // Parse any cookie format
 * const cookies = Utils.parseCookies(rawData)
 *
 * // Convert cookies to header string
 * const header = Utils.toCookieString(cookies)
 * ```
 */
export class Utils extends null {
    /**
     * Parse cookies from various formats
     * Automatically detects the format and parses accordingly
     *
     * @param input - Cookie data in any supported format
     * @returns Parsed cookies object
     */
    static parseCookies(input: string | CookieObject[] | Record<string, string>): Cookies {
        // Already a simple object
        if (typeof input === "object" && !Array.isArray(input)) {
            return input as Cookies;
        }

        // Array of cookie objects (C3C UFC Utility / EditThisCookie)
        if (Array.isArray(input)) {
            return Utils.fromCookieArray(input);
        }

        // String input - detect format
        if (typeof input === "string") {
            const trimmed = input.trim();

            // Try Base64 first
            if (Utils.isBase64(trimmed)) {
                try {
                    const decoded = Buffer.from(trimmed, "base64").toString("utf-8");
                    return Utils.parseCookies(decoded);
                } catch {
                    // Not valid base64, continue with other formats
                }
            }

            // Try JSON (array or object)
            if (trimmed.startsWith("[") || trimmed.startsWith("{")) {
                try {
                    const parsed = JSON.parse(trimmed);
                    return Utils.parseCookies(parsed);
                } catch {
                    // Not valid JSON, continue with other formats
                }
            }

            // Netscape cookie file format (starts with # or domain)
            if (
                trimmed.includes("\t") &&
                (trimmed.startsWith("#") || trimmed.includes(".facebook.com") || trimmed.includes(".messenger.com"))
            ) {
                return Utils.fromNetscape(trimmed);
            }

            // Cookie header string format (name=value; name2=value2)
            if (trimmed.includes("=")) {
                return Utils.fromCookieString(trimmed);
            }
        }

        throw new Error("Unable to parse cookies: unknown format");
    }

    /**
     * Parse cookies from C3C UFC Utility / EditThisCookie array format
     *
     * @param cookies - Array of cookie objects
     * @returns Parsed cookies object
     *
     * @example
     * ```typescript
     * const cookies = Utils.fromCookieArray([
     *     { name: 'c_user', value: '123456' },
     *     { name: 'xs', value: 'abc...' }
     * ])
     * ```
     */
    static fromCookieArray(cookies: CookieObject[]): Cookies {
        const result: Record<string, string> = {};
        for (const cookie of cookies) {
            if (cookie.name && cookie.value !== undefined) {
                result[cookie.name] = String(cookie.value);
            }
        }
        return result as Cookies;
    }

    /**
     * Parse cookies from cookie header string format
     *
     * @param cookieString - Cookie string (e.g., "name1=value1; name2=value2")
     * @returns Parsed cookies object
     *
     * @example
     * ```typescript
     * const cookies = Utils.fromCookieString('c_user=123456; xs=abc...; datr=xyz...')
     * ```
     */
    static fromCookieString(cookieString: string): Cookies {
        const result: Record<string, string> = {};
        const pairs = cookieString.split(/;\s*/);

        for (const pair of pairs) {
            const [name, ...valueParts] = pair.split("=");
            if (name && valueParts.length > 0) {
                const trimmedName = name.trim();
                const value = valueParts.join("=").trim();
                if (trimmedName && value) {
                    result[trimmedName] = value;
                }
            }
        }

        return result as Cookies;
    }

    /**
     * Parse cookies from Netscape/HTTP cookie file format
     *
     * @param content - Netscape cookie file content
     * @returns Parsed cookies object
     *
     * @example
     * ```typescript
     * const cookies = Utils.fromNetscape(`
     * # Netscape HTTP Cookie File
     * .facebook.com	TRUE	/	TRUE	1234567890	c_user	123456
     * .facebook.com	TRUE	/	TRUE	1234567890	xs	abc...
     * `)
     * ```
     */
    static fromNetscape(content: string): Cookies {
        const result: Record<string, string> = {};
        const lines = content.split("\n");

        for (const line of lines) {
            const trimmed = line.trim();

            // Skip comments and empty lines
            if (!trimmed || trimmed.startsWith("#")) {
                continue;
            }

            // Parse tab-separated values
            const parts = trimmed.split("\t");
            if (parts.length >= 7) {
                const name = parts[5];
                const value = parts[6];
                if (name && value) {
                    result[name] = value;
                }
            }
        }

        return result as Cookies;
    }

    /**
     * Parse cookies from Base64 encoded string
     *
     * @param base64 - Base64 encoded cookie data
     * @returns Parsed cookies object
     */
    static fromBase64(base64: string): Cookies {
        const decoded = Buffer.from(base64, "base64").toString("utf-8");
        return Utils.parseCookies(decoded);
    }

    /**
     * Convert cookies object to cookie header string
     *
     * @param cookies - Cookies object
     * @returns Cookie header string
     *
     * @example
     * ```typescript
     * const header = Utils.toCookieString({ c_user: '123456', xs: 'abc...' })
     * // Returns: "c_user=123456; xs=abc..."
     * ```
     */
    static toCookieString(cookies: Cookies): string {
        return Object.entries(cookies)
            .map(([name, value]) => `${name}=${value}`)
            .join("; ");
    }

    /**
     * Convert cookies object to array format (C3C UFC Utility style)
     *
     * @param cookies - Cookies object
     * @param domain - Cookie domain (default: .facebook.com)
     * @returns Array of cookie objects
     *
     * @example
     * ```typescript
     * const arr = Utils.toCookieArray({ c_user: '123456', xs: 'abc...' })
     * ```
     */
    static toCookieArray(cookies: Cookies, domain = ".facebook.com"): CookieObject[] {
        return Object.entries(cookies)
            .filter(([, value]) => value !== undefined)
            .map(([name, value]) => ({
                name,
                value: value as string,
                domain,
                path: "/",
                secure: true,
                httpOnly: true,
            }));
    }

    /**
     * Convert cookies object to Netscape format
     *
     * @param cookies - Cookies object
     * @param domain - Cookie domain (default: .facebook.com)
     * @returns Netscape cookie file content
     */
    static toNetscape(cookies: Cookies, domain = ".facebook.com"): string {
        const lines = ["# Netscape HTTP Cookie File", "# Generated by meta-messenger.js", ""];

        for (const [name, value] of Object.entries(cookies)) {
            // domain, flag, path, secure, expiration, name, value
            const expiration = Math.floor(Date.now() / 1000) + 365 * 24 * 60 * 60; // 1 year
            lines.push(`${domain}\tTRUE\t/\tTRUE\t${expiration}\t${name}\t${value}`);
        }

        return lines.join("\n");
    }

    /**
     * Convert cookies to Base64 encoded JSON
     *
     * @param cookies - Cookies object
     * @returns Base64 encoded string
     */
    static toBase64(cookies: Cookies): string {
        return Buffer.from(JSON.stringify(cookies)).toString("base64");
    }

    /**
     * Filter cookies to only essential ones for Facebook/Messenger
     *
     * @param cookies - Cookies object
     * @returns Filtered cookies with only essential keys
     */
    static filterEssential(cookies: Cookies): Cookies {
        const essential = ["c_user", "xs", "datr", "fr", "sb", "wd", "presence"];
        const result: Record<string, string> = {};

        for (const key of essential) {
            if (cookies[key]) {
                result[key] = cookies[key] as string;
            }
        }

        return result as Cookies;
    }

    /**
     * Validate that cookies contain required fields
     *
     * @param cookies - Cookies object
     * @returns True if cookies are valid
     */
    static validate(cookies: Cookies): boolean {
        const required = ["c_user", "xs"];
        return required.every(key => cookies[key] && cookies[key].length > 0);
    }

    /**
     * Get missing required cookies
     *
     * @param cookies - Cookies object
     * @returns Array of missing cookie names
     */
    static getMissing(cookies: Cookies): string[] {
        const required = ["c_user", "xs"];
        return required.filter(key => !cookies[key] || cookies[key].length === 0);
    }

    /**
     * Check if a string is valid Base64
     */
    private static isBase64(str: string): boolean {
        if (str.length < 4) return false;
        // Check if string contains only valid base64 characters
        const base64Regex = /^[A-Za-z0-9+/]+=*$/;
        if (!base64Regex.test(str)) return false;
        // Try to detect if it's likely base64 (not a simple word)
        return str.length > 20 && str.length % 4 === 0;
    }
}
