use clap::Parser;
use regex::Regex;
use std::collections::HashMap;
use std::fs;
use std::path::{Path, PathBuf};
use std::process::{Command, Stdio};
use std::sync::OnceLock;
use walkdir::WalkDir;

#[derive(Parser)]
struct Cli {
    run_id: String,
}

static NIX_SYSTEMS: OnceLock<HashMap<&'static str, i32>> = OnceLock::new();
static ASSET_REGEX: OnceLock<Regex> = OnceLock::new();

fn get_nix_systems() -> &'static HashMap<&'static str, i32> {
    NIX_SYSTEMS.get_or_init(|| {
        [
            ("X64-Linux", 1),
            ("ARM64-Linux", 2),
            ("X64-macOS", 2),
            ("ARM64-macOS", 3),
        ]
        .into_iter()
        .collect()
    })
}

fn get_asset_regex() -> &'static Regex {
    ASSET_REGEX.get_or_init(|| Regex::new(r"nixpkgs-review-files-[^/]+").unwrap())
}

fn get_current_asset_name() -> String {
    let arch = match std::env::consts::ARCH {
        "x86_64" => "X64",
        "aarch64" => "ARM64",
        _ => "",
    };
    let os = match std::env::consts::OS {
        "linux" => "Linux",
        "macos" => "macOS",
        _ => "",
    };
    if arch.is_empty() || os.is_empty() {
        String::new()
    } else {
        format!("{}-{}", arch, os)
    }
}

fn extract_asset_name(path: &Path) -> Option<String> {
    let path_str = path.to_string_lossy();
    if let Some(caps) = get_asset_regex().captures(&path_str) {
        let matched = caps.get(0).unwrap().as_str();
        for name in get_nix_systems().keys() {
            if matched.contains(name) {
                return Some(name.to_string());
            }
        }
    }
    None
}

fn get_priority(path: &Path) -> (i32, i32) {
    match extract_asset_name(path) {
        None => (999, 999),
        Some(name) => {
            let systems = get_nix_systems();
            let tier = *systems.get(name.as_str()).unwrap_or(&999);
            let favor = if name.ends_with("Linux") { 0 } else { 1 };
            (tier, favor)
        }
    }
}

fn run_command(cmd: &str, args: &[&str]) -> std::io::Result<()> {
    let status = Command::new(cmd).args(args).status()?;
    if !status.success() {
        return Err(std::io::Error::new(
            std::io::ErrorKind::Other,
            format!("Command failed: {} {:?}", cmd, args),
        ));
    }
    Ok(())
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let cli = Cli::parse();
    let in_runner = std::env::var("GITHUB_ACTIONS").unwrap_or_default() == "true";
    let current_asset = get_current_asset_name();

    if !in_runner {
        let _ = Command::new("gh")
            .args(["run", "watch", &cli.run_id, "--interval", "10"])
            .stdin(Stdio::inherit())
            .stdout(Stdio::inherit())
            .stderr(Stdio::inherit())
            .status();
    }

    let temp_dir = tempfile::Builder::new()
        .prefix(&format!("nixpkgs-reviewing-workspace.run-{}.", cli.run_id))
        .tempdir()?;
    let temp_path = temp_dir.path();

    eprintln!("Downloading artifacts to {}...", temp_path.display());
    run_command(
        "gh",
        &["run", "download", &cli.run_id, "--dir", temp_path.to_str().unwrap()],
    )?;

    let all_files: Vec<PathBuf> = WalkDir::new(temp_path)
        .into_iter()
        .filter_map(|e| e.ok())
        .filter(|e| !e.file_type().is_dir())
        .map(|e| e.path().to_owned())
        .collect();

    eprintln!("Downloaded files:");
    for path in &all_files {
        let is_current = match extract_asset_name(path) {
            Some(n) => n == current_asset,
            None => false,
        };
        if is_current {
            eprintln!("  \x1b[32m{}\x1b[0m", path.display());
        } else {
            eprintln!("  {}", path.display());
        }
    }

    let mut reports: Vec<PathBuf> = all_files
        .into_iter()
        .filter(|p| p.file_name().unwrap_or_default() == "report.md")
        .collect();

    if reports.is_empty() {
        eprintln!("No report.md found in {}", temp_path.display());
        std::process::exit(1);
    }

    reports.sort_by_key(|p| get_priority(p));

    let mut final_report = String::new();
    for (i, path) in reports.iter().enumerate() {
        let content = fs::read_to_string(path)?;
        let parts: Vec<&str> = content.split("---").collect();
        
        if i == 0 {
            if let Some(header) = parts.first() {
                final_report.push_str(header);
            }
        }

        if parts.len() > 1 {
            let body = parts[1..].join("---");
            if i > 0 {
                final_report.push('\n');
            }
            final_report.push_str("---");
            final_report.push_str(&body);
        }
    }

    println!("\n{}", final_report);

    Ok(())
}
