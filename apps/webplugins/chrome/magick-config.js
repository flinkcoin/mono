window.MagickModule = {
    locateFile: (path) => {
        if (path === 'magick.wasm') {
            return chrome.runtime.getURL('lib/magick/magick.wasm');
        }
        return path;
    },
    onRuntimeInitialized: () => {
        window.MagickReady = true;
        console.log('Magick module initialized');
    }
};