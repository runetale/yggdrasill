.PHONY: nix-build
nix-build:
	nix build .#notch --no-link --print-out-paths --print-build-logs --extra-experimental-features flakes