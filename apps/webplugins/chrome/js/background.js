import * as Magick from './lib/magick/magickApi.js';
import {FS, Module} from './lib/pdq/pdq-photo-hasher.js';


function fromHex(hexString) {
    if (typeof hexString !== 'string' || hexString.length % 2 !== 0) {
        throw new Error('Invalid hex string');
    }
    const array = new Uint8Array(hexString.length / 2);
    for (let i = 0; i < hexString.length; i += 2) {
        array[i / 2] = parseInt(hexString.substr(i, 2), 16);
    }
    return array;
}

function base32Encode(input) {
    const alphabet = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ234567';
    let output = '';
    let bits = 0;
    let value = 0;

    for (let i = 0; i < input.length; i++) {
        value = (value << 8) | input[i];
        bits += 8;

        while (bits >= 5) {
            output += alphabet[(value >>> (bits - 5)) & 31];
            bits -= 5;
        }
    }

    if (bits > 0) {
        output += alphabet[(value << (5 - bits)) & 31];
    }

    // Optionally add padding '=' characters to reach a length multiple of 8
    // while (output.length % 8 !== 0) {
    //     output += '=';
    // }

    return output;
}


// Handle messages from content scripts
chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {

    if (request.type === 'PROCESS_IMAGE') {

        (async () => {
            try {

                console.log('Processing image:');

                for (let file of request.files) {
                    file.content = new Uint8Array(file.content).buffer;
                }

                let processedFiles = await Magick.Call(request.files, request.command);
                let firstOutputImage = processedFiles[0].outputFiles[0];

                const data = new Uint8Array(await firstOutputImage['blob'].arrayBuffer());

                let filename ="out.pnm";

                let stream = FS.open(filename, 'w+');
                FS.write(stream, data, 0, data.length, 0);
                FS.close(stream);

                //let filename="output.jpg";


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

                sendResponse({success: true, result: base32r});
            } catch (e) {
                console.log(e);
                sendResponse({success: false, error: e.message});
            }
        })();
        return true;
    }

    sendResponse({ success: false, error: 'Unknown message type' });
    return true; // Keep message channel open for async response
});


