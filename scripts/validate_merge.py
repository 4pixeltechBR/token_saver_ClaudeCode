#!/usr/bin/env python3
import os
import sys
import json
import shutil
import glob
import subprocess
import unittest

class TestMergeJson(unittest.TestCase):
    def setUp(self):
        self.test_dir = os.path.abspath(os.path.join(os.path.dirname(__file__), "temp_test_run"))
        os.makedirs(self.test_dir, exist_ok=True)
        self.merge_script = os.path.abspath(os.path.join(os.path.dirname(__file__), "merge_json.py"))
        self.python_cmd = sys.executable

    def tearDown(self):
        if os.path.exists(self.test_dir):
            shutil.rmtree(self.test_dir)

    def run_merge(self, file_path, args_list):
        cmd = [self.python_cmd, self.merge_script, "--file", file_path] + args_list
        result = subprocess.run(cmd, capture_output=True, text=True)
        return result

    def test_1_file_missing(self):
        """Test Case 1: File does not exist."""
        test_file = os.path.join(self.test_dir, "non_existent.json")
        if os.path.exists(test_file):
            os.remove(test_file)

        result = self.run_merge(test_file, [
            "--model", "claude-sonnet-4-6",
            "--update-env", "ENABLE_TOOL_SEARCH=true", "KEY1=VAL1"
        ])
        
        self.assertEqual(result.returncode, 0, f"Script failed with: {result.stderr}")
        self.assertTrue(os.path.exists(test_file), "Output file was not created")

        with open(test_file, "r", encoding="utf-8") as f:
            data = json.load(f)

        self.assertEqual(data.get("model"), "claude-sonnet-4-6")
        self.assertIn("env", data)
        self.assertEqual(data["env"].get("ENABLE_TOOL_SEARCH"), "true")
        self.assertEqual(data["env"].get("KEY1"), "VAL1")
        print("Test Case 1 (File Missing) passed.")

    def test_2_file_empty(self):
        """Test Case 2: File is empty (0 bytes)."""
        test_file = os.path.join(self.test_dir, "empty_file.json")
        with open(test_file, "w", encoding="utf-8") as f:
            f.write("")

        result = self.run_merge(test_file, [
            "--model", "claude-sonnet-4-6",
            "--update-env", "ENABLE_TOOL_SEARCH=true"
        ])

        self.assertEqual(result.returncode, 0, f"Script failed with: {result.stderr}")
        
        with open(test_file, "r", encoding="utf-8") as f:
            data = json.load(f)

        self.assertEqual(data.get("model"), "claude-sonnet-4-6")
        self.assertEqual(data.get("env", {}).get("ENABLE_TOOL_SEARCH"), "true")
        print("Test Case 2 (File Empty) passed.")

    def test_3_valid_json_merge(self):
        """Test Case 3: Syntactically valid JSON with existing settings to preserve."""
        test_file = os.path.join(self.test_dir, "valid.json")
        initial_config = {
            "model": "claude-3-5-sonnet-20241022",
            "theme": "dark",
            "env": {
                "SOME_CUSTOM_VAR": "preserved_val",
                "CLAUDE_CODE_SUBAGENT_MODEL": "claude-3-5-haiku-20241022"
            },
            "mcpServers": {
                "gcal": {
                    "command": "npx",
                    "args": ["-y", "@modelcontextprotocol/server-gcal"]
                }
            }
        }
        with open(test_file, "w", encoding="utf-8") as f:
            json.dump(initial_config, f, indent=2)

        result = self.run_merge(test_file, [
            "--update-env", "ENABLE_TOOL_SEARCH=true"
        ])

        self.assertEqual(result.returncode, 0, f"Script failed with: {result.stderr}")

        with open(test_file, "r", encoding="utf-8") as f:
            data = json.load(f)

        self.assertEqual(data.get("model"), "claude-sonnet-4-6")
        self.assertEqual(data.get("theme"), "dark")
        self.assertEqual(data.get("env", {}).get("CLAUDE_CODE_SUBAGENT_MODEL"), "claude-haiku-4-5-20251001")
        self.assertEqual(data.get("env", {}).get("SOME_CUSTOM_VAR"), "preserved_val")
        self.assertEqual(data.get("env", {}).get("ENABLE_TOOL_SEARCH"), "true")
        self.assertIn("mcpServers", data)
        self.assertIn("gcal", data["mcpServers"])
        print("Test Case 3 (Valid JSON Merge & Model Correction) passed.")

    def test_4_corrupted_json(self):
        """Test Case 4: Malformed/corrupted JSON file (should recover salvageable keys & backup)."""
        test_file = os.path.join(self.test_dir, "corrupted.json")
        corrupted_content = """
        {
            // Old configuration comments
            "model": "claude-3-5-sonnet-20241022",
            "theme": "light",
            "env": {
                "PRESERVE_ME": "yes",
                "CLAUDE_CODE_SUBAGENT_MODEL": "claude-3-5-haiku-20241022",
            },
            "mcpServers": {
                "test": {
                    "command": "node",
                }
            }
        }
        """
        with open(test_file, "w", encoding="utf-8") as f:
            f.write(corrupted_content)

        result = self.run_merge(test_file, [
            "--update-env", "ENABLE_TOOL_SEARCH=true"
        ])

        self.assertEqual(result.returncode, 0, f"Script failed with: {result.stderr}")

        backups = glob.glob(f"{test_file}.*.bak")
        self.assertTrue(len(backups) > 0, "No backup file was created for corrupted configuration")

        with open(test_file, "r", encoding="utf-8") as f:
            data = json.load(f)

        self.assertEqual(data.get("model"), "claude-sonnet-4-6")
        self.assertEqual(data.get("theme"), "light")
        self.assertEqual(data.get("env", {}).get("CLAUDE_CODE_SUBAGENT_MODEL"), "claude-haiku-4-5-20251001")
        self.assertEqual(data.get("env", {}).get("PRESERVE_ME"), "yes")
        self.assertEqual(data.get("env", {}).get("ENABLE_TOOL_SEARCH"), "true")
        print("Test Case 4 (Corrupted JSON Salvage & Backup) passed.")

if __name__ == "__main__":
    unittest.main()
