import { Home, Connect, _404 } from "./views.js";

export const APIendpoint = "http://localhost:8081";

const canActivate = async () => {
    try {
        const authorized = await fetch(APIendpoint, {
            credentials: "include",
        }).then((resp) => resp.ok);
        return authorized;
    } catch (reason) {
        console.log(reason);
    }
    return false;
};

const router = async () => {
    const routes = [
        { path: "/", view: Home, canActivate },
        { path: "/connect", view: Connect, canActivate: async () => true },
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
            route: { path: "/", view: _404, canActivate: async () => true },
            isMatch: false,
        };
    }

    if (!(await match.route.canActivate())) {
        match = {
            route: { path: "/connect", view: Connect },
            isMatch: true,
        };
        history.pushState(null, null, "/connect");
    }

    const view = new match.route.view();
    document.getElementById("root").innerHTML = await view.getHtml();
    view.setCSS();
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
    fetch(`${APIendpoint}/register`, {
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
            document.getElementById("register-server-error").textContent =
                data.message;
        });
};

window.HandleLoginSubmit = (event) => {
    event.preventDefault();
    const data = new FormData(event.target);
    fetch(`${APIendpoint}/login`, {
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
            document.getElementById("login-server-error").textContent =
                data.message;
        });
};

// window.Post = (event) => {
//     event.preventDefault();
//     const data = new FormData(event.target);
//     fetch(`${APIendpoint}/post`, {
//         method: "post",
//         body: data,
//         credentials: "include",
//     })
//         .then((resp) => resp.json())
//         .then((data) => console.log(data));
// };
