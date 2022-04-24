{ pkgs, pname, version, ... }:

let
  inherit (pkgs)
    buildGoModule
    lib
    scdoc;
in
buildGoModule {
  inherit pname;
  version = "v${version}";

  src = lib.cleanSource ../.;

  nativeBuildInputs = with pkgs; [ scdoc ];

  vendorSha256 = "sha256-N3+gaqCJOp5xGOvcJd3OnhPpC1qY1hGzJkZUg7UNrIQ=";

  makeFlags = [
    "PREFIX=$(out)"
  ];

  postBuild = ''
    make $makeFlags
  '';

  preInstall = ''
    make $makeFlags install
  '';

  meta = with lib; {
    homepage = "https://sr.ht/~glorifiedgluer/mata";
    description = "A CLI tool for mataroa / mataroa.blog";
    license = licenses.mit;
    maintainers = with maintainers; [ ratsclub ];
  };
}
