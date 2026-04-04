{
  description = "Nix flake for the Werkzeugkasten Go CLI";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      forAllSystems = f:
        nixpkgs.lib.genAttrs systems (system:
          let
            pkgs = import nixpkgs { inherit system; };
            goPkg = pkgs.go_1_26 or pkgs.go;
          in
          f {
            inherit pkgs goPkg;
          });
    in {
      packages = forAllSystems ({ pkgs, goPkg }:
        let
          werkzeugkasten = (pkgs.buildGoModule.override { go = goPkg; }) {
            pname = "werkzeugkasten";
            version = self.shortRev or "dev";
            src = ./.;

            vendorHash = "sha256-JApJKIHmdUbNBUDjSO5qYw1LffGzcXWQE3aGSYYF3zM=";

            ldflags = [
              "-s"
              "-w"
              "-X github.com/allaman/werkzeugkasten/cli.Version=${self.shortRev or "dev"}"
            ];

            meta = with pkgs.lib; {
              description = "CLI for downloading pre-built binaries";
              homepage = "https://github.com/allaman/werkzeugkasten";
              license = licenses.mit;
              mainProgram = "werkzeugkasten";
              platforms = platforms.unix;
            };
          };
        in {
          default = werkzeugkasten;
          inherit werkzeugkasten;
        });

      apps = forAllSystems ({ pkgs, ... }:
        {
          default = {
            type = "app";
            program = "${self.packages.${pkgs.system}.default}/bin/werkzeugkasten";
          };
        });

      devShells = forAllSystems ({ pkgs, goPkg }:
        {
          default = pkgs.mkShell {
            packages = with pkgs; [
              goPkg
              gopls
              golangci-lint
              govulncheck
              go-tools
              go-task
            ];
          };
        });
    };
}
