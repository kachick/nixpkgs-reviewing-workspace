{
  lib,
  stdenv,
  gleam,
  erlang,
  rebar3,
  makeWrapper,
  gh,
  coreutils,
  git,
}:
let
  binPath = lib.makeBinPath [
    gh
    coreutils
    erlang
    rebar3
    gleam
    git
  ];
in
stdenv.mkDerivation {
  pname = "resume-gleam";
  version = "0.1.0";
  src = ./.;

  nativeBuildInputs = [
    makeWrapper
  ];

  buildPhase = "true";

  installPhase = ''
        mkdir -p $out/share/resume-gleam
        cp -r . $out/share/resume-gleam
        
        mkdir -p $out/bin
        cat <<EOF > $out/bin/resume-gleam
    #!/bin/sh
    set -e
    export PATH="${binPath}:\$PATH"

    # Pass the original PWD to Gleam so it can run gh commands in the correct context
    export GLEAM_WAKE_PATH="\$PWD"

    # Gleam needs a writable directory for build artifacts and dependencies
    TEMP_PROJECT=\$(mktemp -d)
    trap 'rm -rf "\$TEMP_PROJECT"' EXIT

    cp -r "$out/share/resume-gleam/." "\$TEMP_PROJECT/"
    chmod -R +w "\$TEMP_PROJECT"
    cd "\$TEMP_PROJECT"

    # Gleam needs a HOME directory for its cache
    export HOME=\$TEMP_PROJECT

    exec gleam run "\$@"
    EOF
        chmod +x $out/bin/resume-gleam
  '';

  meta = {
    description = "Resume to track a running job in console (Gleam version)";
    mainProgram = "resume-gleam";
  };
}
