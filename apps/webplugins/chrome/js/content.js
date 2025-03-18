async function hashImage(url) {
    const response = await fetch(url);
    if (!response.ok) {
        throw new Error('Failed to fetch image');
    }
    const arrayBuffer = await response.arrayBuffer();
    // const copiedArrayBuffer = structuredClone(arrayBuffer);
    // const inputData = new Uint8Array(arrayBuffer);

    const extension = url.split('.').pop().toLowerCase();
    let inputFilename = `input.${extension}`;

    const tempFilename = 'output_temp.pnm';

    console.log('Processing inputName:', inputFilename);

    inputFilename="input.jpg";

    // Convert image to .pnm using wasm-imagemagick
    const inputFile = {name: inputFilename, content: Array.apply(null, new Uint8Array(arrayBuffer))};
    const files = [inputFile];
    const command = ["convert", inputFilename, "-density", "400x400", tempFilename];

    const hash = await new Promise((resolve, reject) => {
        chrome.runtime.sendMessage({
            type: 'PROCESS_IMAGE',
            files: files,
            command: command,
            tempFilename: tempFilename
        }, response => {
            // Check if response exists and is properly structured
            if (!response) {
                reject(new Error('No response received from background script'));
                return;
            }

            if (response.success) {
                console.log('Image hash:', response.result);
                resolve(response.result);
            } else {
                const error = response.error || 'Unknown error';
                console.error('Error:', error);
                reject(new Error(error));
            }
        });
    });

    return hash;
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
    const img = event.target;
    if (img.tagName.toLowerCase() !== 'img') return;

    let hash = hashCache.get(img);
    if (!hash) {
        try {

            hash = await hashImage(img.src);
            // hash = await hashImageSHA256(img.src);
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
        sendResponse({status: 'ok', enabled: pluginEnabled});
    }
});
