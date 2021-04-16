const copyToClipboard = (s) => {};

const getTopicId = (t) => t.replace(/[.#]/g, '-');

const removeElement = (id) => {
    const elem = document.getElementById(id);
    elem.parentNode?.removeChild(elem);
};

export {
    copyToClipboard,
    getTopicId,
    removeElement,
};
