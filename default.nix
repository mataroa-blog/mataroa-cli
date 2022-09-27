{ pkgs, pname, version, ... }:

let
  inherit (pkgs)
    buildGoModule
    lib
    pandoc;
in
buildGoModule {
  inherit pname;
  version = "v${version}";

  src = lib.cleanSource ./.;

  vendorSha256 = "sha256-N3+gaqCJOp5xGOvcJd3OnhPpC1qY1hGzJkZUg7UNrIQ=";

  subPackages = [ "cmd/mata" ];

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
    description = "A CLI tool for mataroa.blog";
    license = licenses.mit;
    maintainers = with maintainers; [ ratsclub ];
  };
}
