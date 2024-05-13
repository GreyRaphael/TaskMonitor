import { Greet } from '../wailsjs/go/main/App';

window.greet = function () {
    Greet()
        .then((url) => {
            window.location.replace(url);
        })
        .catch((err) => {
            console.error(err);
        });
}

greet();