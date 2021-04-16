import { copyToClipboard, getTopicId, removeElement } from './util.js';
import { unsubscribe } from './socket.js';
import { showMessage } from './bag-modeless.js';

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
    btn1.classList.add('btn', 'btn-secondary', 'text-warning');
    btn1.onclick = () => copyToClipboard(t.value);
    btn1.appendChild(icon1);

    const icon2 = document.createElement('SPAN');
    icon2.classList.add('icon-refresh');
    const btn2 = document.createElement('A');
    btn2.classList.add('btn', 'btn-secondary', 'text-warning');
    btn2.onclick = () => clearCard(t.value);
    btn2.appendChild(icon2);
    const btnGroup = document.createElement('div');
    btnGroup.classList.add('btn-group');
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

const addMessageToCardList = (m, l) => {
    if (!tracker.has(m.topic)) {
        return;
    }
    const actionId = (a) => `${getTopicId(a.topic)}-${a.partition}-${m.offset}-${m.key !== '' ? m.key : '0'}`;

    const node = document.createElement('DIV');
    node.setAttribute('id', actionId(m));
    node.classList.add('list-group-item', 'list-group-item-action');
    node.ondblclick = () => showMessage(m);

    const textnode = document.createTextNode(`{ partition: ${m.partition}, offset: ${m.offset}, key: ${m.key} }`);
    node.appendChild(textnode);

    l.appendChild(node);
    const msgs = tracker.get(m.topic);
    msgs.push(m);
    while (msgs.length > 100) {
        const e = msgs.shift();
        removeElement(actionId(e));
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

const removeCard = (topic) => {
    removeElement(getTopicId(topic));
    unsubscribe([topic]);
    tracker.delete(topic);
};

export {
    addConsumerCard,
    addMessageToCardList,
    clearAllCards,
    getListGroupElement,
    removeAllCards,
    removeCard,
};
