import Toastify from 'toastify-js/src/toastify-es.js';
import 'toastify-js/src/toastify.css';
import * as _ from 'lodash/lodash.js';

const INFO = 'info';
const WARNING = 'warning';
const ERROR = 'error';

const getText = (t) => {
    let { toastType } = t;
    if (![INFO, WARNING, ERROR].includes(toastType)) {
        toastType = INFO;
    }
    return `<h5>${t.title ?? _.startCase(toastType)}</h5><br<hr><p>${t.message}</p>`;
};

const showToast = (t) => {
    const text = getText(t);
    switch (t.toastType) {
        case ERROR:
            Toastify({
                text,
                escapeMarkup: false,
                duration: -1,
                close: true,
                // style: {
                //     background: 'linear-gradient(to right, #7c69ef, #007bff)',
                // },
            }).showToast();
            break;
        case WARNING:
            Toastify({
                text,
                escapeMarkup: false,
                duration: 10000,
                close: true,
                // style: {
                //     background: 'linear-gradient(to right, #303030, #505050)',
                // },
            }).showToast();
            break;
        default:
            // INFO
            Toastify({
                text,
                escapeMarkup: false,
                duration: 5000,
                close: true,
                // style: {
                //     background: 'linear-gradient(to right, #303030, #505050)',
                // },
            }).showToast();
            break;
    }
};

export {
    ERROR,
    INFO,
    WARNING,
    showToast,
};
