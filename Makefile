.PHONY: nix-build
nix-build:
	nix build .#yggdrasill --no-link --print-out-paths --print-build-logs --extra-experimental-features flakes