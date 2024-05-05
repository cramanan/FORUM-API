import { APIendpoint } from "./index.js";

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
        super().setTitle("Login");
    }

    async getHtml() {
        try {
            const html = await fetch("/static/html/connect.html");
            return await html.text();
        } catch (reason) {
            console.log(reason);
        }
        return "<h1>00PS... Something went wrong ://</h1>";
    }

    setCSS() {
        document.querySelector("#viewcss").href = "/static/css/connect.css";
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
    }

    async getHtml() {
        try {
            const response = await fetch(`${APIendpoint}/getposts`);
            const datas = await response.json();
            let postsHTML = "";
            datas.data.posts.forEach((post) => {
                postsHTML += `<a class="post" href="/post/${post.UserID}"><h2>${post.Username}</h2><p>${post.Content}</p></a>`;
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
}

export { Home, Connect, _404 };
