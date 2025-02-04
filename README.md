# Hugo XIV Module

This is a Hugo module that provides some Final Fantasy XIV related shortcodes and data.

## Installation

In your SCSS file, generally `assets/main.scss`, check your theme setting, add the following line:

```scss
@import "scss/xiv.scss";
```

for the toolip to work, add the following somewhere in the body of your HTML:

```go
{{- partials "xiv/eorzea_db.html" . -}}
```

## Features

### Currency

Draw a simple box with the current icon. Take 2 arguments, the amount and the currency name. current default to `gil` if not provided. check the [shortcodes file](layouts/shortcodes/currency.html) for the full list of currency.

```markdown
{{< currency 000 "poetics" >}}
```

### Game Reference

A simple reference to a quest, duty or anything in the Eorzea DB. It can be inline or wrapping.

```markdown
{{< xiv quest=msq db="quest=7a0da925036" >}}Dawntrail{{< /xiv >}}
```
If `db` is specified, it will automatically link to the [Eorzea Database](https://eu.finalfantasyxiv.com/lodestone/playguide/db/).

## Roadmap

Functionality are added they come up. Here are some of the things I want to add:

- [x] Get the [FFXIV Tooltip](https://eu.finalfantasyxiv.com/lodestone/special/fankit/tooltip/) working
- [ ] "Standardize" link references
- [ ] Better currency/token support using data files
- [ ] Map support (?)
- [ ] Support for PhotoSwipe or something similar (No jquery)
- [ ] SBOM for SQEX assets
- [ ] Improve gallery shortcode

## disclaimer

FINAL FANTASY is a registered trademark of Square Enix Holdings Co., Ltd. \
Some images, contents, and assets © SQUARE ENIX \
Used under the [Material Usage License](https://support.na.square-enix.com/rule.php?id=5382&tag=authc).
