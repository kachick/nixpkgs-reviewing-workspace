open Printf

module NixSystem = struct
  type t =
    { _platform : string
    ; tier : int
    }

  let systems =
    [ "X64-Linux", { _platform = "x86_64-linux"; tier = 1 }
    ; "ARM64-Linux", { _platform = "aarch64-linux"; tier = 2 }
    ; "X64-macOS", { _platform = "x86_64-darwin"; tier = 2 }
    ; "ARM64-macOS", { _platform = "aarch64-darwin"; tier = 3 }
    ]
  ;;
end

type priority =
  | Known of int * int
  | Undefined

let compare_priority a b =
  match a, b with
  | Known (t1, f1), Known (t2, f2) ->
    (match Int.compare t1 t2 with
     | 0 -> Int.compare f1 f2
     | n -> n)
  | Known _, Undefined -> -1
  | Undefined, Known _ -> 1
  | Undefined, Undefined -> 0
;;

let contains haystack needle =
  try
    ignore (Str.search_forward (Str.regexp_string needle) haystack 0);
    true
  with
  | Not_found -> false
;;

let extract_asset_name path =
  (* Identify the segment like "nixpkgs-review-files-..." *)
  String.split_on_char '/' path
  |> List.find_map (fun segment ->
    if String.starts_with ~prefix:"nixpkgs-review-files-" segment
    then
      List.find_map
        (fun (name, _) -> if contains segment name then Some name else None)
        NixSystem.systems
    else None)
;;

let get_priority path =
  match extract_asset_name path with
  | None -> Undefined
  | Some name ->
    let tier =
      match List.assoc_opt name NixSystem.systems with
      | Some s -> s.tier
      | None -> 999 (* This case is actually impossible due to find_map logic *)
    in
    let my_favor = if String.ends_with ~suffix:"Linux" name then 0 else 1 in
    Known (tier, my_favor)
;;

let read_file path =
  let ic = open_in path in
  let s = really_input_string ic (in_channel_length ic) in
  close_in ic;
  s
;;

let concat_reports paths =
  paths
  |> List.mapi (fun i path ->
    let content = read_file path in
    match Str.split_delim (Str.regexp "---") content with
    | [] -> ""
    | header :: bodies ->
      let body = String.concat "---" bodies in
      let header_part = if i = 0 then header else "" in
      let body_part =
        if body = "" then "" else if i = 0 then "---" ^ body else "\n---" ^ body
      in
      header_part ^ body_part)
  |> String.concat ""
;;

let run_command cmd args =
  let command = String.concat " " (cmd :: List.map Filename.quote args) in
  if Sys.command command <> 0 then failwith (sprintf "Command failed: %s" command)
;;

let rec list_files acc dir =
  Sys.readdir dir
  |> Array.fold_left
       (fun acc item ->
          let path = Filename.concat dir item in
          if Sys.is_directory path then list_files acc path else path :: acc)
       acc
;;

let () =
  let argv = Sys.argv |> Array.to_list in
  match List.tl argv with
  | [] ->
    eprintf "Usage: %s <run_id>\n" (List.hd argv);
    exit 1
  | run_id :: _ ->
    let in_runner =
      match Sys.getenv_opt "GITHUB_ACTIONS" with
      | Some "true" -> true
      | _ -> false
    in
    let current_asset =
      let get_arch () =
        let ic = Unix.open_process_in "uname -m" in
        let m = input_line ic in
        ignore (Unix.close_process_in ic);
        if m = "x86_64"
        then "X64"
        else if m = "arm64" || m = "aarch64"
        then "ARM64"
        else ""
      in
      let get_os () =
        let ic = Unix.open_process_in "uname -s" in
        let s = input_line ic in
        ignore (Unix.close_process_in ic);
        if s = "Linux" then "Linux" else if s = "Darwin" then "macOS" else ""
      in
      try
        let a = get_arch () in
        let o = get_os () in
        if a = "" || o = "" then "" else a ^ "-" ^ o
      with
      | _ -> ""
    in
    if not in_runner then run_command "gh" [ "run"; "watch"; run_id; "--interval"; "10" ];
    let temp_dir = Filename.temp_dir "nixpkgs-review-run-" ("." ^ run_id) in
    eprintf "Downloading artifacts to %s...\n" temp_dir;
    run_command "gh" [ "run"; "download"; run_id; "--dir"; temp_dir ];
    let all_files = list_files [] temp_dir in
    eprintf "Downloaded files:\n";
    all_files
    |> List.iter (fun path ->
      let is_current =
        match extract_asset_name path with
        | Some n -> String.equal n current_asset
        | None -> false
      in
      if is_current then eprintf "  \x1b[32m%s\x1b[0m\n" path else eprintf "  %s\n" path);
    let reports =
      all_files
      |> List.filter (fun p -> Filename.basename p = "report.md")
      |> List.sort (fun a b -> compare_priority (get_priority a) (get_priority b))
    in
    (match reports with
     | [] ->
       eprintf "No report.md found in %s\n" temp_dir;
       exit 1
     | _ ->
       print_endline "";
       print_endline (concat_reports reports))
;;
