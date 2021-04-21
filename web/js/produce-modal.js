import { getActiveEnv } from './env.js';
import { produce, subscribe } from './socket.js';
import { addConsumerCard } from './cards.js';

const resetProducer = () => {
    document.getElementById('producer-form').reset();
    const icon = document.querySelector('#produce-toggle-compact-payload-btn span');
    icon.classList.remove('icon-layers');
    icon.classList.add('icon-crop');
};

const validateProducerPayload = () => {};

const produceMessage = () => {
    const topic = document.getElementById('produce-topic').value;
    const key = document.getElementById('produce-key').value;
    const headers = document.getElementById('produce-headers').value;
    const payload = document.getElementById('produce-payload').value;

    const { topics } = getActiveEnv();
    const t = topics.find((t) => t.value === topic);
    addConsumerCard(t);
    subscribe([t.key]);
    produce(topic, payload, key, headers);
};

export {
    resetProducer,
    validateProducerPayload,
    produceMessage,
};
