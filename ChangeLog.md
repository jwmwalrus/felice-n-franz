ChangeLog
=========

All noticeable changes in the project  are documented in this file.

Format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

This project uses [semantic versions](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.59.1] 2022-08-27

Upgrade dependencies

## [0.59.0] 2022-07-23

Docker option

### Fixed

* Linter-related issues

### Added

* Option to install with docker-compose

## [0.58.6] 2022-04-22

Update README

## [0.58.5] 2022-04-22

Update Go version and dependencies

## [0.58.4] 2022-01-12

Update Go version and dependencies

## [0.58.3] 2021-12-03

Update bumpy-ride

## [0.58.2] 2021-10-03

Quick fixes

### Fixed

* False positives when getting consumer
* White color on white background issue

## [0.58.1] 2021-09-30

Cleanup

## [0.58.0] 2021-09-30

Reimplement quit mechanism

### Added

* Fully functional stop button for lookup

### Modified

* Unregistering mechanism for topic consumers

## [0.57.0] 2021-09-28

Update dependencies

### Updated

* package.json's dependencies

### Modified

* autoComplete.js's configuration to account for breaking changes.

## [0.56.0] 2021-09-27

New features

### Added

* Predefined headers selector
* Accept user-provided config.json

### Modified

* Use RegExp in filter

## [0.55.1] 2021-09-27

WIP

## [0.55.0] 2021-09-26

Features and improvements

### Added

* Lookup topic messages

### Modified

* Reset/Refresh/Clear handling in bag

### Removed

* Layout footer

### Updated

* Dependencies

## [0.54.2] 2021-07-21

Fix

### Fixed

* Typo in map key

## [0.54.1] 2021-07-21

Fix

### Fixed

* groupId vars expansion for topics

## [0.54.0] 2021-07-19

Features and improvements

### Added

* Options to refresh and reset registry
* Toasts in bag
* Filter available topics
* Group ID per topic

### Modified

* Review toasts backend

## [0.53.3] 2021-06-27

### Modified

* Disabled the Enter key for input

## [0.53.2] 2021-06-20

Minor changes

### Updated

* Imports

### Added

* Tag-related tasks

## [0.53.1] 2021-06-16

Quick fix.

### Fixed

* Subscriber skipping messages

## [0.53.0] 2021-06-14

New features, fixes and improvements

### Fixed

* Toast not displaying properly

### Added

* AssignConsumer
- Filter badge

### Modified

* Rename Search to Filter
* Assets now have version suffix

### Removed

* Messages limit for cards

## [0.52.0] 2021-05-20

Improvements and Fixes

### Fixed

* Environment name for generated config.json

* Webpack's public path

### Modified

* Use fontawesome icons

### Removed

* Dependencies not needed anymore

## [0.51.0] 2021-05-06

New features and improvements

### Added

* Consumer group categories

* Search box in navbar

### Modified

* Show only payload in the message details
    - Use table and pre instead of textarea

### Removed

* Envelope information in details

* Filter modal

## [0.50.2] 2021-05-03

More fixes

## Fixed

* Trim filter values
    - Also, swap envelope and payload

* CSS-related adjustments
    - Increase textarea fonts
    - Remove unused selectors

* Add TODO section to README

* CSS-related issues with color and background-color

## [0.50.1] 2021-05-03

CSS-related fixes

## Fixed

* CSS-related issues with color and background-color

* Typo in the HTML file

## [0.50.0] 2021-05-01

Basic, fully functional, implementation.

## [0.49.0] 2021-04-21

Initial release.
