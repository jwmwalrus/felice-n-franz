import '../css/index.css';
import { copyToClipboard, toggleCompactButton } from './util.js';
import { setActiveEnv } from './env.js';
import { showToast } from './toasts.js';
import { loadSocket } from './socket.js';
import {
    populateAvailable,
    clearTopicsAndGroups,
} from './consume-modal.js';
import {
    addSelectedCards,
    addMessageToCardList,
    clearAllCards,
    getListGroupElement,
    removeAllCards,
} from './cards';
import {
    resetProducer,
    validateProducerPayload,
    produceMessage,
} from './produce-modal.js';
import { resetBag } from './bag-modal.js';

const apiUrl = window.location.origin;

const checkIfValidEnvironment = async (sel) => {
    const { value } = sel.target;
    // showToast({ toastType: INFO, message: 'HELLO HELLO' });
    // showToast({ toastType: ERROR, message: 'HELLO HELLO' });

    if (value) {
        try {
            const res = await fetch(`${apiUrl}/envs/${value}`);
            const active = await res.json();
            setActiveEnv(active);

            populateAvailable();
        } catch (e) {
            console.error(e);
        }
        document.getElementById('clear-contents-btn').disabled = false;
    } else {
        document.getElementById('clear-contents-btn').disabled = true;

        clearTopicsAndGroups();
        setActiveEnv({});
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
        e.innerHTML = '<div class="icon-control-play"></div>';
    } else {
        e.classList.remove('pause', 'play');
        e.innerHTML = '<div class="icon-control-pause"></div>';
    }
};

window.onload = () => {
    document.getElementById('envs-sel').onchange = checkIfValidEnvironment;
    document.getElementById('playpause-btn').onclick = togglePlayPause;
    document.getElementById('clear-contents-btn').onclick = clearAllCards;

    document.getElementById('reset-consumers-btn').onclick = resetConsumers;
    document.getElementById('add-selected-cards-btn').onclick = addSelectedCards;

    document.getElementById('reset-producer-btn').onclick = resetProducer;
    document.getElementById('validate-producer-payload-btn').onclick = validateProducerPayload;
    document.getElementById('produce-message-btn').onclick = produceMessage;

    document.getElementById('bag-topic-btn').onclick = () => copyToClipboard('bag-topic');
    document.getElementById('bag-partition-btn').onclick = () => copyToClipboard('bag-partition');
    document.getElementById('bag-key-btn').onclick = () => copyToClipboard('bag-key');
    document.getElementById('bag-offset-btn').onclick = () => copyToClipboard('bag-offset');
    document.getElementById('bag-timestamp-btn').onclick = () => copyToClipboard('bag-timestamp');
    document.getElementById('bag-headers-btn').onclick = () => copyToClipboard('bag-headers');
    document.getElementById('bag-headers-toggle-compact-btn').onclick = () => toggleCompactButton('bag-headers-toggle-compact-btn', 'bag-headers');
    document.getElementById('bag-payload-btn').onclick = () => copyToClipboard('bag-payload');
    document.getElementById('bag-payload-toggle-compact-btn').onclick = () => toggleCompactButton('bag-payload-toggle-compact-btn', 'bag-payload');
    document.getElementById('reset-bag-btn').onclick = resetBag;

    document.getElementById('showmsg-envelope-btn').onclick = () => copyToClipboard('showmsg-envelope-btn');
    document.getElementById('showmsg-payload-btn').onclick = () => copyToClipboard('showmsg-payload');
    document.getElementById('showmsg-payload-toggle-compact-btn').onclick = () => toggleCompactButton('showmsg-payload-toggle-compact-btn', 'showmsg-payload');

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
