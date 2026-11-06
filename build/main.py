# main.py
import os
import sys
import argparse

from app.android import AndroidBuilder
from app.apple_go import AppleGoBuilder
from app.apple_gomobile import AppleGoMobileBuilder
from app.linux import LinuxBuilder
from app.windows import WindowsBuilder

LOCAL_ARG = "local"


def build_dir_path():
    file_dir = os.path.dirname(__file__)
    dir_path = os.path.abspath(file_dir)
    return dir_path


def parse_local_arg(args: list[str]) -> bool:
    if not args:
        return False
    if args == [LOCAL_ARG]:
        return True
    raise Exception(f"unsupported args: {args}")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(add_help=False)
    parser.add_argument("--xray-core", help="Explicit Xray-core source directory")
    parser.add_argument("--skip-geo", action="store_true", help="Build libraries without downloading runtime geodata")
    options, remaining = parser.parse_known_args()
    if options.xray_core:
        os.environ["LIBXRAY_CORE_DIR"] = os.path.abspath(options.xray_core)
    if options.skip_geo:
        os.environ["LIBXRAY_SKIP_GEO"] = "1"
    sys.argv = [sys.argv[0], *remaining]
    print(sys.argv)
    platform = sys.argv[1]

    if platform == "apple":
        tool = sys.argv[2]
        use_local_xray_core = parse_local_arg(sys.argv[3:])
        if tool == "go":
            builder = AppleGoBuilder(build_dir_path(), use_local_xray_core)
            builder.build()
        elif tool == "gomobile":
            builder = AppleGoMobileBuilder(build_dir_path(), use_local_xray_core)
            builder.build()
        else:
            raise Exception(f"platform {platform} tool {tool} not supported")

    elif platform == "android":
        use_local_xray_core = parse_local_arg(sys.argv[2:])
        builder = AndroidBuilder(build_dir_path(), use_local_xray_core)
        builder.build()

    elif platform == "linux":
        use_local_xray_core = parse_local_arg(sys.argv[2:])
        builder = LinuxBuilder(build_dir_path(), use_local_xray_core)
        builder.build()

    elif platform == "windows":
        use_local_xray_core = parse_local_arg(sys.argv[2:])
        builder = WindowsBuilder(build_dir_path(), use_local_xray_core)
        builder.build()

    else:
        raise Exception(f"platform {platform} not supported")
