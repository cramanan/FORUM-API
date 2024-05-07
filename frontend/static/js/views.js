import { navigateTo, APIendpoint } from "./index.js";

class AbstractView {
    constructor() {}

    setTitle(title) {
        document.title = title;
    }

    setCSS(css) {
        document.querySelector("#viewcss").href = css;
    }

    async getHtml() {
        return "";
    }

    bindListeners() {}
}

class Connect extends AbstractView {
    constructor() {
        super();
        this.setTitle("Connect");
        this.setCSS("/static/css/connect.css");
    }

    async getHtml() {
        return `<div id="connect">
    <form id="login-form">
        <h1>Login</h1>
        <div id="login-server-error"></div>
        <label for="login-email">Email</label>
        <input type="email" id="login-email" name="login-email" />
        <label for="login-password">Password</label>
        <input type="password" id="login-password" name="login-password" />
        <button type="submit">Login</button>
    </form>
    <span id="sep"></span>
    <form id="register-form">
        <h1>Register</h1>
        <div id="register-server-error"></div>
        <label for="register-email">Email</label>
        <input type="email" id="register-email" name="register-email" />
        <label for="register-username">Username</label>
        <input type="text" id="register-username" name="register-username" />
        <label for="register-password">Password</label>
        <input type="password" id="register-password" name="register-password" />
        <label for="register-gender">Gender:</label>
        <select name="register-gender" id="register-gender">
            <option value="M">M</option>
            <option value="F">F</option>
            <option value="O">Other</option>
        </select>
        <label for="register-age">Age</label>
        <input type="number" name="register-age" id="register-age">
        <label for="register-first-name">First Name</label>
        <input type="text" name="register-first-name" id="register-first-name">
        <label for="register-last-name">Last Name</label>
        <input type="text" name="register-last-name" id="register-last-name">
        <button type="submit">Register</button>
    </form>
</div>`;
    }

    bindListeners() {
        const loginform = document.getElementById("login-form");
        const registerform = document.getElementById("register-form");
        loginform?.addEventListener("submit", this.HandleLoginSubmit);
        registerform?.addEventListener("submit", this.HandleRegisterSubmit);
    }

    async HandleRegisterSubmit(event) {
        event.preventDefault();
        try {
            const data = new FormData(event.target);
            const response = await fetch(`${APIendpoint}/register`, {
                method: "post",
                body: data,
                credentials: "include",
            });
            if (response.ok) {
                navigateTo("/");
            } else {
                const server_msg = await response.json();
                document.getElementById("login-server-error").textContent =
                    server_msg.message;
            }
        } catch (reason) {
            console.log(reason);
        }
    }

    async HandleLoginSubmit(event) {
        event.preventDefault();
        try {
            const data = new FormData(event.target);
            const response = await fetch(`${APIendpoint}/login`, {
                method: "post",
                body: data,
                credentials: "include",
            });
            if (response.ok) {
                navigateTo("/");
            } else {
                const server_msg = await response.json();
                document.getElementById("login-server-error").textContent =
                    server_msg.message;
            }
        } catch (reason) {
            console.log(reason);
        }
    }
}

class Home extends AbstractView {
    constructor() {
        super();
        this.setTitle("Real-Time Forum");
        this.setCSS("/static/css/home.css");
    }

    async getHtml() {
        const html = `<nav class="header">
            <h3><a href="/" id="main-title">REAL-TIME FORUM</a></h3>
        </nav>
        <main>
        <form id="post-form">
            <label for="post-content">Create a P0ST</label>
            <textarea name="post-content" id="post-content"></textarea>
            <button type="submit">P0ST</button>
        </form>
        <div id="all-posts">
            ${await this.fetchPosts()}
        </div>
        </main>
        <footer>
        </footer>`;
        return html;
    }

    bindListeners() {
        const postform = document.getElementById("post-form");
        postform?.addEventListener("submit", this.Post);
    }

    async Post(event) {
        event.preventDefault();
        try {
            const data = new FormData(event.target);
            const response = await fetch(`${APIendpoint}/post`, {
                method: "post",
                body: data,
                credentials: "include",
            });
            if (!response.ok) {
                console.log(response);
            }
        } catch (reason) {
            console.log(reason);
        }
    }

    async fetchPosts() {
        let postsHTML = "";
        try {
            const response = await fetch(`${APIendpoint}/getposts`);
            const datas = await response.json();
            datas.data.forEach((post) => {
                postsHTML += `<div class="post"><h2>${post.Username}</h2><p>${post.Content}</p></div>`;
            });
        } catch (error) {
            console.log(error);
        }
        return postsHTML;
    }
}

class _404 extends AbstractView {
    constructor() {
        super();
        this.setTitle("404 Not Found");
    }

    async getHtml() {
        return "<h1>404 NOT FOUND</h1>";
    }
}
export { AbstractView, Home, Connect, _404 };
