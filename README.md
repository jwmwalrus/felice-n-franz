Felice & Franz
==================

A generic, web-based producer and consumer for Apache Kafka.

## Table of Contents
* [Requirements](#requirements)
* [Installation](#installation)
* [Usage](#usage)
* [Configuration](#configuration)

## Requirements

* Go version 1.16 or higher. 
You can probably install it through your system's package manager (`apt`, `brew`, etc.). 
For general instructions, go [here](https://golang.org/doc/install) --no pun intended.

* Any web browser with decent support for ES2020 (e.g., Mozilla Firefox or any Chromium-based browser).

## Installation

To install, open a terminal and execute the following
```bash
go get -u github.com/jwmwalrus/felice-n-franz
```

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

<dl>
    <dt>version</dt>
    <dd>
        <dl>
            <dt>type</dt>
            <dd>int</dd>
            <dt>description</dt>
            <dd>Configuration file version</dd>
            <dt>default</dt>
            <dd><code>1</code></dd>
            <dt>required/optional</dt>
            <dd>required</dd>
        </dl>
    </dd>
    <dt>firstRun</dt>
    <dd>
        <dl>
            <dt>type</dt>
            <dd>bool</dd>
            <dt>description</dt>
            <dd>If true, overwrite the <code>config.json</code> file upon reading its contents</dd>
            <dt>default</dt>
            <dd><code>false</code></dd>
            <dt>required/optional</dt>
            <dd>optional</dd>
        </dl>
    </dd>
    <dt>port</dt>
    <dd>
        <dl>
            <dt>type</dt>
            <dd>int</dd>
            <dt>description</dt>
            <dd>Application's port</dd>
            <dt>default</dt>
            <dd><code>9191</code></dd>
            <dt>required/optional</dt>
            <dd>required</dd>
        </dl>
    </dd>
    <dt>maxTailOffset</dt>
    <dd>
        <dl>
            <dt>type</dt>
            <dd>int</dd>
            <dt>description</dt>
            <dd>Keep only the last <code>n</code> messages per consumer, with <code>n</code> given by this value</dd>
            <dt>default</dt>
            <dd><code>100</code></dd>
            <dt>required/optional</dt>
            <dd>required</dd>
        </dl>
    </dd>
    <dt>consumeForward</dt>
    <dd>
        <dl>
            <dt>type</dt>
            <dd>bool</dd>
            <dt>description</dt>
            <dd>Ignore messages sent by the broker, if timestamp is older than consumer's creation</dd>
            <dt>default</dt>
            <dd><code>false</code></dd>
            <dt>required/optional</dt>
            <dd>optional</dd>
        </dl>
    </dd>
    <dt>envs</dt>
    <dd>
        <dl>
            <dt>type</dt>
            <dd>array</dd>
            <dt>description</dt>
            <dd>Configured environments</dd>
            <dt>default</dt>
            <dd><code>[]</code></dd>
            <dt>required/optional</dt>
            <dd>required</dd>
        </dl>
    </dd>
    <dt>env.name</dt>
    <dd>
        <dl>
            <dt>type</dt>
            <dd>string</dd>
            <dt>description</dt>
            <dd>Environment's name</dd>
            <dt>default</dt>
            <dd></dd>
            <dt>required/optional</dt>
            <dd>required</dd>
        </dl>
    </dd>
    <dt>env.active</dt>
    <dd>
        <dl>
            <dt>type</dt>
            <dd>bool</dd>
            <dt>description</dt>
            <dd>If false, this environent's configuration will be ignored</dd>
            <dt>default</dt>
            <dd>false</dd>
            <dt>required/optional</dt>
            <dd>optional</dd>
        </dl>
    </dd>
    <dt>env.configuration</dt>
    <dd>
        <dl>
            <dt>type</dt>
            <dd>object</dd>
            <dt>description</dt>
            <dd>librdkafka configuration for this environment. See [here]() for further details </dd>
            <dt>default</dt>
            <dd>{}</dd>
            <dt>required/optional</dt>
            <dd>required</dd>
        </dl>
    </dd>
    <dt>env.vars</dt>
    <dd>
        <dl>
            <dt>type</dt>
            <dd>object</dd>
            <dt>description</dt>
            <dd>Environment variables used for topic values</dd>
            <dt>example</dt>
            <dd><code>{"env": ""}</code></dd>
            <dt>default</dt>
            <dd><code>{}</code></dd>
            <dt>required/optional</dt>
            <dd>optional</dd>
        </dl>
    </dd>
    <dt>env.headerPrefix</dt>
    <dd>
        <dl>
            <dt>type</dt>
            <dd>string</dd>
            <dt>description</dt>
            <dd>Common prefix for topic's headers</dd>
            <dt>default</dt>
            <dd></dd>
            <dt>required/optional</dt>
            <dd>optional</dd>
        </dl>
    </dd>
    <dt>env.topicsFrom</dt>
    <dd>
        <dl>
            <dt>type</dt>
            <dd>string</dd>
            <dt>description</dt>
            <dd>Copy <code>env.schemas</code>, <code>env.topics</code> and <code>env.groups</code> from the given environment</dd>
            <dt>default</dt>
            <dd></dd>
            <dt>required/optional</dt>
            <dd>optional</dd>
        </dl>
    </dd>
    <dt>env.schemas</dt>
    <dd>
        <dl>
            <dt>type</dt>
            <dd>object</dd>
            <dt>description</dt>
            <dd>TBD</dd>
            <dt>default</dt>
            <dd>{}</dd>
            <dt>required/optional</dt>
            <dd>optional</dd>
        </dl>
    </dd>
    <dt>env.topics</dt>
    <dd>
        <dl>
            <dt>type</dt>
            <dd>array</dd>
            <dt>description</dt>
            <dd>Configured topics for this environment</dd>
            <dt>default</dt>
            <dd><code>[]</code></dd>
            <dt>required/optional</dt>
            <dd>required</dd>
        </dl>
    </dd>
    <dt>env.topic.name</dt>
    <dd>
        <dl>
            <dt>type</dt>
            <dd>string</dd>
            <dt>description</dt>
            <dd>Topic name</dd>
            <dt>default</dt>
            <dd></dd>
            <dt>required/optional</dt>
            <dd>optional</dd>
        </dl>
    </dd>
    <dt>env.topic.description</dt>
    <dd>
        <dl>
            <dt>type</dt>
            <dd>string</dd>
            <dt>description</dt>
            <dd>Topic description</dd>
            <dt>default</dt>
            <dd></dd>
            <dt>required/optional</dt>
            <dd>optional</dd>
        </dl>
    </dd>
    <dt>env.topic.key</dt>
    <dd>
        <dl>
            <dt>type</dt>
            <dd>string</dd>
            <dt>description</dt>
            <dd>Topic key. Must be a unique identifier for this environment</dd>
            <dt>default</dt>
            <dd></dd>
            <dt>required/optional</dt>
            <dd>required</dd>
        </dl>
    </dd>
    <dt>env.topic.value</dt>
    <dd>
        <dl>
            <dt>type</dt>
            <dd>string</dd>
            <dt>description</dt>
            <dd>Topic value. Actual topic definition. Shall be unique for this environment.</dd>
            <dt>default</dt>
            <dd></dd>
            <dt>required/optional</dt>
            <dd>required</dd>
        </dl>
    </dd>
    <dt>env.topic.headers</dt>
    <dd>
        <dl>
            <dt>type</dt>
            <dd>array</dd>
            <dt>description</dt>
            <dd>Array of key-value objects defining the headers associated to this topic</dd>
            <dt>example</dt>
            <dd><code>[{"key": "Content-Type","value": "application/json"}]</code></dd>
            <dt>default</dt>
            <dd><code>[]</code></dd>
            <dt>required/optional</dt>
            <dd>optional</dd>
        </dl>
    </dd>
    <dt>env.topic.schema</dt>
    <dd>
        <dl>
            <dt>type</dt>
            <dd>object</dd>
            <dt>description</dt>
            <dd>TBD</dd>
            <dt>default</dt>
            <dd><code>{}</code></dd>
            <dt>required/optional</dt>
            <dd>optional</dd>
        </dl>
    </dd>
    <dt>env.groups</dt>
    <dd>
        <dl>
            <dt>type</dt>
            <dd>array</dd>
            <dt>description</dt>
            <dd>Groups of topics</dd>
            <dt>default</dt>
            <dd><code>[]</code></dd>
            <dt>required/optional</dt>
            <dd>optional</dd>
        </dl>
    </dd>
    <dt>env.group.name</dt>
    <dd>
        <dl>
            <dt>type</dt>
            <dd>string</dd>
            <dt>description</dt>
            <dd>Group's display name</dd>
            <dt>default</dt>
            <dd></dd>
            <dt>required/optional</dt>
            <dd>optional</dd>
        </dl>
    </dd>
    <dt>env.group.description</dt>
    <dd>
        <dl>
            <dt>type</dt>
            <dd>array</dd>
            <dt>description</dt>
            <dd>Group'ss description</dd>
            <dt>default</dt>
            <dd></dd>
            <dt>required/optional</dt>
            <dd>optional</dd>
        </dl>
    </dd>
    <dt>env.group.category</dt>
    <dd>
        <dl>
            <dt>type</dt>
            <dd>string</dd>
            <dt>description</dt>
            <dd>TBD</dd>
            <dt>default</dt>
            <dd></dd>
            <dt>required/optional</dt>
            <dd>optional</dd>
        </dl>
    </dd>
    <dt>env.group.id</dt>
    <dd>
        <dl>
            <dt>type</dt>
            <dd>string</dd>
            <dt>description</dt>
            <dd>Group's identifier. Must be unique for this environment and not to collide with topic keys</dd>
            <dt>default</dt>
            <dd></dd>
            <dt>required/optional</dt>
            <dd>required</dd>
        </dl>
    </dd>
    <dt>env.group.keys</dt>
    <dd>
        <dl>
            <dt>type</dt>
            <dd>array</dd>
            <dt>description</dt>
            <dd>List of topic keys that belong to this group. Must not be empty</dd>
            <dt>default</dt>
            <dd></dd>
            <dt>required/optional</dt>
            <dd>required</dd>
        </dl>
    </dd>
</dl>

