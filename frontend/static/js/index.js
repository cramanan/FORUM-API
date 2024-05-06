import { Home, Connect, _404 } from "./views.js";

export const APIendpoint = "http://localhost:8081";

const Authorized = async () => {
    try {
        const authorized = await fetch(APIendpoint, { credentials: "include" });
        return authorized.ok;
    } catch (reason) {
        console.log(reason);
    }
    return false;
};

const router = async () => {
    const routes = [
        { path: "/", view: Home, Authorized },
        { path: "/connect", view: Connect, Authorized: async () => true },
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
            route: { path: "/", view: _404, Authorized: async () => true },
            isMatch: false,
        };
    }

    if (!(await match.route.Authorized())) {
        match = {
            route: { path: "/connect", view: Connect },
            isMatch: true,
        };
        history.pushState(null, null, "/connect");
    }

    const view = new match.route.view();
    document.getElementById("root").innerHTML = await view.getHtml();
};

window.addEventListener("popstate", router);

export const navigateTo = (url) => {
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
