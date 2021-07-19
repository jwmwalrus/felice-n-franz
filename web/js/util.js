const copyStringToClipboard = (str) => {
    navigator
        .clipboard
        .writeText(str)
        .catch(console.error);
};

const createBtnSm = async ({
    id,
    icon,
    classList,
    title,
    onclick,
}) => {
    const cl = classList && classList.length > 0 ? classList : ['btn', 'btn-sm', 'bg-transparent', 'text-warning'];

    const btn = document.createElement('a');
    btn.setAttribute('title', title);
    if (id) {
        btn.setAttribute('id', id);
    }
    for await (const c of cl) {
        btn.classList.add(c);
    }
    btn.onclick = onclick;
    btn.innerHTML = icon;

    return btn;
};

const createBtnGroupSm = async (list = []) => {
    const btnGroup = document.createElement('div');
    btnGroup.classList.add('btn-group-sm');

    for await (const l of list) {
        const btn = await createBtnSm({
            id: l.id,
            icon: l.icon,
            title: l.title,
            classList: l.classList,
            onclick: l.onclick,
        });
        btnGroup.appendChild(btn);
    }

    return btnGroup;
};

const createParagraph = (lines) => {
    const p = document.createElement('p');
    for (let i = 0; i < lines.length; i += 1) {
        if (i > 0) {
            p.appendChild(
                document.createElement('br'),
            );
        }
        p.appendChild(
            document.createTextNode(lines[i]),
        );
    }
    return p;
};

const copyToClipboard = (id) => {
    try {
        const e = document.getElementById(id);
        copyStringToClipboard(e.value);
    } catch (e) {
        console.error(e);
    }
};

const getTopicId = (t) => t.replace(/[.#]/g, '-');

const getActionId = (a) => `${getTopicId(a.topic)}-${a.partition}-${a.offset}-${a.key ? a.key : '0'}`;

const removeElement = (id) => {
    const elem = document.getElementById(id);
    if (!elem) {
        return;
    }
    elem.parentNode?.removeChild(elem);
};

const toggleCompactBtn = (elemId) => {
    const elem = document.getElementById(elemId);

    if (!elem.value) {
        return;
    }

    if (elem.value.includes('\n')) {
        elem.value = JSON.stringify(JSON.parse(elem.value));
    } else {
        elem.value = JSON.stringify(JSON.parse(elem.value), null, 2);
    }
};

export {
    copyStringToClipboard,
    copyToClipboard,
    createBtnSm,
    createBtnGroupSm,
    createParagraph,
    getTopicId,
    getActionId,
    removeElement,
    toggleCompactBtn,
};
