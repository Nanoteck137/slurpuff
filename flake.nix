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
          version = "0.0.1"; 
          src = ./.;
          vendorHash = "sha256-LQrXsuv+t1KhQarh3HFGNmoq7mPiaXskbtC2T7g6T4Q=";
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
