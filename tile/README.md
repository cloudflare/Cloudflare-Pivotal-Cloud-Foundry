# Pivotal Cloud Foundry Tile Generation

## Documentation
http://cf-platform-eng.github.io/isv-portal/tile-generator/

## Setup

The `tile.yml` and `tile-history.yml` files live inside the `src/` directory in this repository but in order to build the tile the structure needs to look like this:

```
tree
.
├── src
│   ├── README.md
│   ├── api
│   ├── broker
│   ├── main
│   ├── main.go
│   └── vendor
└── tile
    ├── resources
    │   └── logo.png
    ├── tile-history.yml
    ├── tile.yml
    └── README.md
```

This is because `tile build` looks at the `packages` `manifest` `path` recursively to generate the tile and it will break if the tile code is inside of the src code. 

After you get mimic this directory structure run `tile build` and the .pivotal file in `tile/product` will be what you upload.