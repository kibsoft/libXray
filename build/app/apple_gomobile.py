import os
import subprocess

from app.build import Builder


class AppleGoMobileBuilder(Builder):
    def before_build(self):
        super().before_build()
        self.clean_lib_dirs(["LibXray.xcframework"])
        self.prepare_gomobile()

    def build(self):
        self.before_build()

        os.chdir(self.lib_dir)
        ret = subprocess.run(
            [
                "gomobile",
                "bind",
                "-target",
                "ios/arm64,iossimulator/arm64,iossimulator/amd64,macos/arm64,macos/amd64",
                "-ldflags=-checklinkname=0",
                "-iosversion",
                "15.0",
            ]
        )
        if ret.returncode != 0:
            raise Exception("build failed")

        self.after_build()

        self.revert_go_env()
