const copyStringToClipboard = (str) => {};

const copyToClipboard = (elem) => {};

const getTopicId = (t) => t.replace(/[.#]/g, '-');

const getActionId = (a) => `${getTopicId(a.topic)}-${a.partition}-${a.offset}-${a.key ? a.key : '0'}`;

const removeElement = (id) => {
    const elem = document.getElementById(id);
    elem.parentNode?.removeChild(elem);
};

const toggleCompactButton = (btnId, elemId) => {
    const elem = document.getElementById(elemId);
    const icon = document.querySelector(`#${btnId} span`);

    if (icon.classList.contains('icon-crop')) {
        elem.value = JSON.stringify(JSON.parse(elem.value));
        icon.classList.replace('icon-crop', 'icon-layers');
    } else {
        elem.value = JSON.stringify(JSON.parse(elem.value), null, 2);
        icon.classList.replace('icon-layers', 'icon-crop');
    }
};

export {
    copyStringToClipboard,
    copyToClipboard,
    getTopicId,
    getActionId,
    removeElement,
    toggleCompactButton,
};
