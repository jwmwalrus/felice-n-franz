import * as _ from 'lodash/lodash.js';
import AutoComplete from '@tarekraafat/autocomplete.js';
import { getActiveEnv } from './env.js';
import { produce, subscribe } from './socket.js';
import { ERROR, showToast } from './toasts.js';
import { addConsumerCard } from './cards.js';

import '@tarekraafat/autocomplete.js/dist/css/autoComplete.02.css';

let acTopic = null;

const resetProducer = () => document.getElementById('producer-form').reset();

const validateProducerPayload = () => {};

const addHeaderToList = (h) => {
    let { value } = document.getElementById('produce-headers');
    if (!value) {
        value = '[]';
    }

    const headers = JSON.parse(value);

    for (const x of headers) {
        if (x.key === h.key && x.value === h.value) {
            return;
        }
    }

    headers.push(h);
    const out = JSON.stringify(headers, null, 2);
    document.getElementById('produce-headers').value = out;
};

const addHeader = () => {
    const e = document.getElementById('produce-predef-headers');
    const { value } = e.options[e.selectedIndex];

    if (!value) {
        return;
    }

    addHeaderToList(JSON.parse(value));
};

const produceMessage = () => {
    const topic = document.getElementById('produce-topic').value;
    const payload = document.getElementById('produce-payload').value;
    if (!topic || !payload) {
        showToast({ toastType: ERROR, message: 'Topic and payload cannot be empty' });
        return;
    }

    let headers = document.getElementById('produce-headers').value;
    if (!headers) {
        headers = '[]';
    }

    const key = document.getElementById('produce-key').value;

    const { topics } = getActiveEnv();
    const t = topics.find((t) => t.value === topic);
    if (_.isEmpty(t)) {
        showToast({ toastType: ERROR, message: `Topic doest not belong to current environment: ${topic}` });
        return;
    }

    addConsumerCard(t);
    subscribe([t.key]);
    produce(topic, payload, key, headers);
};

const removeHeaderFromList = (h) => {
    const { value } = document.getElementById('produce-headers');
    if (!value) {
        return;
    }

    const list = JSON.parse(value).filter((x) => x.key !== h.key || x.value !== h.value);

    const out = JSON.stringify(list, null, 2);
    document.getElementById('produce-headers').value = out;
};

const removeHeader = () => {
    const e = document.getElementById('produce-predef-headers');
    const { value } = e.options[e.selectedIndex];

    if (!value) {
        return;
    }

    removeHeaderFromList(JSON.parse(value));
};

const setPredefinedHeaders = () => {
    const sel = document.getElementById('produce-predef-headers');
    sel.innerHTML = '<option value="" selected>Predefined headers...</option>';
    document.getElementById('produce-headers').value = '[]';

    const topic = document.getElementById('produce-topic').value;
    const headers = getActiveEnv().topics?.find((t) => t.value === topic)?.headers;

    if (headers) {
        headers.forEach((h) => {
            const o = document.createElement('option');
            o.value = JSON.stringify(h);
            o.innerText = `${h.key} -> ${h.value}`;
            sel.appendChild(o);
            if (h.key === 'Content-Type' && h.value === 'application/json') {
                addHeaderToList(h);
            }
        });
    }
};

const setAutoCompleteTopicForProducer = () => {
    acTopic = new AutoComplete({
        selector: '#produce-topic',
        placeHolder: 'Start typing and select...',
        data: {
            src: async () => getActiveEnv().topics ?? [],
            keys: ['value'],
        },
        threshold: 2,
        resultsList: {
            maxResults: 10,
            noResults: true,
        },
        resultItem: {
            highlight: true,
        },
    });

    document.querySelector('#produce-topic').addEventListener(
        'selection',
        (event) => {
            const { selection } = event.detail;
            document.getElementById('produce-topic').value = selection.value[selection.key];
        },
    );

    document.querySelector('#produce-topic').addEventListener(
        'close',
        () => setPredefinedHeaders(),
    );
};

export {
    addHeader,
    resetProducer,
    validateProducerPayload,
    produceMessage,
    removeHeader,
    setAutoCompleteTopicForProducer,
};
