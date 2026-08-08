{
  description = "A drop-in cross-platform reimplementation of macOS dot_clean";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    { self, nixpkgs }:
    let
      inherit (nixpkgs.lib) genAttrs;
      supportedSystems = [
        "aarch64-darwin"
        "aarch64-linux"
        "x86_64-darwin"
        "x86_64-linux"
      ];
      forAllSystems =
        f:
        genAttrs supportedSystems (
          system:
          f {
            inherit system;
            pkgs = import nixpkgs { inherit system; };
          }
        );
    in
    {
      formatter = forAllSystems ({ pkgs, ... }: pkgs.nixfmt);

      packages = forAllSystems (
        { pkgs, ... }:
        {
          default = pkgs.buildGoModule {
            pname = "dotclean";
            version = self.shortRev or self.dirtyShortRev or "dev";
            src = self;

            vendorHash = "sha256-XNj4dHQko/ECAACrD7qcgNq5pgwm5hHqSKpkhdcBD/U=";

            # Go and Make are already provided by buildGoModule macro
            nativeBuildInputs = [ pkgs.pandoc ];

            # Custom Make phases skip buildGoModule's go-test helpers.
            doCheck = false;

            buildPhase = ''
              runHook preBuild
              make build
              runHook postBuild
            '';

            installPhase = ''
              runHook preInstall
              make PREFIX=$out install
              runHook postInstall
            '';

            meta = {
              description = "A drop-in cross-platform reimplementation of macOS dot_clean";
              homepage = "https://github.com/imcrazytwkr/dotclean";
              license = pkgs.lib.licenses.asl20;
              mainProgram = "dotclean";
            };
          };
        }
      );

      devShells = forAllSystems (
        { pkgs, system }:
        {
          default = pkgs.mkShell {
            inputsFrom = [ self.packages.${system}.default ];
            packages = [
              pkgs.just
              pkgs.nixd
            ];
            shellHook = ''
              PS1='\u@dotclean-dev:\w/ > '
              export PS1
            '';
          };
        }
      );
    };
}
