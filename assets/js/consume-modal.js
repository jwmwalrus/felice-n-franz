import { getActiveEnv } from './env.js';
import { removeElement } from './util.js';

const addGroupToList = (g, l, cb) => {
    const node = document.createElement('DIV');
    node.setAttribute('id', g.id);
    node.classList.add('list-group-item', 'list-group-item-action');
    node.ondblclick = cb;

    const textnode = document.createTextNode(g.description);
    node.appendChild(textnode);

    l.appendChild(node);
};

const addTopicToList = (t, l, cb) => {
    const node = document.createElement('DIV');
    node.setAttribute('id', t.key);
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

const populateAvailable = () => {
    clearTopicsAndGroups();

    const { groups, topics } = getActiveEnv();

    const groupsList = document.getElementById('available-groups');
    groups.forEach((g) => addGroupToList(g, groupsList, selectGroup));

    const topicsList = document.getElementById('available-topics');
    topics.forEach((t) => addTopicToList(t, topicsList, selectTopic));
};

const selectGroup = (event) => {
    const id = event.target.id;
    removeElement(id);

    const { groups } = getActiveEnv();

    const groupsList = document.getElementById('selected-groups');
    const g = groups.find((g) => g.id === id);
    addGroupToList(g, groupsList, unselectGroup);

    g.keys.forEach((k) => {
        selectTopic({ target: { id: k } });
    });
};

const selectTopic = (event) => {
    const id = event.target.id;
    removeElement(id);

    const { topics } = getActiveEnv();

    const list = document.getElementById('selected-topics');
    const t = topics.find((t) => t.key === id);
    addTopicToList((t), list, unselectTopic);
};

const unselectGroup = (event) => {
    const id = event.target.id;
    removeElement(id);

    const { groups } = getActiveEnv();

    const groupsList = document.getElementById('available-groups');
    const g = groups.find((g) => g.id === id);
    addGroupToList(g, groupsList, selectGroup);

    g.keys.forEach((k) => {
        unselectTopic({ target: { id: k } });
    });
};

const unselectTopic = (event) => {
    const id = event.target.id;
    removeElement(id);

    const { topics } = getActiveEnv();

    const list = document.getElementById('available-topics');
    const t = topics.find((t) => t.key === id);
    addTopicToList(t, list, selectTopic);
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
