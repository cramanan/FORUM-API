import { Home, Connect, _404, Profile, Chat } from "./views.js";

const APIendpoint = "localhost:8080/api";

const Protected = async () => {
  try {
    const authorized = await fetch(`http://${APIendpoint}/auth`, {
      credentials: "include",
    });
    return authorized.ok;
  } catch (reason) {
    console.log(reason);
  }
  return false;
};

const router = async () => {
  const routes = [
    { path: "/", view: Home, Protected },
    { path: "/connect", view: Connect, Protected: async () => true },
    { path: "/profile", view: Profile, Protected },
    { path: "/chat", view: Chat, Protected },
  ];

  const potentialMatches = routes.map((route) => ({
    route: route,
    isMatch: location.pathname === route.path,
  }));

  let match = potentialMatches.find((potentialMatch) => potentialMatch.isMatch);

  if (!match) {
    match = {
      route: { path: "/", view: _404, Protected: async () => true },
      isMatch: false,
    };
  }

  if (!(await match.route.Protected())) {
    match = {
      route: { path: "/connect", view: Connect },
      isMatch: true,
    };
    history.pushState(null, null, "/connect");
  }

  const view = new match.route.view();
  document.getElementById("root").innerHTML = await view.getHtml();
  view.bindListeners();
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

export { APIendpoint, navigateTo };
