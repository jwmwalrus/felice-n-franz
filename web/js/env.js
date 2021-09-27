import * as _ from 'lodash/lodash.js';

const EXCEPT = ['application/json'];

let activeEnv = {};
let lookupEnv = {};

const getActiveEnv = () => activeEnv;

const setActiveEnv = (val) => {
    activeEnv = val;
};

const getActiveLookup = () => lookupEnv;

const setActiveLookup = (val) => {
    lookupEnv = val;
};

const getOutstandingHeader = (m) => {
    let oh = '';
    if (_.isEmpty(activeEnv)) {
        return oh;
    }

    const mh = m.headers ?? [];
    const th = activeEnv
        .topics
        .find((t) => t.value === m.topic)
        .headers ?? [];

    try {
        for (const i of th) {
            for (const j of mh) {
                if (!(EXCEPT.includes(i.value))
                   && i.key === j.key
                   && (i.value === j.value)
                ) {
                    oh = i.value.replace(new RegExp('^' + activeEnv.headerPrefix), '');
                    break;
                }
            }
        }
    } catch (e) {
        // pass
    }
    return oh;
};

const getTopicName = (m) => {
    const h = getOutstandingHeader(m);

    if (!h) {
        if (!_.isEmpty(activeEnv)) {
            const name = activeEnv.topics.find((t) => t.value === m.topic)?.name ?? '';
            if (name) {
                return name;
            }
        }
        if (!_.isEmpty(lookupEnv)) {
            const name = lookupEnv.topics.find((t) => t.value === m.topic)?.name ?? '';
            if (name) {
                return name;
            }
        }
    }

    return h;
};

export {
    getActiveEnv,
    setActiveEnv,
    getActiveLookup,
    setActiveLookup,
    getOutstandingHeader,
    getTopicName,
};
