#!/usr/bin/env python3
"""Run real terminal launch tests in an isolated controlling PTY (Unix only)."""
import os
from pathlib import Path
import pty
import select
import signal
import subprocess
import tempfile
import time


def main():
    root = Path(__file__).resolve().parents[1]
    with tempfile.TemporaryDirectory(prefix="parley-terminal-") as scratch:
        binary = str(Path(scratch) / "runner.test")
        subprocess.run(["go", "test", "-c", "-o", binary, "./internal/runner"], cwd=root, check=True)
        pid, fd = pty.fork()
        if pid == 0:
            os.chdir(root)
            os.environ["PARLEY_TTY_TEST"] = "1"
            os.execv(binary, [binary, "-test.run=^TestInteractiveRealTerminal", "-test.timeout=45s", "-test.v"])
        output = bytearray()
        sent_child = sent_parent = 0
        deadline = time.monotonic() + 55
        status = None
        try:
            while time.monotonic() < deadline:
                if select.select([fd], [], [], 0.1)[0]:
                    try:
                        output.extend(os.read(fd, 65536))
                    except OSError:
                        pass
                    children = output.count(b"TTY_CHILD_READY")
                    parents = output.count(b"TTY_PARENT_READY")
                    while sent_child < children:
                        os.write(fd, b"child\n")
                        sent_child += 1
                    while sent_parent < parents:
                        os.write(fd, b"parent\n")
                        sent_parent += 1
                reaped, result = os.waitpid(pid, os.WNOHANG)
                if reaped:
                    status = result
                    break
            if status is None:
                os.kill(pid, signal.SIGKILL)
                _, status = os.waitpid(pid, 0)
            print(output.decode(errors="replace"))
            code = os.waitstatus_to_exitcode(status)
            if code == 0 and (sent_child != 1 or sent_parent != 5):
                raise SystemExit("PTY tests did not execute every expected terminal interaction")
            raise SystemExit(0 if code == 0 else 1)
        finally:
            os.close(fd)


if __name__ == "__main__":
    main()
