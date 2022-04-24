{ pkgs, ... }:

with pkgs; mkShell {
  buildInputs = [
    gnumake
    go
    golangci-lint
  ];
}
