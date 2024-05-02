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
        const html = await fetch("/static/html/connect.html")
            .then((resp) => {
                if (resp.ok) {
                    return resp.text();
                }
                throw new Error("Resource not found");
            })
            .then((txt) => {
                return txt;
            })
            .catch((reason) => console.log(reason));

        return html;
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
        const html = await fetch("/static/html/404.html")
            .then((resp) => {
                if (resp.ok) {
                    return resp.text();
                }
                throw new Error("Resource not found");
            })
            .then((txt) => {
                return txt;
            })
            .catch((reason) => console.log(reason));

        return html;
    }
}

class Home extends AbstractView {
    constructor() {
        super().setTitle("Real-Time Forum");
    }

    async getHtml() {
        const html = await fetch("/static/html/index.html")
            .then((resp) => {
                if (resp.ok) {
                    return resp.text();
                }
                throw new Error("Resource not found");
            })
            .then((txt) => {
                return txt;
            })
            .catch((reason) => console.log(reason));

        return html;
    }

    setCSS() {
        document.querySelector("#viewcss").href = "/static/css/home.css";
    }
}

export { Home, Connect, _404 };
