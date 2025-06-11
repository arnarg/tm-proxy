{
  lib,
  buildGoApplication,
}: let
  version = "0.0.0-alpha.1";
in
  buildGoApplication {
    inherit version;
    pname = "tm-proxy";

    src =
      builtins.filterSource (
        path: type:
          type
          == "directory"
          || baseNameOf path == "go.mod"
          || baseNameOf path == "go.sum"
          || lib.hasSuffix ".go" path
      )
      ./.;

    modules = ./gomod2nix.toml;

    subPackages = ["cmd/tm-proxy"];
  }
