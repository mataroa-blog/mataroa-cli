{
  description = "Mataroa CLI";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils }:
    let
      pname = "mata";
      version = "0.2.0";
    in
    {
      overlays.default = final: prev: {
        mata = final.callPackage ./default.nix {
          inherit final pname version;
        };
      };
    } //
    utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };

        inherit (pkgs)
          callPackage
          mkShell

          gnumake
          go

          # https://github.com/golang/vscode-go/blob/master/docs/tools.md
          delve
          go-outline
          golangci-lint
          gomodifytags
          gopls
          gopkgs
          gotests
          impl
          ;
      in
      rec {
        # `nix build`
        packages."${pname}" = callPackage ./default.nix {
          inherit pkgs pname version;
        };
        packages.default = packages."${pname}";

        # `nix run`
        apps."${pname}" = utils.lib.mkApp {
          drv = packages."${pname}";
        };
        apps.default = apps."${pname}";

        # `nix develop`
        devShells = {
          default = mkShell {
            buildInputs = [
              gnumake
              go

              # https://github.com/golang/vscode-go/blob/master/docs/tools.md
              delve
              go-outline
              golangci-lint
              gomodifytags
              gopls
              gopkgs
              gotests
              impl
            ];
          };

          ci = mkShell {
            buildInputs = [
              gnumake
              go
              golangci-lint
            ];
          };
        };
      });
}
