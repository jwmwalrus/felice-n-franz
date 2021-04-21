import {
    copyStringToClipboard,
    getActionId,
    removeElement,
} from './util.js';

import { getOutstandingHeader } from './env.js';

const bag = [];

const removeFromBag = (m) => removeElement(getActionId(m));

const addMessageToList = (m, l, id, cb = null, dbl = null) => {
    const node = document.createElement('DIV');
    node.setAttribute('id', id);
    node.classList.add('list-group-item', 'list-group-item-action');

    if (cb) {
        node.onclick = cb;
    }
    if (dbl) {
        node.ondblclick = dbl;
    }

    const left = document.createElement('DIV');
    left.innerHTML = `<small>offset: ${m.offset}</small>`;
    const middle = document.createElement('DIV');
    const h = getOutstandingHeader(m);
    if (h) {
        middle.innerHTML = `<small>type: ${h}</small>`;
    } else {
        middle.innerHTML = `</small>key: ${m.key ? m.key : '0'}</small>`;
    }
    const right = document.createElement('SMALL');
    right.innerHTML = `</small>ts: ${m.timestamp}</small>`;

    const text = document.createElement('DIV');
    node.classList.add('d-flex');
    text.appendChild(left);
    text.appendChild(middle);
    text.appendChild(right);

    node.appendChild(text);

    l.appendChild(node);
};

const showBagMessage = (m) => {
    document.getElementById('bag-copy-raw-btn=b').onclick = () => copyStringToClipboard(JSON.stringify(m));
    document.getElementById('bag-remove-msg-btn').onclick = () => removeFromBag(m);

    const {
        topic,
        partition,
        key,
        offset,
        timestamp,
        timestampType,
    } = m;
    document.getElementById('bag-topic').value = topic;
    document.getElementById('bag-partition').value = partition;
    document.getElementById('bag-key').value = key;
    document.getElementById('bag-offset').value = offset;
    document.getElementById('bag-timestamp').value = timestamp;
    document.getElementById('bag-timestamptype').value = timestampType;

    const headers = JSON.stringify(JSON.parse(m.headers), null, 2);
    document.getelementbyid('bag-headers').value = headers;

    const payload = JSON.stringify(JSON.parse(m.payload), null, 2);
    document.getElementById('bag-payload').value = payload;
};

const addToBag = async (m) => {
    bag.push(m);
    await addMessageToList(
        m,
        document.getElementById('bag-list'),
        `bag-${getActionId(m)}`,
        () => showBagMessage(m),
    );
};

const resetBag = () => {
    const e = document.getElementById('bag-list');
    e.innerHTML = '';
    // bag.length = 0;
};

const showMessage = (m) => {
    document.getElementById('showmsg-raw-copy-btn').onclick = () => copyStringToClipboard(JSON.stringify(m));
    document.getElementById('showmsg-add-to-bag-btn').onclick = async () => addToBag(m);

    const {
        topic,
        partition,
        key,
        offset,
        timestamp,
        timestampType,
        headers,
    } = m;
    const envelope = JSON.stringify({
        topic,
        partition,
        key,
        offset,
        timestamp,
        timestampType,
        headers,
    }, null, 2);
    document.getElementById('showmsg-envelope').value = envelope;

    const payload = JSON.stringify(JSON.parse(m.value), null, 2);
    document.getElementById('showmsg-payload').value = payload;

    document.getElementById('showmsg-btn').click();
};

export {
    addToBag,
    addMessageToList,
    resetBag,
    showMessage,
};
