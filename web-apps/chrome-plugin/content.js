(() => {
    console.log("Content script loaded");

    // Variables to manage state
    let lastImg = null;         // Tracks the last detected image
    let frameRequested = false; // Prevents multiple frame requests
    let lastX = -1;             // Stores the latest mouse X position
    let lastY = -1;             // Stores the latest mouse Y position

    function handleMouseMove(event) {
        // Capture the current mouse coordinates
        lastX = event.clientX;
        lastY = event.clientY;

        // Only request a frame if one isn’t already scheduled
        if (!frameRequested) {
            frameRequested = true;
            requestAnimationFrame(() => {
                // Reset the flag so a new frame can be requested
                frameRequested = false;

                // Find elements under the cursor at the latest position
                const elements = document.elementsFromPoint(lastX, lastY);
                let currentImg = null;

                // Look for the first <img> element under the cursor
                for (const el of elements) {
                    if (el.tagName.toLowerCase() === 'img') {
                        currentImg = el;
                        break;
                    }
                }

                // If a new image is detected, send its source
                if (currentImg && currentImg !== lastImg) {
                    chrome.runtime.sendMessage({ type: 'IMAGE_SRC', src: currentImg.src });
                    lastImg = currentImg;
                }
                // If the mouse leaves an image, reset the last image
                else if (!currentImg && lastImg) {
                    lastImg = null;
                    // Optional: Notify the background script when leaving an image
                    // chrome.runtime.sendMessage({ type: 'IMAGE_SRC', src: null });
                }
            });
        }
    }

    // Attach the throttled handler to the document
    document.addEventListener('mousemove', handleMouseMove);
})();