import * as _ from 'lodash/lodash.js';
import {
    copyStringToClipboard,
    createBtnGroupSm,
    getActionId,
    getTopicId,
    removeElement,
} from './util.js';
import { getActiveEnv, getOutstandingHeader } from './env.js';
import { subscribe, unsubscribe } from './socket.js';
import { populateAvailable } from './consume-modal.js';
import { addToBag } from './bag-modal.js';
import {
    broomIcon,
    compressIcon,
    copyIcon,
    ellipsisHIcon,
    shoppingBagIcon,
} from './icons.js';

const tracker = new Map();

const FILTER_TYPE_TEXT = 0;
const FILTER_TYPE_KEY = 1;
const FILTER_TYPE_RE = 2;

const filter = {
    type: FILTER_TYPE_TEXT,
    ignoreCase: false,
    value: '',
};

const messageMatchesRegExp = (m) => m.value.match(new RegExp(filter.value)) !== null;

const messageMatchesKey = (m) => filter.value in JSON.parse(m.value);

const messageMatchesText = (m) => m.value.includes(filter.value);

const messageMatchesFilter = (m) => {
    switch (filter.type) {
        case FILTER_TYPE_RE:
            return messageMatchesRegExp(m);
        case FILTER_TYPE_KEY:
            return messageMatchesKey(m);
        default:
            return messageMatchesText(m);
    }
};

const applyFilter = () => {
    filter.value = document.getElementById('search-input').value.trim();

    if (filter.value === '') {
        resetFilter();
        return;
    }

    for (const k of tracker.keys()) {
        tracker.get(k).forEach((m) => {
            document.getElementById(`details-${getActionId(m)}`).classList.remove('is-visible');
            const matches = messageMatchesFilter(m);
            const e = document.getElementById(getActionId(m));
            if (matches) {
                e.classList.add('is-visible');
            } else {
                e.classList.remove('is-visible');
            }
        });
    }

    document.getElementById('search-btn').classList.replace('text-warning', 'text-danger');
    document.getElementById('reset-search-btn').style.display = 'block';
};

const resetFilter = () => {
    filter.key = '';
    filter.value = '';
    filter.and = false;

    for (const k of tracker.keys()) {
        tracker.get(k).forEach((m) => {
            const e = document.getElementById(getActionId(m));
            e.classList.add('is-visible');
        });
    }

    document.getElementById('search-input').value = '';
    document.getElementById('search-btn').classList.replace('text-danger', 'text-warning');
    document.getElementById('reset-search-btn').style.display = 'none';
};

const getListGroupElement = (topic) => {
    try {
        return document.querySelector(`#${getTopicId(topic)} .card-body .list-group`);
    } catch (e) {
        return null;
    }
};

const updateCardBadge = (topic) => {
    if (!tracker.has(topic)) {
        return;
    }
    const cardId = getTopicId(topic);
    const badge = document.querySelector(`#${cardId} .card-title span.badge`);
    badge.innerText = tracker.get(topic).length.toString();
};

const clearCard = (topic) => {
    if (!tracker.has(topic)) {
        return;
    }

    const e = getListGroupElement(topic);
    if (e !== null) {
        e.innerHTML = '';
    }

    tracker.set(topic, []);
    updateCardBadge(topic);
};

const addConsumerCard = async (t) => {
    if (tracker.has(t.value)) {
        return;
    }
    tracker.set(t.value, []);

    const parent = document.getElementById('consumer-cards');

    const title = document.createElement('h5');
    title.classList.add('card-title', 'mr-auto');
    title.innerHTML = `${t.name}&nbsp;&nbsp;<span class="badge badge-secondary">0</span>`;

    const subtitle = document.createElement('h6');
    subtitle.classList.add('card-subtitle');
    subtitle.innerText = t.value;

    const btnGroup = await createBtnGroupSm([
        {
            icon: copyIcon,
            title: 'Copy topic',
            onclick: () => copyStringToClipboard(t.value),
        },
        {
            icon: broomIcon,
            title: 'Clear messages',
            onclick: () => clearCard(t.value),
        },
    ]);

    const controls = document.createElement('div');
    controls.classList.add('card-controls');
    controls.appendChild(title);
    controls.appendChild(btnGroup);

    const header = document.createElement('div');
    header.classList.add('card-header');
    header.appendChild(controls);
    header.appendChild(subtitle);

    const list = document.createElement('div');
    list.classList.add('list-group');

    const body = document.createElement('div');
    body.classList.add('card-body');
    body.appendChild(list);

    const card = document.createElement('div');
    card.setAttribute('id', getTopicId(t.value));
    card.classList.add('card', 'bg-dark');
    card.appendChild(header);
    card.appendChild(body);

    parent.appendChild(card);
};

const createMessageNode = async (m) => {
    const truncate = (s) => _.truncate(s, { length: 40 });

    let small = `offset: ${m.offset}`;
    const h = getOutstandingHeader(m);
    if (h) {
        small += `<br>type: ${truncate(h)}`;
    } else {
        small += `<br>key: ${m.key ? truncate(m.key) : '0'}`;
    }
    small += `<br>ts: ${m.timestamp}`;
    const text = document.createElement('div');
    text.innerHTML = `<p>${small}</p>`;

    const btnGroup = await createBtnGroupSm([
        {
            icon: ellipsisHIcon,
            title: 'Toggle details',
            onclick: () => {
                const details = document.getElementById(`details-${getActionId(m)}`);
                details.classList.toggle('is-visible');
            },
        },
        {
            icon: shoppingBagIcon,
            title: 'Add message to bag',
            onclick: () => addToBag(m),
        },
    ]);

    const flex = document.createElement('div');
    flex.classList.add('d-flex');
    btnGroup.classList.add('ml-auto');
    flex.appendChild(text);
    flex.appendChild(btnGroup);

    const node = document.createElement('div');
    node.setAttribute('id', getActionId(m));
    node.classList.add('list-group-item', 'toggle-content', 'is-visible');
    node.appendChild(flex);

    if (!messageMatchesFilter(m)) {
        node.classList.remove('is-visible');
    }

    return node;
};

const createMessageDetails = async (m) => {
    const payload = JSON.stringify(JSON.parse(m.value), null, 2);

    const pBtnGroup = await createBtnGroupSm([
        {
            title: 'Copy payload',
            icon: copyIcon,
            onclick: () => copyStringToClipboard(payload),
        },
        {
            title: 'Copy compact payload',
            icon: compressIcon,
            onclick: () => copyStringToClipboard(m.value),
        },
    ]);
    pBtnGroup.classList.add('btn-group-vertical');

    const code = document.createElement('code');
    code.innerText = payload;

    const pre = document.createElement('pre');
    pre.classList.add('scrollable-pre');
    pre.appendChild(code);

    const td1 = document.createElement('td');
    td1.classList.add('scrollable-td');
    td1.appendChild(pre);

    const td2 = document.createElement('td');
    td2.appendChild(pBtnGroup);

    const tr = document.createElement('tr');
    tr.appendChild(td1);
    tr.appendChild(td2);

    const table = document.createElement('table');
    table.appendChild(tr);

    const form = document.createElement('form');
    form.appendChild(table);

    const body = document.createElement('div');
    body.classList.add('card-body');
    body.appendChild(form);

    const card = document.createElement('div');
    card.classList.add('card');
    card.appendChild(body);

    const details = document.createElement('div');
    details.setAttribute('id', `details-${getActionId(m)}`);
    details.classList.add('list-group-item', 'list-group-item-details', 'toggle-content');
    details.appendChild(card);

    return details;
};

const addMessageToCardList = async (m, l) => {
    if (!tracker.has(m.topic)) {
        return;
    }

    const node = await createMessageNode(m);
    const details = await createMessageDetails(m);

    l.appendChild(node);
    l.appendChild(details);

    const maxTailOffset = getActiveEnv();
    const msgs = tracker.get(m.topic);
    msgs.push(m);
    while (msgs.length > maxTailOffset) {
        const e = msgs.shift();
        await removeElement(getActionId(e));
        await removeElement(`details-${getActionId(e)}`);
    }
    tracker.set(m.topic, msgs);
    updateCardBadge(m.topic);
};

const clearAllCards = () => {
    const keys = Array.from(tracker.keys());
    keys.forEach((k) => clearCard(k));
};

const removeAllCards = () => {
    const parent = document.getElementById('consumer-cards');
    parent.innerHTML = '';

    const keys = Array.from(tracker.keys());
    unsubscribe(keys);
    tracker.clear();
};

const removeCard = async (topic) => {
    await removeElement(getTopicId(topic));
    unsubscribe([topic]);
    tracker.delete(topic);
};

const addSelectedCards = async () => {
    const parent = document.getElementById('selected-topics');
    let list = parent.querySelectorAll('div');
    list = Array.from(list);

    const { topics } = getActiveEnv();
    const topicList = [];
    list.forEach((l) => {
        const t = topics.find((t) => t.key === l.id);
        addConsumerCard(t);
        topicList.push(t);
    });

    subscribe(topicList.map((l) => l.key));

    const keys = Array.from(tracker.keys());
    keys.forEach((k) => {
        const v = topicList.find((l) => l.value === k);
        if (!v) {
            removeCard(k);
        }
    });
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
    updateCardBadge,
    applyFilter,
    resetFilter,
};
