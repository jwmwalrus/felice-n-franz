import * as _ from 'lodash/lodash.js';
import {
    copyStringToClipboard,
    getActionId,
    removeElement,
} from './util.js';
import { getTopicName } from './env.js';

const bag = [];

const updateBagBadge = () => {
    const badge = document.querySelector('#bag-btn span.badge');
    badge.innerText = bag.length.toString();
};

const removeFromBag = (m) => {
    const id = getActionId(m);
    let i = 0;
    for (;;) {
        if (i >= bag.length) {
            break;
        }
        if (id === getActionId(bag[i])) {
            bag.splice(i, 1);
            break;
        }
        i += 1;
    }

    removeElement(`bag-${id}`);
    updateBagBadge();
};

const showBagMessage = (m) => {
    document.getElementById('bag-copy-raw-btn').onclick = () => copyStringToClipboard(JSON.stringify(m));
    document.getElementById('bag-remove-msg-btn').onclick = () => removeFromBag(m);

    const {
        topic,
        partition,
        key,
        offset,
        timestamp,
        timestampType,
        timestampTypeString,
    } = m;

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

const addToBag = async (m) => {
    bag.push(m);

    const node = document.createElement('div');
    node.setAttribute('id', `bag-${getActionId(m)}`);
    node.classList.add('list-group-item', 'list-group-item-action');
    if (bag.length === 1) {
        node.classList.add('selected');
    }
    node.onclick = () => {
        showBagMessage(m);
        try {
            document.querySelector('#bag-list .active')?.classList.remove('active');
        } catch (e) {
            // pass
        }
        document.getElementById(`bag-${getActionId(m)}`).classList.add('active');
    };

    const truncate = (s) => _.truncate(s, { length: 30 });

    const name = getTopicName(m);
    let inner = '';
    if (name) {
        inner += `${name}<br>`;
    }
    inner += `offset: ${m.offset}<br>`;
    if (m.key) {
        inner += `key: ${truncate(m.key)}<br>`;
    }
    inner += `ts: ${m.timestamp}`;
    const child = document.createElement('p');
    child.innerHTML = `${inner}`;

    node.appendChild(child);

    document.getElementById('bag-list').appendChild(node);
    updateBagBadge();
};

const resetBag = () => {
    const e = document.getElementById('bag-list');
    e.innerHTML = '';
    bag.length = 0;
    document.getElementById('bag-form').reset();
    updateBagBadge();
};

export {
    addToBag,
    resetBag,
};
