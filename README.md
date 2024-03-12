# The backend for SlugQuest

## Running the server

The server can be compiled into an executable and ran locally:
```bash
# Compile into an executable
go build -o server
./server
```

## Nix

An easy way to develop and run. Install from [here](https://nixos.org/download#nix-install-linux).

Set up flakes with:
```bash
mkdir ~/.config/nix
echo "experimental-features = nix-command flakes" > ~/.config/nix/nix.conf
```

### Commands:
```bash
# Get a development shell
nix develop

# Build the program
nix build

# Run the program
nix run
```

## Resources

* [Go Website](https://go.dev/)
