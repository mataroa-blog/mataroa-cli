# Dependencies

## Dependency Policy

The mataroa-cli project has an unusually strict yet usually unclear dependency policy.

Vague rules include:

* No third-party Django apps.
* All Go packages should be individually vetted.
    * Packages should be published from community-trusted organisations or developers.
    * Packages should be actively maintained (though not necessarily actively developed).
    * Packages should hold a high quality of coding practices.

## Adding a new dependency

After approving a dependency, the process to add it is:

`go get example.com/dependency-x`

## Upgrading dependencies

When a new libraries versions are out it’s a good idea to upgrade everything.

Steps:

1. `go get -u`
2. `go mod tidy`
