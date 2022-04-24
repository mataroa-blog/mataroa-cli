{ pkgs, ... }:

with pkgs; mkShell {
  buildInputs = [
    gnumake
    go
    scdoc

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
}
