import {
    copyStringToClipboard,
    getActionId,
    getTopicId,
    removeElement,
} from './util.js';
import { getActiveEnv } from './env.js';
import { subscribe, unsubscribe } from './socket.js';
import { populateAvailable } from './consume-modal.js';
import { addToBag, addMessageToList, showMessage } from './bag-modal.js';

const tracker = new Map();

const getListGroupElement = (topic) => document.querySelector(`#${getTopicId(topic)} .card-body .list-group`);

const clearCard = (topic) => {
    if (!tracker.has(topic)) {
        return;
    }

    const e = getListGroupElement(topic);
    if (e !== null) {
        e.innerHTML = '';
    }

    tracker.set(topic, []);
};

const addConsumerCard = (t) => {
    if (tracker.has(t.value)) {
        return;
    }
    tracker.set(t.value, []);

    const parent = document.getElementById('consumer-cards');

    const title = document.createElement('H5');
    title.classList.add('card-title');
    title.innerText = t.value;

    const icon1 = document.createElement('SPAN');
    icon1.classList.add('icon-docs');
    const btn1 = document.createElement('A');
    btn1.setAttribute('title', 'Copy topic');
    btn1.classList.add('btn', 'btn-sm', 'bg-transparent', 'text-warning');
    btn1.onclick = () => copyStringToClipboard(t.value);
    btn1.appendChild(icon1);

    const icon2 = document.createElement('SPAN');
    icon2.classList.add('icon-reload');
    const btn2 = document.createElement('A');
    btn2.setAttribute('title', 'Clear messages');
    btn2.classList.add('btn', 'btn-sm', 'bg-transparent', 'text-warning');
    btn2.onclick = () => clearCard(t.value);
    btn2.appendChild(icon2);
    const btnGroup = document.createElement('div');
    btnGroup.classList.add('btn-group-sm');
    btnGroup.appendChild(btn1);
    btnGroup.appendChild(btn2);

    const controls = document.createElement('DIV');
    controls.classList.add('card-controls');
    controls.appendChild(btnGroup);

    const header = document.createElement('DIV');
    header.classList.add('card-header');
    header.appendChild(controls);
    header.appendChild(title);

    const list = document.createElement('DIV');
    list.classList.add('list-group');

    const body = document.createElement('DIV');
    body.classList.add('card-body');
    body.appendChild(list);

    const card = document.createElement('DIV');
    card.setAttribute('id', getTopicId(t.value));
    card.classList.add('card', 'bg-dark');
    card.appendChild(header);
    card.appendChild(body);

    parent.appendChild(card);
};

const addMessageToCardList = async (m, l) => {
    if (!tracker.has(m.topic)) {
        return;
    }

    await addMessageToList(
        m,
        l,
        getActionId(m),
        () => showMessage(m),
        () => addToBag(m),
    );

    const msgs = tracker.get(m.topic);
    msgs.push(m);
    while (msgs.length > 100) {
        const e = msgs.shift();
        await removeElement(getActionId(e));
    }
    tracker.set(m.topic, msgs);
};

const clearAllCards = () => {
    const keys = Array.of(tracker.keys());
    keys.forEach((k) => clearCard(k));
};

const removeAllCards = () => {
    const parent = document.getElementById('consumer-cards');
    parent.innerHTML = '';

    const keys = Array.of(tracker.keys());
    unsubscribe(keys);
    tracker.clear();
};

const removeCard = async (topic) => {
    await removeElement(getTopicId(topic));
    unsubscribe([topic]);
    tracker.delete(topic);
};

const addSelectedCards = async () => {
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

const resetConsumers = () => {
    removeAllCards();
    populateAvailable();
};

export {
    addSelectedCards,
    addConsumerCard,
    addMessageToCardList,
    clearAllCards,
    getListGroupElement,
    removeAllCards,
    removeCard,
    resetConsumers,
};
