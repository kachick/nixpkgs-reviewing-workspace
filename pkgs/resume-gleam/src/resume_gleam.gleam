import argv
import gleam/erlang/os
import gleam/int
import gleam/io
import gleam/list
import gleam/option.{type Option, None, Some}
import gleam/result
import gleam/string
import shellout
import simplifile

type NixSystem {
  NixSystem(platform: String, tier: Int)
}

const systems = [
  #("X64-Linux", NixSystem("x86_64-linux", 1)),
  #("ARM64-Linux", NixSystem("aarch64-linux", 2)),
  #("X64-macOS", NixSystem("x86_64-darwin", 2)),
  #("ARM64-macOS", NixSystem("aarch64-darwin", 3)),
]

fn get_current_asset_name(work_dir: String) -> String {
  let arch = case shellout.command("uname", ["-m"], work_dir, []) {
    Ok(m) -> {
      let m = string.trim(m)
      case m {
        "x86_64" -> "X64"
        "arm64" | "aarch64" -> "ARM64"
        _ -> ""
      }
    }
    Error(_) -> ""
  }

  let os_name = case os.family() {
    os.Linux -> "Linux"
    os.Darwin -> "macOS"
    _ -> ""
  }

  case arch == "" || os_name == "" {
    True -> ""
    False -> arch <> "-" <> os_name
  }
}

fn extract_asset_name(path: String) -> Option(String) {
  let seg =
    path
    |> string.split("/")
    |> list.find(fn(segment) {
      string.starts_with(segment, "nixpkgs-review-files-")
    })

  case seg {
    Ok(segment) -> {
      systems
      |> list.find_map(fn(sys) {
        let #(name, _) = sys
        case string.contains(segment, name) {
          True -> Ok(name)
          False -> Error(Nil)
        }
      })
      |> option.from_result
    }
    Error(_) -> None
  }
}

pub fn main() {
  case argv.load().arguments {
    [run_id] -> run(run_id)
    _ -> io.println("Usage: resume_gleam <run_id>")
  }
}

fn list_all_files(dir: String) -> List(String) {
  case simplifile.read_directory(dir) {
    Ok(items) -> {
      items
      |> list.flat_map(fn(item) {
        let path = dir <> "/" <> item
        case simplifile.is_directory(path) {
          Ok(True) -> list_all_files(path)
          _ -> [path]
        }
      })
    }
    Error(_) -> []
  }
}

fn run(run_id: String) {
  let in_runner = os.get_env("GITHUB_ACTIONS") == Ok("true")
  // The original working directory where the user invoked the command
  let work_dir = os.get_env("GLEAM_WAKE_PATH") |> result.unwrap(".")
  let current_asset = get_current_asset_name(work_dir)

  let _ = case in_runner {
    False -> {
      io.println("Watching run " <> run_id <> "...")
      let _ =
        shellout.command(
          "gh",
          ["run", "watch", run_id, "--interval", "10"],
          work_dir,
          [],
        )
      Nil
    }
    True -> Nil
  }

  let temp_dir = "/tmp/nixpkgs-review-run-" <> run_id
  let _ = simplifile.create_directory_all(temp_dir)

  io.println("Downloading artifacts to " <> temp_dir <> "...")
  let res =
    shellout.command(
      "gh",
      ["run", "download", run_id, "--dir", temp_dir],
      work_dir,
      [],
    )
  case res {
    Ok(out) -> io.println(out)
    Error(#(status, err)) ->
      io.println("Error " <> int.to_string(status) <> ": " <> err)
  }

  let files = list_all_files(temp_dir)

  io.println("Downloaded files:")
  files
  |> list.each(fn(path) {
    let is_current = case extract_asset_name(path) {
      Some(n) -> n == current_asset
      _ -> False
    }
    case is_current {
      True -> io.println("  \u{001b}[32m" <> path <> "\u{001b}[0m")
      False -> io.println("  " <> path)
    }
  })

  let reports =
    files
    |> list.filter(fn(p) { string.ends_with(p, "report.md") })
    |> list.sort(fn(a, b) {
      let #(t1, f1) = get_priority(a)
      let #(t2, f2) = get_priority(b)
      case t1 == t2 {
        True -> int.compare(f1, f2)
        False -> int.compare(t1, t2)
      }
    })

  case reports {
    [] -> io.println("No report.md found")
    _ -> {
      let final_report = concat_reports(reports)
      io.println("\n" <> final_report)
    }
  }
}

fn get_priority(path: String) -> #(Int, Int) {
  case extract_asset_name(path) {
    None -> #(999, 999)
    Some(name) -> {
      let tier =
        list.key_find(systems, name)
        |> result.map(fn(sys) { sys.tier })
        |> result.unwrap(999)
      let favor = case string.ends_with(name, "Linux") {
        True -> 0
        False -> 1
      }
      #(tier, favor)
    }
  }
}

fn concat_reports(paths: List(String)) -> String {
  paths
  |> list.index_map(fn(path, i) {
    let content = simplifile.read(path) |> result.unwrap("")
    let parts = string.split(content, "---")

    let header = list.first(parts) |> result.unwrap("")
    let bodies = list.drop(parts, 1) |> string.join("---")

    let header_part = case i {
      0 -> header
      _ -> ""
    }
    let body_part = case bodies == "" {
      True -> ""
      False ->
        case i {
          0 -> "---" <> bodies
          _ -> "\n---" <> bodies
        }
    }
    header_part <> body_part
  })
  |> string.join("")
}
