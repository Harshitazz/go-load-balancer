# Dynamic Backend Load Balancer

This project provides a **dynamic load balancer** using Docker Compose, which can automatically adjust to multiple backend services and distribute traffic based on configurable weights.

## Features
- **Dynamic Backend Configuration**: Add, update, and remove backend services at runtime using a simple command-line interface (CLI).
- **Automated `docker-compose.yml` Generation**: The CLI tool generates or updates the `docker-compose.yml` and `backends.json` config files automatically.
- **Automatic Docker Compose Build & Up**: The `init` and `add` commands automatically trigger `docker-compose up --build` to build and start the containers, ensuring everything is up-to-date.
- **Health Checks**: The load balancer performs periodic health checks to ensure only healthy backends serve requests.
- **Weighted Round-Robin Load Balancing**: Requests are distributed to backends based on their weight, simulating a weighted round-robin algorithm.

## Requirements

Before using the project, make sure you have the following:
- **Docker**: [Install Docker](https://www.docker.com/get-started)
- **Docker Compose**: [Install Docker Compose](https://docs.docker.com/compose/install/)
- **Python 3.x**: [Download Python](https://www.python.org/downloads/)

## Project Structure

- `cli.py`: A Python script that handles backend management and automatic `docker-compose.yml` generation.
- `Dockerfile.backend`: Dockerfile for building the backend service containers.
- `docker-compose.yml`: Automatically generated or updated by the `cli.py` script to define services, including the load balancer and backends.
- `config/backends.json`: Config file containing backend details, including URLs and weights.

## Setup

1. Clone the repository to your local machine:

   ```bash
   git clone https://github.com/yourusername/load-balancer.git
   cd load-balancer

## CLI Commands

### `init` Command

The `init` command initializes the backend services, generates the necessary configuration files (`docker-compose.yml` and `backends.json`), and starts the services using `docker-compose up --build`.

###backend1:9001:2: Defines a backend service backend1 running on port 9001 with weight 2.

#### Usage

    ```bash
    python cli.py init backend1:9001:2 backend2:9002:3

### `add` Command

The add command allows you to add new backends to the existing configuration and restart the services.

    ```bash
    python cli.py add backend3:9003:1


### `remove` Command
Removes a backend from the configuration by either its name or port.

Usage:

```bash

python cli.py remove <backend_name_or_port>