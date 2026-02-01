/*
 * meta-messenger.js, Unofficial Meta Messenger Chat API for NodeJS
 * Copyright (c) 2026 Elysia and contributors
 */

import { execSync } from "node:child_process";

export function detectPlatform() {
    const { platform } = process;
    const { arch } = process;
    const isMusl = detectMusl();

    const libc = platform === "linux" ? (isMusl ? "musl" : "gnu") : "";
    const triplet = platform === "linux" ? `${platform}-${arch}-${libc}` : `${platform}-${arch}`;

    const ext = platform === "win32" ? "dll" : platform === "darwin" ? "dylib" : "so";

    return { platform, arch, libc, triplet, ext };
}

function detectMusl() {
    try {
        if (process.platform !== "linux") return false;
        if (process.report && typeof process.report.getReport === "function") {
            const rep = process.report.getReport();
            const glibc = rep.header && rep.header.glibcVersionRuntime;
            return !glibc;
        }
<<<<<<< Updated upstream
    }
    catch {
=======
    } catch {
>>>>>>> Stashed changes
        //
    }
    try {
        const out = execSync("ldd --version 2>&1 || true", { encoding: "utf8" });
        return /musl/i.test(out);
<<<<<<< Updated upstream
    }
    catch {
=======
    } catch {
>>>>>>> Stashed changes
        //
    }
    return false;
}

export default detectPlatform;
