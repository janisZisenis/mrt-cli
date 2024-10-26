const core = require('@actions/core');
const fetch = require('node-fetch');
const fs = require('fs');
const path = require('path');
const os = require('os');

async function run() {
    try {
        const version = core.getInput('version');
        const platform = os.platform();
        let url;

        if (platform === 'darwin') {
            url = `https://github.com/janisZisenis/multi-repo-tool/releases/download/v${version}/mrt-darwin-amd64`;
        } else if (platform === 'linux') {
            url = `https://github.com/janisZisenis/multi-repo-tool/releases/download/v${version}/mrt-linux-amd64`;
        } else {
            throw new Error(`Unsupported platform: ${platform}`);
        }

        const binDir = path.join(os.homedir(), 'bin');
        const binPath = path.join(binDir, 'mrt');

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

        core.info(`MRT version ${version} installed successfully`);
    } catch (error) {
        core.setFailed(`Action failed with error ${error}`);
    }
}

run();
