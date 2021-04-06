// assets/js/index.js
var apiUrl = window.location.origin;
var activeEnv = {};
var addConsumerCard = (t) => {
  const parent = document.getElementById("consumer-cards");
  const list = document.createElement("DIV");
  list.classList.add("list-group");
  const title = document.createElement("H5");
  title.classList.add("card-title");
  title.innerText = t.value;
  const body = document.createElement("DIV");
  body.classList.add("card-body");
  const card = document.createElement("DIV");
  card.classList.add("card");
  card.classList.add("bg-dark");
  body.appendChild(title);
  body.appendChild(list);
  card.appendChild(body);
  parent.appendChild(card);
};
var addGroupToList = (g, l, cb) => {
  const node = document.createElement("DIV");
  node.setAttribute("id", g.name);
  node.classList.add("list-group-item");
  node.classList.add("list-group-item-action");
  node.ondblclick = cb;
  const textnode = document.createTextNode(g.description);
  node.appendChild(textnode);
  l.appendChild(node);
};
var addTopicToList = (t, l, cb) => {
  const node = document.createElement("DIV");
  node.setAttribute("id", t.key);
  node.classList.add("list-group-item");
  node.classList.add("list-group-item-action");
  node.ondblclick = cb;
  const textnode = document.createTextNode(t.value);
  node.appendChild(textnode);
  l.appendChild(node);
};
var clearTopicsAndGroups = () => {
  document.getElementById("available-groups").innerHTML = "";
  document.getElementById("available-topics").innerHTML = "";
  document.getElementById("selected-groups").innerHTML = "";
  document.getElementById("selected-topics").innerHTML = "";
};
var removeElement = (id) => {
  const elem = document.getElementById(id);
  elem.parentNode.removeChild(elem);
};
var removeAllCards = () => {
  const parent = document.getElementById("consumer-cards");
  parent.innerHTML = "";
};
var selectGroup = (event) => {
  const id = event.target.id;
  removeElement(id);
  const groupsList = document.getElementById("selected-groups");
  const g = activeEnv.groups.find((g2) => g2.name === id);
  addGroupToList(g, groupsList, unselectGroup);
  g.keys.forEach((k) => {
    selectTopic({target: {id: k}});
  });
};
var selectTopic = (event) => {
  const id = event.target.id;
  removeElement(id);
  const list = document.getElementById("selected-topics");
  const t = activeEnv.topics.find((t2) => t2.key === id);
  addTopicToList(t, list, unselectTopic);
};
var unselectGroup = (event) => {
  const id = event.target.id;
  removeElement(id);
  const groupsList = document.getElementById("available-groups");
  const g = activeEnv.groups.find((g2) => g2.name === id);
  addGroupToList(g, groupsList, selectGroup);
  g.keys.forEach((k) => {
    unselectTopic({target: {id: k}});
  });
};
var unselectTopic = (event) => {
  const id = event.target.id;
  removeElement(id);
  const list = document.getElementById("available-topics");
  const t = activeEnv.topics.find((t2) => t2.key === id);
  addTopicToList(t, list, selectTopic);
};
var populateAvailable = () => {
  clearTopicsAndGroups();
  const groupsList = document.getElementById("available-groups");
  activeEnv.groups.forEach((g) => addGroupToList(g, groupsList, selectGroup));
  const topicsList = document.getElementById("available-topics");
  activeEnv.topics.forEach((t) => addTopicToList(t, topicsList, selectTopic));
};
window.addSelectedCards = async () => {
  removeAllCards();
  const parent = document.getElementById("selected-topics");
  let list = parent.querySelectorAll("div");
  list = Array.from(list);
  list.forEach((l) => {
    const t = activeEnv.topics.find((t2) => t2.key === l.id);
    addConsumerCard(t);
  });
  const topicKeys = list.map((l) => l.id);
  try {
    const method = "POST";
    const body = JSON.stringify({env: activeEnvironment, topicKeys});
    await fetch(`${apiUrl}/consumers`, {method, body});
  } catch (e) {
    console.error(e);
  }
};
window.checkIfValidEnvironment = async (sel) => {
  const {value} = sel;
  if (value !== "") {
    try {
      const res = await fetch(`${apiUrl}/envs/${value}`);
      activeEnv = await res.json();
      console.log(activeEnv);
      populateAvailable();
    } catch (e) {
      console.error(e);
    }
    document.getElementById("clear-contents-btn").disabled = false;
  } else {
    document.getElementById("clear-contents-btn").disabled = true;
    clearTopicsAndGroups();
    activeEnv = {};
  }
};
window.resetConsumers = () => {
  removeAllCards();
  populateAvailable();
};
var js_default = populateAvailable;
export {
  js_default as default
};
//# sourceMappingURL=index.js.map
