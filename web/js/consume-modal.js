import { removeElement } from './util.js';
import { getActiveEnv } from './env.js';

const addGroupToList = (g, l, cb) => {
    const node = document.createElement('div');
    node.setAttribute('id', g.id);
    if (g.description) {
        node.setAttribute('title', g.description);
    }
    node.classList.add('list-group-item', 'list-group-item-action');
    node.ondblclick = cb;

    const textnode = document.createTextNode(g.name);
    node.appendChild(textnode);

    l.appendChild(node);
};

const addTopicToList = (t, l, cb) => {
    const node = document.createElement('div');
    node.setAttribute('id', t.key);
    if (t.description) {
        node.setAttribute('title', t.description);
    }
    node.classList.add('list-group-item', 'list-group-item-action');
    node.ondblclick = cb;

    const textnode = document.createTextNode(t.value);
    node.appendChild(textnode);

    l.appendChild(node);
};

const clearTopicsAndGroups = () => {
    document.getElementById('available-groups').innerHTML = '';
    document.getElementById('available-topics').innerHTML = '';
    document.getElementById('selected-groups').innerHTML = '';
    document.getElementById('selected-topics').innerHTML = '';
};

const populateAvailable = async () => {
    clearTopicsAndGroups();

    const { groups, topics } = getActiveEnv();

    const groupsList = document.getElementById('available-groups');
    for await (const g of groups) {
        addGroupToList(g, groupsList, selectGroup);
    }

    const topicsList = document.getElementById('available-topics');
    for await (const t of topics) {
        addTopicToList(t, topicsList, selectTopic);
    }
};

const selectGroup = async (event) => {
    const id = event.target.id;
    await removeElement(id);

    const { groups } = getActiveEnv();

    const groupsList = document.getElementById('selected-groups');
    const g = groups.find((g) => g.id === id);
    await addGroupToList(g, groupsList, unselectGroup);

    for await (const k of g.keys) {
        selectTopic({ target: { id: k } });
    }
};

const selectTopic = async (event) => {
    const id = event.target.id;
    await removeElement(id);

    const { topics } = getActiveEnv();

    const list = document.getElementById('selected-topics');
    const t = topics.find((t) => t.key === id);
    await addTopicToList((t), list, unselectTopic);
};

const unselectGroup = async (event) => {
    const id = event.target.id;
    await removeElement(id);

    const { groups } = getActiveEnv();

    const groupsList = document.getElementById('available-groups');
    const g = groups.find((g) => g.id === id);
    await addGroupToList(g, groupsList, selectGroup);

    for await (const k of g.keys) {
        unselectTopic({ target: { id: k } });
    }
};

const unselectTopic = async (event) => {
    const id = event.target.id;
    await removeElement(id);

    const { topics } = getActiveEnv();

    const list = document.getElementById('available-topics');
    const t = topics.find((t) => t.key === id);
    await addTopicToList(t, list, selectTopic);
};

export {
    addGroupToList,
    addTopicToList,
    clearTopicsAndGroups,
    populateAvailable,
    selectGroup,
    selectTopic,
    unselectGroup,
    unselectTopic,
};
