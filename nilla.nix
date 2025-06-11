let
  pins = import ./npins;

  nilla = import pins.nilla;

  systems = ["x86_64-linux" "aarch64-linux"];
in
  nilla.create ({config}: {
    includes = [
      "${pins.nilla-utils}/modules"
    ];

    config = {
      # Load all pins from npins and generate nilla inputs.
      generators.inputs.pins = pins;

      # Add gomod2nix using overlay
      inputs.nixpkgs.settings.overlays = [
        (final: prev: let
          callPackage = final.callPackage;
        in {
          inherit (callPackage "${pins.gomod2nix}/builder" {}) buildGoApplication mkGoEnv mkVendorEnv;
          gomod2nix = callPackage "${pins.gomod2nix}/default.nix" {};
        })
      ];

      packages.default = config.packages.tm-proxy;
      packages.tm-proxy = {
        inherit systems;

        builder = "nixpkgs";
        settings.pkgs = config.inputs.nixpkgs.result;

        package = import ./default.nix;
      };

      shells.default = {
        inherit systems;

        builder = "nixpkgs";
        settings.pkgs = config.inputs.nixpkgs.result;

        shell = {
          mkShellNoCC,
          gomod2nix,
          ...
        }:
          mkShellNoCC {
            packages = [gomod2nix];
          };
      };
    };
  })
