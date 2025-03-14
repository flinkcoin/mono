import { getPDQMD5Hash } from "../pdq/pdq.js";
import * as pdq from "../pdq/pdq-photo-hasher.js";
import * as Magick from '../node_modules/wasm-imagemagick/dist/magickApi.js';

const imgElement = document.getElementById('hoveredImage');
const imageHashElement = document.getElementById('imageHash');
const scoreElement = document.getElementById('score');
const tooltip = document.getElementById('tooltip');
const spinner = document.getElementById('spinner');

chrome.runtime.onMessage.addListener(async (message, sender, sendResponse) => {
    if (message.type === 'IMAGE_SRC') {

        imgElement.src = "";
        imageHashElement.textContent = "";
        scoreElement.textContent = "";
        tooltip.style.display = 'none';
        spinner.style.display = 'block';

        let blobUrl;
        try {
            // Fetch image data
            const response = await fetch(message.src);
            if (!response.ok) {
                throw new Error('Failed to fetch image');
            }
            const arrayBuffer = await response.arrayBuffer();
            const imageData = new Uint8Array(arrayBuffer);

            // Create Blob and URL for image display
            const blob = new Blob([imageData], { type: response.headers.get('Content-Type') || 'image/jpeg' });
            blobUrl = URL.createObjectURL(blob);

            // Promise for image loading
            const imageLoadPromise = new Promise((resolve, reject) => {
                imgElement.onload = () => resolve();
                imgElement.onerror = () => reject(new Error('Image failed to load'));
                imgElement.src = blobUrl;
            });

            // Prepare parameters for PDQ hash computation
            const extension = message.src.split('.').pop().toLowerCase() || 'jpg'; // Default to 'jpg' if no other extension
            const fname = `input.${extension}`;
            const tempfname = 'output.pnm';
            const hashPromise = getPDQMD5Hash(imageData, fname, tempfname, true, false, null);

            // Wait for both image load and hash computation
            const [_, hashResult] = await Promise.all([imageLoadPromise, hashPromise]);
            imageHashElement.textContent = "Image hash: " + hashResult;

        } catch (error) {
            console.error('Error:', error);
            imageHashElement.textContent = "Error computing hash: " + error.message;
        } finally {
            spinner.style.display = 'none';
            if (blobUrl) {
                URL.revokeObjectURL(blobUrl);
            }
        }
    } else if (message.type === 'CLEAR_IMAGE') {
        // Clear the display if needed
        // imgElement.src = "";
        // imageHashElement.textContent = "";
        // scoreElement.textContent = "";
        // tooltip.style.display = 'none';
        // spinner.style.display = 'none';
    }
});


async function testMagick(message) {
    const response = await fetch(message.src);
    if (!response.ok) {
        throw new Error('Failed to fetch image');
    }
    const arrayBuffer = await response.arrayBuffer();
    const inputData = new Uint8Array(arrayBuffer);

    // Determine file extension.
    const extension = message.src.split('.').pop().toLowerCase();
    const inputFilename = `input.${extension}`;
    const tempFilename = 'output.pnm';

    // Convert image to .pnm using wasm-imagemagick.
    const inputFile = { name: inputFilename, content: inputData };
    const command = ["convert", inputFilename, "-density", "400x400", tempFilename];
    const processedFiles = await Magick.Call([inputFile], command);
    if (!processedFiles.length) {
        throw new Error('Image conversion failed');
    }

    // Get the converted .pnm data (if needed for further processing)
    const outputFile = processedFiles[0];
    const outputData = new Uint8Array(await outputFile.blob.arrayBuffer());
    console.log(outputData)
}