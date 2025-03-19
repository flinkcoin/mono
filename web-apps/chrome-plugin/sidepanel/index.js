import * as Magick from '../node_modules/wasm-imagemagick/dist/magickApi.js';
import {FS, Module} from '../pdq/pdq-photo-hasher.js';
import {fromHex, base32Encode} from '../libs/helper.js';
import {getImageExtension} from '../libs/extension.js';

const imgElement = document.getElementById('hoveredImage');
const imageHashElement = document.getElementById('imageHash');
const scoreElement = document.getElementById('score');
const tooltip = document.getElementById('tooltip');
const spinner = document.getElementById('spinner');


chrome.runtime.onMessage.addListener(async (message, sender, sendResponse) => {
    try {
        imgElement.src = message.src;
        imageHashElement.textContent = "";
        scoreElement.textContent = "";
        tooltip.style.display = 'none';
        spinner.style.display = 'none';

        console.log('Fetching image from:', message.src);

        // Fetch image from URL
        const response = await fetch(message.src);
        if (!response.ok) {
            throw new Error('Failed to fetch image');
        }
        const arrayBuffer = await response.arrayBuffer();
        // const copiedArrayBuffer = structuredClone(arrayBuffer);
        // const inputData = new Uint8Array(arrayBuffer);

        const extension = getImageExtension(message.src, response);
        let inputFilename = `input.${extension}`;
        const tempFilename = 'output_temp.pnm';

        console.log('Processing inputName:', inputFilename);

        //inputFilename="input.jpg";

        // Convert image to .pnm using wasm-imagemagick
        const inputFile = {name: inputFilename, content: Array.apply(null, new Uint8Array(arrayBuffer))};
        const files = [inputFile];
        const command = ["convert", inputFilename, "-density", "400x400", tempFilename];

        console.log('Processing image:');
        for (let file of files) {
            file.content = new Uint8Array(file.content).buffer;
        }

        let processedFiles = await Magick.Call(files, command);
        let firstOutputImage = processedFiles[0];

        const data = new Uint8Array(await firstOutputImage['blob'].arrayBuffer());

        let filename = "out.pnm";

        let stream = FS.open(filename, 'w+');
        FS.write(stream, data, 0, data.length, 0);
        FS.close(stream);

        var result = Module.ccall(
            'getHash',	// name of C function
            'string',	// return type
            ['string'],	// argument types
            [filename]	// arguments
        );

        let binary = fromHex(result);
        let base32r = base32Encode(binary)
        console.log(base32r)


        // Remove the file so that we can free some memory.
        FS.unlink(filename);

        // Display the hash
        imageHashElement.textContent = result;
        console.log("Final hash:", result);
    } catch (error) {
        console.error('Error processing image:', error);
    }
    return true;
});


async function testMagick(message) {
    const response = await fetch(message.src);
    if (!response.ok) {
        throw new Error('Failed to fetch image');
    }
    const arrayBuffer = await response.arrayBuffer();
    const imageData = new Uint8Array(arrayBuffer);

    // Create Blob and URL for image display
    const blob = new Blob([imageData], {type: response.headers.get('Content-Type') || 'image/jpeg'});
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


}