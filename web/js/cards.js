import * as _ from 'lodash/lodash.js';
import {
    copyStringToClipboard,
    copyToClipboard,
    createBtnSm,
    createBtnGroupSm,
    getActionId,
    getTopicId,
    removeElement,
    toggleCompactBtn,
} from './util.js';
import { getActiveEnv, getOutstandingHeader } from './env.js';
import { subscribe, unsubscribe } from './socket.js';
import { populateAvailable } from './consume-modal.js';
import { addToBag } from './bag-modal.js';

const tracker = new Map();

const filter = {
    key: '',
    and: false,
    value: '',
};

const messageMatchesFilter = (m) => {
    let matches = true;

    const parsed = JSON.parse(m.value);
    const valueRe = new RegExp(`\b${filter.value}\b`);
    if (
        (filter.key && !filter.and && !filter.value)
        || (filter.key && filter.and && !filter.value)
    ) {
        matches = filter.key in parsed;
    } else if (
        (filter.key && !filter.and && filter.value)
    ) {
        matches = filter.key in parsed || m.value.match(valueRe) !== null;
    } else if (
        (filter.key && filter.and && filter.value)
    ) {
        matches = filter.key in parsed && m.value.match(valueRe) !== null;
    } else if (
        (!filter.key && !filter.and && filter.value)
        || (!filter.key && filter.and && filter.value)
    ) {
        matches = m.value.match(valueRe) !== null;
    }

    return matches;
};

const applyFilter = () => {
    filter.key = document.getElementById('filter-key').value;
    filter.value = document.getElementById('filter-value').value;
    filter.and = document.getElementById('filter-and').checked;

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

    document.getElementById('filter-btn').classList.replace('text-warning', 'text-danger');
};

const resetFilter = () => {
    filter.key = '';
    filter.value = '';
    filter.and = false;
    document.getElementById('filter-form').reset();

    for (const k of tracker.keys()) {
        tracker.get(k).forEach((m) => {
            const e = document.getElementById(getActionId(m));
            e.classList.add('is-visible');
        });
    }

    document.getElementById('filter-btn').classList.replace('text-danger', 'text-warning');
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
            iconClass: 'icon-docs',
            title: 'Copy topic',
            onclick: () => copyStringToClipboard(t.value),
        },
        {
            iconClass: 'icon-reload',
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
            iconClass: 'icon-options',
            title: 'Toggle details',
            onclick: () => {
                const details = document.getElementById(`details-${getActionId(m)}`);
                details.classList.toggle('is-visible');
            },
        },
        {
            iconClass: 'icon-bag',
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
    const eText = document.createElement('span');
    eText.classList.add('input-group-text');
    eText.innerHTML = '<small><b>Envelope:</b></small>';

    const envelopeId = `details-${getActionId(m)}-envelope`;
    const eTextArea = document.createElement('textarea');
    eTextArea.setAttribute('rows', 3);
    eTextArea.setAttribute('id', envelopeId);
    eTextArea.readonly = true;
    eTextArea.classList.add('form-control');

    const {
        topic,
        partition,
        key,
        offset,
        timestamp,
        timestampType,
        headers,
    } = m;
    eTextArea.value = JSON.stringify({
        topic,
        partition,
        key,
        offset,
        timestamp,
        timestampType,
        headers,
    }, null, 2);

    const eBtn = await createBtnSm({
        title: 'Copy envelope',
        iconClass: 'icon-docs',
        classList: ['btn', 'btn-sm', 'bg-transparent', 'text-warning'],
        onclick: () => copyToClipboard(envelopeId),
    });

    const eBox = document.createElement('div');
    eBox.classList.add('input-group', 'mb-3');
    eBox.appendChild(eText);
    eBox.appendChild(eTextArea);
    eBox.appendChild(eBtn);

    const pText = document.createElement('span');
    pText.classList.add('input-group-text');
    pText.innerHTML = '<small></b>Payload:</b></small>';

    const payloadId = `details-${getActionId(m)}-payload`;
    const pTextArea = document.createElement('textarea');
    pTextArea.setAttribute('id', payloadId);
    pTextArea.setAttribute('rows', 5);
    pTextArea.readonly = true;
    pTextArea.classList.add('form-control');
    pTextArea.value = JSON.stringify(JSON.parse(m.value), null, 2);

    const pBtnGroup = await createBtnGroupSm([
        {
            title: 'Toggle compact JSOM',
            iconClass: 'icon-crop',
            onclick: () => toggleCompactBtn(payloadId),
        },
        {
            title: 'Copy payload',
            iconClass: 'icon-docs',
            onclick: () => copyToClipboard(payloadId),
        },
    ]);

    const pBox = document.createElement('div');
    pBox.classList.add('input-group', 'mb-3');
    pBox.appendChild(pText);
    pBox.appendChild(pTextArea);
    pBox.appendChild(pBtnGroup);

    const form = document.createElement('form');
    form.appendChild(eBox);
    form.appendChild(pBox);

    const body = document.createElement('div');
    body.classList.add('card-body');
    body.appendChild(form);

    const card = document.createElement('div');
    card.classList.add('card', 'bg-dark');
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
