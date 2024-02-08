# The backend for slugquest

## Resources

* [Go Website](https://go.dev/)

## Commands

```bash
# Run the program
go run main.go
```

## Nix

An easy way to develop and run

Install from [here](https://nixos.org/download#nix-install-linux)

Set up flakes with:
```bash
mkdir ~/.config/nix
echo "experimental-features = nix-command flakes" > ~/.config/nix/nix.conf
```

Commands:
```bash
# Get a development sheel
nix develop
# Build the program
nix build
# Run the program
nix run
```

