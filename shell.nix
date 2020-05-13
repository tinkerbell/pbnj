let _pkgs = import <nixpkgs> { };
in { pkgs ? import (_pkgs.fetchFromGitHub {
  owner = "NixOS";
  repo = "nixpkgs-channels";
  #branch@date: nixpkgs-unstable@2020-04-17
  rev = "10100a97c8964e82b30f180fda41ade8e6f69e41";
  sha256 = "011f36kr3c1ria7rag7px26bh73d1b0xpqadd149bysf4hg17rln";
}) { } }:

with pkgs;

mkShell { buildInputs = [ go golangci-lint ]; }
