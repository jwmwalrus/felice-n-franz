import * as _ from 'lodash/lodash.js';
import {
    copyStringToClipboard,
    createBtnGroupSm,
    createParagraph,
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

const filter = {
    ignoreCase: false,
    value: '',
};

const messageMatchesFilter = (m) => m.value.includes(filter.value);

const updateCardBadges = (topic) => {
    if (!tracker.has(topic)) {
        return;
    }
    const cardId = getTopicId(topic);
    const badge = document.querySelector(`#${cardId} .card-title > span.badge`);
    badge.innerText = tracker.get(topic).length.toString();

    if (filter.value) {
        document.querySelector(`#${cardId} .card-title span.filter-badge`).classList.add('active-filter');
        const badge = document.querySelector(`#${cardId} .card-title span.filter-badge > span.badge`);
        badge.innerText = document.querySelectorAll(`#${cardId} > div.card-body > div.list-group .list-group-item.list-group-item-node.is-visible`).length;
    } else {
        document.querySelector(`#${cardId} .card-title span.filter-badge`).classList.remove('active-filter');
    }
};

const clearFilter = () => {
    filter.value = '';

    for (const k of tracker.keys()) {
        tracker.get(k).forEach((m) => {
            const e = document.getElementById(getActionId(m));
            e.classList.add('is-visible');
        });
        updateCardBadges(k);
    }

    document.getElementById('filter-input').value = '';
    document.getElementById('filter-btn').classList.replace('text-danger', 'text-warning');
    document.getElementById('clear-filter-btn').style.display = 'none';
};

const applyFilter = () => {
    filter.value = document.getElementById('filter-input').value.trim();

    if (filter.value === '') {
        clearFilter();
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
        updateCardBadges(k);
    }

    document.getElementById('filter-btn').classList.replace('text-warning', 'text-danger');
    document.getElementById('clear-filter-btn').style.display = 'block';
};

const getListGroupElement = (topic) => {
    try {
        return document.querySelector(`#${getTopicId(topic)} > div.card-body > div.list-group`);
    } catch (e) {
        return null;
    }
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
    updateCardBadges(topic);
};

const addConsumerCard = async (t) => {
    if (tracker.has(t.value)) {
        return;
    }
    tracker.set(t.value, []);

    const parent = document.getElementById('consumer-cards');

    const filterBadge = document.createElement('span');
    filterBadge.classList.add('badge', 'badge-danger');
    filterBadge.appendChild(
        document.createTextNode('0'),
    );

    const filterSpan = document.createElement('span');
    filterSpan.classList.add('filter-badge', 'toggle-content');
    filterSpan.appendChild(filterBadge);
    filterSpan.appendChild(
        document.createTextNode('\u002F'), // Solidus, HTML code &#47;
    );

    const badge = document.createElement('span');
    badge.classList.add('badge', 'badge-secondary');
    badge.appendChild(
        document.createTextNode('0'),
    );

    const title = document.createElement('h5');
    title.classList.add('card-title', 'mr-auto');
    title.appendChild(
        document.createTextNode(`${t.name}`),
    );
    title.appendChild(
        document.createElement('br'),
    );
    title.appendChild(filterSpan);
    title.appendChild(badge);

    const subtitle = document.createElement('h6');
    subtitle.classList.add('card-subtitle');
    subtitle.appendChild(
        document.createTextNode(t.value),
    );

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

    const text = document.createElement('div');

    const lines = [];
    lines.push(`offset: ${m.offset}`);

    const h = getOutstandingHeader(m);
    if (h) {
        lines.push(`type: ${truncate(h)}`);
    } else {
        lines.push(`key: ${m.key ? truncate(m.key) : '0'}`);
    }

    lines.push(`ts: ${m.timestamp}`);

    text.appendChild(
        createParagraph(lines),
    );

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
    node.classList.add('list-group-item', 'list-group-item-node', 'toggle-content', 'is-visible');
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

    const msgs = tracker.get(m.topic);
    msgs.push(m);
    tracker.set(m.topic, msgs);
    updateCardBadges(m.topic);
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
    updateCardBadges,
    applyFilter,
    clearFilter,
};
