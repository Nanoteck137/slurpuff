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
          vendorHash = "sha256-yC5g9Y4tCtHFPfeHTtip/xSoiR9yAWf0Ys/rAZ1C2+I=";

          nativeBuildInputs = [ pkgs.makeWrapper ];

          postFixup = ''
          wrapProgram $out/bin/slurpuff --prefix PATH : ${pkgs.lib.makeBinPath [ pkgs.ffmpeg ]}
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
          ];
        };
      }
    );
}
