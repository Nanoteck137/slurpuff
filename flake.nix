{
  description = "Simple music convertion tool";

  inputs = {
    nixpkgs.url      = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url  = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        overlays = [];
        pkgs = import nixpkgs {
          inherit system overlays;
        };

        program = pkgs.buildGoModule {
          pname = "slurpuff";
          version = self.shortRev or "dirty";
          src = ./.;
          vendorHash = "sha256-y5McAeXQWLO2VkS5AG1nUY0LpiwLCztKjiqOm0vgWm0=";

          buildInputs = [
            pkgs.ffmpeg
          ];
        };
      in
      {
        packages.default = program;
        packages.slurpuff = program;

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            ffmpeg
          ];
        };
      }
    );
}
