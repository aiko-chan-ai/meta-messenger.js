/*
 * meta-messenger.js, Unofficial Meta Messenger Chat API for NodeJS
 * Copyright (c) 2026 Elysia and contributors
 */

import { createWriteStream } from "node:fs";
import { mkdir, rename, unlink } from "node:fs/promises";
import https from "node:https";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

import { detectPlatform } from "./detect-platform.mjs";
import { packageJson as pkg } from "./package.mjs";

const __dirname = dirname(fileURLToPath(import.meta.url));

function defaultRepoSlug() {
    const repo = pkg.repository;
    if (typeof repo === "string") {
        const m = repo.match(/github:(.+)/i);
        if (m) return m[1];
        if (/^[\w-]+\/[\w.-]+$/.test(repo)) return repo;
    } else if (repo && typeof repo === "object" && repo.url) {
        const m = repo.url.match(/github\.com[:/]+([^#]+?)(?:\.git)?$/i);
        if (m) return m[1];
    }
    return "yumi-team/meta-messenger.js";
}

function buildBaseURL() {
    const repo = defaultRepoSlug();
    const versionTag = `v${pkg.version}`;
    const baseURL = `https://github.com/${repo}/releases/download/${versionTag}`;
    return baseURL.replace(/\/$/, "");
}

function httpGet(url, redirectCount = 0) {
    console.log(`[${pkg.name}] HTTP GET: ${url}${redirectCount > 0 ? ` (redirect #${redirectCount})` : ""}`);
    const reqStart = Date.now();
    return new Promise((resolve, reject) => {
        https
            .get(url, res => {
                console.log(`[${pkg.name}] Response: HTTP ${res.statusCode} in ${Date.now() - reqStart}ms`);
                if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
                    console.log(`[${pkg.name}] Following redirect to: ${res.headers.location}`);
                    return resolve(httpGet(res.headers.location, redirectCount + 1));
                }
                if (res.statusCode !== 200) {
                    return reject(new Error(`HTTP ${res.statusCode} for ${url}`));
                }
                resolve(res);
            })
            .on("error", err => {
                console.log(`[${pkg.name}] HTTP error after ${Date.now() - reqStart}ms:`, err?.message);
                reject(err);
            });
    });
}

async function downloadTo(url, dstPath) {
    console.log(`[${pkg.name}] Downloading to: ${dstPath}`);
    await mkdir(dirname(dstPath), { recursive: true });
    const tmp = `${dstPath}.download`;
    try {
        await unlink(tmp);
    } catch {
        //
    }
    const res = await httpGet(url);
    console.log(`[${pkg.name}] Starting file write...`);
    const writeStart = Date.now();
    let bytesWritten = 0;
    await new Promise((resolve, reject) => {
        const out = createWriteStream(tmp);
        res.on("data", chunk => {
            bytesWritten += chunk.length;
        });
        res.pipe(out);
        res.on("error", reject);
        out.on("error", reject);
        out.on("finish", () => {
            console.log(`[${pkg.name}] File write completed: ${bytesWritten} bytes in ${Date.now() - writeStart}ms`);
            // Destroy the response stream to release the socket
            res.destroy();
            resolve();
        });
    });
    await rename(tmp, dstPath);
    console.log(`[${pkg.name}] File renamed to final destination`);
}

export async function downloadPrebuilt() {
    console.log(`[${pkg.name}] downloadPrebuilt() started`);
    const { triplet, ext } = detectPlatform();
    const baseURL = buildBaseURL();
    const filename = `messagix-${triplet}.${ext}`;
    const url = `${baseURL}/${filename}`;
    console.log(`[${pkg.name}] Target URL: ${url}`);

    const out = join(__dirname, "..", "build", `messagix.${ext}`);
    const totalStart = Date.now();
    try {
        await downloadTo(url, out);
        console.log(`[${pkg.name}] Downloaded prebuilt from ${url} in ${Date.now() - totalStart}ms`);
        return true;
    } catch (err) {
        console.warn(
            `[${pkg.name}] No remote prebuilt found at ${url} (after ${Date.now() - totalStart}ms):`,
            err?.message || String(err),
        );
        return false;
    }
}

// node scripts/download-prebuilt.mjs
const isMain = process.argv[1] && fileURLToPath(import.meta.url) === process.argv[1];
if (isMain) {
    downloadPrebuilt()
        .then(ok => {
            if (!ok) process.exit(1);
        })
        .catch(err => {
            console.error(`[${pkg.name}] download-prebuilt failed:`, err?.message || String(err));
            process.exit(1);
        });
}
