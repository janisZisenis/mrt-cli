const cache = require('@actions/cache');
const core = require('@actions/core');
const exec = require('@actions/exec');
const fetch = require('node-fetch');
const fs = require('fs');
const path = require('path');
const os = require('os');

async function run() {
    try {
        const version = core.getInput('version') || 'v0.1.0';
        const platform = os.platform();
        const cacheKey = `mrt-${version}-${platform}`;
        const binDir = path.join(os.homedir(), 'bin');
        const binPath = path.join(binDir, 'mrt');

        // Restore cache
        const cacheId = await cache.restoreCache([binDir], cacheKey);
        if (cacheId) {
            core.info(`Cache hit for key: ${cacheKey}`);
            core.addPath(binDir);
            await exec.exec('mrt', ['version']);
            return;
        }

        core.info(`Cache miss for key: ${cacheKey}`);
        let url;

        if (platform === 'darwin') {
            url = `https://github.com/janisZisenis/multi-repo-tool/releases/download/v${version}/mrt-darwin-amd64`;
        } else if (platform === 'linux') {
            url = `https://github.com/janisZisenis/multi-repo-tool/releases/download/v${version}/mrt-linux-amd64`;
        } else {
            throw new Error(`Unsupported platform: ${platform}`);
        }

        if (!fs.existsSync(binDir)){
            fs.mkdirSync(binDir);
        }

        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`Failed to download MRT binary: ${response.statusText}`);
        }

        const fileStream = fs.createWriteStream(binPath);
        await new Promise((resolve, reject) => {
            response.body.pipe(fileStream);
            response.body.on("error", reject);
            fileStream.on("finish", resolve);
        });

        fs.chmodSync(binPath, '0755');
        core.addPath(binDir);

        // Save cache
        await cache.saveCache([binDir], cacheKey);

        core.info(`MRT version ${version} installed successfully`);
    } catch (error) {
        core.setFailed(`Action failed with error ${error}`);
    }
}

run();
