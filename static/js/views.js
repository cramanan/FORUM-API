import { navigateTo, APIendpoint } from "./index.js";

class View {
  constructor() {}

  /**
   * Sets the document's title to the inputed string.
   *
   * @param {string} title - The title to set
   */

  setTitle(title) {
    document.title = title;
  }

  /**
   * Sets the document's css to the inputed string.
   * The document will change the <link> tag with id="viewcss"
   *
   * @param {string} css - The title to set
   */

  setCSS(css) {
    document.getElementById("viewcss").href = css;
  }

  /**
   * Returns HTML as a string
   *
   * @returns {string}
   */

  async getHtml() {
    return "";
  }

  /**
   * Executes JS after rendering HTML from getHTML()
   */

  async bindListeners() {}
}

class Connect extends View {
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
      const response = await fetch(`http://${APIendpoint}/register`, {
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
      const response = await fetch(`http://${APIendpoint}/login`, {
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

class Home extends View {
  constructor() {
    super();
    this.setTitle("Real-Time Forum");
    this.setCSS("/static/css/home.css");
  }

  async getHtml() {
    const html = `
    <header class="header">
        <h3><a href="/" id="main-title" data-link>REAL-TIME FORUM</a></h3>
        <button id="logout-button" title="Log Out"><img src="/static/images/logout.svg" width="34"/></button>
    </header>
    <aside id="sidebar" class="closed">
      <h2>Users :</h2>
      <ul id="users"></ul>
    </aside>
    <main>
      <form id="post-form">
          <label for="post-content">Create a P0ST</label>
          <textarea name="post-content" id="post-content"></textarea>
          <button type="submit">P0ST</button>
      </form>
      <div id="all-posts"></div>
    </main>
    <footer>
    </footer>`;
    return html;
  }

  async fetchUsers() {}

  async bindListeners() {
    const allposts = document.getElementById("all-posts");
    allposts?.append(...(await this.fetchPosts()));

    const postform = document.getElementById("post-form");
    postform?.addEventListener("submit", (event) => this.Post(event));

    // const logout = document.getElementById("logout-button");
    // logout?.addEventListener("click", (event) => this.Logout(event));

    const sidebar = document.getElementById("sidebar");
    const closeBtn = document.createElement("button");
    closeBtn.id = "closeBtn";
    closeBtn.textContent = "<";
    closeBtn.onclick = () => {
      sidebar.classList.toggle("closed");
      closeBtn.textContent = closeBtn.textContent === "<" ? ">" : "<";
    };
    sidebar.prepend(closeBtn);

    const conn = new WebSocket(`ws://${APIendpoint}/ws`);
    conn.onmessage = (msg) => console.log(JSON.parse(msg.data));
    const users = document.getElementById("users");
    console.log(users);
  }

  async Post(event) {
    event.preventDefault();
    try {
      const data = new FormData(event.target);
      console.log(data);
      const response = await fetch(`http://${APIendpoint}/post`, {
        method: "post",
        body: data,
        credentials: "include",
      });

      if (response.ok) {
        document.getElementById("post-content").value = "";
        const allposts = document.getElementById("all-posts");
        allposts.innerHTML = "";
        allposts.append(...(await this.fetchPosts()));
      }
    } catch (reason) {
      console.log(reason);
    }
  }

  async fetchPosts() {
    const postsHTML = [];
    try {
      const response = await fetch(`http://${APIendpoint}/getposts`);
      const datas = await response.json();
      datas.data.forEach((post) => {
        const div = document.createElement("div");
        div.className = "post";
        const h2 = document.createElement("h2");
        h2.textContent = post.Username;
        const p = document.createElement("p");
        p.textContent = post.Content;
        div.append(h2, p);
        postsHTML.push(div);
      });
    } catch (error) {
      console.log(error);
    }
    return postsHTML;
  }

  async Logout() {
    try {
      const response = await fetch(`http://${APIendpoint}/logout`, {
        credentials: "include",
      });
      navigateTo("/connect");
    } catch (reason) {
      console.log(reason);
    }
  }
}

class _404 extends View {
  constructor() {
    super();
    this.setTitle("404 Not Found");
  }

  async getHtml() {
    return "<h1>404 NOT FOUND</h1>";
  }
}

class Profile extends View {}

export { View, Home, Connect, _404, Profile };
