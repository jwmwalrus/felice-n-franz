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
    clearFilter,
} from './cards';
import {
    resetProducer,
    validateProducerPayload,
    produceMessage,
    setAutoCompleteTopic,
} from './produce-modal.js';
import { resetBag } from './bag-modal.js';
import { pauseIcon, playIcon } from './icons.js';

import '../css/index.css';

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
        e.innerHTML = playIcon;
    } else {
        e.classList.remove('pause', 'play');
        e.innerHTML = pauseIcon;
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

    document.getElementById('filter-btn').onclick = applyFilter;
    document.getElementById('clear-filter-btn').onclick = clearFilter;


    window.addEventListener('keydown',(e) => {
        if(e.keyIdentifier == 'U+000A'
            || e.keyIdentifier=='Enter'
            || e.keyCode == 13) {
            if(e.target.nodeName == 'INPUT' && e.target.type == 'text') {
                e.preventDefault();
                return false;
            }
        }
    }, true);

    setAutoCompleteTopic();

    loadSocket({
        open: () => console.info('Web socket started!'),
        close: () => console.info('Web socket closed!'),
        error: () => console.error('Web socket error!'),
        message: async (event) => {
            const messages = event.data.split('\n').map((s) => JSON.parse(s));
            const pp = document.getElementById('playpause-btn');
            if (!pp.classList.contains('play')) {
                return;
            }
            for await (const m of messages) {
                if ('toastType' in m) {
                    showToast(m);
                } else if ('topic' in m) {
                    const l = getListGroupElement(m.topic);
                    console.log(l);
                    if (l !== null) {
                        await addMessageToCardList(m, l);
                    }
                }
            }
        },
    });
};

export default apiUrl;
