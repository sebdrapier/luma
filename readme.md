# LUMA

LUMA is a DMX lighting controller with a Go backend and a React frontend. It lets you manage fixtures, presets and shows, and controls a DMX USB interface. The frontend is compiled and embedded directly into the Go binary using `go:embed` for easy deployment.

## Features

- Manage fixtures, presets and shows
- Real-time DMX output over USB
- REST and WebSocket APIs for remote control
- Single binary distribution with the React UI included

## Requirements

- Go 1.24 or newer
- Node.js and npm (for building the frontend)
- A DMX USB interface

## Quick Start

Use the provided script to build the frontend and start the backend:

```bash
./dev.sh
```

Pass `-n` to skip rebuilding the frontend if it has already been built.

Open `http://localhost:3000` in your browser to access the UI.

## Building the Frontend Manually

```bash
cd src/frontend
npm install
npm run build
```

The build output is written to `src/frontend/dist` and will be embedded in the Go binary.

## Running the Backend

```bash
go run ./src/backend
```

Configuration is done via environment variables:

- `SERVER_PORT` – HTTP port (default `:3000`)
- `DMX_PORT` – serial port used for DMX (default `/dev/cu.usbserial-A10QIXZO`)
- `DATA_FILE` – path to the project file (default `.data/project.yaml`)
- `ENABLE_DMX` – set to `false` to disable DMX output

After starting the server, open `http://localhost:3000` in your browser to use
the web interface.

## Docker

A `dockerfile` is provided to build a self-contained image:

```bash
docker build -t luma:latest .
docker run --rm -p 3000:3000 luma
```

## API

The backend exposes REST endpoints under `/api` and a WebSocket endpoint at `/ws/control` for real-time DMX commands. See the source in `src/backend/api` and `src/backend/ws` for details.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
