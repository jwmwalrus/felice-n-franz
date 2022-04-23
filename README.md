Felice & Franz
==================

A generic, web-based producer and consumer for Apache Kafka.

## Contents
* [Requirements](#requirements)
* [Installation](#installation)
* [Usage](#usage)
    * [Logging Flags](#logging-flags)
* [Configuration](#configuration)
* [TODO](#todo)

## Requirements

* Go version 1.16 or higher. 
You can probably install it through your system's package manager (`apt`, `brew`, etc.). 
For general instructions, go [here](https://golang.org/doc/install) --no pun intended.

* Any web browser with decent support for ES2020 (e.g., Mozilla Firefox or any Chromium-based browser).

## Installation

To install, open a terminal and execute the following
```bash
go install github.com/jwmwalrus/felice-n-franz@latest
```

The same command can be used for subsequent updates.

## Usage

Assuming `$GOBIN` is in your PATH, open a terminal and execute the following
```bash
felice-n-franz
```

Then click on the displayed URL (which by default should be `http://localhost:9191`).

The web interface is pretty self-explanatory (select environment, consume or produce messages by topic, filter, etc.).

### Logging Flags

Passing `-h` upon executing the `felice-n-franz` command, will display options related to logging.

The application's logs can be found at `${XDG_DATA_HOME}/felice-n-franz/`.

## Configuration

In order to operate properly (or at all), a `config.json` file is required.

Upon first run, one is generated at `${XDG_CONFIG_HOME}/felice-n-franz/`. The expected contents are as follows:

| Path | Type | Description | Default | Required |
| :--- | :--- | :---------- | :------ | :------- |
| version | int | Configuration file version | 1 | required |
| firstRun | bool | If true, overwrite the <code>${XDG_CONFIG_HOME}/felice-n-franz/config.json</code> file with in-memory values | <code>false</code> | optional |
| port | int | Application's port | <code>9191</code> | required |
| envs | array | Configured environments | <code>[]</code> | required |
| env.name | string | Environment's name |  | required |
| env.active | bool | If false, this environent's configuration will be ignored | false | optional |
| env.configuration | object | librdkafka configuration for this environment. See [here](https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md) for further details  | {} | required |
| env.vars | object | Environment variables used for topic values | <code>{"env": ""}</code> | <code>{}</code> | optional |
| env.headerPrefix | string | Common prefix for topic's headers |  | optional |
| env.topicsFrom | string | Copy <code>env.schemas</code>, <code>env.topics</code> and <code>env.groups</code> from the given environment |  | optional
| env.schemas | object | TBD | {} | optional
| env.topics | array | Configured topics for this environment | <code>[]</code> | required | env.topic.name | string | Topic name |  | optional |
| env.topic.description | string | Topic description |  | optional |
| env.topic.key | string | Topic key. Must be a unique identifier for this environment |  | required |
| env.topic.value | string | Topic value. Actual topic definition. Shall be unique for this environment. |  | required |
| env.topic.headers | array | Array of key-value objects defining the headers associated to this topic | <code>[{"key": "Content-Type","value": "application/json"}]</code> | <code>[]</code> | optional |
| env.topic.schema | object | TBD | <code>{}</code> | optional |
| env.groups | type | description | default | required/optional |
| env.group.name | string | Group's display name |  | optional |
| env.group.description | array | Group'ss description |  | optional |
| env.group.category | string | Group's category |  | optional |
| env.group.id | string | Group's identifier. Must be unique for this environment and not to collide with topic keys |  | required |
| env.group.keys | array | List of topic keys that belong to this group. Must not be empty |  | required |

## TODO

- [ ] Add a way of reordering consumer cards.
- [ ] Add topic schema validation.


