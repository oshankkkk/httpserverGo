{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell {
  packages = with pkgs; [
    go
    gopls
  ];

  shellHook = ''
    echo "Dev environment loaded"

    export GOPATH="$PWD/.gopath"
    export PATH="$PWD/.gopath/bin:$PATH"

    echo "go:   $(go version)"
  '';
}


