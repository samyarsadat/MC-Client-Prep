<h1 align="center">Minecraft Server Manifest Specification</h1>
<p align="center"><i>Version 1.0 - <strong>Written by Samyar Sadat Akhavi</strong> - Doc. Ver. 2025.4.3</i></p>

----
This document specifies the format of the Minecraft Server Manifest File used by the Minecraft Client Preparation Program to prepare a Minecraft client based upon the server's requirements. The file should be in the JSON format and publicly accessible at a static URL.

The file should provide the following information regarding the server:
1. The server's address (either IP address or domain name)
2. The server's port number
3. The Minecraft version of the server
4. The mod-loader type (if any) used by the server
5. The mod-loader version (if any) used by the server
6. A list of mods (if any) to be downloaded onto the client
7. A list of resource packs (if any) to be downloaded on the client
8. A list of shader packs (if any) to be downloaded on the client

More information regarding these parameters and valid values is provided below.

<br>

## Overall Structure
The overall structure of the JSON file is as follows:
```json
{
  "server_manifest_version": 1.0,
  "datetime": "2025-04-03T17:05:00Z",
  "server_info": {
    "address": null,
    "port": 25565,
    "name": "Samyar's Minecraft Server",
    "version": "1.20.1",
    "modloader": "forge",
    "modloader_version": "latest"
  },
  "mod_manifest": {
    "curseforge": [
      "example-project-slug"
    ],
    "modrinth": [
      "example-project-slug:beta"
    ],
    "url_download": [
      {
        "name": "Example Mod",
        "url": "https://example.com",
        "sha1": ""
      }
    ]
  },
  "resourcepack_manifest": {
    "vanillatweaks": [],
    "curseforge": [],
    "modrinth": [],
    "url_download": []
  },
  "shaderpack_manifest": {
    "curseforge": [],
    "modrinth": [],
    "url_download": []
  }
}
```

<br>

## Detailed Description
This is a detailed description of the JSON file. After the name of each key, there is information regarding the expected type, and whether the key must be populated with a value (type = type, req = required). Keys which are not required to be populated can be left empty, populated with `null` or omitted entirely.
- `server_manifest_version` [type: `float`, req: **yes**]: this is the specification version that the file conforms to.
- `datetime` [type: `string`, req: **yes**]: the time and date at which the file was last updated (formatted as `YYYY-MM-DDTHH:MM:SSZ`, ISO 8601).
- `server_info` [type: `dict`, req: **yes**]: basic information regarding the server.
    - `address` [type: `string`, req: **no**]: the server's address, either IP address or domain name.
    - `port` [type: `int`, req: **yes**]: the server's port number.
    - `name` [type: `string`, req: **no**]: optional display name for server.
    - `version` [type: `string`, req: **yes**]: the minecraft version used by the server (formatted as `1.xx.x`).
    - `modloader` [type: `string`, req: **yes**]: the mod-loader (if any) used by the server (supported values listed below).
    - `modloader_version` [type: `string`, req: **yes** - if mod-loader is non-vanilla]: the mod-loader version used by the server (see note below).
- `mod_manifest` [type: `dict`, req: **no**]: list of mods to be installed on the client.
    - `curseforge` [type: `list[string]`, req: **no**]: list of mods to be downloaded from CurseForge (see format below).
    - `modrinth` [type: `list[string]`, req: **no**]: list of mods to be downloaded from Modrinth (see format below).
    - `url_download` [type: `list[dict]`, req: **no**]: list of mods to be downloaded a given file URL.
        - `name` [type: `string`, req: **no**]: optional display name for file.
        - `url` [type: `string`, req: **yes**]: the URL to the file.
        - `sha1` [type: `string`, req: **no**]: optional `SHA1` hash of the file for verification.
- `resourcepack_manifest` [type: `dict`, req: **no**]: list of resource packs to be installed on the client.
    - `vanillatweaks` [type: `list[string]`, req: **no**]: list VanillaTweaks resource pack share-codes to be downloaded.
    - `curseforge` [type: `list[string]`, req: **no**]: list of resource packs to be downloaded from CurseForge (see format below).
    - `modrinth` [type: `list[string]`, req: **no**]: list of resource packs to be downloaded from Modrinth (see format below).
    - `url_download` [type: `list[dict]`, req: **no**]: list of resource packs to be downloaded a given file URL.
        - `name` [type: `string`, req: **no**]: optional display name for file.
        - `url` [type: `string`, req: **yes**]: the URL to the file.
        - `sha1` [type: `string`, req: **no**]: optional `SHA1` hash of the file for verification.
- `shaderpack_manifest` [type: `dict`, req: **no**]: list of shader packs to be installed on the client.
    - `curseforge` [type: `list[string]`, req: **no**]: list of shader packs to be downloaded from CurseForge (see format below).
    - `modrinth` [type: `list[string]`, req: **no**]: list of shader packs to be downloaded from Modrinth (see format below).
    - `url_download` [type: `list[dict]`, req: **no**]: list of shader packs to be downloaded a given file URL.
        - `name` [type: `string`, req: **no**]: optional display name for file.
        - `url` [type: `string`, req: **yes**]: the URL to the file.
        - `sha1` [type: `string`, req: **no**]: optional `SHA1` hash of the file for verification.

> [!NOTE]
> You can set `modloader_version` to `latest`/`recommended` (for Forge and NeoForge) or `latest`/`snapshot` (for Fabric) instead of manually specifying a version number.

<br>

### Supported Mod-loaders
Supported mod-loaders as of version 1.0 are:
1. Forge (value: `"forge"`)
2. NeoForge (value: `"neoforge"`)
3. Fabric (value: `"fabric"`)
4. Vanilla (no mod-loader, value: `"vanilla"`)

<br>

### CurseForge Resource Format
Under the hood, [mc-image-helper](https://github.com/itzg/mc-image-helper) is used to download files from CurseForge. As a result, the required format is the same as with the `mc-image-helper curseforge-files` command.

CurseForge projects (mods, resource packs, shader packs) can be specified with one of the following formats:
- `<project ID>|<slug>:<file ID>`
- `<project ID>|<slug>@<filename matcher>`
- `<project ID>|<slug>`
- Project page URL

Note that unless a specific version is requested, the latest applicable version is automatically selected.

<br>

### Modrinth Resource Format
Under the hood, [mc-image-helper](https://github.com/itzg/mc-image-helper) is used to download files from Modrinth. As a result, the required format is the same as with the `mc-image-helper modrinth` command.

Modrinth projects (mods, resource packs, shader packs) can be specified with one of the following formats:
- `<project ID>|<slug>:<release type>`
- `<project ID>|<slug>:<version name/number>`
- `<project ID>|<slug>`

Note that `<release type>` can only be `release`, `beta`, or `alpha`. Unless a specific version is requested, the latest applicable version is automatically selected.

<br>

Copyright © 2025, Samyar Sadat Akhavi.\
Document version 2025.4.3 for specification version 1.0.