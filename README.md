# Coverage

[![GitHub](https://img.shields.io/badge/project-Data_Together-487b57.svg?style=flat-square)](http://github.com/datatogether)
[![Slack](https://img.shields.io/badge/slack-Archivers-b44e88.svg?style=flat-square)](https://archivers-slack.herokuapp.com/)
[![License](https://img.shields.io/github/license/datatogether/coverage.svg?style=flat-square)](./LICENSE)
[![Codecov](https://img.shields.io/codecov/c/github/datatogether/coverage.svg?style=flat-square)](https://codecov.io/gh/datatogether/coverage)
[![CircleCI](https://img.shields.io/circleci/project/github/datatogether/coverage.svg?style=flat-square)](https://circleci.com/gh/datatogether/coverage)

**Coverage** is a project for visualizing the status of digital data archiving efforts across various data repositories run by different initiatives. Its current scope covers data within the epa.gov top-level domain.

This code repo provides the JSON back-end: [`https://api.archivers.co/coverage`](https://api.archivers.co/coverage)

The [`datatogether/webapp` repo](https://github.com/datatogether/webapp) provides the visual front-end: [`https://archivers.co/coverage`](https://archivers.co/coverage)


## Current Data Repositories

Actual source datasets can be found in each [`/repositories/*` directory](/repositories). It currently includes the following:

  * [Archivers 2](https://alpha.archivers.space/)
  * [archivers.space](https://archivers.space/)
  * [EDGI Nomination Tool](https://chrome.google.com/webstore/detail/nominationtool/abjpihafglmijnkkoppbookfkkanklok?hl=en) Uncrawlables
  * [The Internet Archive](https://archive.org/)
  * [Project Svalbard](https://github.com/datproject/svalbard) JSON-LD crawl

Requests for new data repositories are tracked under the [`data-repository`](https://github.com/datatogether/coverage/labels/data-repository) issue label.


##  How It Works

It takes a list of urls and associated archiving information, and turns that into a tree of url paths with associated coverage information.

The output is cached in [`cache.json`](cache.json). Because this is a large file, we provide incremental pieces of the cached tree as a web server. To dynamically calculate coverage completion to can work with the `cache.json` file.


## Routes

* `/healthcheck` - server status
* `/repositories` - listing of all data repositories
* `/repositories/:repository_uuid` - details of one data repository
* `/fulltree` - 
* `/coverage` - 
* `/tree` - alias for `/coverage`


## Getting Involved

We would love involvement from more people! If you notice any errors or would like to submit changes, please see our [Contributing Guidelines](./github/CONTRIBUTING.md).

We use GitHub issues for [tracking bugs and feature requests](./issues) and Pull Requests (PRs) for [submitting changes](./pulls)


## Installation

Running this project can be done either directly on your workstation system, or in a "container" via Docker.

For people comfortable with Docker, or who are excited to learn about it, it can be the best way to get going.

### Docker Install

Running this project via Docker requires:

  * [Docker](https://docs.docker.com/engine/installation/)
  * [`docker-compose`](https://docs.docker.com/compose/install/)

Running the project in a Docker container should be as simple as:

```
make setup
make run
```

If you get an error about a port "address already in use", you can change the `PORT` environment variable in your local `.env` file.

Barring any changes, you may now visit a JSON endpoint at: `http://localhost:8080/repositories`

### Local System Install

Running this project directly on your system requires:

  * Go 1.7+
  * Postgresql

(Setting up these services is outside the scope of this README.)

```
cd path/to/coverage
createdb datatogether_coverage
go build
go get ./
# Set a free port on which to serve JSON
export PORT=8080
# Your postgresql instance may be running on a different port
export POSTGRES_DB_URL=postgres://localhost:5432/datatogether_coverage
$GOPATH/bin/coverage
```

Barring any changes, you may now visit a JSON endpoint at: `http://localhost:8080/repositories`



## Development

Please follow the install instructions above! Inclusion of tests are appreciated!


## License & Copyright

Copyright (C) 2017 Data Together
This program is free software: you can redistribute it and/or modify it under
the terms of the GNU Affero General Public License as published by the Free Software
Foundation, version 3.0.

This program is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE.

See the [`LICENSE`](./LICENSE) file for details.
