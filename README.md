# MC-Server-Monitor

Minecraft Server Monitor is a Golang web server application to monitor one's own Minecraft server (thru RCON). This is targeted towards a small Minecraft server size (<10) to do basic adminstrative tasks.  
A terraform script (for GCP) is provided to host the server on the internet quickly (thanks to [Futurice](https://github.com/futurice/terraform-examples) for help)

## Installation

Use [git](https://git-scm.com/) version control to retrieve the Minecraft Server Monitor repository.

```bash
git pull https://github.com/itzsBananas/mc-server-monitor.git
```

## Usage (Local)

### Quickstart

To get started quickly (w/ [Docker](https://www.docker.com/)), use the docker command at the repository root

```bash
docker compose up
```

This command boots up a minecraft server and the web server preemptively serving ports 25565 and 25575 respectively.

### Advanced

If you have a local server already running, you can run the command (w/ [Golang](https://go.dev/) at the repository root

```bash
go run ./cmd/web
```

Alternatively, with Docker

```bash
docker compose up mc-server-monitor
```

Both commands boot up the web server at port 25575.

## Usage (hosted on GCP)

IN PROGRESS; use the scripts in the deployments directory (for advanced users)

## Acknowledgements

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](https://choosealicense.com/licenses/mit/)
