window.PDQModule = {
    locateFile: (path) => {
        if (path === 'pdq-photo-hasher.wasm') {
            return chrome.runtime.getURL('lib/pdq/pdq-photo-hasher.wasm');
        }
        return path;
    },
    onRuntimeInitialized: () => {
        window.PDQReady = true;
        console.log('PDQ module initialized');
    }
};