import * as _ from 'lodash/lodash.js';
import { v4 as uuidv4 } from 'uuid';
import { ERROR, showToast } from './toasts.js';
import {
    copyStringToClipboard,
    createParagraph,
    removeElement,
} from './util.js';
import { getTopicName } from './env.js';

let activeMsgId = '';
const bag = {
    messages: new Map(),
    toasts: new Map(),
};

const updateBagBadges = () => {
    const badge = document.querySelector('#bag-btn span.badge');
    const msgBadge = document.querySelector('#v-pills-messages-tab span.badge');
    const toastBadge = document.querySelector('#v-pills-toasts-tab span.badge');

    let text = bag.messages.size.toString();
    if (Array.from(bag.toasts.values())
        .filter((t) => t.toastType === ERROR && !t.canBeIgnored)
        .length > 0) {
        text += ' (*)';
    }
    badge.innerText = text;

    msgBadge.innerText = bag.messages.size.toString();
    toastBadge.innerText = bag.toasts.size.toString();
};

const removeFromBag = (id) => {
    if (bag.messages.has(id)) {
        bag.messages.delete(id);
    }

    removeElement(`bag-msg-${id}`);
    document.getElementById('bag-message-form').reset();
    updateBagBadges();
};

const showBagMessage = (id) => {
    activeMsgId = id;
    const m = bag.messages.get(id);

    document.getElementById('bag-copy-raw-btn').onclick = () => copyStringToClipboard(JSON.stringify(m));
    document.getElementById('bag-remove-msg-btn').onclick = () => removeFromBag(id);

    const {
        label,
        topic,
        partition,
        key,
        offset,
        timestamp,
        timestampType,
        timestampTypeString,
    } = m;

    document.getElementById('bag-label').value = label;
    document.getElementById('bag-topic').value = topic;
    document.getElementById('bag-partition').value = partition;
    document.getElementById('bag-key').value = key;
    document.getElementById('bag-offset').value = offset;
    document.getElementById('bag-timestamp').value = timestamp;
    document.getElementById('bag-timestamp-help-btn').setAttribute('title', `type: ${timestampType} (${timestampTypeString})`);

    const headers = JSON.stringify(m.headers, null, 2);
    document.getElementById('bag-headers').value = headers;

    const payload = JSON.stringify(JSON.parse(m.value), null, 2);
    document.getElementById('bag-payload').value = payload;
};

const createMessageSignature = (m) => {
    const text = document.createElement('div');

    if (m.label !== '') {
        text.appendChild(
            createParagraph([`${m.label}`]),
        );
        return text;
    }

    const lines = [];
    const name = getTopicName(m);
    if (name) {
        lines.push(`${name}`);
    }
    lines.push(`offset: ${m.offset}`);
    lines.push(`ts: ${m.timestamp}`);

    text.appendChild(
        createParagraph(lines),
    );

    return text;
};

const addMessageToBag = async (m) => {
    const id = uuidv4();
    const newMsg = { label: '', ...m };
    bag.messages.set(id, newMsg);

    const node = document.createElement('div');
    node.setAttribute('id', `bag-msg-${id}`);
    node.classList.add('list-group-item', 'list-group-item-action');
    if (bag.messages.size === 1) {
        node.classList.add('active');
    }
    node.onclick = () => {
        showBagMessage(id);
        try {
            document.querySelector('#bag-message-list .active')?.classList.remove('active');
        } catch (e) {
            // pass
        }
        document.getElementById(`bag-msg-${id}`).classList.add('active');
    };

    const text = createMessageSignature(newMsg);
    node.appendChild(text);

    document.getElementById('bag-message-list').appendChild(node);
    updateBagBadges();
};

const addToastToBag = (t) => {
    const id = uuidv4();
    bag.toasts.set(id, t);
    if (!t.canBeIgnored) {
        showToast(t);
    }

    const node = document.createElement('div');
    node.setAttribute('id', `toast-${id}`);
    node.classList.add('list-group-item');

    const truncate = (s) => _.truncate(s, { length: 50 });

    const text = document.createElement('div');

    const lines = [];
    lines.push(`${t.title}`);
    lines.push('---');
    lines.push(`type: ${t.toastType}`);
    lines.push(`message: ${t.message}`);
    if (t.topic !== '') {
        lines.push(`topic: ${truncate(t.topic)}`);
    }
    if (t.partition) {
        lines.push(`partition: ${t.partition}`);
    }
    if (t.offset) {
        lines.push(`offset: ${t.offset}`);
    }
    const sentAt = new Date(t.sentAt * 1000);
    lines.push(`sent at: ${sentAt.toLocaleString()}`);

    text.appendChild(
        createParagraph(lines),
    );

    node.appendChild(text);

    document.getElementById('bag-toast-list').appendChild(node);
    updateBagBadges();
};

const addToBag = (x) => {
    if ('toastType' in x) {
        addToastToBag(x);
        return;
    }

    addMessageToBag(x);
};

const clearBagMessages = (withBadges = false) => {
    const e = document.getElementById('bag-message-list');
    e.innerHTML = '';
    bag.messages.size = 0;
    document.getElementById('bag-message-form').reset();

    if (withBadges) {
        updateBagBadges();
    }
};
const clearBagToasts = (withBadges = false) => {
    const e = document.getElementById('bag-toast-list');
    e.innerHTML = '';
    bag.toasts.size = 0;

    if (withBadges) {
        updateBagBadges();
    }
};

const resetBag = () => {
    clearBagMessages(false);
    clearBagToasts(false);
    updateBagBadges();
};

const updateMessageSignature = () => {
    const label = document.getElementById('bag-label').value;

    const id = activeMsgId !== '' ? activeMsgId : document.querySelector('#bag-message-list .active').id.substr(8);

    if (!id) {
        return;
    }

    const node = document.getElementById(`bag-msg-${id}`);

    const m = bag.messages.get(id);
    m.label = label;
    bag.messages.set(id, m);
    const text = createMessageSignature(m);

    node.innerHTML = '';
    node.appendChild(text);
};

export {
    addToBag,
    clearBagMessages,
    clearBagToasts,
    resetBag,
    updateMessageSignature,
};
