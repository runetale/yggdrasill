{
  description = "notch";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    let
      version = builtins.substring 0 8 self.lastModifiedDate;
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];
      eachSystem = flake-utils.lib.eachSystem supportedSystems;
    in
    eachSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };
      in
    {
      packages = {
        notch = pkgs.buildGoModule {
          pname = "notch";
          inherit version;
          src = self;
          buildPhase = ''
            export GOFLAGS="-mod=mod"
            export GOPROXY=https://proxy.golang.org,direct
            go build -o notch cli/cli.go
          '';
          installPhase = ''
            mkdir -p $out/bin
            cp notch $out/bin/
          '';
          allowGoReference = true;
          vendorHash = "sha256-Gw+9NMT5r1mMl9KHxKdVlv2sPTeJNImcHXU6eU/mERE=";
          meta = with pkgs.lib; {
            description = "notch project";
            license = licenses.mit;
            maintainers = with maintainers; [ "shinta@gx14ac.com" ];
          };
        };
      };
      devShells.default = pkgs.mkShell {
        buildInputs = with pkgs;
          [
            go
          ];
      };
      defaultPackage = self.packages.${system}.notch;
    });
}
