import AutoComplete from '@tarekraafat/autocomplete.js';
import { DateTime } from 'luxon';
import * as _ from 'lodash/lodash.js';
import { v4 as uuidv4 } from 'uuid';

import { ERROR, showToast } from './toasts.js';
import {
    copyStringToClipboard,
    createParagraph,
    removeElement,
    createBtnSm,
} from './util.js';
import {
    getActiveLookup,
    getTopicName,
    setActiveLookup,
} from './env.js';
import { lookup, stopAllLookup } from './socket.js';
import { shoppingBagIcon } from './icons.js';

const REPLAY_FROM_BEGINNING = 'beginning';
const REPLAY_FROM_TIMESTAMP = 'timestamp';

const apiUrl = window.location.origin;

const bag = {
    activeMsgId: '',
    lookup: {
        env: '',
        replayType: '',
        searchIds: [],
    },
    messages: new Map(),
    toasts: new Map(),
};
let acTopic = null;

const updateBagBadges = () => {
    const badge = document.querySelector('#bag-btn span.badge');
    const msgBadge = document.querySelector('#pills-messages-tab span.badge');
    const toastBadge = document.getElementById('toasts-badge');

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
    bag.activeMsgId = id;
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

const addSearchResult = async (m) => {
    if (!bag.lookup.searchIds.includes(m.searchId)) {
        return;
    }

    const node = document.createElement('div');
    node.classList.add('d-flex');
    node.appendChild(
        await createBtnSm({
            icon: shoppingBagIcon,
            title: 'Add message to bag',
            onclick: async () => addMessageToBag(m),
        }),
    );
    node.appendChild(
        createParagraph([JSON.stringify(m)]),
    );

    document.getElementById('bag-lookup-list').appendChild(node);
};

const addToastToBag = (t) => {
    const id = uuidv4();
    bag.toasts.set(id, t);
    if (!t.canBeIgnored) {
        showToast(t);
    }

    const node = document.createElement('div');
    node.setAttribute('id', `toast-${id}`);
    node.classList.add('list-group-item', 'list-group-item-noclick');

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

const addToBag = async (x) => {
    if ('toastType' in x) {
        addToastToBag(x);
        return;
    }

    await addMessageToBag(x);
};

const clearBagMessages = (withBadges = false) => {
    const e = document.getElementById('bag-message-list');
    e.innerHTML = '';
    bag.messages.clear();
    document.getElementById('bag-message-form').reset();

    if (withBadges) {
        updateBagBadges();
    }
};

const clearBagToasts = (withBadges = false) => {
    const e = document.getElementById('bag-toast-list');
    e.innerHTML = '';
    bag.toasts.clear();

    if (withBadges) {
        updateBagBadges();
    }
};

const clearLookup = async () => {
    document.getElementById('bag-lookup-list').innerHTML = '';
};

const enableLookupGo = async () => {
    if (bag.lookup.env !== ''
        && bag.lookup.replayType !== ''
        && document.getElementById('bag-lookup-topic').value !== ''
        && document.getElementById('bag-lookup-offset').value !== ''
        && document.getElementById('bag-lookup-pattern').value !== ''
    ) {
        document.getElementById('bag-lookup-go-btn').removeAttribute('disabled');
        return;
    }

    document.getElementById('bag-lookup-go-btn').setAttribute('disabled', 'disabled');
};

const fireLookup = async () => {
    const searchId = uuidv4();
    bag.lookup.searchIds.push(searchId);

    const offset = document.getElementById('bag-lookup-offset').value;
    const parsed = Date.parse(offset);
    if (parsed === 0 || Number.isNaN(parsed)) {
        showToast({
            toastType: ERROR,
            title: 'Invalid DateTime',
            message: 'Unable to parse the provided offset as a DateTime',
        });
        return;
    }

    const payload = {
        type: bag.lookup.replayType,
        pattern: document.getElementById('bag-lookup-pattern').value,
        offset,
        searchId,
    };

    lookup(
        bag.lookup.env,
        document.getElementById('bag-lookup-topic').value,
        [JSON.stringify(payload)],
    );
};

const setAutoCompleteTopicForLookup = () => {
    acTopic = new AutoComplete({
        selector: '#bag-lookup-topic',
        placeHolder: 'Start typing and select...',
        data: {
            src: async () => getActiveLookup().topics ?? [],
            keys: ['value'],
        },
        threshold: 2,
        resultsList: {
            maxResults: 10,
            noResults: true,
        },
        resultItem: {
            highlight: true,
        },
    });

    document.querySelector('#bag-lookup-topic').addEventListener(
        'selection',
        (event) => {
            const { selection } = event.detail;
            document.getElementById('bag-lookup-topic').value = selection.value[selection.key];
        },
    );
};

const setLookupEnvironment = async (sel) => {
    bag.lookup.env = sel?.target.value ?? '';
    if (!bag.lookup.env) {
        setActiveLookup({});
        return;
    }
    try {
        const res = await fetch(`${apiUrl}/envs/${bag.lookup.env}`);
        const payload = await res.json();
        setActiveLookup(payload);
    } catch (e) {
        console.error(e);
    }
};

const setLookupType = async (sel) => {
    bag.lookup.replayType = sel?.target.value ?? '';
    switch (bag.lookup.replayType) {
        case REPLAY_FROM_BEGINNING:
            document.getElementById('bag-lookup-offset').value = '0';
            break;
        case REPLAY_FROM_TIMESTAMP:
            document.getElementById('bag-lookup-offset').value = DateTime.local().toISO();
            break;
        default:
            document.getElementById('bag-lookup-offset').value = '';
    }
    await enableLookupGo();
};

const resetLookup = async () => {
    document.getElementById('bag-lookup-form').reset();
    await setLookupEnvironment();
    await setLookupType();
};

const stopLookup = async () => {
    stopAllLookup(bag.lookup.searchIds);
    bag.lookup.searchIds = [];
};

const updateMessageSignature = () => {
    const label = document.getElementById('bag-label').value;

    const id = bag.activeMsgId !== '' ? bag.activeMsgId : document.querySelector('#bag-message-list .active').id.substr(8);

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
    addSearchResult,
    addToBag,
    clearBagMessages,
    clearBagToasts,
    enableLookupGo,
    fireLookup,
    resetLookup,
    stopLookup,
    clearLookup,
    setAutoCompleteTopicForLookup,
    setLookupEnvironment,
    setLookupType,
    updateMessageSignature,
};
