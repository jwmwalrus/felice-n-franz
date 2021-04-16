const bag = [];

const resetBag = () => {
    const e = document.getElementById('bag-list');
    e.innerHTML = '';
    bag.length = 0;
};

const showMessage = (m) => {
    // TODO: implement
};

export {
    resetBag,
    showMessage,
};
