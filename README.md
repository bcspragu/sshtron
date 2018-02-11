# SSHTron

SSHTron is a multiplayer lightcycle game that runs through SSH. Just run the command below and you'll be playing in seconds:

    $ ssh sshtron.zachlatta.com

_Controls: WASD or vim keybindings to move (**do not use your arrow keys**). Escape or Ctrl+C to exit._

**Code quality disclaimer:** _SSHTron was built in ~20 hours at [BrickHack 2](https://brickhack.io/). Here be dragons._

## Want to choose color yourself?

There are total 7 colors to choose from: Red, Green, Yellow, Blue, Magenta, Cyan and White

    $ ssh red@sshtron.zachlatta.com

If the color you picked is already taken in all open games, you'll randomly be assigned a color.

## Running Your Own Copy (in Docker)

Clone the project and `cd` into its directory.

```sh
# Build the SSHTron Docker image
$ docker build -t sshtron .

# Spin up the container with always-restart policy
$ docker run -t -p 2022:2022 --restart=always --name sshtron sshtron
```

For Raspberry Pi, change the base image in `Dockerfile` from `golang:latest` to `apicht/rpi-golang:latest`.

## License

SSHTron is licensed under the MIT License. See the full license text in [`LICENSE`](LICENSE).
