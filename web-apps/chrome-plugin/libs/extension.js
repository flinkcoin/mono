// Mapping of common image MIME types to their standard extensions
const mimeToExt = {
    'image/jpeg': 'jpg',
    'image/png': 'png',
    'image/gif': 'gif',
    'image/bmp': 'bmp',
    'image/webp': 'webp',
    'image/svg+xml': 'svg'
};

/**
 * Extracts the MIME type from a data URL
 * @param {string} src - The data URL (e.g., "data:image/png;base64,...")
 * @returns {string|null} - The MIME type or null if not a valid data URL
 */
function getMimeFromDataUrl(src) {
    if (!src.startsWith('data:')) return null;
    const header = src.split(',')[0];
    const mimePart = header.split(';')[0];
    return mimePart.slice(5); // Remove 'data:' prefix
}

/**
 * Extracts the file extension from a regular URL
 * @param {string} src - The URL (e.g., "https://example.com/image.jpg")
 * @returns {string|null} - The extension or null if not found
 */
function getExtensionFromUrl(src) {
    try {
        const url = new URL(src);
        const pathname = url.pathname;
        const filename = pathname.split('/').pop();
        const parts = filename.split('.');
        if (parts.length > 1) {
            const ext = parts.pop().toLowerCase();
            // Normalize 'jpeg' to 'jpg'
            return ext === 'jpeg' ? 'jpg' : ext;
        }
        return null;
    } catch (e) {
        return null;
    }
}

/**
 * Determines the image file extension based on message.src and response
 * @param {Object} message - Object with a 'src' property (URL or data URL)
 * @param {Response|null} response - Fetch response object or null
 * @returns {string} - The image extension (e.g., 'jpg', 'png') or 'unknown'
 */
function getImageExtension(src, response) {

    // Handle data URLs
    if (src.startsWith('data:')) {
        const mime = getMimeFromDataUrl(src);
        if (mime && mime in mimeToExt) {
            return mimeToExt[mime];
        }
        return 'unknown';
    }

    // Handle regular URLs
    if (response && response.headers) {
        const contentType = response.headers.get('Content-Type');
        if (contentType) {
            const mime = contentType.split(';')[0].trim();
            if (mime in mimeToExt) {
                return mimeToExt[mime];
            }
        }
    }

    // Fallback to URL extension
    const ext = getExtensionFromUrl(src);
    const knownExtensions = ['jpg', 'png', 'gif', 'bmp', 'webp', 'svg'];
    if (ext && knownExtensions.includes(ext)) {
        return ext;
    }

    return 'unknown';
}

export {getImageExtension}