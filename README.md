# Prefine

_**prefine is a CLI tool designed to optimize various file types commonly found in web development projects including images (PNG, JPG, WebP), configuration files (JSON, YAML), and more.**_

## Installation

### Build from source

```
git clone https://github.com/priory-io/prefine
cd prefine
go mod tidy
go build -o prefine cmd/prefine/main.go
./prefine --version
```

### AUR

```
paru -Sy prefine

# or
yay -Sy prefine

# or or
git clone https://aur.archlinux.org/prefine.git
cd prefine
makepkg -si
```

### License

We license all our projects under the [MIT License](./LICENSE)

---

Made with ♥  by [Priory](https://github.com/priory-io)

