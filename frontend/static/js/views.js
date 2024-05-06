import { APIendpoint, navigateTo } from "./index.js";

class AbstractView {
    constructor() {}

    setTitle(title) {
        document.title = title;
    }

    async getHtml() {
        return "";
    }

    setCSS() {}
}

class Connect extends AbstractView {
    constructor() {
        super().setTitle("Connect");
        window.HandleLoginSubmit = this.HandleLoginSubmit;
        window.HandleRegisterSubmit = this.HandleRegisterSubmit;
    }

    async getHtml() {
        return `<div id="connect">
    <form onsubmit="HandleLoginSubmit(event)">
        <h1>Login</h1>
        <div id="login-server-error"></div>
        <label for="login-email">Email</label>
        <input type="email" id="login-email" name="login-email" />
        <label for="login-password">Password</label>
        <input type="password" id="login-password" name="login-password" />
        <button type="submit">Login</button>
    </form>
    <span id="sep"></span>
    <form onsubmit="HandleRegisterSubmit(event)">
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

    setCSS() {
        document.querySelector("#viewcss").href = "/static/css/connect.css";
    }

    HandleRegisterSubmit(event) {
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
    }

    HandleLoginSubmit(event) {
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
    }
}

class _404 extends AbstractView {
    constructor() {
        super().setTitle("404 Not Found");
    }

    async getHtml() {
        return "<h1>404 NOT FOUND</h1>";
    }
}

class Home extends AbstractView {
    constructor() {
        super().setTitle("Real-Time Forum");
        window.Post = this.Post;
    }

    async getHtml() {
        try {
            const response = await fetch(`${APIendpoint}/getposts`);
            const datas = await response.json();
            let postsHTML = "";
            datas.data.forEach((post) => {
                postsHTML += `<div class="post"><h2>${post.Username}</h2><p>${post.Content}</p></div>`;
            });
            const html = `<nav class="header">
                <h3><a href="/" id="main-title">REAL-TIME FORUM</a></h3>
            </nav>
            <main>
            <form id="post-form" onsubmit="Post(event)">
                <label for="post-content">Create a P0ST</label>
                <textarea name="post-content" id="post-content"></textarea>
                <button type="submit">P0ST</button>
            </form>
            <div id="all-posts">
                ${postsHTML}
            </div>
        </main>
        <footer>

        </footer>`;
            return html;
        } catch (error) {
            console.log(error);
        }
        return "<h1>00PS... Something went wrong ://</h1>";
    }

    setCSS() {
        document.querySelector("#viewcss").href = "/static/css/home.css";
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

            console.log(response.ok);
        } catch (reason) {
            console.log(reason);
        }
    }
}

export { Home, Connect, _404 };
