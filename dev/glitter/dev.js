import './dev.css'
import 'htmx.org';
import Alpine from 'alpinejs'
window.Alpine = Alpine
Alpine.start()

document.addEventListener("htmx:configRequest", (e) => {
    const token = document.cookie
        .split("; ")
        .find(r => r.startsWith("csrf_="))
        ?.split("=")[1];
    if (token) e.detail.headers["X-Csrf-Token"] = token;
});
