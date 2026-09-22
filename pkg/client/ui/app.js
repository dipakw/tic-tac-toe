class Backend {
    constructor({ listen }) {
        this.sse = new EventSource("/sse");

        this.sse.onopen = () => {
            listen("open");
        };

        this.sse.onerror = () => {
            listen("error");
        };

        this.sse.onmessage = (event) => {
            try {
                listen("message", JSON.parse(event.data));
            } catch {
                // Nothing.
            }
        };
    }
}

const modes = {
    "wait": ["wait"],
    "ask": ["me", "form"],
    "play": ["me", "board", "peer"],
};

const blocks = [
    "wait",
    "me",
    "form",
    "board",
    "peer",
];

class App {
    state = {
        mode: "wait",
    };

    connected = false;

    getCell(x, y) {
        const row = document.querySelector(`#board > div:nth-child(${x + 1})`);
        const cell = row.querySelector(`button:nth-child(${y + 1})`);
        return cell;
    }

    onCellClick(x, y) {
        if (!this.connected) {
            return;
        }

        fetch(`/click?x=${x}&y=${y}`);
    }

    setupCells() {
        for (let x = 0; x < 3; x++) {
            for (let y = 0; y < 3; y++) {
                const cell = this.getCell(x, y);

                cell.addEventListener("click", () => {
                    this.onCellClick(x, y);
                });
            }
        }
    }

    setConnected() {
        this.connected = true;
        const span = document.querySelector("#status span");
        span.classList.add("connected");
        span.textContent = "Connected";
    }

    setDisconnected() {
        this.connected = false;
        const span = document.querySelector("#status span");
        span.classList.remove("connected");
        span.textContent = "Disconnected";
    }

    async startGame() {
        const peerId = document.querySelector("#play_peer_id").value.trim();

        if (peerId == "") {
            return;
        }

        await fetch(`/start-game?peer_id=${peerId}`);
    }

    setupEventListeners() {
        document.querySelector("#play").addEventListener("click", () => {
            this.startGame();
        });
    }

    setupBackend() {
        this.backend = new Backend({
            listen: (event, data) => {
                switch (event) {
                    case "open":
                        this.setConnected();
                        break;

                    case "error":
                        this.setDisconnected();
                        break;

                    case "message":
                        this.state = data;
                        this.update();
                        break;

                    default:
                        // Nothing.
                        break;
                }
            }
        });
    }

    updateBoard() {
        for (let x = 0; x < 3; x++) {
            for (let y = 0; y < 3; y++) {
                const stateCell = this.state?.rows?.[x]?.[y] ?? {};
                const value = stateCell?.value ?? "&nbsp;";
                const cell = this.getCell(x, y);
                cell.innerHTML = value;
            }
        }
    }

    updateMeInfo() {
        document.querySelector("#me_name").value = this.state?.profile?.name || "";
        document.querySelector("#me_id").value = this.state?.profile?.id || "";
    }

    updatePeerInfo() {
        document.querySelector("#peer_name").value = this.state?.peer?.name || "";
        document.querySelector("#peer_id").value = this.state?.peer?.id || "";
    }

    update() {
        // Hide all first.
        blocks.forEach((id) => {
            document.querySelector(`#${id}`).classList.add("none");
        });

        // Show necessary.
        (modes[this.state.mode] || []).forEach((id) => {
            document.querySelector(`#${id}`).classList.remove("none");
        });

        this.updateMeInfo();
        this.updatePeerInfo();
        this.updateBoard();
    }

    async run() {
        this.setupEventListeners();
        this.setupBackend();
        await fetch("/register");
        this.setupCells();
        this.update();
    }
}

(async () => {
    (new App()).run();
})();