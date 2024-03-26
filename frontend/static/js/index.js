console.log("JS is loaded");

import Login from "./views/Login.js";
import Home from "./views/Home.js";
import ERROR404 from "./views/ERROR404.js";
import Register from "./views/Register.js";

const navigateTo = (url) => {
    history.pushState(null, null, url);
    router();
};

const router = async () => {
    const routes = [
        { path: "/", view: Home },
        { path: "/login", view: Login },
        { path: "/register", view: Register },
    ];

    const potentialMatches = routes.map((route) => {
        return {
            route: route,
            isMatch: location.pathname === route.path,
        };
    });

    let match = potentialMatches.find(
        (potentialMatch) => potentialMatch.isMatch
    );

    if (!match) {
        match = {
            route: { path: "/", view: ERROR404 },
            isMatch: false,
        };
    }
    const view = new match.route.view();
    view.setCSS();
    document.getElementById("root").innerHTML = await view.getHtml();
};

window.addEventListener("popstate", router);

document.addEventListener("DOMContentLoaded", () => {
    document.body.addEventListener("click", (e) => {
        if (e.target.matches("[data-link]")) {
            e.preventDefault();
            navigateTo(e.target.href);
        }
    });
    router();
});

// const socket = new WebSocket("ws://localhost:8081/");
// socket.onopen = () => console.log("connected");
// socket.onclose = () => console.log("disconnected");

//Login
window.HandleRegisterSubmit = (event) => {
    event.preventDefault();
    const data = new FormData(event.target);
    fetch("http://localhost:8081/register", { method: "post", body: data })
        .then((resp) => {
            if (resp.ok) {
                navigateTo("/");
            }
            return resp.json();
        })
        .then((data) => {
            document.getElementById("server-error").textContent = data.message;
        });
};

window.HandleLoginSubmit = (event) => {
    event.preventDefault();
    const data = new FormData(event.target);
    fetch("http://localhost:8081/login", { method: "post", body: data })
        .then((resp) => {
            if (resp.ok) {
                navigateTo("/");
            }
            return resp.json();
        })
        .then((data) => {
            document.getElementById("server-error").textContent = data.message;
        })
        .catch(() => {});
};
