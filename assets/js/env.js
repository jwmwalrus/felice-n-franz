let activeEnv = {};

const getActiveEnv = () => activeEnv;

const setActiveEnv = (val) => {
    activeEnv = val;
};

export {
    getActiveEnv,
    setActiveEnv,
};
