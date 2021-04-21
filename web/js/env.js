import * as _ from 'lodash/lodash.js';

const EXCEPT = ['application/json'];

let activeEnv = {};

const getActiveEnv = () => activeEnv;

const setActiveEnv = (val) => {
    activeEnv = val;
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
        .headers ?? {};

    try {
        for (const i of th) {
            for (const j in mh) {
                if (!(EXCEPT.includes(i.value))
                   && i.key === j
                   && (i.value === mh[j]
                       || activeEnv.headerPrefix + i.value === mh[j])
                ) {
                    oh = i.value;
                    break;
                }
            }
        }
    } catch (e) {
        // pass
    }
    return oh;
};

export {
    getActiveEnv,
    setActiveEnv,
    getOutstandingHeader,
};
