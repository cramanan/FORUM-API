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
        <input type="email" id="login-email" name="email" />
        <label for="login-password">Password</label>
        <input type="password" id="login-password" name="password" />
        <button type="submit">Login</button>
      </form>
      <span id="sep"></span>
      <form id="register-form">
        <h1>Register</h1>
        <div id="register-server-error"></div>
        <label for="register-email">Email</label>
        <input type="email" id="register-email" name="email" />
        <label for="name">Username</label>
        <input type="text" id="name" name="name" />
        <label for="password">Password</label>
        <input type="password" id="register-password" name="password" />
        <label for="gender">Gender:</label>
        <select name="gender" id="gender">
          <option value="M">M</option>
          <option value="F">F</option>
          <option value="O">Other</option>
        </select>
        <label for="age">Age</label>
        <input type="number" name="age" id="age" />
        <label for="firstName">First Name</label>
        <input type="text" name="firstName" id="firstName" />
        <label for="lastName">Last Name</label>
        <input type="text" name="lastName" id="lastName" />
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
      const data = Object.fromEntries(new FormData(event.target).entries());
      console.log(data);
      const response = await fetch(`http://${APIendpoint}/register`, {
        credentials: "include",
        method: "post",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      });
      if (response.ok) navigateTo("/");
      else {
        document.getElementById("register-server-error").textContent =
          await response.json();
      }
    } catch (reason) {
      console.log(reason);
    }
  }

  async HandleLoginSubmit(event) {
    event.preventDefault();
    try {
      const data = Object.fromEntries(new FormData(event.target).entries());
      const response = await fetch(`http://${APIendpoint}/login`, {
        credentials: "include",
        method: "post",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      });
      if (response.ok) navigateTo("/");
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
    return ` <header class="header">
        <h3><a href="/" id="main-title" data-link>REAL-TIME FORUM</a></h3>
        <div id="user-box">
          <a href="/profile" data-link>
            <div id="username"></div>
            <div id="id"></div>
          </a>
          <button id="logout-button" title="Log Out">
            <img src="/static/images/logout.svg" width="34" />
          </button>
        </div>
      </header>
      <main>
        <form id="post-form">
          <label for="post-content">Create a P0ST</label>
          <textarea name="content" id="post-content"></textarea>
          <button type="submit">P0ST</button>
        </form>
        <div id="all-posts"></div>
      </main>
      <footer></footer>`;
  }

  async bindListeners() {
    const username = document.getElementById("username");
    username.textContent = await fetch(`http://${APIendpoint}/`)
      .then((resp) => resp.json())
      .then((v) => v.name);

    const allposts = document.getElementById("all-posts");
    allposts?.append(...(await this.fetchPosts()));

    const postform = document.getElementById("post-form");
    postform?.addEventListener("submit", (event) => this.Post(event));

    const logout = document.getElementById("logout-button");
    logout?.addEventListener("click", (event) => this.Logout(event));
  }

  async Post(event) {
    event.preventDefault();
    try {
      const data = Object.fromEntries(new FormData(event.target).entries());
      const response = await fetch(`http://${APIendpoint}/post`, {
        credentials: "include",
        method: "post",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
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
      datas.forEach((post) => {
        // create elements to avoid injections
        const div = document.createElement("div");
        div.className = "post";
        const h2 = document.createElement("h2");
        h2.textContent = post.username;
        const p = document.createElement("p");
        p.textContent = post.content;
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
      console.log(response);
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

class Profile extends View {
  constructor() {
    super();
  }

  async getHtml() {
    const info = await fetch(`http://${APIendpoint}/`).then((resp) =>
      resp.json()
    );
    return `
      <h1>Profile</h1>
      <div>${Object.entries(info)
        .map(([k, v]) => {
          const div = document.createElement("div");
          div.textContent = `${k} : ${v}`;
          return div.outerHTML;
        })
        .join("")}</div>
    `;
  }
}

class Chat extends View {
  constructor() {
    super();
    this.setTitle("Chat");
  }

  async getHtml() {
    return `
      <header>
        <h1>Chat</h1>
      </header>
      <main>
        <ul id="users">Flemme lol.</ul>
      </main>
    `;
  }
}

export { View, Home, Connect, _404, Profile, Chat };
