#!/usr/bin/env python3
"""Build standalone release packages. Development only; no third-party packages."""
import argparse
import hashlib
import os
from pathlib import Path
import subprocess
import zipfile

VERSION = "3.0.0"
ROOT = Path(__file__).resolve().parents[1]
TARGETS = [(system, arch) for system in ("windows", "darwin", "linux") for arch in ("amd64", "arm64")]


def skill_entries(prefix):
    entries = {}
    for file in (ROOT / "skill").rglob("*"):
        if file.is_file():
            entries[f"{prefix}/{file.relative_to(ROOT / 'skill').as_posix()}"] = file.read_bytes()
    return entries


def write_zip(target, root, entries):
    with zipfile.ZipFile(target, "w", zipfile.ZIP_DEFLATED, compresslevel=9) as archive:
        for path, content in entries.items():
            item = zipfile.ZipInfo(root + "/" + path)
            item.create_system = 3
            item.external_attr = 0o100644 << 16
            item.compress_type = zipfile.ZIP_DEFLATED
            archive.writestr(item, content)


def make_plugin_packages(out):
    plugin_manifest = ('{"$schema":"https://agent-plugins.org/schemas/1.0.0/plugin.schema.json",'
                       f'"name":"token-saver","version":"{VERSION}",'
                       '"description":"Safe, read-only context auditing across coding harnesses",'
                       '"license":"MIT"}\n').encode()
    minimax_manifest = (f'{{"name":"token-saver","version":"{VERSION}",'
                        '"skills":["./skills/SKILL.md"]}\n').encode()
    plugin_readme = f"""# Token Saver plugin {VERSION}

This package contains the portable Agent Skill for hosts that support the
Agent Skills or Agent Plugins format. It provides discovery and read-only
guidance. Install the matching Token Saver release package to add the native
companion CLI for deterministic JSON audit, plan, apply, and rollback.

The plugin never changes model, reasoning effort, permissions, or context
automatically. Read the repository README and skill references before enabling
shell execution for a third-party plugin.
""".encode()
    base = {
        "plugin.json": plugin_manifest,
        ".claude-plugin/plugin.json": minimax_manifest,
        "README.md": plugin_readme,
        "LICENSE": (ROOT / "LICENSE").read_bytes(),
    }
    base.update(skill_entries("skills/token-saver"))
    for source in (ROOT / "skill").rglob("*"):
        if source.is_file():
            relative = source.relative_to(ROOT / "skill").as_posix()
            base["skills/" + relative] = source.read_bytes()
    target = out / f"token-saver-{VERSION}-agent-plugin.zip"
    write_zip(target, f"token-saver-{VERSION}-agent-plugin", base)
    return target


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--out", type=Path, default=ROOT / "dist")
    parser.add_argument("--go", default="go")
    parser.add_argument("--target", help="Optional single target, e.g. windows-amd64")
    args = parser.parse_args()
    out = args.out.resolve()
    out.mkdir(parents=True, exist_ok=True)
    package_hashes = []
    for system, arch in TARGETS:
        if args.target and args.target != f"{system}-{arch}":
            continue
        name = f"token-saver-{VERSION}-{system}-{arch}"
        stage = out / ".build" / f"{system}-{arch}"
        stage.mkdir(parents=True, exist_ok=True)
        if not stage.resolve().is_relative_to(out):
            raise RuntimeError("Build staging path escaped the output directory")
        binary_name = "token-saver.exe" if system == "windows" else "token-saver"
        binary = Path(stage) / binary_name
        env = dict(os.environ, GOOS=system, GOARCH=arch, CGO_ENABLED="0", GOTOOLCHAIN="local")
        subprocess.run([args.go, "build", "-trimpath", "-buildvcs=false", "-ldflags=-s -w", "-o", str(binary), "."], cwd=ROOT, env=env, check=True)
        data = binary.read_bytes()
        checksum = hashlib.sha256(data).hexdigest()
        entries = {"bin/" + binary_name: data,
                   "SHA256SUMS.txt": f"{checksum}  bin/{binary_name}\n".encode(),
                   "README.md": (ROOT / "README.md").read_bytes(),
                   "README.en.md": (ROOT / "README.en.md").read_bytes(),
                   "LICENSE": (ROOT / "LICENSE").read_bytes()}
        for file in ("install.ps1", "instalar.cmd") if system == "windows" else ("install.sh",):
            entries[file] = (ROOT / file).read_bytes()
        entries.update(skill_entries("skill"))
        target = out / (name + ".zip")
        with zipfile.ZipFile(target, "w", zipfile.ZIP_DEFLATED, compresslevel=9) as archive:
            for path, content in entries.items():
                item = zipfile.ZipInfo(name + "/" + path)
                item.create_system = 3
                item.external_attr = (0o100755 if path.startswith("bin/") or path.endswith(".sh") else 0o100644) << 16
                item.compress_type = zipfile.ZIP_DEFLATED
                archive.writestr(item, content)
        package_hashes.append(f"{hashlib.sha256(target.read_bytes()).hexdigest()}  {target.name}")
        print(target.name, flush=True)
    if not package_hashes:
        raise SystemExit("Unknown target")
    plugin = make_plugin_packages(out)
    package_hashes.append(f"{hashlib.sha256(plugin.read_bytes()).hexdigest()}  {plugin.name}")
    (out / "SHA256SUMS.txt").write_text("\n".join(package_hashes) + "\n", encoding="utf-8")


if __name__ == "__main__":
    main()
