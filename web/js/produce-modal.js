import * as _ from 'lodash/lodash.js';
import AutoComplete from '@tarekraafat/autocomplete.js';
import { getActiveEnv } from './env.js';
import { produce, subscribe } from './socket.js';
import { ERROR, showToast } from './toasts.js';
import { addConsumerCard } from './cards.js';

import '@tarekraafat/autocomplete.js/dist/css/autoComplete.css';

let acTopic = null;

const resetProducer = () => document.getElementById('producer-form').reset();

const validateProducerPayload = () => {};

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

const setAutoCompleteTopic = () => {
    acTopic = new AutoComplete({
        selector: '#produce-topic',
        placeHolder: 'Start typing and select...',
        data: {
            src: async () => getActiveEnv().topics ?? [],
            key: ['value'],
        },
        threshold: 2,
        onSelection: (feedback) => {
            document.getElementById('produce-topic').value = feedback.selection.value[feedback.selection.key];
        },
    });
};

export {
    resetProducer,
    validateProducerPayload,
    produceMessage,
    setAutoCompleteTopic,
};
