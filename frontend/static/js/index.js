import Home from "./views/Home.js";
import ERROR404 from "./views/ERROR404.js";
import Connect from "./views/Connect.js";

const router = async () => {
    const routes = [
        { path: "/", view: Home },
        { path: "/connect", view: Connect },
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

    const authorized = await fetch("http://localhost:8081/", {
        credentials: "include",
    }).then((resp) => resp.ok);
    console.log(authorized);

    const view = new match.route.view();
    view.setCSS();
    document.getElementById("root").innerHTML = await view.getHtml();
};

window.addEventListener("popstate", router);

const navigateTo = (url) => {
    history.pushState(null, null, url);
    router();
};

document.addEventListener("DOMContentLoaded", () => {
    document.body.addEventListener("click", (e) => {
        if (e.target.matches("[data-link]")) {
            e.preventDefault();
            navigateTo(e.target.href);
        }
    });
    router();
});

// Register
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
    fetch("http://localhost:8081/login", {
        method: "post",
        body: data,
        credentials: "include",
    })
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
