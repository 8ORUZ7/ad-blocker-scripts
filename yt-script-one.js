(function () {
    'use strict';

    const logAction = (action) => console.log(`[Ad Blocker] ${action}`);

    const showNotification = (message, type = "info") => {
        const notification = document.createElement("div");
        Object.assign(notification.style, {
            position: "fixed",
            bottom: "10px",
            right: "10px",
            zIndex: "9999",
            backgroundColor: type === "error" ? "#ff4c4c" : "#4caf50",
            color: "#fff",
            padding: "10px",
            borderRadius: "4px",
            fontFamily: "Arial, sans-serif",
            fontSize: "12px",
            maxWidth: "300px",
            textAlign: "center",
            boxShadow: "0 2px 10px rgba(0,0,0,0.2)",
        });
        notification.textContent = message;
        document.body.appendChild(notification);
        setTimeout(() => notification.remove(), 5000);
    };

    const skipAds = () => {
        const adSelectors = [
            ".ytp-ad-skip-button",
            ".ytp-ad-overlay-close-button",
            "div.ytp-ad-module",
        ];
        adSelectors.forEach((selector) => {
            document.querySelectorAll(selector).forEach((el) => {
                el.click();
                el.remove();
                logAction(`Skipped ad: ${selector}`);
            });
        });
    };

 
    const bypassAdBlocker = () => {
        const observer = new MutationObserver((mutations) => {
            mutations.forEach((mutation) => {
                mutation.addedNodes.forEach((node) => {
                    if (node.nodeType === 1 && node.tagName === "YTD-POPUP-CONTAINER" && node.textContent.includes("ad blocker")) {
                        node.remove();
                        logAction("Removed ad blocker detection popup.");
                    }
                });
            });
        });
        observer.observe(document.body, { childList: true, subtree: true });
        logAction("Ad blocker bypass active.");
    };

  
    const adaptToMobile = () => {
        if (!document.querySelector('meta[name="viewport"]')) {
            const viewport = document.createElement("meta");
            viewport.name = "viewport";
            viewport.content = "width=device-width, initial-scale=1";
            document.head.appendChild(viewport);
            logAction("Added viewport meta tag for mobile.");
        }
    };


    const checkScriptHealth = async () => {
        try {
            const response = await fetch("https://github.com/8ORUZ7/ad-blocker-scripts/blob/main/yt-script-one.js");
            if (response.ok) {
                const data = await response.json();
                if (data.latestVersion && data.latestVersion !== "2.3") {
                    showNotification("New version available. Update now!", "error");
                    logAction("Update available notification displayed.");
                }
            }
        } catch (error) {
            logAction("Failed to check for updates: " + error.message);
            showNotification("Script health check failed. Check updates manually.", "error");
        }
    };

    const init = () => {
        logAction("Initializing Ad Blocker...");
        showNotification("Ad Blocker activated!", "info");
        adaptToMobile();
        bypassAdBlocker();
        const observer = new MutationObserver(skipAds);
        observer.observe(document.body, { childList: true, subtree: true });
        setInterval(checkScriptHealth, 60000);
    };

    init();
})();
