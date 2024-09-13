# MC-Server-Monitor

Minecraft Server Monitor is a Golang web server application to monitor one's own Minecraft server (thru RCON). This is targeted towards a small Minecraft server size (<10) to do basic adminstrative tasks.  

## Dependencies

### Overall
* Go w/ STD
* SQLite
* Air (live-reloading for Go)
* Docker (for Minecraft server and cloud deployment)
* Tailwind CSS w/ DaisyUI (for simplified CSS)
* Prettier (for code formatting)

Look at ``package.json`` and ``go.mod`` for all dependencies in this project

## Usage (Local)

### Quickstart

To get started quickly (w/ [Docker](https://www.docker.com/) installed), use the docker command at the repository root

```bash
make server
```

This command boots up a minecraft server and the web server preemptively serving ports 25565 and 25575 respectively.

## Usage (hosted on GCP)

⚠️ Incurs cost

Cost Breakdown

-   Static IP - $1.49 / month
-   Compute Engine (Preemptible VM) - $0.01 / hr
-   Persistent Disk - $0.50 / month
-   VM-VM egress pricing - <=$0.15 per GB


## Acknowledgements

Thanks to [itzg](https://github.com/itzg) and [Futurice](https://github.com/futurice) for the docker image and prior terraform script respectively.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](https://choosealicense.com/licenses/mit/)
