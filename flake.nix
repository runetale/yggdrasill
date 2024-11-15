{
  description = "yggdrasill";

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
        yggdrasill = pkgs.buildGoModule {
          pname = "yggdrasill";
          inherit version;
          src = self;
          buildPhase = ''
            export GOFLAGS="-mod=mod"
            export GOPROXY=https://proxy.golang.org,direct
            go build -o yggdrasill cli/cli.go
          '';
          installPhase = ''
            mkdir -p $out/bin
            cp yggdrasill $out/bin/
          '';
          allowGoReference = true;
          vendorHash = "sha256-9smmluvWI6CvcT0pzXLaJDe88sN8yB/DjRiUxWC1bAA=";
          meta = with pkgs.lib; {
            description = "yggdrasill project";
            license = licenses.mit;
            maintainers = with maintainers; [ "shinta@runetale.com" ];
          };
        };
      };
      devShells.default = pkgs.mkShell {
        buildInputs = with pkgs;
          [
            go
            git
          ];
      };
      defaultPackage = self.packages.${system}.yggdrasill;
    });
}
