workspace(name = "com_zombiezen_go_capnproto2")

http_archive(
    name = "io_bazel_rules_go",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.7.1/rules_go-0.7.1.tar.gz",
    sha256 = "341d5eacef704415386974bc82a1783a8b7ffbff2ab6ba02375e1ca20d9b031c",
)
load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains", "go_repository")
go_rules_dependencies()
go_register_toolchains()

go_repository(
    name = "com_github_kylelemons_godebug",
    importpath = "github.com/kylelemons/godebug",
    sha256 = "4415b09bae90e41695bc17e4d00d0708e1f6bbb6e21cc22ce0146a26ddc243a7",
    strip_prefix = "godebug-a616ab194758ae0a11290d87ca46ee8c440117b0",
    urls = [
        "https://github.com/kylelemons/godebug/archive/a616ab194758ae0a11290d87ca46ee8c440117b0.zip",
    ],
)

go_repository(
    name = "org_golang_x_net",
    importpath = "golang.org/x/net",
    sha256 = "880dc04d0af397dce6875ee2349bbb4295fe5a47352f7a4da4270456f726edd4",
    strip_prefix = "net-f5079bd7f6f74e23c4d65efa0f4ce14cbd6a3c0f",
    urls = [
        "https://github.com/golang/net/archive/f5079bd7f6f74e23c4d65efa0f4ce14cbd6a3c0f.zip",
    ],
)
