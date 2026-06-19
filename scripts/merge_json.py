#!/usr/bin/env python3
import os
import sys
import json
import re
import argparse
from datetime import datetime

# Helper to log to stderr
def log_warning(msg):
    sys.stderr.write(f"WARNING: {msg}\n")

def log_info(msg):
    sys.stderr.write(f"INFO: {msg}\n")

def clean_comments_and_commas(content):
    # Remove javascript-style comments // and /* ... */
    content_no_comments = re.sub(r'//.*', '', content)
    content_no_comments = re.sub(r'/\*.*?\*/', '', content_no_comments, flags=re.DOTALL)
    # Remove trailing commas before closing braces/brackets
    content_clean = re.sub(r',\s*([\]}])', r'\1', content_no_comments)
    return content_clean

def salvage_corrupted_json(content):
    salvaged = {}
    # Search for "model": "something"
    model_match = re.search(r'"model"\s*:\s*"([^"]+)"', content)
    if model_match:
        salvaged["model"] = model_match.group(1)
    
    # Search for "env" block
    env_match = re.search(r'"env"\s*:\s*\{([^}]+)\}', content)
    if env_match:
        env_content = env_match.group(1)
        env_dict = {}
        for kv_match in re.finditer(r'"([^"]+)"\s*:\s*"([^"]+)"', env_content):
            env_dict[kv_match.group(1)] = kv_match.group(2)
        if env_dict:
            salvaged["env"] = env_dict

    # Other top-level simple key-value pairs
    for kv_match in re.finditer(r'"([^"]+)"\s*:\s*("[^"]*"|\d+|true|false|null)', content):
        key = kv_match.group(1)
        val_str = kv_match.group(2)
        if key not in ["env", "mcpServers", "model"]:
            try:
                salvaged[key] = json.loads(val_str)
            except Exception:
                pass
                
    return salvaged

def correct_model_names(config):
    # Model update mappings — always keep current
    model_mappings = {
        "claude-3-5-sonnet-20241022": "claude-sonnet-4-6",
        "claude-3-5-sonnet-latest": "claude-sonnet-4-6",
        "claude-3.5-sonnet": "claude-sonnet-4-6",
    }
    subagent_mappings = {
        "claude-3-5-haiku-20241022": "claude-haiku-4-5-20251001",
        "claude-3-5-haiku-latest": "claude-haiku-4-5-20251001",
        "claude-3.5-haiku": "claude-haiku-4-5-20251001",
    }
    
    # Auto-correct main model
    if "model" in config and isinstance(config["model"], str):
        old_model = config["model"]
        if old_model in model_mappings:
            config["model"] = model_mappings[old_model]
            log_info(f"Updated main model from '{old_model}' to '{config['model']}'")
        
    # Auto-correct subagent model in env
    if "env" in config and isinstance(config["env"], dict):
        subagent_key = "CLAUDE_CODE_SUBAGENT_MODEL"
        if subagent_key in config["env"] and isinstance(config["env"][subagent_key], str):
            old_sub = config["env"][subagent_key]
            if old_sub in subagent_mappings:
                config["env"][subagent_key] = subagent_mappings[old_sub]
                log_info(f"Updated subagent model from '{old_sub}' to '{config['env'][subagent_key]}'")
                
    return config

def main():
    parser = argparse.ArgumentParser(description="Merge and correct Claude Code JSON configuration files.")
    parser.add_argument("--file", default="~/.claude.json", help="Path to the JSON settings file")
    parser.add_argument("--config-type", default="settings", choices=["settings", "mcp", "generic"], help="Configuration type")
    parser.add_argument("--model", help="Set the main model name")
    parser.add_argument("--update-env", nargs="*", help="Updates for the env dictionary in KEY=VALUE format")
    parser.add_argument("--auto-correct", action="store_true", help="Auto-correct/upgrade older model names")
    parser.add_argument("--dry-run", action="store_true", help="Print the resulting JSON to stdout instead of writing to disk")
    
    args = parser.parse_args()
    
    file_path = os.path.abspath(os.path.expanduser(args.file))
    
    config = {}
    file_exists = os.path.exists(file_path)
    file_empty = False
    
    if file_exists:
        if os.path.getsize(file_path) == 0:
            file_empty = True
            log_info(f"File {file_path} is empty. Initializing new configuration.")
        else:
            try:
                with open(file_path, "r", encoding="utf-8") as f:
                    content = f.read()
            except Exception as e:
                log_warning(f"Failed to read file {file_path}: {e}. Starting fresh.")
                content = ""
                file_empty = True
            
            if not file_empty:
                # Try parsing direct JSON
                try:
                    config = json.loads(content)
                    if not isinstance(config, dict):
                        log_warning("JSON root is not an object. Starting fresh.")
                        config = {}
                except json.JSONDecodeError:
                    log_warning("Failed to parse JSON directly. File is malformed/corrupted.")
                    # Backup the corrupted file immediately
                    try:
                        backup_path = f"{file_path}.{datetime.now().strftime('%Y%m%d%H%M%S')}.bak"
                        with open(backup_path, "w", encoding="utf-8") as bf:
                            bf.write(content)
                        log_info(f"Backed up corrupted file to {backup_path}")
                    except Exception as be:
                        log_warning(f"Could not write backup file: {be}")
                    
                    # Try cleaning comments/commas
                    log_warning("Attempting to clean comments and commas...")
                    try:
                        clean_content = clean_comments_and_commas(content)
                        config = json.loads(clean_content)
                        if not isinstance(config, dict):
                            config = {}
                    except json.JSONDecodeError:
                        log_warning("Cleaned JSON parsing failed. Attempting to salvage keys...")
                        config = salvage_corrupted_json(content)
                            
    else:
        log_info(f"File {file_path} does not exist. Will create it.")
        
    # Ensure config is a dictionary
    if not isinstance(config, dict):
        config = {}
        
    # Apply updates
    if args.model:
        config["model"] = args.model
        
    if args.update_env:
        if "env" not in config or not isinstance(config["env"], dict):
            config["env"] = {}
        for item in args.update_env:
            if "=" in item:
                k, v = item.split("=", 1)
                config["env"][k.strip()] = v.strip()
            else:
                log_warning(f"Ignoring invalid env update format: {item}. Expected KEY=VALUE.")
                
    # Auto-correct model names if requested or if it is settings and auto-correct is implicitly active
    if args.auto_correct or args.config_type == "settings":
        config = correct_model_names(config)
        
    # Output / Save
    json_output = json.dumps(config, indent=2, ensure_ascii=False)
    
    if args.dry_run:
        try:
            print(json_output)
        except UnicodeEncodeError:
            if hasattr(sys.stdout, "buffer"):
                sys.stdout.buffer.write(json_output.encode("utf-8"))
                sys.stdout.buffer.write(b"\n")
            else:
                print(json_output.encode("utf-8", errors="replace").decode(errors="replace"))
    else:
        # Create parent directory if missing
        dir_name = os.path.dirname(file_path)
        if dir_name and not os.path.exists(dir_name):
            try:
                os.makedirs(dir_name, exist_ok=True)
            except Exception as e:
                log_warning(f"Failed to create directory {dir_name}: {e}")
                
        try:
            with open(file_path, "w", encoding="utf-8") as f:
                f.write(json_output)
            log_info(f"Successfully updated configuration file: {file_path}")
        except Exception as e:
            log_warning(f"Failed to write to file {file_path}: {e}")
            sys.exit(1)

if __name__ == "__main__":
    main()
