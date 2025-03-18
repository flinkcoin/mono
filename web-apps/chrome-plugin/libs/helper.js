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

export { fromHex, base32Encode };