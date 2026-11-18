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


class GreetingsTests(unittest.TestCase):
    def test_default_greetings_include_every_line(self):
        self.assertEqual(
            hello.greetings(),
            ["Hello, Python!", "Hello, every!"],
        )

    def test_named_greetings_include_every_line(self):
        self.assertEqual(
            hello.greetings("World"),
            ["Hello, World!", "Hello, every!"],
        )

    def test_every_greeting_constant(self):
        self.assertEqual(hello.EVERY_GREETING, "Hello, every!")


class MainTests(unittest.TestCase):
    def _run(self, argv):
        buf = io.StringIO()
        with redirect_stdout(buf):
            code = hello.main(argv)
        return code, buf.getvalue().strip()

    def test_main_prints_default_greetings(self):
        code, out = self._run([])
        self.assertEqual(code, 0)
        self.assertEqual(out, "Hello, Python!\nHello, every!")

    def test_main_prints_given_name(self):
        code, out = self._run(["World"])
        self.assertEqual(code, 0)
        self.assertEqual(out, "Hello, World!\nHello, every!")

    def test_main_accepts_dash_name(self):
        code, out = self._run(["--", "-x"])
        self.assertEqual(code, 0)
        self.assertEqual(out, "Hello, -x!\nHello, every!")


if __name__ == "__main__":
    unittest.main()
