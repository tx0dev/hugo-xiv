# Hugo XIV Module

This is a Hugo module that provides some Final Fantasy XIV related shortcodes and data.

## Installation

In your SCSS file, generally `assets/main.scss`, check your theme setting, add the following line:

```scss
@import "scss/xiv.scss";
```
There's no styling yet.


for the toolip to work, add the following somewhere in the body of your HTML:

```go
{{- partials "xiv/eorzea_db.html" . -}}
```
The partial will only inject the [Eorzea Database javascript](https://na.finalfantasyxiv.com/lodestone/special/fankit/tooltip/) when a `xiv` shortcode is present.

## Features

### Currency

Draw a simple box with the current icon. Take 2 arguments, the amount and the currency name. current default to `gil` if not provided. Check the [Currency data file](data/xiv_currencies.json) for the full list of currency.

The Society currency are also added, check the [Societies data file](data/xiv_societies.json). Use the *simple* name of the society.

```markdown
{{< currency 000 vanuvanu >}}
{{< tomestone 100 poetics >}}
{{< scrip 100 purple crafter >}}
```

### Game Reference

A simple reference to a quest, duty or anything in the Eorzea DB. It can be inline or wrapping.

```markdown
{{< xiv quest=msq db="quest:7a0da925036" >}}Dawntrail{{< /xiv >}}
```
If `db` is specified, it will automatically link to the [Eorzea Database](https://eu.finalfantasyxiv.com/lodestone/playguide/db/).

### Map

You can embed a map, *with custom marker* using the `map` shortcode. Two way to use it:
* Just the map with a self closing tag
* or with markers.
Aetheryte and shards are automatically added.

The marker needs the *in game* `x` and `y` coordinates, `text` is optional (but recommended), and an optional `img` attribute with the name of a marker. see the `maker` dict in [XIV Icons](data/xiv_icons.json).

```tpl
{{< map "Tuliyollal" />}}
{{< map "The Rak'tika Greatwood" >}}
{{< marker x=10,3 y=5.4 text="My Marker" >}}
{{< /map >}}
```

The map data is from TeamCraft and XIVAPI. The map image is from XIVAPI.

## Roadmap

Functionality are added they come up. Here are some of the things I want to add:

- [x] Get the [FFXIV Tooltip](https://eu.finalfantasyxiv.com/lodestone/special/fankit/tooltip/) working
- [x] "Standardize" link references
- [x] Better currency/token support using data files
- [x] Map support
  - [ ] Better marker for fate and mob area
- [ ] Support for PhotoSwipe or something similar (No jquery)
- [x] SBOM for SQEX assets
- [ ] Improve gallery shortcode

## Square Enix Copyright

FINAL FANTASY is a registered trademark of Square Enix Holdings Co., Ltd. \
Some images, contents, and assets © SQUARE ENIX \
Used under the [Material Usage License](https://support.na.square-enix.com/rule.php?id=5382&tag=authc).
