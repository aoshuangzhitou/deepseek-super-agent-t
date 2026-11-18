#!/usr/bin/env python3
"""Tests for hello.py."""

import io
import unittest
from contextlib import redirect_stdout

import hello


class GreetingTests(unittest.TestCase):
    def test_default_greeting(self):
        self.assertEqual(hello.greeting(), "Hello, Python!")

    def test_greeting_with_name(self):
        self.assertEqual(hello.greeting("小明"), "Hello, 小明!")

    def test_blank_name_falls_back_to_default(self):
        self.assertEqual(hello.greeting("   "), "Hello, Python!")


class MainTests(unittest.TestCase):
    def _run(self, argv):
        buf = io.StringIO()
        with redirect_stdout(buf):
            code = hello.main(argv)
        return code, buf.getvalue().strip()

    def test_main_prints_default_greeting(self):
        code, out = self._run([])
        self.assertEqual(code, 0)
        self.assertEqual(out, "Hello, Python!")

    def test_main_prints_given_name(self):
        code, out = self._run(["World"])
        self.assertEqual(code, 0)
        self.assertEqual(out, "Hello, World!")

    def test_main_accepts_dash_name(self):
        code, out = self._run(["--", "-x"])
        self.assertEqual(code, 0)
        self.assertEqual(out, "Hello, -x!")


if __name__ == "__main__":
    unittest.main()
