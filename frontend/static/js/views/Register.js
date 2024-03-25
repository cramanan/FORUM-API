import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
    constructor() {
        super().setTitle("Register");
    }

    async getHtml() {
        const html = await fetch("/static/js/templates/register.html")
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
