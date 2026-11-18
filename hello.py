#!/usr/bin/env python3
"""Print a friendly hello greeting.

Usage:
    python3 hello.py            # Hello, Python!
    python3 hello.py 小明       # Hello, 小明!
"""

from __future__ import annotations

import argparse
import sys
from typing import List, Optional

DEFAULT_NAME = "Python"


def greeting(name: Optional[str] = None) -> str:
    """Return the greeting line for ``name``.

    An empty or whitespace-only name falls back to :data:`DEFAULT_NAME` so the
    message always addresses someone.
    """
    cleaned = (name or "").strip()
    return "Hello, {}!".format(cleaned or DEFAULT_NAME)


def main(argv: Optional[List[str]] = None) -> int:
    """Print the greeting and return a process exit code."""
    parser = argparse.ArgumentParser(description="Print a friendly hello greeting.")
    parser.add_argument(
        "name",
        nargs="?",
        default=DEFAULT_NAME,
        help='name to greet (default "{}")'.format(DEFAULT_NAME),
    )
    args = parser.parse_args(argv)
    print(greeting(args.name))
    return 0


if __name__ == "__main__":
    sys.exit(main())
