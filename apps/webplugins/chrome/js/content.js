// import { Call as MagickCall } from './lib/magick/magickApi.js';
// const WASM_MAGICK_URL = chrome.runtime.getURL('./lib/magick/magickApi.js');

// const loadWasmModule = async () => {
//     s
//     const mod = await import(WASM_MAGICK_URL);

//     // default export is an init function
//     const isOk = await mod.default().catch((e) => {
//         console.warn('Failed to init wasm module in content script. Probably CSP of the page has restricted wasm loading.', e);
//         return null;
//     });

//     return isOk ? mod : null;
// };

async function hashImage(url) {
    try {
        // loadWasmModule();
        // Fetch the image
        // const response = await fetch(url);
        // if (!response.ok) {
        //     throw new Error('Failed to fetch image');
        // }
        // const arrayBuffer = await response.arrayBuffer();
        // const inputData = new Uint8Array(arrayBuffer);

        // const extension = url.split('.').pop().toLowerCase();
        // const inputFilename = `input.${extension}`;
        // const tempFilename = 'output.pnm';

        // // Convert image to .pnm using wasm-imagemagick
        // const inputFile = { name: inputFilename, content: inputData };
        // const command = ["convert", inputFilename, "-density", "400x400", tempFilename];
        // const processedFiles = await MagickCall([inputFile], command);
        // if (!processedFiles.length) {
        //     throw new Error('Image conversion failed');
        // }

        // // Get the converted .pnm data
        // const outputFile = processedFiles[0];
        // const outputData = new Uint8Array(await outputFile.blob.arrayBuffer());

        // console.log(outputData);
        // Write to PDQ module's Emscripten FS
        // pdqModule.FS.writeFile(tempFilename, outputData);

        // // Compute PDQ hash
        // const hash = pdqModule.ccall('getHash', 'string', ['string'], [tempFilename]);

        // // Clean up
        // pdqModule.FS.unlink(tempFilename);

        return hash;
    } catch (e) {
        console.error(e);
        return 'Error: ' + e.message;
    }
}


async function hashImageSHA256(url) {
    const response = await fetch(url);
    const arrayBuffer = await response.arrayBuffer();
    const hashBuffer = await crypto.subtle.digest('SHA-256', arrayBuffer);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    return hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
}


// PLUGIN UI LOGIC
let pluginEnabled = false;
let tooltipDiv = null;
const hashCache = new WeakMap();

function createTooltip() {
    tooltipDiv = document.createElement('div');
    tooltipDiv.style.position = 'fixed';
    tooltipDiv.style.background = 'rgba(0, 0, 0, 0.75)';
    tooltipDiv.style.color = '#fff';
    tooltipDiv.style.padding = '5px';
    tooltipDiv.style.borderRadius = '3px';
    tooltipDiv.style.zIndex = '9999';
    tooltipDiv.style.fontSize = '12px';
    tooltipDiv.style.pointerEvents = 'none';
    tooltipDiv.style.display = 'none';
    document.body.appendChild(tooltipDiv);
}

function showTooltip(text, x, y) {
    if (!tooltipDiv) createTooltip();
    tooltipDiv.textContent = text;
    tooltipDiv.style.left = (x + 10) + 'px';
    tooltipDiv.style.top = (y + 10) + 'px';
    tooltipDiv.style.display = 'block';
}

function hideTooltip() {
    if (tooltipDiv) {
        tooltipDiv.style.display = 'none';
    }
}

async function imageMouseOverHandler(event) {
    console.log("Image mouse over handler");
    const img = event.target;
    if (img.tagName.toLowerCase() !== 'img') return;

    let hash = hashCache.get(img);
    if (!hash) {
        try {
            hash = await hashImage(img.src);
            hashCache.set(img, hash);
        } catch (e) {
            hash = 'Error: ' + e.message;
        }
    }
    showTooltip(hash, event.clientX, event.clientY);
}

function imageMouseOutHandler(event) {
    const img = event.target;
    if (img.tagName.toLowerCase() !== 'img') return;
    hideTooltip();
}

function enablePlugin() {
    if (!pluginEnabled) {
        pluginEnabled = true;
        document.addEventListener('mouseover', imageMouseOverHandler, true);
        document.addEventListener('mouseout', imageMouseOutHandler, true);
        console.log('Image hash plugin enabled');
    }
}

function disablePlugin() {
    if (pluginEnabled) {
        pluginEnabled = false;
        document.removeEventListener('mouseover', imageMouseOverHandler, true);
        document.removeEventListener('mouseout', imageMouseOutHandler, true);
        hideTooltip();
        console.log('Image hash plugin disabled');
    }
}

chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
    if (message.action === 'togglePlugin') {
        if (message.enabled) {
            enablePlugin();
        } else {
            disablePlugin();
        }
        sendResponse({ status: 'ok', enabled: pluginEnabled });
    }
});
