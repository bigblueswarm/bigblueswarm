<p align="center">
<img src="assets/header.png" alt="Architecture" />
</p>

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/44dddb06956746e98d474324a1dbbe5a)](https://www.codacy.com/gh/bigblueswarm/bigblueswarm/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=bigblueswarm/bigblueswarm&amp;utm_campaign=Badge_Grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/44dddb06956746e98d474324a1dbbe5a)](https://www.codacy.com/gh/bigblueswarm/bigblueswarm/dashboard?utm_source=github.com&utm_medium=referral&utm_content=bigblueswarm/bigblueswarm&utm_campaign=Badge_Coverage)
[![Code linting](https://github.com/bigblueswarm/bigblueswarm/actions/workflows/lint.yml/badge.svg)](https://github.com/bigblueswarm/bigblueswarm/actions/workflows/lint.yml)
[![Unit tests and coverage](https://github.com/bigblueswarm/bigblueswarm/actions/workflows/unit_test.yml/badge.svg)](https://github.com/bigblueswarm/bigblueswarm/actions/workflows/unit_test.yml)
[![Integration tests](https://github.com/bigblueswarm/bigblueswarm/actions/workflows/integration_test.yml/badge.svg)](https://github.com/bigblueswarm/bigblueswarm/actions/workflows/integration_test.yml)\
[![Docker build](https://github.com/bigblueswarm/bigblueswarm/actions/workflows/docker_image.yml/badge.svg)](https://github.com/bigblueswarm/bigblueswarm/actions/workflows/docker_image.yml)
![Docker Image](https://img.shields.io/docker/v/sledunois/bigblueswarm?label=Docker)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/bigblueswarm/bigblueswarm)
![GitHub](https://img.shields.io/github/license/bigblueswarm/bigblueswarm)


BigBlueSwarm is an open source metric-based multi-tenant load balancer that manages a pool of [BigBlueButton](https://bigbluebutton.org/) servers, an open source web conferencing system for online learning. It works as a proxy and makes the server pool appear as a single server. Send standard BigBlueButton API requests and BigBlueSwarm distributes these requests to the least loaded BigBlueButton server in the pool.

## Documentation

- [Introduction](docs/introduction.md)
- [First steps](docs/first_steps/readme.md)
  - [Installation](docs/first_steps/installation.md)
  - [Configuration](docs/first_steps/configuration.md)
  - [Initialize your cluster](docs/first_steps/initialization.md)
- [API](docs/api/readme.md)
  - [Custom errors](docs/api/CustomErrors.md)
  - [InstanceList](docs/api/InstanceList.md)
  - [Tenant](docs/api/Tenant.md)


## Manage BigBlueSwarm
Manage your BigBlueSwarm cluster using the [bbsctl](https://github.com/bigblueswarm/bbsctl) cli tool.

## Roadmap
Checkout [BigBlueSwarm public roadmap](https://github.com/users/SLedunois/projects/4).

## Contributors

![GitHub Contributors Image](https://contrib.rocks/image?repo=bigblueswarm/bigblueswarm)
