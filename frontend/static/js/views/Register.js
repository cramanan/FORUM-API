import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
    constructor() {
        super().setTitle("Register");
    }

    async getHtml() {
        const html = await fetch("/static/html/register.html")
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
        document.querySelector('link[rel="stylesheet"]').href =
            "/static/css/register.css";
    }
}
