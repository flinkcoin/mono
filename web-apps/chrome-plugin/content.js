(async () => {
    console.log("Content script loaded");

    async function imageMouseOverHandler(event) {
        const img = event.target;
        if (img.tagName.toLowerCase() !== 'img') return;
        chrome.runtime.sendMessage({ type: 'IMAGE_SRC', src: img.src });
    }

    function imageMouseOutHandler(event) {
        const img = event.target;
        if (img.tagName.toLowerCase() !== 'img') return;
        chrome.runtime.sendMessage({ type: 'CLEAR_IMAGE' });
    }

    function scanImages() {
        const images = document.querySelectorAll('img');
        images.forEach(img => {
            img.addEventListener('mouseover', imageMouseOverHandler);
            img.addEventListener('mouseout', imageMouseOutHandler);
        });
    }

    scanImages();
})();
