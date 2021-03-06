import { getActiveEnv, getActiveLookup } from './env.js';

const socketUrl = 'ws://' + document.location.host + '/ws';
let conn;

const loadSocket = ({
    open = null,
    message = null,
    error = null,
    close = null,
}) => {
    conn = new WebSocket(socketUrl);
    if (open) {
        conn.onopen = open;
    }
    if (message) {
        conn.onmessage = message;
    }
    if (error) {
        conn.onerror = error;
    }
    if (close) {
        conn.onclose = close;
    }
};

const lookup = (env, topic, payload) => {
    const msg = {
        type: 'lookup',
        env,
        topic,
        payload,
    };
    conn.send(JSON.stringify(msg));
};

const produce = (topic, payload, key = '', headers = []) => {
    const { name } = getActiveEnv();
    const msg = {
        type: 'produce',
        env: name,
        topic,
        key,
        headers: JSON.parse(headers),
        payload: [JSON.stringify(JSON.parse(payload))],
    };
    conn.send(JSON.stringify(msg));
};

const refreshSubscriptions = () => {
    const msg = {
        type: 'refresh',
    };
    conn.send(JSON.stringify(msg));
};

const resetSubscriptions = () => {
    const msg = {
        type: 'reset',
    };
    conn.send(JSON.stringify(msg));
};

const stopAllLookup = (searchIds) => {
    const { name } = getActiveLookup();
    const msg = {
        type: 'stopLookup',
        env: name,

        payload: searchIds,
    };
    conn.send(JSON.stringify(msg));
};

const subscribe = (topicKeys) => {
    const { name } = getActiveEnv();
    const msg = {
        type: 'consume',
        env: name,
        payload: topicKeys,
    };
    conn.send(JSON.stringify(msg));
};

const unsubscribe = (topics) => {
    const { name } = getActiveEnv();
    const msg = {
        type: 'unsubscribe',
        env: name,

        payload: topics,
    };
    conn.send(JSON.stringify(msg));
};

export {
    loadSocket,
    lookup,
    produce,
    refreshSubscriptions,
    resetSubscriptions,
    stopAllLookup,
    subscribe,
    unsubscribe,
};
