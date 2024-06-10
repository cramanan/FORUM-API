document.addEventListener("DOMContentLoaded", () => {
    // Handle initial load
    handleRoute();
    // Handle hash change
    window.addEventListener('hashchange', handleRoute);
});

function handleRoute() {
    const hash = window.location.hash || '#home';
    switch (hash) {
        case '#home':
            loadHome();
            break;
        case '#about':
            loadAbout();
            break;
        default:
            loadHome();
            break;
    }
}

function navigate(event, route) {
    event.preventDefault();
    window.location.hash = route;
}

function loadHome() {
    fetch('/api/home')
        .then(response => response.json())
        .then(data => {
            document.getElementById('content').innerHTML = `
                <h1>Home</h1>
                <p>${data.message}</p>
            `;
        })
        .catch(error => console.error('Error fetching data:', error));
}

function loadAbout() {
    fetch('/api/about')
        .then(response => response.json())
        .then(data => {
            document.getElementById('content').innerHTML = `
                <h1>About</h1>
                <p>${data.message}</p>
            `;
        })
        .catch(error => console.error('Error fetching data:', error));
}
