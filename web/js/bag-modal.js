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
    toats: new Map(),
};

const updateBagBadges = () => {
    const badge = document.querySelector('#bag-btn span.badge');
    let text = bag.messages.length.toString();
    if (bag.toasts.filter((t) => t.toastType === ERROR && !t.canBeIgnored).length > 0) {
        text += ' (*)';
    }
    badge.innerText = text;
};

const removeFromBag = (id) => {
    if (bag.messages.contains(id)) {
        bag.messages.remove(id);
    }

    removeElement(`bag-msg-${id}`);
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
            createParagraph(`${m.label}`),
        );
        return text;
    }

    const name = getTopicName(m);
    if (name) {
        text.appendChild(
            createParagraph(`${name}`),
        );
    }
    text.appendChild(
        createParagraph(`offset: ${m.offset}`),
    );
    text.appendChild(
        createParagraph(`ts: ${m.timestamp}`),
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
    if (bag.messages.length === 1) {
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
    node.classList.add('list-group-item', 'list-group-item-action');
    if (bag.toasts.length === 1) {
        node.classList.add('active');
    }

    const truncate = (s) => _.truncate(s, { length: 30 });

    const text = document.createElement('div');
    text.appendChild(
        createParagraph(`${t.title}`),
    );
    text.appendChild(
        createParagraph(`type: ${t.type}`),
    );
    text.appendChild(
        createParagraph(`message: ${t.message}`),
    );
    if (t.topic !== '') {
        text.appendChild(
            createParagraph(`topic: ${truncate(t.topic)}`),
        );
    }
    if (t.partition) {
        text.appendChild(
            createParagraph(`partition: ${t.partition}`),
        );
    }
    if (t.offset) {
        text.appendChild(
            createParagraph(`offset: ${t.offset}`),
        );
    }
    text.appendChild(
        createParagraph(`sent at: ${t.sentAt}`),
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
    bag.messages.length = 0;
    document.getElementById('bag-message-form').reset();

    if (withBadges) {
        updateBagBadges();
    }
};
const clearBagToasts = (withBadges = false) => {
    const e = document.getElementById('bag-toast-list');
    e.innerHTML = '';
    bag.toasts.length = 0;

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
    const label = document.getElementById('bag-label');

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
