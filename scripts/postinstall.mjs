/*
 * meta-messenger.js, Unofficial Meta Messenger Chat API for NodeJS
 * Copyright (c) 2026 Elysia and contributors
 */

import { spawnSync } from "node:child_process";
import { existsSync } from "node:fs";
import { copyFile, mkdir, rm } from "node:fs/promises";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

import { detectPlatform } from "./detect-platform.mjs";
import { downloadPrebuilt } from "./download-prebuilt.mjs";
import { packageJson as pkg } from "./package.mjs";

const __dirname = dirname(fileURLToPath(import.meta.url));

async function copyIfExists(src, dst) {
    try {
        await mkdir(dirname(dst), { recursive: true });
        await copyFile(src, dst);
        return true;
<<<<<<< Updated upstream
    }
    catch (err) {
=======
    } catch (err) {
>>>>>>> Stashed changes
        if (err?.code === "ENOENT") return false;
        throw err;
    }
}

async function run() {
<<<<<<< Updated upstream
    console.log(`[${pkg.name}] postinstall started`);
    const startTime = Date.now();

=======
>>>>>>> Stashed changes
    if (process.env.MESSAGIX_SKIP_POSTINSTALL === "true") {
        console.log(`[${pkg.name}] Skipping postinstall (MESSAGIX_SKIP_POSTINSTALL=true)`);
        return;
    }

<<<<<<< Updated upstream
    console.log(`[${pkg.name}] Detecting platform...`);
    const { triplet, ext } = detectPlatform();
    console.log(`[${pkg.name}] Platform: ${triplet}, ext: ${ext}`);

    const buildOut = join(__dirname, "..", "build", `messagix.${ext}`);
    console.log(`[${pkg.name}] Build output path: ${buildOut}`);
=======
    const { triplet, ext } = detectPlatform();
    const buildOut = join(__dirname, "..", "build", `messagix.${ext}`);
>>>>>>> Stashed changes

    // 1) Prefer local prebuilt shipped in npm tarball
    const prebuiltDir = join(__dirname, "..", "prebuilt", triplet);
    const prebuilt = join(prebuiltDir, `messagix.${ext}`);
<<<<<<< Updated upstream
    console.log(`[${pkg.name}] Checking local prebuilt at: ${prebuilt}`);
    if (await copyIfExists(prebuilt, buildOut)) {
        console.log(`[${pkg.name}] Using local prebuilt for ${triplet}`);
        if (process.env.MESSAGIX_KEEP_PREBUILT !== "true") {
            try { await rm(prebuiltDir, { recursive: true, force: true }); }
            catch {
                //
            }
        }
        console.log(`[${pkg.name}] postinstall completed in ${Date.now() - startTime}ms`);
=======
    if (await copyIfExists(prebuilt, buildOut)) {
        console.log(`[${pkg.name}] Using local prebuilt for ${triplet}`);
        if (process.env.MESSAGIX_KEEP_PREBUILT !== "true") {
            try {
                await rm(prebuiltDir, { recursive: true, force: true });
            } catch (error) {
                console.error(error);
            }
        }
>>>>>>> Stashed changes
        return;
    }
    console.log(`[${pkg.name}] No local prebuilt found`);

    // 2) Try remote prebuilt from GitHub Releases
<<<<<<< Updated upstream
    console.log(`[${pkg.name}] Attempting to download remote prebuilt...`);
    const downloadStart = Date.now();
    try {
        if (await downloadPrebuilt()) {
            console.log(`[${pkg.name}] Download completed in ${Date.now() - downloadStart}ms`);
            console.log(`[${pkg.name}] postinstall completed in ${Date.now() - startTime}ms`);
            return;
        }
    }
    catch (err) {
        console.log(`[${pkg.name}] Download failed after ${Date.now() - downloadStart}ms:`, err?.message || String(err));
=======
    try {
        if (await downloadPrebuilt()) return;
    } catch (error) {
        console.error(error);
>>>>>>> Stashed changes
    }

    // 3) No prebuilt available. Try local build if allowed
    if (process.env.MESSAGIX_BUILD_FROM_SOURCE === "true") {
        console.log(`[${pkg.name}] No prebuilt found. Attempting local build...`);
        try {
<<<<<<< Updated upstream
            const res = spawnSync(
                process.platform === "win32" ? "npm.cmd" : "npm",
                ["run", "build:go"],
                { cwd: join(__dirname, ".."), stdio: "inherit", env: process.env },
            );
            if (res.status === 0 && existsSync(buildOut)) return;
            console.warn(`[${pkg.name}] Local build did not produce the native library.`);
        }
        catch (err) {
=======
            const res = spawnSync(process.platform === "win32" ? "npm.cmd" : "npm", ["run", "build:go"], {
                cwd: join(__dirname, ".."),
                stdio: "inherit",
                env: process.env,
            });
            if (res.status === 0 && existsSync(buildOut)) return;
            console.warn(`[${pkg.name}] Local build did not produce the native library.`);
        } catch (err) {
>>>>>>> Stashed changes
            console.warn(`[${pkg.name}] Local build failed:`, err?.message || String(err));
        }
    }

    console.warn(`[${pkg.name}] No prebuilt found (local/remote) and no local build completed.`);
    console.warn(`[${pkg.name}] Expected triplet: ${triplet}, file: messagix-${triplet}.${ext}`);
    console.warn(
        `[${pkg.name}] You can:\n` +
<<<<<<< Updated upstream
        "  - set MESSAGIX_BUILD_FROM_SOURCE=true and re-run install\n" +
        "  - or build manually with: npm run build:go",
    );
}

run()
    .then(() => process.exit(0))
    .catch(err => {
        console.error(`[${pkg.name}] postinstall failed:`, err?.message || String(err));
        process.exit(1);
    });
=======
            "  - set MESSAGIX_BUILD_FROM_SOURCE=true and re-run install\n" +
            "  - or build manually with: npm run build:go",
    );
}

run().catch(err => {
    console.error(`[${pkg.name}] postinstall failed:`, err?.message || String(err));
});
>>>>>>> Stashed changes
