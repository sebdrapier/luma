# LUMA

LUMA is a DMX lighting controller composed of a Go backend and a React
frontend. The application lets you manage fixtures, presets and shows and
controls a DMX interface over USB. The frontend is embedded in the Go binary
using `go:embed`.

## Requirements

- Go 1.24 or newer
- Node.js and npm (for building the frontend)
- a DMX USB interface

## Building the Frontend

```bash
cd src/frontend
npm install
npm run build
```

The build output is placed in `src/frontend/dist` and is embedded in the Go
binary.

## Running the Backend

```bash
go run ./src/backend
```

Environment variables can be used to configure the server:

- `SERVER_PORT` – HTTP port (default `:3000`)
- `DMX_PORT` – serial port used for DMX (default `/dev/cu.usbserial-A10QIXZO`)
- `DATA_FILE` – path to the project file (default `.data/project.yaml`)
- `ENABLE_DMX` – set to `false` to disable DMX output

After starting the server, open `http://localhost:3000` in your browser to use
the web interface.

## API

The backend exposes REST endpoints under `/api` and a WebSocket endpoint at
`/ws/control` for real‑time DMX commands. See the source in `src/backend/api`
and `src/backend/ws` for details.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file
for details.
