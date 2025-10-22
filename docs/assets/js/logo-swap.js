(function () {
    function currentScheme() {
        // Material sets data-md-color-scheme on <html>. Usually "default" or "slate"
        var scheme = document.documentElement.getAttribute('data-md-color-scheme');
        if (scheme) return scheme;
        // Fallback to OS preference
        return window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches
            ? 'slate'
            : 'default';
    }

    function updateLogo() {
        var img = document.querySelector('img.brand-logo');
        if (!img) return;
        var isDark = currentScheme() === 'slate';
        var darkSrc = img.getAttribute('data-dark');
        var lightSrc = img.getAttribute('data-light');
        var next = isDark ? darkSrc : lightSrc;
        if (next && img.getAttribute('src') !== next) {
            img.setAttribute('src', next);
        }
    }

    // Run on load
    document.addEventListener('DOMContentLoaded', updateLogo);

    // Observe theme changes (Material toggler updates this attribute on <html>)
    var obs = new MutationObserver(function (mutations) {
        for (var m of mutations) {
            if (m.type === 'attributes' && m.attributeName === 'data-md-color-scheme') {
                updateLogo();
            }
        }
    });
    obs.observe(document.documentElement, { attributes: true });

    // Also listen for OS scheme changes as a fallback
    if (window.matchMedia) {
        var mq = window.matchMedia('(prefers-color-scheme: dark)');
        if (mq.addEventListener) mq.addEventListener('change', updateLogo);
        else if (mq.addListener) mq.addListener(updateLogo); // Safari <14
    }
})();
