#!/usr/bin/env python3
"""Print friendly hello greetings.

Usage:
    python3 hello.py            # Hello, Python!
                                # Hello, every!
    python3 hello.py 小明       # Hello, 小明!
                                # Hello, every!
"""

from __future__ import annotations

import argparse
import sys
from typing import List, Optional

DEFAULT_NAME = "Python"
EVERY_GREETING = "Hello, every!"


def greeting(name: Optional[str] = None) -> str:
    """Return the greeting line for ``name``.

    An empty or whitespace-only name falls back to :data:`DEFAULT_NAME` so the
    message always addresses someone.
    """
    cleaned = (name or "").strip()
    return "Hello, {}!".format(cleaned or DEFAULT_NAME)


def greetings(name: Optional[str] = None) -> List[str]:
    """Return every output line for ``name``.

    The personal greeting comes first, followed by the "hello every" line.
    """
    return [greeting(name), EVERY_GREETING]


def main(argv: Optional[List[str]] = None) -> int:
    """Print the greetings and return a process exit code."""
    parser = argparse.ArgumentParser(description="Print friendly hello greetings.")
    parser.add_argument(
        "name",
        nargs="?",
        default=DEFAULT_NAME,
        help='name to greet (default "{}")'.format(DEFAULT_NAME),
    )
    args = parser.parse_args(argv)
    for line in greetings(args.name):
        print(line)
    return 0


if __name__ == "__main__":
    sys.exit(main())
