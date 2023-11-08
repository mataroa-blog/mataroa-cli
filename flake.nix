{
  description = "Sharpie monorepo";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    devenv = {
      url = "github:cachix/devenv";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = inputs@{ self, devenv, nixpkgs, ... }:
  let
    # System types to support.
    supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-darwin" ];

    # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
    forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

    # Nixpkgs instantiated for supported system types.
    nixpkgsFor = forAllSystems (system: import nixpkgs {
      inherit system;
    });
  in
  {
    devShells = forAllSystems (system:
    let
      pkgs = nixpkgsFor."${system}";
    in
    {
      # `nix develop`
      default = devenv.lib.mkShell {
        inherit inputs pkgs;
        modules = [
          ({ pkgs, lib, ... }: {
            languages.go.enable = true;
            pre-commit.hooks = {
              gofmt.enable = true;
              gotest.enable = true;
              govet.enable = true;
              staticcheck.enable = true;
            };
          })
        ];
      };
    });
  };
}
