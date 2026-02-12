{
  ocamlPackages,
  makeWrapper,
  lib,
  gh,
  coreutils,
}:
ocamlPackages.buildDunePackage {
  pname = "resume-ocaml";
  version = "0.1.0";
  src = ./.;

  duneVersion = "3";

  buildInputs = with ocamlPackages; [
    ocamlformat
  ];

  nativeBuildInputs = [
    makeWrapper
  ];

  postInstall = ''
    wrapProgram "$out/bin/resume-ocaml" --prefix PATH : ${
      lib.makeBinPath [
        gh
        coreutils
      ]
    }
  '';

  meta = {
    description = "Resume to track a running job in console (OCaml version)";
    mainProgram = "resume-ocaml";
  };
}
