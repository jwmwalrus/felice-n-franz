import '../css/index.css';
import { copyToClipboard, toggleCompactBtn } from './util.js';
import { setActiveEnv } from './env.js';
import { showToast } from './toasts.js';
import { loadSocket } from './socket.js';
import {
    populateAvailable,
    clearTopicsAndGroups,
    filterGroups,
} from './consume-modal.js';
import {
    addSelectedCards,
    addMessageToCardList,
    clearAllCards,
    getListGroupElement,
    removeAllCards,
    applyFilter,
    resetFilter,
} from './cards';
import {
    resetProducer,
    validateProducerPayload,
    produceMessage,
    setAutoCompleteTopic,
} from './produce-modal.js';
import { resetBag } from './bag-modal.js';

const apiUrl = window.location.origin;

const resetEnvironment = () => {
    removeAllCards();
    clearTopicsAndGroups();
    resetProducer();
};

const checkIfValidEnvironment = async (sel) => {
    const { value } = sel.target;

    if (value) {
        try {
            const res = await fetch(`${apiUrl}/envs/${value}`);
            const active = await res.json();
            resetEnvironment();
            setActiveEnv(active);

            populateAvailable();
            document.getElementById('produce-message-btn').disabled = false;
        } catch (e) {
            console.error(e);
        }
        document.getElementById('clear-contents-btn').disabled = false;
    } else {
        resetEnvironment();
        document.getElementById('clear-contents-btn').disabled = true;

        setActiveEnv({});
        document.getElementById('produce-message-btn').disabled = true;
    }
};

const resetConsumers = () => {
    removeAllCards();
    populateAvailable();
};

const togglePlayPause = () => {
    const e = document.getElementById('playpause-btn');
    if (e.classList.contains('play')) {
        e.classList.replace('play', 'pause');
        const i = document.querySelector('#playpause-btn div');
        i.classList.replace('icon-control-pause', 'icon-control-play');
    } else {
        e.classList.remove('pause', 'play');
        const i = document.querySelector('#playpause-btn div');
        i.classList.replace('icon-control-play', 'icon-control-pause');
    }
};

window.onload = () => {
    document.getElementById('envs-sel').onchange = checkIfValidEnvironment;
    document.getElementById('playpause-btn').onclick = togglePlayPause;
    document.getElementById('clear-contents-btn').onclick = clearAllCards;

    document.getElementById('group-categories').onchange = filterGroups;
    document.getElementById('reset-consumers-btn').onclick = resetConsumers;
    document.getElementById('add-selected-cards-btn').onclick = addSelectedCards;

    document.getElementById('produce-payload-toggle-compact-btn').onclick = () => toggleCompactBtn('produce-payload');
    document.getElementById('reset-producer-btn').onclick = resetProducer;
    document.getElementById('validate-producer-payload-btn').onclick = validateProducerPayload;
    document.getElementById('produce-message-btn').onclick = produceMessage;

    document.getElementById('bag-copy-raw-btn').onclick = () => {};
    document.getElementById('bag-remove-msg-btn').onclick = () => {};

    document.getElementById('bag-topic-btn').onclick = () => copyToClipboard('bag-topic');
    document.getElementById('bag-partition-btn').onclick = () => copyToClipboard('bag-partition');
    document.getElementById('bag-key-btn').onclick = () => copyToClipboard('bag-key');
    document.getElementById('bag-offset-btn').onclick = () => copyToClipboard('bag-offset');
    document.getElementById('bag-timestamp-btn').onclick = () => copyToClipboard('bag-timestamp');
    document.getElementById('bag-headers-btn').onclick = () => copyToClipboard('bag-headers');
    document.getElementById('bag-headers-toggle-compact-btn').onclick = () => toggleCompactBtn('bag-headers');
    document.getElementById('bag-payload-btn').onclick = () => copyToClipboard('bag-payload');
    document.getElementById('bag-payload-toggle-compact-btn').onclick = () => toggleCompactBtn('bag-payload');
    document.getElementById('reset-bag-btn').onclick = resetBag;

    document.getElementById('reset-filter-btn').onclick = resetFilter;
    document.getElementById('apply-filter-btn').onclick = applyFilter;

    setAutoCompleteTopic();

    loadSocket({
        open: () => console.info('Web socket started!'),
        close: () => console.info('Web socket closed!'),
        error: () => console.info('Web socket error!'),
        message: async (event) => {
            const messages = event.data.split('\n').map((s) => JSON.parse(s));
            const pp = document.getElementById('playpause-btn');
            if (!pp.classList.contains('play')) {
                return;
            }
            for await (const m of messages) {
                if ('toastType' in m) {
                    if (getListGroupElement(m.topic) !== null) {
                        showToast(m);
                    }
                } else if ('topic' in m) {
                    const l = getListGroupElement(m.topic);
                    if (l !== null) {
                        await addMessageToCardList(m, l);
                    }
                }
            }
        },
    });
};

export default apiUrl;
