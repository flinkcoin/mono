let pluginEnabled = false;

document.getElementById('toggleButton').addEventListener('click', () => {
    pluginEnabled = !pluginEnabled;

    chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
        chrome.tabs.sendMessage(tabs[0].id, { action: 'togglePlugin', enabled: pluginEnabled }, (response) => {
            document.getElementById('toggleButton').textContent = pluginEnabled ? 'Disable Plugin' : 'Enable Plugin';
        });
    });
});
