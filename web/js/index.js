import '../css/index.css';
import { getActiveEnv, setActiveEnv } from './env.js';
import { loadSocket, subscribe } from './socket.js';
import { populateAvailable, clearTopicsAndGroups } from './consume-modal.js';
import { INFO, ERROR, showToast } from './toasts.js';
import {
    addConsumerCard,
    addMessageToCardList,
    clearAllCards,
    getListGroupElement,
    removeAllCards,
} from './cards';
import { resetBag } from './bag-modeless.js';

const apiUrl = window.location.origin;

window.onload = () => {
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

window.addSelectedCards = async () => {
    removeAllCards();
    const parent = document.getElementById('selected-topics');
    let list = parent.querySelectorAll('div');
    list = Array.from(list);
    const { topics } = getActiveEnv();
    list.forEach((l) => {
        const t = topics.find((t) => t.key === l.id);
        addConsumerCard(t);
    });

    const topicKeys = list.map((l) => l.id);
    subscribe(topicKeys);
};

window.checkIfValidEnvironment = async (sel) => {
    const { value } = sel;
    showToast({ toastType: INFO, message: 'HELLO HELLO' });
    showToast({ toastType: ERROR, message: 'HELLO HELLO' });

    if (value !== '') {
        try {
            const res = await fetch(`${apiUrl}/envs/${value}`);
            const active = await res.json();
            setActiveEnv(active);

            populateAvailable();
        } catch (e) {
            console.error(e);
        }
        // document.getElementById("produce-btn").disabled = false;
        document.getElementById('clear-contents-btn').disabled = false;
    } else {
        // document.getElementById("produce-btn").disabled = true;
        document.getElementById('clear-contents-btn').disabled = true;

        clearTopicsAndGroups();
        setActiveEnv({});
    }
};

window.clearAllContents = clearAllCards;

window.produceMessage = () => {};

window.resetConsumers = () => {
    removeAllCards();
    populateAvailable();
};

window.resetProducer = () => {};

window.resetBag = resetBag;

window.togglePlayPause = () => {
    const e = document.getElementById('playpause-btn');
    if (e.classList.contains('play')) {
        e.classList.remove('play');
        e.classList.add('pause');
        e.innerHTML = '<div class="icon-control-play"></div>';
    } else {
        e.classList.remove('pause');
        e.classList.add('play');
        e.innerHTML = '<div class="icon-control-pause"></div>';
    }
};

window.validateProducerMessage = () => {};

export default apiUrl;
