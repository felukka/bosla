{
  description = "A serious, well-scoped versioning tool.";

  inputs.nixpkgs.url = "nixpkgs/nixos-25.11";

  outputs =
    { self, nixpkgs }:
    let
      inherit (nixpkgs) lib;

      systems = [
        "x86_64-linux"
        "aarch64-linux"
      ];

      overlay = final: _: {
        bosla = final.callPackage ./package.nix { };
      };

      pkgsFor = lib.genAttrs systems (
        system:
        import nixpkgs {
          inherit system;
          overlays = [ overlay ];
        }
      );

      forAllSystems = f: lib.genAttrs systems (system: f pkgsFor.${system});
    in
    {
      overlays.default = overlay;

      packages = forAllSystems (pkgs: {
        default = pkgs.bosla;
      });

      devShells = forAllSystems (pkgs: {
        default = pkgs.mkShellNoCC {
          packages = with pkgs; [
            cobra-cli
            goimports-reviser
            golangci-lint
            golines
            sqlc
            (python3.withPackages (
              p: with p; [
                mkdocs-material
              ]
            ))
          ];
        };
      });
    };
}
