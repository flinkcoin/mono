window.Module = {
    locateFile: (path) => {
        if (path === 'magick.wasm') {
            return chrome.runtime.getURL('lib/magick/magick.wasm');
        }

        if (path === 'pdq-photo-hasher.wasm') {
            return chrome.runtime.getURL('lib/pdq/pdq-photo-hasher.wasm');
        }
        return path;
    },
};