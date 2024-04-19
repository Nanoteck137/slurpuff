{
  description = "Simple music convertion tool";

  inputs = {
    nixpkgs.url      = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url  = "github:numtide/flake-utils";
    opusimage.url    = "github:nanoteck137/opusimage";
  };

  outputs = { self, nixpkgs, flake-utils, opusimage, ... }:
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
          vendorHash = "sha256-R9eot+MwJtsY0arPlvzVyelajmgHi9KLAPSqsvVS/2Y=";

          nativeBuildInputs = [ pkgs.makeWrapper ];

          postFixup = ''
          wrapProgram $out/bin/slurpuff --prefix PATH : ${pkgs.lib.makeBinPath [ pkgs.ffmpeg opusimage.packages.${system}.default ]}
          '';
        };
      in
      {
        packages.default = program;
        packages.slurpuff = program;

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            ffmpeg
            opusimage.packages.${system}.default
          ];
        };
      }
    );
}
